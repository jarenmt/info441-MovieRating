package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/assignments-fixed-ssunni12/servers/gateway/handlers"
)

func TestCors(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/v1/users/", nil)
	var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("IWDLAMWLDN"))
	})
	corsTest := handlers.NewHeaderHandler(testHandler)
	corsTest.ServeHTTP(recorder, request)
	response := recorder.Result()
	if response.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Header 1 didn't get added; Access-Control-Allow-Origin is not *")
	}
	if response.Header.Get("Access-Control-Allow-Methods") != "GET, PUT, POST, PATCH, DELETE" {
		t.Errorf("Header 2 didn't get added; Access-Control-Allow-Methods are not GET, PUT, POST, PATCH, DELETE")
	}
	if response.Header.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("Header 3 didn't get added; ccess-Control-Allow-Headers are not Content-Type, Authorization")
	}
	if response.Header.Get("Access-Control-Expose-Headers") != "Authorization" {
		t.Errorf("Header 4 didn't get added; Access-Control-Expose-Headers is not Authorization")
	}
	if response.Header.Get("Access-Control-Max-Age") != "600" {
		t.Errorf("Header 5 didn't get added; Access-Control-Max-Age is not 600")
	}
}
