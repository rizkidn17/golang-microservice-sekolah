package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang-microservice-sekolah/internal/database/handler"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	//grouping routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/", s.HelloWorldHandler)
		r.Route("/v1", func(r chi.Router) {
			r.Get("/", s.HelloWorldHandler)
			r.Get("/health", s.healthHandler)
			
			r.Route("/schools", func(r chi.Router) {
				r.Get("/", handler.GetSchoolsHandler)
				r.Get("/:{uuid}", handler.GetSchoolByUuidHandler)
				r.Post("/", handler.CreateSchoolHandler)
				r.Put("/:{uuid}", handler.UpdateSchoolByUuidHandler)
				r.Delete("/:{uuid}", handler.DeleteSchoolByUuidHandler)
			})
			
			r.Route("/users", func(r chi.Router) {
				r.Post("/register", handler.RegisterUserHandler)
				r.Post("/get-token", handler.GenerateUserTokenHandler)
			})
		})
	})
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
