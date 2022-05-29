package main

import (
	"encoding/json"
)

func jsonDump(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
