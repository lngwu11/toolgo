package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
	"time"
)

//Limiter 限流器对象
type Limiter struct {
	value int64
	max   int64
	ts    int64
}

//NewLimiter 产生一个限流器
func NewLimiter(cnt int64) *Limiter {
	return &Limiter{
		value: 0,
		max:   cnt,
		ts:    time.Now().Unix(),
	}
}

//Ok 是否可以通过
func (l *Limiter) Ok() bool {
	ts := time.Now().Unix()
	tsOld := atomic.LoadInt64(&l.ts)
	if ts != tsOld {
		atomic.StoreInt64(&l.ts, ts)
		atomic.StoreInt64(&l.value, 1)
		return true
	}
	return atomic.AddInt64(&(l.value), 1) < l.max
}

//SetMax 设置最大限制
func (l *Limiter) SetMax(m int64) {
	l.max = m
}

//MaxAllowed 限流器
func MaxAllowed(limitValue int64) func(c *gin.Context) {
	limiter := NewLimiter(limitValue)
	// 返回限流逻辑
	return func(c *gin.Context) {
		if !limiter.Ok() {
			// 超过每秒limitValue，返回429错误码
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
