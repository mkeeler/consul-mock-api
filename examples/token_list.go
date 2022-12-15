package examples

import (
	"testing"

	"github.com/hashicorp/consul/api"
	mockconsul "github.com/mkeeler/consul-mock-api"
)

func TestTokenList(t *testing.T) {
	m := mockconsul.NewConsul(t)
	defer m.Close()

	// http.Get will add both of the headers, but we don't want to care about them.
	m.SetFilteredHeaders([]string{
		"Accept-Encoding",
		"User-Agent",
	})

	// set up an expectation for a GET /v1/acl/tokens
	m.ACLTokenList(200, []*api.ACLTokenListEntry{
		{
			CreateIndex: 1,
			ModifyIndex: 2,
			AccessorID:  "cd45b25a-07cf-4a60-ae71-2a227eaace8e",
			Description: "fake token",
			NodeIdentities: []*api.ACLNodeIdentity{
				{
					NodeName:   "foo",
					Datacenter: "dc1",
				},
			},
		},
	}).Once() // expect this call to be made exactly once

	cfg := api.DefaultConfig()
	cfg.Address = m.URL()

	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("error when creating a new HTTP client: %v", err)
	}

	tokens, _, err := client.ACL().TokenList(nil)
	if err != nil {
		t.Fatalf("error when retrieving token list: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("wrong number of token returned - expected: 1, got: %d", len(tokens))
	}
}
