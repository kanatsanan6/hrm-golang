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
	Service service.Service
}

func NewServer(q queries.Queries, s service.Service) *Server {
	server := &Server{
		Queries: q,
		Service: s,
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
	v1.Post("/invite", s.inviteUser)

	company := v1.Group("/company")

	company.Post("/", s.createCompany)
	company.Get("/users", s.getUsers)
	company.Delete("/users/:id", s.deleteUser)

	company.Get("/leaves", s.getLeaves)
	company.Post("/leaves", s.createLeave)

	company.Get("/leave_types", s.getLeaveTypes)

	s.Router = app
}

func (s *Server) Start(addr string) error {
	port := fmt.Sprintf(":%s", addr)
	return s.Router.Listen(port)
}
