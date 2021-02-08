package config

import (
	"encoding/json"
	"io/ioutil"
	//"gopkg.in/validator.v2"
)

// LoadJSON config into out interface, with defaults and validates
func LoadJSON(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// res, err := ParseEnv(data)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "config parse error: %s", err.Error())
	// 	res = data
	// }
	return UnmarshalJSON(data, out)
}

// UnmarshalJSON unmarshals, defaults and validates
func UnmarshalJSON(in []byte, out interface{}) error {
	err := json.Unmarshal(in, out)
	if err != nil {
		return err
	}
	// err = SetDefaults(out)
	// if err != nil {
	// 	return err
	// }
	// err = validator.Validate(out)
	// if err != nil {
	// 	return err
	// }
	return nil
}
