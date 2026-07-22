package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/internal/model"
	fixture "github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRegisterEndpoint(t *testing.T) {
    testCases := []struct {
        name string

        reqMethod string
        reqPath string
        reqBody string

        useUserCommonFixture bool
        closeDBBeforeReq bool

        expectedStatusCode int
        expectedMessage string
        expectedResponse string
    }{
        {
            name: "201 - registers user successfully",
            reqMethod: http.MethodPost,
            reqPath: "/v1/users/register",
            reqBody: `{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
            expectedStatusCode: http.StatusCreated,
            expectedMessage: "User registered successfully!",
        },
        {
            name: "400 - invalid input (empty body)",
            reqMethod: http.MethodPost,
            reqPath: "/v1/users/register",
            reqBody: `{}`,
            expectedStatusCode: http.StatusBadRequest,
            expectedMessage: "Invalid input",
        },
        {
            name: "409 - username already exists",
            reqMethod: http.MethodPost,
            reqPath: "/v1/users/register",
            reqBody: `{"username":"test1","display_name":"Test User","password":"Password123@","email":"new@example.com"}`,
            useUserCommonFixture: true,
            expectedStatusCode: http.StatusConflict,
            expectedMessage: "username already exists",
        },
        {
            name: "409 - email already exists",
            reqMethod: http.MethodPost,
            reqPath: "/v1/users/register",
            reqBody: `{"username":"newname","display_name":"Test User","password":"Password123@","email":"test2@example.com"}`,
            useUserCommonFixture: true,
            expectedStatusCode: http.StatusConflict,
            expectedMessage: "email already exists",
        },
        {
            name: "500 - database error",
            reqMethod: http.MethodPost,
            reqPath: "/v1/users/register",
            reqBody: `{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
            closeDBBeforeReq: true,
            expectedStatusCode: http.StatusInternalServerError,
            expectedResponse: `{"message":"Processing error"}`,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
            t.Setenv("SERVICE_NAME", "bookmark_service_test")

            // register custom validations used by handlers
            require.NoError(t, validation.RegisterValidation())

            cfg, err := api.NewConfig()
            require.NoError(t, err)

            redisClient := redisPkg.InitMockRedis(t)

            var db *gorm.DB
            if tc.useUserCommonFixture {
                db = fixture.NewFixture(t, &fixture.UserCommonTest{})
            } else {
                db = sqldb.InitMockDB(t)
                // migrate user model
                require.NoError(t, db.AutoMigrate(&model.User{}))
            }

            if tc.closeDBBeforeReq {
                sqlDB, derr := db.DB()
                require.NoError(t, derr)
                require.NoError(t, sqlDB.Close())
            }

            testAPI := api.NewEngine(cfg, redisClient, db)

            req := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBufferString(tc.reqBody))
            req.Header.Set("Content-Type", "application/json")
            recorder := httptest.NewRecorder()

            testAPI.ServeHTTP(recorder, req)

            assert.Equal(t, tc.expectedStatusCode, recorder.Code)

            if tc.expectedMessage != "" {
                var body struct{ Message string `json:"message"` }
                require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
                assert.Equal(t, tc.expectedMessage, body.Message)
                return
            }

            if tc.expectedResponse != "" {
                assert.JSONEq(t, tc.expectedResponse, recorder.Body.String())
                return
            }

            if tc.expectedStatusCode == http.StatusCreated {
                var resp struct{
                    Data struct{ Username string `json:"username"` } `json:"data"`
                    Message string `json:"message"`
                }
                require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
                assert.Equal(t, "johndoe", resp.Data.Username)
                assert.Equal(t, "User registered successfully!", resp.Message)
            }
        })
    }
}

