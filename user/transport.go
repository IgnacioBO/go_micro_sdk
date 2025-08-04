package user

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

	//Interfac Trasnport que tiene LOS METODOS A MANEJAR; en este caso SOLO el get, se eniv el id, y devivel el user y error
	Transport interface {
		Get(id string) (*domain.User, error)
	}

	//El cliente va hacer el Tranpsoert que definimos en el repo go_http_cliente
	//clientGoHttp es el CLIENT del repo_go_http cloente pero le pusimos un alias para no confundir con otro client
	//Ojo ClientHTTP es un struct QUE IMPLEMENTE EL INTERFACE Transport DE ESTE PACKAGE! Dentro de este tiene client que IMPLENTE EL TRANPOSERT DEL REPO go_http_client
	clientHTTP struct {
		client clientGoHttp.Transport
	}
)

// Generamremos una funcion NewHttpClient
// Funcion recibira la baseURL y un token (si se necesita)
func NewHttpClient(baseURL, token string) Transport {
	header := http.Header{}

	if token != "" {
		header.Set("Authorization", token)
	}

	//Retornamos un clientHttp que tiene un Transport usando New (que recibe header, baseurl, timeout, y log)
	//Ese es un struct que tiene metodso (get, post, etc) mas otros campos como los header, body, etc  -> ESO LO HCIMOS con el package go_http_cloente que
	return &clientHTTP{
		client: clientGoHttp.New(header, baseURL, 5000*time.Millisecond, true),
	}

}

// Metodo Get
// Recibo un id y devivle user y un id
func (c *clientHTTP) Get(id string) (*domain.User, error) {

	//Aqui definimos una stuct que sera DtaREspone y dentro de fata le pasaremos un User vacio
	//Esto para que la trnapofracion del marsh fucnoia bien
	dataResponse := DataResponse{Data: &domain.User{}}

	//Aqui isaremos el packge url de go
	u := url.URL{}
	//para agregar el path, query params, etc
	//Aqui "hardcodeamos" el path del user + id
	u.Path += fmt.Sprintf("/users/%s", id)

	//Ahora usaremo el get del client (que el el GET del go_htt_client)
	reps := c.client.Get(u.String())

	//si la resp hay error (el package go_htt client) se deuevle
	if reps.Err != nil {
		return nil, reps.Err
	}

	//Fillup le pasamos el datarespone para transformar la response de struct a json (o xml)
	//Esto ahce que dataResponse se llene con los datos de respuesta
	if err := reps.FillUp(&dataResponse); err != nil {
		return nil, err
	}

	//Si deuvle 400 es notfound
	if reps.StatusCode == 404 {
		return nil, ErrNotFound{fmt.Sprintf("%s", dataResponse.Message)}
	}

	//Si deuvle si es 300 o superior esu nerror no controlado
	//poudeo poner mas si queiro
	if reps.StatusCode > 299 {
		return nil, fmt.Errorf("%s", dataResponse.Message)
	}

	//Casteo del Data de la respuesta dle serviocn, con el objeto user
	//Esto porque devovlere eso (dataREsponse.Data es un interface{}, por eso casteamos a user)
	return dataResponse.Data.(*domain.User), nil
}
