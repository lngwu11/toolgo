package connpool

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// returned from a pool connection when the
	// maximum number of connections in the pool has been reached.
	errPoolExhausted = errors.New("the pool has been exhausted")
	errPoolClosed    = errors.New("the pool has been closed")
	errFactoryNotSet = errors.New("the factory not set")
)

type Pool interface {
	ID() string
	Get() (IConn, error)
	GetContext(ctx context.Context) (IConn, error)
	Put(conn IConn) error
	Close()
}

type pool struct {
	id           string        // key of poolMap
	config       *Config       // configs for external
	mu           sync.Mutex    // mu protects the following fields
	closed       bool          // set to true when the pool is closed.
	active       int           // the number of open connections in the pool
	initOnce     sync.Once     // the init ch once func
	ch           chan struct{} // limits open connections when p.Wait is true
	idle         idleList      // idle connections
	waitCount    int64         // total number of connections waited for.
	waitDuration time.Duration // total time waited for new connections.
}

// ID return the pool id.
// it's also the poolMap key
func (p *pool) ID() string {
	return p.id
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (p *pool) Get() (IConn, error) {
	// GetContext returns errorConn in the first argument when an error occurs.
	return p.GetContext(context.Background())
}

// GetContext gets a connection using the provided context.
// The provided Context must be non-nil. If the context expires before the
// connection is complete, an error is returned. Any expiration on the context
// will not affect the returned connection.
// If the function completes without error, then the application must close the
// returned connection.
func (p *pool) GetContext(ctx context.Context) (IConn, error) {
	// Wait until there is a vacant connection in the pool.
	waited, err := p.waitVacantConn(ctx)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()

	if waited > 0 {
		p.waitCount++
		p.waitDuration += waited
	}

	// Prune stale connections at the back of the idle list.
	n := p.idle.count
	for i := 0; i < n && p.idle.back != nil && p.idle.back.Closed(); i++ {
		p.idle.popBack()
		p.active--
	}

	// Get idle connection from the front of idle list.
	for p.idle.front != nil {
		pc := p.idle.front
		p.idle.popFront()
		p.mu.Unlock()
		if pc.Reusable() && !pc.Closed() {
			return pc, nil
		}
		p.mu.Lock()
		p.active--
	}

	// Check for pool closed before dialing a new connection.
	if p.closed {
		p.mu.Unlock()
		return nil, errPoolClosed
	}

	// Handle limit for p.Wait == false.
	if !p.config.Wait && p.config.MaxActive > 0 && p.active >= p.config.MaxActive {
		p.mu.Unlock()
		return nil, errPoolExhausted
	}

	p.active++
	p.mu.Unlock()
	var conn IConn
	conn, err = p.dial(ctx)
	if err != nil {
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		p.mu.Unlock()
		return nil, err
	}
	return conn, nil
}

// Put a connection to the pool. If the idle count more than MaxIdle,
// the connection of list will be pop from back and close.
func (p *pool) Put(conn IConn) error {
	p.mu.Lock()
	if !p.closed {
		p.idle.pushFront(conn)
		if p.idle.count > p.config.MaxIdle {
			conn = p.idle.back
			p.idle.popBack()
		} else {
			conn = nil
		}
	}

	if conn != nil {
		p.mu.Unlock()
		_ = conn.Close()
		p.mu.Lock()
		p.active--
	}

	if p.ch != nil && !p.closed {
		p.ch <- struct{}{}
	}
	p.mu.Unlock()
	return nil
}

// Close the connection and release the resources used by the pool.
func (p *pool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.active -= p.idle.count
	pc := p.idle.front
	p.idle.count = 0
	p.idle.front, p.idle.back = nil, nil
	if p.ch != nil {
		close(p.ch)
	}
	p.mu.Unlock()
	for ; pc != nil; pc = pc.getNext() {
		_ = pc.Close()
	}
}

func newPool(key string, config *Config) Pool {
	return &pool{id: key, config: config}
}

// waitVacantConn waits for a vacant connection in pool if waiting
// is enabled and pool size is limited, otherwise returns instantly.
// If ctx expires before that, an error is returned.
//
// If there were no vacant connection in the pool right away it returns the time spent waiting
// for that connection to appear in the pool.
func (p *pool) waitVacantConn(ctx context.Context) (waited time.Duration, err error) {
	if !p.config.Wait || p.config.MaxActive <= 0 {
		// No wait or no connection limit.
		return 0, nil
	}

	p.lazyInit()

	// wait indicates if we believe it will block, so it's not 100% accurate
	// however for stats it should be good enough.
	wait := len(p.ch) == 0
	var start time.Time
	if wait {
		start = time.Now()
	}

	select {
	case <-p.ch:
		// Additionally, check that context hasn't expired while we were waiting,
		// because `select` picks a random `case` if several of them are "ready".
		select {
		case <-ctx.Done():
			p.ch <- struct{}{}
			return 0, ctx.Err()
		default:
		}
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	if wait {
		return time.Since(start), nil
	}
	return 0, nil
}

func (p *pool) lazyInit() {
	p.initOnce.Do(func() {
		p.ch = make(chan struct{}, p.config.MaxActive)
		if p.closed {
			close(p.ch)
		} else {
			for i := 0; i < p.config.MaxActive; i++ {
				p.ch <- struct{}{}
			}
		}
	})
}

func (p *pool) dial(_ context.Context) (IConn, error) {
	if p.config.Factory == nil {
		return nil, errFactoryNotSet
	}

	conn, err := p.config.Factory(p.config.Param)
	if err != nil {
		return nil, err
	}
	go conn.Detector()
	return conn, nil
}
