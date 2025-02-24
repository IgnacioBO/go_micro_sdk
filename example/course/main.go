package main

import (
	"errors"
	"fmt"
	"os"

	courseSDK "github.com/IgnacioBO/go_micro_sdk/course"
)

func main() {
	courseTrans := courseSDK.NewHttpClient("http://localhost:8002", "")

	course, err := courseTrans.Get("da538f62-6499-4d51-bd14-ee51ca2c70ea")
	if err != nil {
		//Con este errors.As nos servira en ernollment paar pode manejra el status code y indetnificar
		if errors.As(err, &courseSDK.ErrNotFound{}) {
			fmt.Println("Not found:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Internal Server Error", err.Error())
		os.Exit(1)
	}
	fmt.Println(course)

	fmt.Println(course.Name)

	//******Error*****//

	fmt.Println("\n**TEST ERROR**")
	courseTrans = courseSDK.NewHttpClient("http://localhost:8002", "")
	//Aqui aremos un get
	course, err = courseTrans.Get("boexite")
	if err != nil {
		if errors.As(err, &courseSDK.ErrNotFound{}) {
			fmt.Println("Not found: ", err.Error())
			//os.Exit(1)
		}
		//fmt.Println("Interanl Server Error", err.Error())
		//os.Exit(1)
	}

	fmt.Println("\n**TEST ERROR BASE URL NO EXISTE**")
	courseTrans = courseSDK.NewHttpClient("http://localhost:7001", "")
	//Aqui aremos un get
	course, err = courseTrans.Get("boexite")
	if err != nil {
		if errors.As(err, &courseSDK.ErrNotFound{}) {
			fmt.Println("Not found: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Interanl Server Error", err.Error())
		os.Exit(1)
	}
	fmt.Println(course)

	//Aqui puedo usar propieaddes el usuario
	//Asi pdeo accder a todo
	fmt.Println(course.Name)
}
