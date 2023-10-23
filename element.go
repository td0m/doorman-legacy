package doorman

import "strings"

type Element string

func (e Element) Type() string {
	return strings.SplitN(string(e), ":", 2)[0]
}
