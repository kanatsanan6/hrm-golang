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
)

func TestServer_createLeave(t *testing.T) {
	email := utils.RandomEmail()
	user := model.User{Email: email}

	description := utils.RandomString(10)
	startDate := "2023-04-24"
	endDate := "2023-04-25"
	leaveType := "vacation_leave"

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
			name: "Bad Request",
			body: fiber.Map{},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, user, email)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "Cannot create leave",
			body: fiber.Map{
				"description": description,
				"start_date":  startDate,
				"end_date":    endDate,
				"leave_type":  leaveType,
			},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, user, email)
				q.EXPECT().
					CreateLeave(gomock.Any()).
					Return(model.Leave{}, errors.New("not_found"))
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
			},
		},
		{
			name: "OK",
			body: fiber.Map{
				"description": description,
				"start_date":  startDate,
				"end_date":    endDate,
				"leave_type":  leaveType,
			},
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, user, email)
				q.EXPECT().
					CreateLeave(gomock.Any()).
					Return(model.Leave{
						Description: description,
						Status:      "pending",
						LeaveType:   leaveType,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data api.LeaveType `json:"data"`
				}

				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusCreated, res.StatusCode)
				assert.Equal(t, "pending", result.Data.Status)
				assert.Equal(t, leaveType, result.Data.LeaveType)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			p := mock_service.NewMockPolicyInterface(ctrl)
			tc.buildStub(q)

			server := api.NewServer(q, p)

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/company/leaves", bytes.NewBuffer(data))
			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}

func TestServer_getLeaves(t *testing.T) {
	user := GenerateUser(nil)
	leave1 := queries.LeaveType{
		ID:          1,
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2023, 01, 02, 0, 0, 0, 0, time.UTC),
		LeaveType:   "vacation_leave",
		CreatedAt:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	leave2 := queries.LeaveType{
		ID:          2,
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2023, 01, 02, 0, 0, 0, 0, time.UTC),
		LeaveType:   "vacation_leave",
		CreatedAt:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, email string)
		buildStub     func(q *mock_queries.MockQueries, p *mock_service.MockPolicyInterface)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name:      "Unauthorized",
			setupAuth: func(t *testing.T, req *http.Request, email string) {},
			buildStub: func(q *mock_queries.MockQueries, p *mock_service.MockPolicyInterface) {},
			checkResponse: func(t *testing.T, res *http.Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries, p *mock_service.MockPolicyInterface) {
				MockMe(q, *user, user.Email)
				q.EXPECT().
					GetLeaves(gomock.Any()).
					Times(1).
					Return([]queries.LeaveType{leave1, leave2})
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				var result struct {
					Data []queries.LeaveType `json:"data"`
				}
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, 1, int(result.Data[0].ID))
				assert.Equal(t, 2, int(result.Data[1].ID))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			q := mock_queries.NewMockQueries(ctrl)
			p := mock_service.NewMockPolicyInterface(ctrl)
			tc.buildStub(q, p)

			server := api.NewServer(q, p)

			req := httptest.NewRequest("GET", "/api/v1/company/leaves", nil)
			tc.setupAuth(t, req, user.Email)

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}
