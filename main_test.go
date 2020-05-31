package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/vinhut/graphservice/mocks"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks.NewMockAuthService(ctrl)

	router := setupRouter(mock_auth)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
