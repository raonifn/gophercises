package urlshort

import "net/http"

type Server struct {
	handler http.Handler
}

type HandlerStacker func(http.Handler) http.Handler

func NewServer() *Server {
	return &Server{handler: defaultHandler()}
}

func (s *Server) Start(addr string) {
	http.ListenAndServe(addr, s.handler)
}

func (s *Server) StackHandler(stacker HandlerStacker) {
	nh := stacker(s.handler)
	s.handler = nh
}

func defaultHandler() http.Handler {
	return http.HandlerFunc(notFound)
}
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
