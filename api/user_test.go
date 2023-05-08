package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/model"
	mock_queries "github.com/kanatsanan6/hrm/queries/mock"
	mock_service "github.com/kanatsanan6/hrm/service/mock"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func GenerateUser(companyId *int64) *model.User {
	return &model.User{
		ID:                int64(utils.RandomNumber(1, 10)),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		Email:             utils.RandomEmail(),
		CompanyID:         companyId,
		EncryptedPassword: utils.RandomString(16),
	}

}

func TestServer_signUp(t *testing.T) {
	email := utils.RandomEmail()
	password := utils.RandomString(10)
	firstName := utils.RandomString(10)
	lastName := utils.RandomString(10)
	companyName := utils.RandomString(10)

	testCases := []struct {
		Name          string
		body          fiber.Map
		buildStub     func(q *mock_queries.MockQueries, s *mock_service.MockService)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			Name:      "BadRequest",
			body:      fiber.Map{},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			Name: "Cannot create user",
			body: fiber.Map{
				"email":        email,
				"password":     password,
				"first_name":   firstName,
				"last_name":    lastName,
				"company_name": companyName,
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					CreateUser(gomock.Any()).
					Times(1).
					Return(model.User{}, errors.New("cannot_create"))
				q.EXPECT().
					CreateCompany(gomock.Any()).
					Times(1).
					Return(model.Company{}, nil)

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			Name: "OK",
			body: fiber.Map{
				"email":        email,
				"password":     password,
				"first_name":   firstName,
				"last_name":    lastName,
				"company_name": companyName,
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					CreateUser(gomock.Any()).
					Times(1).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateCompany(gomock.Any()).
					Times(1).
					Return(model.Company{}, nil)
				q.EXPECT().
					CreateLeaveType(gomock.Any()).
					Times(4).
					Return(model.LeaveType{}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusCreated, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q, s)

			server := api.NewServer(q, s)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/sign_up", bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestServer_signIn(t *testing.T) {
	password := utils.RandomString(10)
	hash, _ := utils.Encrypt(password)

	user := &model.User{Email: utils.RandomEmail(), EncryptedPassword: hash}
	email := user.Email

	testCases := []struct {
		Name          string
		body          fiber.Map
		buildStub     func(q *mock_queries.MockQueries)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			Name:      "BadRequest",
			body:      fiber.Map{},
			buildStub: func(q *mock_queries.MockQueries) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			Name: "User cannot be found",
			body: fiber.Map{"email": email, "password": password},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(email).
					Times(1).
					Return(model.User{}, errors.New("not_found"))

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Error string `json:"errors"`
				}

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
				assert.Equal(t, "email or password is incorrect", result.Error)
			},
		},
		{
			Name: "Incorrect password",
			body: fiber.Map{"email": email, "password": "incorrect"},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(email).
					Times(1).
					Return(*user, nil)

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Error string `json:"errors"`
				}

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
				assert.Equal(t, "email or password is incorrect", result.Error)
			},
		},
		{
			Name: "correct password",
			body: fiber.Map{"email": email, "password": password},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(email).
					Times(1).
					Return(*user, nil)

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q)

			server := api.NewServer(q, s)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/sign_in", bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			res, _ := server.Router.Test(req, -1)

			tc.checkResponse(t, res)
		})
	}
}

func TestServer_me(t *testing.T) {
	email := utils.RandomEmail()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, email string)
		buildStub     func(q *mock_queries.MockQueries, s *mock_service.MockService)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Cannot find user",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(email).
					Return(model.User{}, errors.New("not_found"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
			},
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(email).
					Return(model.User{}, nil)
				s.EXPECT().
					Export(gomock.Any()).
					Return([]map[string]string{{"key": "value"}})
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q, s)

			server := api.NewServer(q, s)
			app := server.Router
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

			ctx.Locals("email", email)

			req := httptest.NewRequest("GET", "/api/v1/me", nil)

			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)

			tc.checkResponse(t, res)
		})
	}
}

