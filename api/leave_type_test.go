package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/model"
	mock_queries "github.com/kanatsanan6/hrm/queries/mock"
	mock_service "github.com/kanatsanan6/hrm/service/mock"
	"github.com/stretchr/testify/assert"
)

func TestServer_getLeaveTypes(t *testing.T) {
	user := GenerateUser(nil)
	email := user.Email

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
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, email string) {
				AddAuth(t, req, email)
			},
			buildStub: func(q *mock_queries.MockQueries) {
				MockMe(q, *user, email)
				q.EXPECT().
					GetUserLeaveTypes(gomock.Any()).
					Return([]model.LeaveType{{ID: 1}, {ID: 2}}, nil)
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

			req := httptest.NewRequest("GET", "/api/v1/company/leave_types", nil)
			tc.setupAuth(t, req, email)

			res, _ := server.Router.Test(req, -1)
			tc.checkResponse(t, res)
		})
	}
}
