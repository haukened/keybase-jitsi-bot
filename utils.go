package main

import (
	"encoding/json"
)

// this JSON pretty prints errors and debug
func p(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "  ")
	return string(s)
}
