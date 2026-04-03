package tests

import "testing"

// TestAny word %any% allow to have any string
func TestAny(t *testing.T) {
	CompareJSONBody(
		t,
		[]byte(`{"nc_issuer_id":"10010949", "nc_issuer_status":"Pending"}`),
		[]byte(`{"nc_issuer_id":"%any%", "nc_issuer_status":"Pending"}`),
	)
}

func TestAnyInsideArray(t *testing.T) {
	CompareJSONBody(
		t,
		[]byte(`{"items":[{"id":21,"created_at":"2026-03-30 11:30:00 +0200 CEST","nested":["keep","dynamic"]}]}`),
		[]byte(`{"items":[{"id":"%any%","created_at":"%any%","nested":["keep","%any%"]}]}`),
	)
}
