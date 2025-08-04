package course

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	clientGoHttp "github.com/IgnacioBO/go_http_client/client" //clientGoHttp es el CLIENT del repo_go_http cloente pero le pusimos un alias para no confundir con otro client
	"github.com/IgnacioBO/gomicro_domain/domain"
)

//Aqui definiremos struct e intrface

type (
	//Struct Data respone, es lo que RECIBIREMSO AL CLIENTE que le pego (User por ejemplo) [El client http y user]
	DataResponse struct {
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Meta    interface{} `json:"meta"`
	}

	Transport interface {
		Get(id string) (*domain.Course, error)
	}

	clientHTTP struct {
		client clientGoHttp.Transport
	}
)

func NewHttpClient(baseURL, token string) Transport {
	header := http.Header{}

	if token != "" {
		header.Set("Authorization", token)
	}

	return &clientHTTP{
		client: clientGoHttp.New(header, baseURL, 5000*time.Millisecond, true),
	}

}

func (c *clientHTTP) Get(id string) (*domain.Course, error) {

	dataResponse := DataResponse{Data: &domain.Course{}}

	u := url.URL{}
	u.Path += fmt.Sprintf("/courses/%s", id)

	reps := c.client.Get(u.String())

	if reps.Err != nil {
		return nil, reps.Err
	}

	//Fillup le pasamos el datarespone para transformar la response de struct a json (o xml)
	//Esto ahce que dataResponse.
	if err := reps.FillUp(&dataResponse); err != nil {
		return nil, err
	}

	//Si deuvle 400 es notfound
	if reps.StatusCode == 404 {
		return nil, ErrNotFound{fmt.Sprintf("%s", dataResponse.Message)}
	}

	if reps.StatusCode > 299 {
		return nil, fmt.Errorf("%s", dataResponse.Message)
	}

	return dataResponse.Data.(*domain.Course), nil
}
