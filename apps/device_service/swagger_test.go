package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSwagger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerSwagger(router)

	for _, test := range []struct {
		path string
		want string
	}{
		{"/swagger", "SwaggerUIBundle"},
		{"/swagger/openapi.yaml", "title: Device Service API"},
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.path, nil))
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), test.want) {
			t.Fatalf("GET %s returned %d: %s", test.path, response.Code, response.Body.String())
		}
	}
}
