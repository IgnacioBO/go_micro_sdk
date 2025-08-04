package course_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/IgnacioBO/go_http_client/client"
	courseSdk "github.com/IgnacioBO/go_micro_sdk/course"
	"github.com/IgnacioBO/gomicro_domain/domain"
)

var header http.Header
var sdk courseSdk.Transport

//Funcion

func TestMain(m *testing.M) {
	//Steamso header
	header = http.Header{}                         //Le asignamos un tupo de dato header
	header.Set("Content-Type", "application/json") //Le definimos el content type
	//Aqui hacemos un newHttpClient para generar el cliente
	sdk = courseSdk.NewHttpClient("base-url", "")
	os.Exit(m.Run())
}

// Primer test
func TestGet_Response404Error(t *testing.T) {
	expectedError := courseSdk.ErrNotFound{Message: "course with id: 1 not found"}
	//Aqui con client.AddMockups agregamos el mock al mapa de mocks, por lo que cuando hagamos get, se teomar como un mock
	//El mock es como una respuesta que le damos al cliente cuando hace un get a un determinado path
	//Cuando uno hace get, el flujo verifica si hay un mock, si hay un mock, lo devuelve, si no hay un mock, hace el get real
	//Aqui estamos creando un mock con data creada dentro de "client.Mock{}"
	err := client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
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

	//Entonces aca llamamos al sdk el get  con id 1 (que es el que mockeamos)
	course, err := sdk.Get("1")

	//Si los errores son distintos generar fallo en el tet
	if !errors.Is(err, expectedError) {
		t.Errorf("expected %v, got %v", expectedError, err)
	}

	//t.Errorf("expected %v, got %v", expectedError, err)

	//Si el curso no es nulo, significa que no lanzo un 404, por lo tanto falla el test
	if course != nil {
		t.Errorf("expected nil, got %v", course)
	}

}

func TestGet_Response500Error(t *testing.T) {
	expectedError := courseSdk.ErrNotFound{Message: "internal server error"}

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
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
	course, err := sdk.Get("1")

	//err = errors.New("Error")
	if err == nil || err.Error() != expectedError.Message {
		t.Errorf("expected err %v, got %v", expectedError, err)
	}

	//Si el curso no es nulo, significa que no lanzo un 404, por lo tanto falla el test
	if course != nil {
		t.Errorf("expected nil, got %v", course)
	}
}

func TestGet_Response4MarshalError(t *testing.T) {
	expectedError := errors.New("unexpected end of JSON input")

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		RespBody:     `{`, //Ponemos un body que no se puede marshalizar (esta mal formado el JSON)
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	course, err := sdk.Get("1")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected err %v, got %v", expectedError, err)
	}

	//Si el curso no es nulo, significa que no lanzo un 404, por lo tanto falla el test
	if course != nil {
		t.Errorf("expected nil, got %v", course)
	}
}

func TestGet_ClientError(t *testing.T) {
	expectedError := errors.New("client error")

	err := client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		Err:          expectedError,
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	course, err := sdk.Get("1")

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected err %v, got %v", expectedError, err)
	}

	if course != nil {
		t.Errorf("expected nil, got %v", course)
	}
}

func TestGet_Success(t *testing.T) {
	startDate, _ := time.Parse("2006-01-02T15:04:05-07:00", "2023-01-01T21:00:00-03:00")
	endDate, _ := time.Parse("2006-01-02T15:04:05-07:00", "2023-01-01T21:00:00-03:00")
	expectedCourse := &domain.Course{
		ID:        "9eec597b-7636-4ef5-81e6-5336b966235f",
		Name:      "Matematicas",
		StartDate: startDate,
		EndDate:   endDate,
	}

	expectedCourseJSON, err := json.Marshal(expectedCourse)
	if err != nil {
		t.Fatalf("failed to marshal expected course: %v", err)
	}

	err = client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		RespBody: fmt.Sprintf(`{
		"message": "success",
		"status": 200,
		"data": %s
		}`, expectedCourseJSON),
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	course, err := sdk.Get("1")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if course == nil {
		t.Errorf("expected a course, got nil")
	}

	if course.ID != expectedCourse.ID {
		t.Errorf("expected course ID %s, got %s", expectedCourse.ID, course.ID)
	}
	if course.Name != expectedCourse.Name {
		t.Errorf("expected course Name %s, got %s", expectedCourse.Name, course.Name)
	}
	if !course.StartDate.Equal(expectedCourse.StartDate) {
		t.Errorf("expected course StartDate %s, got %s", expectedCourse.StartDate, course.StartDate)
	}
	if !course.EndDate.Equal(expectedCourse.EndDate) {
		t.Errorf("expected course EndDate %s, got %s", expectedCourse.EndDate, course.EndDate)
	}

}

func TestGet_CustomToken(t *testing.T) {
	var sdkCustom = courseSdk.NewHttpClient("base-url", "customToken")

	expectedCourse := &domain.Course{
		ID:   "9eec597b-7636-4ef5-81e6-5336b966235f",
		Name: "Matematicas",
	}

	expectedCourseJSON, err := json.Marshal(expectedCourse)
	if err != nil {
		t.Fatalf("failed to marshal expected course: %v", err)
	}

	err = client.AddMockups(&client.Mock{
		URL:          "base-url/courses/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespHeaders:  header,
		RespBody: fmt.Sprintf(`{
		"message": "success",
		"status": 200,
		"data": %s
		}`, expectedCourseJSON),
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	course, err := sdkCustom.Get("1")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if course == nil {
		t.Errorf("expected a course, got nil")
	}

}
