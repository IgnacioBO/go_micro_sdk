package user

//Aqui devovleremos los errores que podemos devoler del request
//Por ejemplo bad request, etc

type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	return e.Message
}
