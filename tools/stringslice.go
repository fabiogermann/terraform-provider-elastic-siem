package tools

import (
	"encoding/json"
	"fmt"
)

type StringSlice []string

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err == nil {
		*s = multi
		return nil
	}

	return fmt.Errorf("invalid value for StringSlice: %s", string(data))
}
