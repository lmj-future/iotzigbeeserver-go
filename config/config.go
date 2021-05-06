package config

import (
	"encoding/json"
	"io/ioutil"
)

// LoadJSON config into out interface, with defaults and validates
func LoadJSON(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return UnmarshalJSON(data, out)
}

// UnmarshalJSON unmarshals, defaults and validates
func UnmarshalJSON(in []byte, out interface{}) error {
	err := json.Unmarshal(in, out)
	if err != nil {
		return err
	}
	return nil
}
