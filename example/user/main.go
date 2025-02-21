package main

//Aqui ejecutaremo el SDK y probaremos si anda bien

import (
	"errors"
	"fmt"
	"os"

	userSDK "github.com/IgnacioBO/go_micro_sdk/user"
)

func main() {
	//Aqui usaremos la funcion NewHttpClient del pacajkge user, le ponemos url base y sin tocket y luego un get
	userTrans := userSDK.NewHttpClient("http://localhost:8001", "")

	//Aqui aremos un get
	user, err := userTrans.Get("ea7ea618-d1e5-48ff-bfc0-6a3e89e28958")
	if err != nil {
		//Con este errors.As nos servira en ernollment paar pode manejra el status code y indetnificar
		if errors.As(err, &userSDK.ErrNotFound{}) {
			fmt.Println("Not found: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Interanl Server Error", err.Error())
		os.Exit(1)
	}
	fmt.Println(user)

	//Aqui puedo usar propieaddes el usuario
	//Asi pdeo accder a todo
	fmt.Println(user.FirstName)

	//******Error*****//

	fmt.Println("\n**TEST ERROR**")
	userTrans = userSDK.NewHttpClient("http://localhost:8001", "")
	//Aqui aremos un get
	user, err = userTrans.Get("boexite")
	if err != nil {
		if errors.As(err, &userSDK.ErrNotFound{}) {
			fmt.Println("Not found: ", err.Error())
			//os.Exit(1)
		}
		//fmt.Println("Interanl Server Error", err.Error())
		//os.Exit(1)
	}

	fmt.Println("\n**TEST ERROR BASE URL NO EXISTE**")
	userTrans = userSDK.NewHttpClient("http://localhost:7001", "")
	//Aqui aremos un get
	user, err = userTrans.Get("boexite")
	if err != nil {
		if errors.As(err, &userSDK.ErrNotFound{}) {
			fmt.Println("Not found: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Interanl Server Error", err.Error())
		os.Exit(1)
	}
	fmt.Println(user)

	//Aqui puedo usar propieaddes el usuario
	//Asi pdeo accder a todo
	fmt.Println(user.FirstName)
}
