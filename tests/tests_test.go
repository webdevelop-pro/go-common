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
