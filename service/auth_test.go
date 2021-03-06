package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"sample-project/storage"
	"sample-project/structs"
)

func (s fakeStoragePostgres) GetUserUUID(u *structs.User) (string, error) {
	return "01234567-8901-2345-6789-012345678901", nil
}

func (s fakeStoragePostgresErr) GetUserUUID(u *structs.User) (string, error) {
	return "", errors.New("no rows")
}

func (s fakeStorageRedis) RegisterAuth(userUUID string, token *structs.AuthToken) error {
	return nil
}

func (s fakeStorageRedisErr) RegisterAuth(userUUID string, token *structs.AuthToken) error {
	return errors.New("something gone wrong")
}

func (s fakeStorageRedis) DeleteAuth(accessUUID string) error {
	return nil
}

func TestLogin(t *testing.T) {
	fsPG := fakeStoragePostgres{}
	fsRD := fakeStorageRedis{}
	fsPGerr := fakeStoragePostgresErr{}
	fsRDerr := fakeStorageRedisErr{}

	testCases := []struct {
		name         string
		username     string
		password     string
		statusCode   int
		responseBody string
		isReqSuccess bool
		reqMessage   string
		storagePG    storage.Storager
		storageRD    storage.RedisStorager
	}{
		{
			name:         "case success auth",
			username:     "foo",
			password:     "foobar",
			isReqSuccess: true,
			reqMessage:   "OK",
			statusCode:   http.StatusOK,
			storagePG:    fsPG,
			storageRD:    fsRD,
		},
		{
			name:         "case empty credentials",
			statusCode:   http.StatusBadRequest,
			isReqSuccess: false,
			reqMessage:   "empty username or password",
			storagePG:    fsPG,
			storageRD:    fsRD,
		},
		{
			name:         "case user not found",
			statusCode:   http.StatusUnauthorized,
			username:     "foo",
			password:     "foobar",
			isReqSuccess: false,
			reqMessage:   "incorrect username or password",
			storagePG:    fsPGerr,
			storageRD:    fsRD,
		},
		{
			name:         "case auth register err",
			statusCode:   http.StatusInternalServerError,
			username:     "foo",
			password:     "foobar",
			isReqSuccess: false,
			reqMessage:   "something gone wrong",
			storagePG:    fsPG,
			storageRD:    fsRDerr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &APIv1{store: tc.storagePG, redisClient: tc.storageRD}

			req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
			req.Form = url.Values{"username": []string{tc.username}, "password": []string{tc.password}}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/login")

			assert.NoError(t, api.login(c))
			assert.Equal(t, tc.statusCode, rec.Code)

			resp := Response{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.isReqSuccess, resp.Success)
			assert.Equal(t, tc.reqMessage, resp.Message)
		})
	}
}

func TestLogout(t *testing.T) {
	fsRD := fakeStorageRedis{}
	fsPG := fakeStoragePostgres{}

	testCases := []struct {
		name         string
		statusCode   int
		responseBody string
		isReqSuccess bool
		reqMessage   string
		storagePG    storage.Storager
		storageRD    storage.RedisStorager
	}{
		{
			name:         "case success logout",
			isReqSuccess: true,
			reqMessage:   "OK",
			statusCode:   http.StatusOK,
			storageRD:    fsRD,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &APIv1{redisClient: tc.storageRD, store: fsPG}

			req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/logout")

			req.Form = url.Values{"username": []string{"foo"}, "password": []string{"foobar"}}
			assert.NoError(t, api.login(c))

			resp := Response{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.isReqSuccess, resp.Success)
			assert.Equal(t, tc.reqMessage, resp.Message)
			assert.NotEmpty(t, resp.Data)

			req.Header.Set("Authorization", fmt.Sprint("Bearer ", resp.Data))

			// reset recorder & context after login to prevent response mixing
			rec = httptest.NewRecorder()
			c = e.NewContext(req, rec)
			c.SetPath("/logout")
			assert.NoError(t, api.logout(c))
			assert.Equal(t, tc.statusCode, rec.Code)

			resp = Response{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.isReqSuccess, resp.Success)
			assert.Equal(t, tc.reqMessage, resp.Message)
		})
	}
}
