# consul-mock-api
Consul Mock HTTP API helpers

## Usage

```go
package somepkg

import (
   "testing"

   "github.com/hashicorp/consul/api"
   mockconsul "github.com/mkeeler/consul-mock-api"
)

func TestTokenList(t *testing.T) {
   m := mockconsul.NewConsul(t)

   // http.Get will add both of the headers but we don't want to care about them.
   m.SetFilteredHeaders([]string{
      "Accept-Encoding",
      "User-Agent",
   })

   // set up an expectation for a GET /v1/acl/tokens
   m.ACLTokenList(200, []*api.ACLTokenListEntry{
      CreateIndex: 1,
      ModifyIndex: 2,
      AccessorID:  "cd45b25a-07cf-4a60-ae71-2a227eaace8e",
      Description: "fake token",
      NodeIdentities: []*api.ACLNodeIdentity{
         {
             NodeName: "foo",
             Datacenter: "dc1",
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
   
   // no need to do anything if using Go 1.14 as the expectations will be checked 
   // and the test HTTP server will be closed automatically. If not using Go 1.14 then
   // you should defer m.Close() after its creation
}
```