package api

func (s *serverState) routes() {
	s.router.HandleFunc("/", s.solveHandler())
}
