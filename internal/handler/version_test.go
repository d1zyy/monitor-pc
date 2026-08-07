package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/version", GetVersion)

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal("Failed to create request: ", err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"version":"unknown","commit":"unknown","build_time":"unknown"}`, w.Body.String())
}
