package api_test

import (
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
	mock_service "github.com/kanatsanan6/hrm/service/mock"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestServer_getUsers(t *testing.T) {
	email := utils.RandomEmail()
	id := int64(utils.RandomNumber(1, 10))
	user := model.User{
		ID:        int64(utils.RandomNumber(1, 10)),
		Email:     email,
		CompanyID: &id,
	}
	otherUser := model.User{
		ID:        int64(utils.RandomNumber(1, 10)),
		Email:     email,
		CompanyID: &id,
	}
	company := model.Company{
		ID:   id,
		Name: utils.RandomString(10),
	}

	testCases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, req *http.Request, email string)
		buildStub     func(q *mock_queries.MockQueries)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			body:      fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Company Not found",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(user, nil)
				q.EXPECT().
					FindCompanyByID(gomock.Eq(id)).
					Times(1).
					Return(model.Company{}, errors.New("not_found"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
			},
		},
		{
			name: "OK",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(user, nil)
				q.EXPECT().
					FindCompanyByID(gomock.Eq(id)).
					Times(1).
					Return(company, nil)
				q.EXPECT().
					GetUsers(gomock.Eq(company.ID)).
					Return([]model.User{user, otherUser}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data []model.User `json:"data"`
				}

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusOK, res.StatusCode)

				assert.Equal(t, 2, len(result.Data))
				assert.Equal(t, user.ID, result.Data[0].ID)
				assert.Equal(t, otherUser.ID, result.Data[1].ID)
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
			app := server.Router
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

			ctx.Locals("email", email)

			req := httptest.NewRequest("GET", "/api/v1/company/users", nil)

			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)

			tc.checkResponse(t, res)
		})
	}
}
