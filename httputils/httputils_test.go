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

	_, code, err := SendRequest(req)
	if err != nil {
		t.Errorf("cannot send request: %s", err.Error())
		t.FailNow()
	}

	if code != 200 {
		t.Errorf("expected 200 code")
		t.FailNow()
	}
}