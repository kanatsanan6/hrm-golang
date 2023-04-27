package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

type userType struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func userResponse(user model.User) userType {
	return userType{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type SignUpBody struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (s *Server) signUp(c *fiber.Ctx) error {
	var body SignUpBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	hash, err := utils.Encrypt(body.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	user, err := s.Queries.CreateUser(queries.CreateUserArgs{
		Email:             body.Email,
		EncryptedPassword: hash,
		FirstName:         body.FirstName,
		LastName:          body.LastName,
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}
	return utils.JsonResponse(c, fiber.StatusCreated, userResponse(user))
}

type SignInBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type tokenResponse struct {
	CompanyID *uint  `json:"company_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func signInResponse(token string, claims utils.CustomClaims, user model.User) tokenResponse {
	return tokenResponse{
		CompanyID: user.CompanyID,
		Token:     token,
		ExpiresAt: claims.ExpiresAt,
	}

}

func (s *Server) signIn(c *fiber.Ctx) error {
	var body SignInBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	user, err := s.Queries.FindUserByEmail(body.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, "email or password is incorrect")
	}

	if matched := utils.ComparePasswords(user.EncryptedPassword, body.Password); !matched {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, "email or password is incorrect")
	}

	signedToken, claims, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	response := signInResponse(signedToken, claims, user)

	return utils.JsonResponse(c, fiber.StatusOK, response)
}

func (s *Server) me(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	user, err := s.Queries.FindUserByEmail(email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusOK, userResponse(user))
}
