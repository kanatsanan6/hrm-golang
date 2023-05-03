package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	mock_queries "github.com/kanatsanan6/hrm/queries/mock"
	mock_service "github.com/kanatsanan6/hrm/service/mock"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func generateCompany() *model.Company {
	return &model.Company{
		ID:        uint(utils.RandomNumber(0, 10)),
		Name:      utils.RandomString(10),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestServer_createCompany(t *testing.T) {
	email := "kanatsanan.j1998@gmail.com"
	company := generateCompany()
	user := GenerateUser(&company.ID)

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
			name: "BadRequest",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "NotFound User",
			body: fiber.Map{"name": company.Name},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
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
			name: "Cannot create company",
			body: fiber.Map{"name": company.Name},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateCompany(gomock.Eq(queries.CreateCompanyArgs{Name: company.Name})).
					Times(1).
					Return(model.Company{}, errors.New("cannot create"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "Cannot update user companyID",
			body: fiber.Map{"name": company.Name},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateCompany(gomock.Eq(queries.CreateCompanyArgs{Name: company.Name})).
					Times(1).
					Return(model.Company{}, nil)
				q.EXPECT().
					UpdateUserCompanyID(gomock.Eq(model.User{}), gomock.Eq(uint(0))).
					Times(1).
					Return(errors.New("cannot update"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "Cannot update user companyID",
			body: fiber.Map{"name": company.Name},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(model.User{}, nil)
				q.EXPECT().
					CreateCompany(gomock.Eq(queries.CreateCompanyArgs{Name: company.Name})).
					Times(1).
					Return(*company, nil)
				q.EXPECT().
					UpdateUserCompanyID(gomock.Eq(model.User{}), gomock.Eq(company.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data api.CompanyType `json:"data"`
				}

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusCreated, res.StatusCode)
				assert.Equal(t, company.ID, result.Data.ID)
				assert.Equal(t, company.Name, result.Data.Name)
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

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/company", bytes.NewBuffer(data))

			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)

			tc.checkResponse(t, res)
		})
	}
}

func TestServer_getUsers(t *testing.T) {
	email := utils.RandomEmail()
	id := uint(utils.RandomNumber(1, 10))
	user := &model.User{
		ID:        uint(utils.RandomNumber(1, 10)),
		Email:     email,
		CompanyID: &id,
	}
	otherUser := &model.User{
		ID:        uint(utils.RandomNumber(1, 10)),
		Email:     email,
		CompanyID: &id,
	}
	company := &model.Company{
		ID:    id,
		Name:  utils.RandomString(10),
		Users: []model.User{*user, *otherUser},
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
			name: "User Not found",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
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
			name: "Company Not found",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(*user, nil)
				q.EXPECT().
					FindCompanyByID(gomock.Eq(*user.CompanyID)).
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
				MockMe(q, *user, email)
				q.EXPECT().
					FindUserByEmail(gomock.Eq(email)).
					Times(1).
					Return(*user, nil)
				q.EXPECT().
					FindCompanyByID(gomock.Eq(*user.CompanyID)).
					Times(1).
					Return(*company, nil)
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
