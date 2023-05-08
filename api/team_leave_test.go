package api_test

// func TestServer_getTeamLeaves(t *testing.T) {
// 	user := GenerateUser(nil)
// 	email := user.Email

// 	testCases := []struct {
// 		name          string
// 		setupAuth     func(t *testing.T, req *http.Request, email string)
// 		buildStub     func(q *mock_queries.MockQueries)
// 		checkResponse func(t *testing.T, res *http.Response)
// 	}{
// 		{
// 			name:      "Unauthorized",
// 			setupAuth: func(t *testing.T, req *http.Request, email string) {},
// 			buildStub: func(q *mock_queries.MockQueries) {},
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
// 			},
// 		},
// 		{
// 			name: "OK",
// 			setupAuth: func(t *testing.T, req *http.Request, email string) {
// 				AddAuth(t, req, email)
// 			},
// 			buildStub: func(q *mock_queries.MockQueries) {
// 				MockMe(q, *user, email)
// 				q.EXPECT().
// 					GetTeamLeaves(gomock.Any()).
// 					Return([]queries.LeaveStruct{}, nil)
// 			},
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				assert.Equal(t, fiber.StatusOK, res.StatusCode)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			q := mock_queries.NewMockQueries(ctrl)
// 			s := mock_service.NewMockService(ctrl)
// 			tc.buildStub(q)

// 			server := api.NewServer(q, s)

// 			req := httptest.NewRequest("GET", "/api/v1/company/team_leaves", nil)
// 			tc.setupAuth(t, req, email)

// 			res, _ := server.Router.Test(req, -1)
// 			tc.checkResponse(t, res)
// 		})
// 	}
// }
