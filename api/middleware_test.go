package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/kanatsanan6/hrm/api"
	mock_queries "github.com/kanatsanan6/hrm/queries/mock"
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

func TestMiddleware(t *testing.T) {
	email := utils.RandomEmail()
	fmt.Println(email)

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
				assert.Equal(t, email, result.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			server := api.NewServer(q)
			app := server.Router

			app.Use(api.AuthMiddleware())
			app.Use(api.MeMiddleware())

			app.Get("/me", func(c *fiber.Ctx) error {
				return utils.JsonResponse(c, fiber.StatusOK, c.Locals("email"))
			})

			req := httptest.NewRequest("GET", "/me", nil)

			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)

			tc.checkResponse(t, res)
		})
	}
}
