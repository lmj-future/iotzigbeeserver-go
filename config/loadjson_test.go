package config

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	param := make(map[string]interface{})
	LoadJSON("./default.json", &param)
	fmt.Println(param)
}
