package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/queries"
)

type Server struct {
	router  *fiber.App
	Queries *queries.Queries
}

func NewServer() *Server {
	server := &Server{}

	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	app := fiber.New()

	app.Post("/sign_up", s.signUp)
	app.Get("/sign_in", s.signIn)

	s.router = app
}

func (s *Server) Start(addr string) error {
	port := fmt.Sprintf(":%s", addr)
	return s.router.Listen(port)
}
