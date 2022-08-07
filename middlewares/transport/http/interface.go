package http

import "github.com/gin-gonic/gin"

type RouterObject struct {
	HttpMethod string
	Pattern    string
	Handler    gin.HandlerFunc
}

type DecodeRequestFunc func(*gin.Context) (interface{}, error)
type EncodeResponseFunc func(*gin.Context, interface{}) error
