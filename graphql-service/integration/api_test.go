package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/riyadennis/sigist/graphql-service/graph/model"
)

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

func getUserQuery(id int16) error {
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
