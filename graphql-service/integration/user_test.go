package integration

import (
	"context"
	"log"
	"testing"

	"github.com/cucumber/godog"
	flag "github.com/spf13/pflag"

	"github.com/riyadennis/sigist/graphql-service/graph/model"
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/riyadennis/sigist/graphql-service/service"
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestUsers(t *testing.T) {
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
	ctx.Step(`^he sign up with details below:$`, u.heSignUpWithDetailsBelow)
	ctx.Step(`^"([^"]*)" is a user$`, u.isAUser)
	ctx.Step(`^there should be a user called "([^"]*)"$`, u.thereShouldBeAUserCalled)
}

func BeforeScenario(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	config := internal.Config{
		Env:            "test",
		Port:           ":8081",
		LogLevel:       "debug",
		DBFile:         "test.db",
		MigrationsPath: "../migrations",
	}
	newService, err := service.NewService(config)
	if err != nil {
		log.Fatal("failed to initialise newService", err)
	}
	err = newService.Start()
	if err != nil {
		log.Fatal("failed to start newService", err)
	}

	return ctx, nil
}

func (u *UserTest) heSignUpWithDetailsBelow(arg1 *godog.Table) error {
	var err error
	u.APIResponse, err = saveUserMutation(
		&model.User{
			FirstName: arg1.Rows[1].Cells[0].Value,
			LastName:  arg1.Rows[1].Cells[1].Value,
			Email:     arg1.Rows[1].Cells[2].Value,
			JobTitle:  arg1.Rows[1].Cells[3].Value,
		},
	)
	return err
}

func (u *UserTest) thereShouldBeAUserCalled(arg1 string) error {
	response, err := getUserQueryByName(arg1)
	if err != nil {
		return err
	}

	for _, rsp := range response.Data.GetUser {
		if rsp.FirstName == arg1 {
			return nil
		}
	}

	return nil
}

func (u *UserTest) isAUser(arg1 string) error {
	response, err := getUserQueryByName(arg1)
	if err != nil {
		return err
	}

	for _, rsp := range response.Data.GetUser {
		if rsp.FirstName == arg1 {
			return nil
		}
	}

	return nil
}
