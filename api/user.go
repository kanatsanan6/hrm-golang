package api

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/service"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/sethvargo/go-password/password"
)

func userResponse(user model.User) model.User {
	return model.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type SignUpBody struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
}

func (s *Server) signUp(c *fiber.Ctx) error {
	var body SignUpBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	company, err := s.Queries.CreateCompany(queries.CreateCompanyArgs{Name: body.CompanyName})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
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
		CompanyID:         &company.ID,
		Role:              "admin",
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	for _, lType := range model.DefaultLeaveType {
		_, err := s.Queries.CreateLeaveType(queries.CreateLeaveTypeArgs{
			Name:   lType["name"].(string),
			Usage:  0,
			Max:    lType["max"].(int),
			UserID: user.ID,
		})
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	return utils.JsonResponse(c, fiber.StatusCreated, userResponse(user))
}

type SignInBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func signInResponse(token string, claims utils.CustomClaims, user model.User) tokenResponse {
	return tokenResponse{
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

type MeType struct {
	User   model.User           `json:"user"`
	Policy []service.PolicyType `json:"policy"`
}

func meResponse(user model.User, policies []map[string]string) *MeType {
	var policyResult []service.PolicyType
	for _, policy := range policies {
		policyResult = append(policyResult, service.PolicyType{
			Subject: policy["subject"],
			Action:  policy["action"],
		})
	}

	return &MeType{
		User:   userResponse(user),
		Policy: policyResult,
	}
}

func (s *Server) me(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

	return utils.JsonResponse(c, fiber.StatusOK, meResponse(user, s.Service.Export(c)))
}

type InviteUserBody struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (s *Server) inviteUser(c *fiber.Ctx) error {
	if authorized := s.Service.Authorize(c, "user_management", "invite"); !authorized {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var body InviteUserBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	user, _ := s.Queries.FindUserByEmail(body.Email)
	if !reflect.DeepEqual(user, model.User{}) {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, "user already exists")
	}

	password, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	hash, err := utils.Encrypt(password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	currentUser := c.Locals("user").(model.User)
	user, err = s.Queries.CreateUser(queries.CreateUserArgs{
		Email:             body.Email,
		EncryptedPassword: hash,
		FirstName:         body.FirstName,
		LastName:          body.LastName,
		CompanyID:         currentUser.CompanyID,
		Role:              "member",
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	token := user.GenerateResetPasswordToken()
	_, err = s.Queries.UpdateUser(queries.UpdateUserArgs{
		ID:                 user.ID,
		ResetPasswordToken: &token,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	// TODO: move this to worker
	messageBody := fmt.Sprintf("Please reset your password from this link: %s/reset-password/%s", os.Getenv("FRONTEND_URL"), token)
	if err := s.Service.Send(body.Email, "[HRM] You are invited", messageBody); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusCreated, userResponse(user))
}

type ForgetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *Server) forgetPassword(c *fiber.Ctx) error {
	var body ForgetPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	user, err := s.Queries.FindUserByEmail(body.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	token := user.GenerateResetPasswordToken()
	_, err = s.Queries.UpdateUser(queries.UpdateUserArgs{
		ID:                 user.ID,
		ResetPasswordToken: &token,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	// TODO: move this to worker
	messsageBody := fmt.Sprintf("Please reset your password from this link: %s/reset-password/%s", os.Getenv("FRONTEND_URL"), token)
	if err := s.Service.Send(body.Email, "[HRM] Reset Password", messsageBody); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusOK, userResponse(user))
}

type ResetPasswordBody struct {
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func (s *Server) resetPassword(c *fiber.Ctx) error {
	var body ResetPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	user, err := s.Queries.FindUserByResetPasswordToken(body.Token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
	}

	hash, err := utils.Encrypt(body.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	_, err = s.Queries.UpdateUser(queries.UpdateUserArgs{
		ID:       user.ID,
		Password: &hash,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusOK, userResponse(user))
}

type DeleteUserBody struct {
	ID int64 `json:"id" validate:"required"`
}

func (s *Server) deleteUser(c *fiber.Ctx) error {
	if authorized := s.Service.Authorize(c, "user_management", "delete"); !authorized {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var body DeleteUserBody
	if err := c.ParamsParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}
	user, err := s.Queries.FindUserByID(body.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
	}

	if err := s.Queries.DeleteUser(user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusNoContent, "")
}
