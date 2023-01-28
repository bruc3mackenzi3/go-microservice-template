package integrationtests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type userRequest struct {
	Name  string
	Email string
	Phone string
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func TestCreateUser(t *testing.T) {
	// Each test needs this to avoid being run when the unit test suite is executed
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("Integration tests not enabled with INTEGRATION_TESTS env flag")
	}
	user := userRequest{
		Name:  "Bruce Wayne",
		Email: "batman@gmail.com",
		Phone: "+15199999999",
	}

	statusCode, userResp := createUserRequest(t, user)
	assert.Equal(t, 201, statusCode, "Got wrong status code")
	assert.Equal(t, user.Name, userResp.Name, "Name doesn't match")
	assert.Equal(t, user.Email, userResp.Email, "Email doesn't match")
	assert.Equal(t, user.Phone, userResp.Phone, "Phone doesn't match")

	statusCode, userResp = createUserRequest(t, user)
	assert.Equal(t, 400, statusCode, "Got wrong status code")
}

func createUserRequest(t *testing.T, user userRequest) (int, userResponse) {
	t.Cleanup(func() {
		deleteUser(user.Email)
	})

	var userResponse userResponse

	// Prepare request data
	body, err := json.Marshal(user)
	if err != nil {
		t.Error("Failed to marshal user request to json")
		return 0, userResponse
	}

	// Perform HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/users/", bytes.NewBuffer(body))
	if err != nil {
		t.Error("Failed to create new http.Request object")
		return 0, userResponse
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Error("Failed to perform create user request")
		return resp.StatusCode, userResponse
	}

	if resp.StatusCode == 201 {
		// Parse response data if a user was created
		var respBody []byte
		defer resp.Body.Close()
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Error("Error reading response body")
			return resp.StatusCode, userResponse
		}

		err = json.Unmarshal(respBody, &userResponse)
		if err != nil {
			t.Errorf("Failed to unmarshal json in response body: `%s`", respBody)
			return resp.StatusCode, userResponse
		}
	}
	return resp.StatusCode, userResponse
}
