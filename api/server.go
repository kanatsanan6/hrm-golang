package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New())

	v1 := app.Group("/api/v1")

	v1.Post("/sign_up", s.signUp)
	v1.Post("/sign_in", s.signIn)

	v1.Use(AuthMiddleware())
	v1.Use(MeMiddleware())

	v1.Get("/me", s.me)
	v1.Post("/company", s.createCompany)

	s.router = app
}

func (s *Server) Start(addr string) error {
	port := fmt.Sprintf(":%s", addr)
	return s.router.Listen(port)
}
