package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	oldVersion := version
	defer func() {
		version = oldVersion
	}()

	version = "test"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rootHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	expected := "Hello, DevOps! version=test\n"
	if rec.Body.String() != expected {
		t.Fatalf("expected %q, got %q", expected, rec.Body.String())
	}
}
