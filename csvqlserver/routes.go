package csvqlserver

import (
	"net/http"
)

func (s *Server) Route(router *http.ServeMux) {
	router.HandleFunc("POST /csvql/submit", s.Submit)
	router.HandleFunc("POST /csvql/query", s.Query)
	router.HandleFunc("GET /csvql/schemas", s.Schemas)
}
