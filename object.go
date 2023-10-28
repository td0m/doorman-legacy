package doorman

import (
	"fmt"
	"strings"
)

type Object string

func (o Object) Validate() error {
	if !strings.Contains(string(o), ":") {
		return fmt.Errorf("invalid object format")
	}

	return nil
}

func (o Object) Type() string {
	return strings.SplitN(string(o), ":", 2)[0]
}

func (o Object) Value() string {
	return strings.SplitN(string(o), ":", 2)[1]
}
