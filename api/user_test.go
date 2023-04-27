package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/model"
	mock_queries "github.com/kanatsanan6/hrm/queries/mock"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func GenerateUser(hash string) *model.User {
	return &model.User{
		ID:                uint(utils.RandomNumber(1, 10)),
		Email:             utils.RandomEmail(),
		EncryptedPassword: hash,
	}

}

func TestServer_signUp(t *testing.T) {
	email := utils.RandomEmail()
	password := utils.RandomString(10)
	firstName := utils.RandomString(10)
	lastName := utils.RandomString(10)

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
			Name: "Cannot create user",
			body: fiber.Map{
				"email":      email,
				"password":   password,
				"first_name": firstName,
				"last_name":  lastName,
			},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					CreateUser(gomock.Any()).
					Times(1).
					Return(model.User{}, errors.New("cannot_create"))

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			Name: "OK",
			body: fiber.Map{
				"email":      email,
				"password":   password,
				"first_name": firstName,
				"last_name":  lastName,
			},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					CreateUser(gomock.Any()).
					Times(1).
					Return(model.User{}, nil)

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
			tc.buildStub(q)

			server := api.NewServer(q)

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

	user := GenerateUser(hash)
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
			tc.buildStub(q)

			server := api.NewServer(q)

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
		buildStub     func(q *mock_queries.MockQueries)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Cannot find user",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
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
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(email).
					Return(model.User{}, nil)
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
			tc.buildStub(q)

			server := api.NewServer(q)
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
