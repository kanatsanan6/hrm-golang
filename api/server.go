package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	router *fiber.App
}

func NewServer() *Server {
	server := &Server{}

	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	app := fiber.New()

	s.router = app
}

func (s *Server) Start(addr string) error {
	port := fmt.Sprintf(":%s", addr)
	return s.router.Listen(port)
}
