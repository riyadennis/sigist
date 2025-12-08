package pacts

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/riyadennis/event-management/rest-service/service"
)

func TestClientPact_Local(t *testing.T) {
	// initialize PACT DSL
	pact := dsl.Pact{
		Consumer: "example-client",
		Provider: "example-server",
	}

	// setup a PACT Mock Server
	pact.Setup(true)

	t.Run("get email by id", func(t *testing.T) {
		id := "1"
		createdAt := time.Now().Format(time.RFC3339)
		email := "alice@gmail.com"
		sourceName := "default"
		pact.
			AddInteraction().                                      // specify PACT interaction
			Given("User email exists").                            // specify Provider state
			UponReceiving("email 'alice@gmail.com' is requested"). // specify test case name
			WithRequest(dsl.Request{                               // specify expected request
				Method: "GET",
				Path:   dsl.Term("/emails", ""), // specify matching for endpoint
			}).
			WillRespondWith(dsl.Response{ // specify minimal expected response
				Status: 200,
				Body: dsl.Like([]*service.EmailResponse{ // specify matching for response body
					{
						ID:         id,
						Email:      email,
						SourceName: sourceName,
						CreatedAt:  createdAt,
					},
				}),
			})
		db, err := service.SetUpDB("test.db", "../migrations")
		if err != nil {
			t.Fatal(err)
		}
		// verify interaction on client side
		err = pact.Verify(func() error {
			_, err = service.SaveEmail(
				db, &service.Request{
					Email:   email,
					Sources: []string{sourceName},
				},
				id,
				createdAt,
			)
			if err != nil {
				t.Fatal(err)
			}
			// execute function
			emails, err := service.FetchEmails(db, otelzap.New(zap.NewExample()))
			if err != nil {
				t.Fatal(err)
			}
			t.Log(emails[0].Email)
			// check if actual emails is equal to expected
			if emails == nil || len(emails) == 0 {
				return fmt.Errorf("expected emails with ID %s but got %v", id, emails)
			}

			return err
		})

		if err != nil {
			t.Fatalf("error %v", err)
		}
	})

	// write Contract into file
	if err := pact.WritePact(); err != nil {
		t.Fatal(err)
	}

	// stop PACT mock server
	pact.Teardown()
}
