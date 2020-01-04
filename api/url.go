package api

import "fmt"

func (s *Server) modelURI(id int) string {
	return fmt.Sprintf("%s/models/%d", s.baseURL, id)
}
