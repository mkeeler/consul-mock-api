package examples

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/consul/api"
	mockconsul "github.com/mkeeler/consul-mock-api"
	"reflect"
	"testing"
)

func TestTokenCreate(t *testing.T) {
	m := mockconsul.NewConsul(t)
	defer m.Close()

	m.SetFilteredHeaders([]string{
		"Accept-Encoding",
		"User-Agent",
		"Content-Length",
		"Content-Type",
	})

	// this body will be checked and compared to the body that is sent to the request
	body := &api.ACLToken{
		AccessorID:  "cd45b25a-07cf-4a60-ae71-2a227eaace8e",
		SecretID:    "cd45b25a-07cf-4a60-ae71-2a227eaace8e",
		Description: "test token",
	}

	// Convert the body into a slice of bytes.
	// You must use the json encoder as it appends a newline to the end of the slice.
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(body); err != nil {
		t.FailNow()
	}

	// mock the call to PUT /v1/acl/token
	m.ACLTokenCreate(buf.Bytes(), 200, body).Once() // expect one call to this endpoint

	cfg := api.DefaultConfig()
	cfg.Address = m.URL()

	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("error when creating a new HTTP client: %v", err)
	}

	token, _, err := client.ACL().TokenCreate(body, nil)
	if err != nil {
		t.Fatalf("failed to get get token: %v", err)
	}

	if !reflect.DeepEqual(token, body) {
		t.Fatal("ACL token got from endpoint does not match expected result")
	}
}
