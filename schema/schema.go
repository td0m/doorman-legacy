package schema

import "github.com/td0m/doorman"

type Schema struct {
	Types map[string]Type
}

type Type map[string]Relation

type Relation struct {
	Computed ComputedRelation
}

type ComputedRelation struct {
}

func (c ComputedRelation) ToSet(contextualElement doorman.Element) doorman.Set {
	return doorman.Set{}
}
