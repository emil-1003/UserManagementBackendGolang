package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/handlers"
	"github.com/emil-1003/UserManagementBackendGolang/pkg/middleware"
)

type Server struct {
	Name    string
	Version string
	Router  *mux.Router
	Port    string
}

func New(name, version, port, path string) (*Server, error) {
	r := mux.NewRouter()

	s := r.PathPrefix(fmt.Sprintf("/%s/%s", path, version)).Subrouter()

	s.Path("/signup").Handler(handlers.Signup()).Methods("POST")
	s.Path("/login").Handler(handlers.Login()).Methods("POST")

	s.Path("/users").Handler(middleware.AuthOnly(handlers.GetUsers())).Methods("GET")
	s.Path("/users/{id}").Handler(middleware.AdminOnly(handlers.DeleteUser())).Methods("DELETE")

	return &Server{name, version, s, port}, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.Port, s.Router)
}
