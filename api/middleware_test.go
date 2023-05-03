package api_test

import (
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
)

func AddAuth(t *testing.T, req *http.Request, email string) {
	jwtToken, _, err := utils.GenerateJWT(email)

	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	authorization := fmt.Sprintf("Bearer %s", jwtToken)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
}

func MockMe(q *mock_queries.MockQueries, user model.User, email string) {
	q.EXPECT().
		FindUserByEmail(gomock.Eq(email)).
		Times(1).
		Return(model.User{Email: email}, nil)
}

func TestAuthMiddleware(t *testing.T) {
	email := utils.RandomEmail()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, email string)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, req *http.Request, email string) {

			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data string `json:"data"`
				}
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, "success", result.Data)
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
			app := server.Router

			app.Use(server.AuthMiddleware())
			app.Get("/test", func(c *fiber.Ctx) error {
				return utils.JsonResponse(c, fiber.StatusOK, "success")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			tc.setupAuth(t, req, email)

			res, _ := app.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestMeMiddleware(t *testing.T) {
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
			name: "User not found",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
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
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(model.User{Email: email}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data api.UserType `json:"data"`
				}
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, email, result.Data.Email)
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
			app := server.Router

			app.Use(server.AuthMiddleware(), server.MeMiddleware())
			app.Get("/test", func(c *fiber.Ctx) error {
				return utils.JsonResponse(c, fiber.StatusOK, c.Locals(("user")))
			})

			tc.buildStub(q)

			req := httptest.NewRequest("GET", "/test", nil)
			tc.setupAuth(t, req, email)

			res, _ := app.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}
