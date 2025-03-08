package httputils

import (
	"context"
	"testing"
)

func TestSendRequest(t *testing.T) {
	ctx := context.Background()
	req, err := CreateDefaultRequest(
		ctx,
		Request{
			Host:   "google.com",
			Scheme: "https",
			Method: "GET",
			Path:   "/",
			Body:   []byte{},
		},
	)
	if err != nil {
		t.Errorf("cannot create default request: %s", err.Error())
		t.FailNow()
	}

	_, headers, code, err := SendRequest(req)
	if err != nil {
		t.Errorf("cannot send request: %s", err.Error())
		t.FailNow()
	}

	if code != 200 {
		t.Errorf("expected 200 code")
		t.FailNow()
	}

	if headers.Get("Content-Type") != "text/html; charset=ISO-8859-1" {
		t.Errorf("expected text/html; charset=ISO-8859-1 from google")
		t.FailNow()
	}
}
