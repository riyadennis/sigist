package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/riyadennis/sigist/graphql-service/graph/model"
)

var (
	opts = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}
	hostUrl = "http://localhost:8081/graphql"
)

type Response struct {
	Data struct {
		SaveUserFeedback struct {
			ID        string `json:"id"`
			FirstName string `json:"firstName omitempty"`
		} `json:"SaveUserFeedback"`
	} `json:"data"`
}

type GetUserResponse struct {
	Data struct {
		GetUserFeedback []struct {
			ID        int    `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Feedback  string `json:"feedback"`
		} `json:"GetUserFeedback"`
	} `json:"data"`
}

func saveUserFeedbackMutation(user *model.User) (*Response, error) {
	mutation := `{"query":"mutation {\n\tSaveUserFeedback(input: {\n\t\tfirstName: \"` + user.FirstName + `\"\n\t\tlastName: \"` + user.LastName + `\"\n\t\temail: \"` + user.Email + `'\"\n\t\tjobTitle: \"` + user.JobTitle + `\"\n\t\tfeedback: \"` + user.Feedback + `\"\n\t}\n\t){\n\t\tid\n\t}\n}"}`
	resp, err := http.Post(hostUrl, "application/json", strings.NewReader(mutation))
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to save user: %s", resp.Status)
	}

	response := &Response{}
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getUserQueryByName(firstName string) (*GetUserResponse, error) {
	body := fmt.Sprintf(`{"query":"\n{\n\tGetUserFeedback(filter: {\n\t\tfirstName: \"` + firstName + `\"\n\t}){\n\t\tfirstName\n\t\tfeedback\n\t}\n}","variables":{}}`)
	resp, err := http.Post(hostUrl, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := &GetUserResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
