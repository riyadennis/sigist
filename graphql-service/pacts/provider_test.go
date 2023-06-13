package pacts

import (
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var lastName = "" // User doesn't exist

func TestProvider(t *testing.T) {

	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Provider: "UserFeedbackProvider",
	}

	// Verify the Provider using the locally saved Pact Files
	_, err := pact.VerifyProvider(t,
		types.VerifyRequest{
			ProviderBaseURL: "http://localhost:4000",
			PactURLs:        []string{"provider-interactions.json"},
			StateHandlers: types.StateHandlers{
				"User foo exists": func() error {
					lastName = "Doe"
					return nil
				},
			},
		})

	assert.NoError(t, err)
}
