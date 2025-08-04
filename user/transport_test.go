package user_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/IgnacioBO/go_http_client/client"
	userSDK "github.com/IgnacioBO/go_micro_sdk/user"
	"github.com/IgnacioBO/gomicro_domain/domain"
)

var header http.Header
var sdk userSDK.Transport

func TestMain(m *testing.M) {
	header = http.Header{}
	header.Set("Content-Type", "application/json") //Le definimos el content type
	sdk = userSDK.NewHttpClient("base-url", "")
	os.Exit(m.Run())
}

func TestGet_Response404Error(t *testing.T) {
	expectedError := userSDK.ErrNotFound{Message: "user with id: 1 not found"}

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 404,
		RespHeaders:  header,
		RespBody: fmt.Sprintf(`{"status": 404,
				"message":"%s"
				}`, expectedError.Error()),
	})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	user, err := sdk.Get("1")

	if !errors.Is(err, expectedError) {
		t.Errorf("expected %v, got %v", expectedError, err)
	}

	if user != nil {
		t.Errorf("expected nil, got %v", user)
	}

}

func TestGet_Response500Error(t *testing.T) {
	expectedError := userSDK.ErrNotFound{Message: "internal server error"}

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 500,
		RespHeaders:  header,
		RespBody: fmt.Sprintf(`{"status": 500,
				"message":"%s"
				}`, expectedError.Error()),
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	//Entonces aca llamamos al sdk el get  con id 1 (que es el que mockeamos)
	user, err := sdk.Get("1")

	//err = errors.New("Error")
	if err == nil || err.Error() != expectedError.Message {
		t.Errorf("expected err %v, got %v", expectedError, err)
	}

	//Si el curso no es nulo, significa que no lanzo un 404, por lo tanto falla el test
	if user != nil {
		t.Errorf("expected nil, got %v", user)
	}
}

func TestGet_Response4MarshalError(t *testing.T) {
	expectedError := errors.New("unexpected end of JSON input")

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		RespBody:     `{`, //Ponemos un body que no se puede marshalizar (esta mal formado el JSON)
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	user, err := sdk.Get("1")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected err '%v', got '%v'", expectedError, err)
	}

	//Si el curso no es nulo, significa que no lanzo un 404, por lo tanto falla el test
	if user != nil {
		t.Errorf("expected nil, got %v", user)
	}
}

func TestGet_ClientError(t *testing.T) {
	expectedError := errors.New("client error")

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		Err:          expectedError,
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	user, err := sdk.Get("1")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected err %v, got %v", expectedError, err)
	}

	if user != nil {
		t.Errorf("expected nil, got %v", user)
	}
}

func TestGet_Success(t *testing.T) {
	expectedUser := &domain.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
	}

	expectedUserJSON, err := json.Marshal(expectedUser)
	if err != nil {
		t.Fatalf("failed to marshal expected user: %v", err)
	}

	err = client.AddMockups(&client.Mock{
		URL:          "base-url/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		RespBody: fmt.Sprintf(`{
		"message": "success",
		"status": 200,
		"data": %s
		}`, expectedUserJSON),
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	user, err := sdk.Get("1")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if user == nil {
		t.Errorf("expected a course, got nil")
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %s, got %s", expectedUser.ID, user.ID)
	}
	if user.FirstName != expectedUser.FirstName {
		t.Errorf("expected user FirstName %s, got %s", expectedUser.FirstName, user.FirstName)
	}
	if user.LastName != expectedUser.LastName {
		t.Errorf("expected user LastName %s, got %s", expectedUser.LastName, user.LastName)
	}
	if user.Email != expectedUser.Email {
		t.Errorf("expected user Email %s, got %s", expectedUser.Email, user.Email)
	}
	if user.Phone != expectedUser.Phone {
		t.Errorf("expected user Phone %s, got %s", expectedUser.Phone, user.Phone)
	}

}
