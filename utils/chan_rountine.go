package utils

import (
	"reflect"
	"sync"
	"time"
)

type ChanRountineProc func(value interface{})

type ChanRountine struct {
	selectCase []reflect.SelectCase
	rountine   []ChanRountineProc
	step       time.Duration
	lock       sync.Mutex
}

func NewChanRountine(step time.Duration) *ChanRountine {
	r := ChanRountine{
		selectCase: make([]reflect.SelectCase, 0, 0),
		rountine:   make([]ChanRountineProc, 0, 0),
		step:       step,
	}
	go r.run()
	return &r
}

func (c *ChanRountine) run() {
	for {
		if len(c.selectCase) <= 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		chosen, receive, ok := reflect.Select(c.selectCase)
		if ok && len(c.rountine) > chosen && c.rountine[chosen] != nil {
			c.rountine[chosen](receive.Interface())
		}
		if c.step > 0 {
			time.Sleep(c.step)
		}
	}
}

func (c *ChanRountine) Register(channel interface{}, rountine ChanRountineProc) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	exist := false
	for _, sc := range c.selectCase {
		if sc.Chan == reflect.ValueOf(channel) {
			exist = true
			break
		}
	}
	if !exist {
		var selectCase reflect.SelectCase
		selectCase.Chan = reflect.ValueOf(channel)
		selectCase.Dir = reflect.SelectRecv
		c.rountine = append(c.rountine, rountine)
		c.selectCase = append(c.selectCase, selectCase)
	}
	return nil
}
