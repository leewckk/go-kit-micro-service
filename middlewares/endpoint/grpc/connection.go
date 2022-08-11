package grpc

import (
	"sync"

	"google.golang.org/grpc"
)

type Connection struct {
	conn sync.Map
}

func init() {
}

func NewConnection() *Connection {
	c := &Connection{}
	return c
}

func (this *Connection) GetConnection(instance string) *grpc.ClientConn {
	if c, ok := this.conn.Load(instance); ok {
		return c.(*grpc.ClientConn)
	}

	c, err := grpc.Dial(instance, grpc.WithInsecure())
	if nil != err {
		panic(err)
	}
	this.conn.Store(instance, c)
	return c
}
