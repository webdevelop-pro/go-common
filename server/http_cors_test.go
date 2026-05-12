package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCORSPreflightAllowsPrivateNetworkRequestHeader pins the fix for the
// Firefox failure: "header 'access-control-request-private-network' is not
// allowed according to header 'Access-Control-Allow-Headers' from CORS
// preflight response". The preflight response's Access-Control-Allow-Headers
// must list `Access-Control-Request-Private-Network`.
func TestCORSPreflightAllowsPrivateNetworkRequestHeader(t *testing.T) {
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "0")

	srv := NewServer()

	req := httptest.NewRequest(http.MethodOptions, "/v1.0/internal/objects/offer/130/", nil)
	req.Header.Set("Origin", "https://www.torque.investments")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "authorization,access-control-request-private-network")
	req.Header.Set("Access-Control-Request-Private-Network", "true")

	rec := httptest.NewRecorder()
	srv.Echo.ServeHTTP(rec, req)

	allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
	if !strings.Contains(strings.ToLower(allowHeaders), "access-control-request-private-network") {
		t.Fatalf("expected Access-Control-Allow-Headers to include access-control-request-private-network, got %q", allowHeaders)
	}

	allowPNA := rec.Header().Get("Access-Control-Allow-Private-Network")
	if allowPNA != "true" {
		t.Fatalf("expected Access-Control-Allow-Private-Network: true, got %q", allowPNA)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "https://www.torque.investments" {
		t.Fatalf("expected Access-Control-Allow-Origin to echo the origin, got %q", allowOrigin)
	}

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rec.Code)
	}
}
