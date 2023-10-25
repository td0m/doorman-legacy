package doorman

type Verb string

type Role struct {
	ID    string
	Verbs []Verb
}

func NewRole(id string, optverbs ...[]Verb) Role {
	verbs := []Verb{}
	if len(optverbs) > 0 {
		verbs = optverbs[0]
	}
	return Role{ID: id, Verbs: verbs}
}

