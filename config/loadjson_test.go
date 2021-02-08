package config

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	fmt.Println("")

	param := make(map[string]interface{})
	LoadJSON("./default.json", &param)
	fmt.Println(param)
}
