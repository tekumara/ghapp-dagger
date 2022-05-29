package main

import (
	"encoding/json"
	"regexp"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

// remove ansi color control codes
// stolen from https://github.com/acarl005/stripansi/blob/2749a05/stripansi.go
func StripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

func jsonDump(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
