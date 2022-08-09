package http

func NewTransportObject(method, pattern string, server *Server) *RouterObject {
	r := RouterObject{
		HttpMethod: method,
		Pattern:    pattern,
		Handler:    server.ServerHTTP,
	}
	return &r
}
