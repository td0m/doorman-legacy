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
	return strings.Split(string(o), ":")[0]
}
