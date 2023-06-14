package integration

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/kafka"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/sigist/graphql-service/graph/model"
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/riyadennis/sigist/graphql-service/service"
)

var kafkaContainer *gnomock.Container

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestUsers(t *testing.T) {
	defer func() {
		assert.NoError(t, gnomock.Stop(kafkaContainer))
	}()

	flag.Parse()
	opts.Paths = flag.Args()
	userTest := &UserTest{}
	status := godog.TestSuite{
		Name:                "user tests",
		ScenarioInitializer: userTest.InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}.Run()
	if status != 0 {
		t.Errorf("falied to run tests got status %d", status)
	}
}

type UserTest struct {
	APIResponse *Response
}

func (u *UserTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(BeforeScenario)
	ctx.Step(`^he add his details and feedback as below:$`, u.heAddedHisDetailsAndFeedbackAsBelow)
	ctx.Step(`^there should be a user called "([^"]*)" saved in the system with feedback "([^"]*)"$`, u.thereShouldBeAUserCalledSavedInTheSystemWithFeedback)
	ctx.Step(`^"([^"]*)" is a user$`, u.isAUser)
}

func BeforeScenario(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	var err error
	kafkaContainer, err = gnomock.Start(
		kafka.Preset(kafka.WithTopics("data-pipe")),
		gnomock.WithDebugMode(), gnomock.WithLogWriter(os.Stdout),
		gnomock.WithContainerName("db-kafka-test"),
	)

	if err != nil {
		log.Fatal("failed to start kafka", err)
		return ctx, err
	}

	config := internal.Config{
		Env:            "test",
		Port:           ":8081",
		LogLevel:       "debug",
		DBFile:         "test.db",
		MigrationsPath: "../migrations",
		KafkaBroker:    kafkaContainer.DefaultAddress(),
		KafkaTopic:     "data-pipe",
	}

	newService, err := service.NewService(config)
	if err != nil {
		log.Fatal("failed to initialise newService", err)
		return ctx, err
	}

	err = newService.Start()
	if err != nil {
		log.Fatal("failed to start newService", err)
		return ctx, err
	}

	return ctx, nil
}

func (u *UserTest) heAddedHisDetailsAndFeedbackAsBelow(arg1 *godog.Table) error {
	var err error
	u.APIResponse, err = saveUserFeedbackMutation(
		&model.User{
			FirstName: arg1.Rows[1].Cells[0].Value,
			LastName:  arg1.Rows[1].Cells[1].Value,
			Email:     arg1.Rows[1].Cells[2].Value,
			JobTitle:  arg1.Rows[1].Cells[3].Value,
		},
	)
	return err
}

func (u *UserTest) thereShouldBeAUserCalledSavedInTheSystemWithFeedback(arg1, arg2 string) error {
	response, err := getUserQueryByName(arg1)
	if err != nil {
		return err
	}

	for _, rsp := range response.Data.GetUserFeedback {
		if rsp.FirstName != arg1 {
			return errors.New("user not found")
		}
		if rsp.Feedback == arg2 {
			return errors.New("feedback not found")
		}
	}

	return nil
}

func (u *UserTest) isAUser(arg1 string) error {
	response, err := getUserQueryByName(arg1)
	if err != nil {
		return err
	}

	for _, rsp := range response.Data.GetUserFeedback {
		if rsp.FirstName == arg1 {
			return nil
		}
	}

	return nil
}
