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
		SaveUser struct {
			ID        int16  `json:"id"`
			FirstName string `json:"firstName omitempty"`
		} `json:"saveUser"`
	} `json:"data"`
}

type GetUserResponse struct {
	Data struct {
		GetUser []struct {
			ID        int    `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
		} `json:"GetUser"`
	} `json:"data"`
}

func saveUserMutation(user *model.User) (*Response, error) {
	mutation := `{"query":"mutation {\n\tsaveUser(input: {\n\t\tfirstName: \"` + user.FirstName + `\"\n\t\tlastName: \"` + user.LastName + `\"\n\t\temail: \"` + user.Email + `\"\n\t\tjobTitle: \"` + user.JobTitle + `\"\n\t}\n\t){\n\t\tid\n\t}\n}"}`
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
	body := fmt.Sprintf(`{"query":"\n{\n\tGetUser(filter: {\n\t\tfirstName: \"%s\"\n\t}){\n\t\tid\n\t\tfirstName\n\t\tlastName\n\t\temail\n\t}\n}","variables":{}}}`, firstName)
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

func getUserQueryWithID(id int64) error {
	body := fmt.Sprintf(`{"query":"\n{\n\tGetUser(filter: {\n\t\tid: %d\n\t}){\n\t\tid\n\t\tfirstName\n\t\tlastName\n\t\temail\n\t}\n}","variables":{}}`, id)
	resp, err := http.Post(hostUrl, "application/json", strings.NewReader(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get user")
	}

	return nil
}
