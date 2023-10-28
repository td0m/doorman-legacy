package doorman

import "encoding/json"

type Change struct {
	Type    string
	Payload json.RawMessage
}
