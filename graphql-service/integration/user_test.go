package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	flag "github.com/spf13/pflag"

	"github.com/riyadennis/sigist/graphql-service/graph/model"
	"github.com/riyadennis/sigist/graphql-service/internal"
	"github.com/riyadennis/sigist/graphql-service/service"
)

var (
	opts = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}
	hostUrl = "http://localhost:8081/graphql"
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

type Response struct {
	Data struct {
		SaveUser struct {
			ID        int16  `json:"id"`
			FirstName string `json:"firstName omitempty"`
		} `json:"saveUser"`
	} `json:"data"`
}

func (u *UserTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(BeforeScenario)
	ctx.Step(`^he sign up with details below:$`, u.heSignUpWithDetailsBelow)
	ctx.Step(`^there should be (\d+) user called "([^"]*)"$`, u.thereShouldBeUserCalled)
	ctx.Step(`^"([^"]*)" is a user$`, u.wantsToUseTheSystem)
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

func (u *UserTest) thereShouldBeUserCalled(arg1 int, arg2 string) error {
	createdUserID := u.APIResponse.Data.SaveUser.ID
	if createdUserID == 0 {
		log.Fatal("no user created")
	}

	return getUserQuery(createdUserID)
}

func (u *UserTest) wantsToUseTheSystem(arg1 string) error {

	return nil
}
