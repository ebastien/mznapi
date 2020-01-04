package api

func (s *Server) routes() {
	s.router.HandleFunc("/", s.solveHandler())
	s.router.HandleFunc("/models", s.createHandler())
}
