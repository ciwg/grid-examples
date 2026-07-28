package grid

import "encoding/json"

func marshalIndentJSON(value any) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}

func unmarshalJSON(body []byte, value any) error {
	return json.Unmarshal(body, value)
}