func TestServer_inviteUser(t *testing.T) {
	email := "kanatsanan.j1998@gmail.com"
	firstName := utils.RandomString(10)
	lastName := utils.RandomString(10)
	user := &model.User{Email: email}

	testCases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, req *http.Request, email string)
		buildStub     func(q *mock_queries.MockQueries, s *mock_service.MockService)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			body:      fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Unauthorized 2",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(false)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Bad Request",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "user already exists",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{Email: email}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "user already exists",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{Email: email}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "cannot create user",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateUser(gomock.Any()).
					Return(model.User{Email: email}, errors.New("cannot create user"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "cannot update user token",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateUser(gomock.Any()).
					Return(model.User{Email: email}, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(model.User{}, errors.New("cannot update"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "Cannot send Email",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateUser(gomock.Any()).
					Return(model.User{Email: email}, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(model.User{Email: email}, nil)
				s.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
			},
		},
		{
			name: "OK",
			body: fiber.Map{"email": email, "first_name": firstName, "last_name": lastName},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "invite").
					Return(true)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateUser(gomock.Any()).
					Return(model.User{Email: email}, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(model.User{Email: email}, nil)
				s.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusCreated, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			server := api.NewServer(q, s)
			c := server.Router.AcquireCtx(&fasthttp.RequestCtx{})

			companyId := int64(utils.RandomNumber(1, 10))
			currentUser := model.User{
				CompanyID: &companyId,
			}
			c.Locals("users", currentUser)

			tc.buildStub(q, s)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/invite", bytes.NewBuffer(data))
			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestServer_forgetPassword(t *testing.T) {
	email := utils.RandomEmail()
	user := GenerateUser(nil)

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStub     func(q *mock_queries.MockQueries, s *mock_service.MockService)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "BadRequest",
			body:      fiber.Map{"email": "invalid"},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "User not found",
			body: fiber.Map{"email": email},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(model.User{}, errors.New("not_found"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
			},
		},
		{
			name: "Cannot update password",
			body: fiber.Map{"email": email},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(*user, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(model.User{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "Cannot send email",
			body: fiber.Map{"email": email},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(*user, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(*user, nil)
				s.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
			},
		},
		{
			name: "OK",
			body: fiber.Map{"email": email},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Return(*user, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(*user, nil)
				s.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q, s)

			server := api.NewServer(q, s)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/forget_password", bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestServer_resetPassword(t *testing.T) {
	password := utils.RandomString(10)
	token := utils.RandomString(16)
	body := fiber.Map{"password": password, "token": token}
	user := GenerateUser(nil)

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStub     func(q *mock_queries.MockQueries)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "BadRequest",
			body:      fiber.Map{},
			buildStub: func(q *mock_queries.MockQueries) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "Invalid Token",
			body: body,
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByResetPasswordToken(gomock.Eq(token)).
					Return(model.User{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Cannot update password",
			body: body,
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByResetPasswordToken(gomock.Eq(token)).
					Return(*user, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(model.User{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "OK",
			body: body,
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByResetPasswordToken(gomock.Eq(token)).
					Return(*user, nil)
				q.EXPECT().
					UpdateUser(gomock.Any()).
					Return(*user, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q)

			server := api.NewServer(q, s)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", "/api/v1/reset_password", bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestServer_deleteUser(t *testing.T) {
	email := "kanatsanan.j1998@gmail.com"
	id := uint(utils.RandomNumber(1, 10))
	user := GenerateUser(nil)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, email string)
		buildStub     func(q *mock_queries.MockQueries, s *mock_service.MockService)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Unauthorized 2",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "delete").
					Return(false)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "User Not found",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "delete").
					Return(true)
				q.EXPECT().
					FindUserByID(gomock.Eq(int64(id))).
					Return(model.User{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
			},
		},
		{
			name: "Cannot delete user",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "delete").
					Return(true)
				q.EXPECT().
					FindUserByID(gomock.Eq(int64(id))).
					Return(*user, nil)
				q.EXPECT().
					DeleteUser(gomock.Eq(int64(user.ID))).
					Return(errors.New("errors"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, s *mock_service.MockService) {
				MockMe(q, *user, email)
				s.EXPECT().
					Authorize(gomock.Any(), "user_management", "delete").
					Return(true)
				q.EXPECT().
					FindUserByID(gomock.Eq(int64(id))).
					Return(*user, nil)
				q.EXPECT().
					DeleteUser(gomock.Eq(int64(user.ID))).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusNoContent, res.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			s := mock_service.NewMockService(ctrl)
			tc.buildStub(q, s)

			server := api.NewServer(q, s)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/company/users/%d", id), nil)
			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}
