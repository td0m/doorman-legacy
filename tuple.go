package doorman

type Tuple struct {
	Subject Object
	Role    string
	Object  Object
}

func NewTuple(sub Object, role string, obj Object) Tuple {
	return Tuple{sub, role, obj}
}

type Connection struct {
	Role   string
	Object Object
}

type Path []Connection

