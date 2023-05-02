package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/service"
)

type Server struct {
	Router  *fiber.App
	Queries queries.Queries
	Policy  service.PolicyInterface
}

func NewServer(q queries.Queries, p service.PolicyInterface) *Server {
	server := &Server{
		Queries: q,
		Policy:  p,
	}

	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	app := fiber.New()

	app.Use(cors.New())

	v1 := app.Group("/api/v1")

	v1.Post("/forget_password", s.forgetPassword)
	v1.Put("/reset_password", s.resetPassword)
	v1.Post("/sign_up", s.signUp)
	v1.Post("/sign_in", s.signIn)

	v1.Use(s.AuthMiddleware(), s.MeMiddleware())

	v1.Get("/me", s.me)
	v1.Post("/company", s.createCompany)
	v1.Get("/company/users", s.getUsers)
	v1.Post("/invite", s.inviteUser)
	v1.Delete("/company/users/:id", s.deleteUser)

	v1.Post("/company/leaves", s.createLeave)

	s.Router = app
}

func (s *Server) Start(addr string) error {
	port := fmt.Sprintf(":%s", addr)
	return s.Router.Listen(port)
}
