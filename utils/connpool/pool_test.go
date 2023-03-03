package connpool

import (
	"log"
	"testing"
	"time"
)

var _ IConn = (*TestConnection)(nil)

type TestConnection struct {
	Connection
}

func (t TestConnection) Close() error {
	log.Println("close")
	return nil
}

func (t TestConnection) Closed() bool {
	return false
}

func (t TestConnection) Reusable() bool {
	return true
}

func (t TestConnection) Detector() {

}

func TestConnPool(t *testing.T) {
	p := GetPool("test-connection", &Config{
		Factory:   New,
		Param:     nil,
		MaxIdle:   10,
		MaxActive: 100,
		Wait:      false,
	})
	c, _ := p.Get()
	// do sth.
	time.Sleep(10 * time.Millisecond)
	_ = p.Put(c)
}

func New(p interface{}) (IConn, error) {
	log.Println("new connection")
	return &TestConnection{}, nil
}
