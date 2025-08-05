package mock

import (
	"testing"

	"github.com/IgnacioBO/go_micro_sdk/course"
)

//Que pasa si un desarrollador genera otra funcion que no esta en el interface?
//Por ejem plo ahoar solo esta Get, pero ahora se le ocurre agregar una funcion GetAll, que no esta en el interface
//Entonces va ha hacer que el mock falle, porque no esta implementando el interface
// Entonces si queremos agregar una nueva funcion, tenemos que agregarla al interface y al mock

//entonces para que esto no pase haremos un test que verifique que el mock implementa el interface, y si no significa que el mock deberia ser modificado
//Esto dara un error de build si el mock no implementa el interface

func TestMock_Course(t *testing.T) {
	t.Run("should implement course.Transport interface", func(t *testing.T) {
		//Lo que hacemmos es asignarle a una variable (_) del tipo course.Transport el mock, y si no falla, significa que implementa el interface
		var _ course.Transport = (*CourseSdkMock)(nil) // Aqui verificamos que CourseSdkMock implementa el interface course.Transport
		//Aqui (*CourseSdkMock)(nil) es una conversión de tipo: toma el valor nil (sin tipo) y lo “etiqueta” como *CourseSdkMock. -> Es lo mismo que cuando hacemos por ejemplo int64(21)
		//En la plactica podemos pro ejemplo usarlo asi user = User(dto) -> Aqui convierto dto a User, y si son identicos (mismos campos, tipos, etc) no dara error
		//Si courseSdkMock no implementa el interface course.Transport, esto causará un error de compilación.

		//Esto es similar a:
		//var _ course.Transport = &CourseSdkMock{}
		//La diferencia es que el primero no hay asignacion de memoria
	})

}
