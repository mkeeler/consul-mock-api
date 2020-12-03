package mockconsul

import "github.com/mkeeler/mock-http-api"

type Consul struct {
	*mockapi.MockAPI
}

func NewConsul(t mockapi.TestingT) *Consul {
	return &Consul{
		MockAPI: mockapi.NewMockAPI(t),
	}
}
