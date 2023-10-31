package parser

import (
	"testing"
)

func TestParseText(t *testing.T) {
	file := `@host = http://127.0.0.1:3000
### 
POST {{host}}/goto/Kixwrc2g HTTP/1.1
Content-Type: text/html; charset=utf-8

{
	"here": "value"
}`

	ParseRequestText(file)

	if len(Variables) != 1 {
		t.Errorf("Not enough variables parsed found %d needed %d", len(Variables), 1)
	}

	if val, ok := Variables["host"]; !ok {
		t.Errorf("Variable %s not present", "host")
	} else {
		url := "http://127.0.0.1:3000"
		if val != url {
			t.Errorf("Variable %s didn't get the right value. Found %s needed %s", "host", val, url)
		}
	}

	if len(Requests) != 1 {
		t.Errorf("Not enough requests")
	}

	if Requests[0].Name != "POST" {
		t.Errorf("Got name %s instead of POST", Requests[0].Name)
	}

	body := `{
	"here": "value"
}`
	if Requests[0].Body != body {
		t.Errorf("Got %s instead of %s as body", Requests[0].Body, body)
	}

	// Test Parse response
	input := `{"this":"true"}`
	expectedOutput := `{
  "this": "true"
}`

	if ParseResponse(input) != expectedOutput {
		t.Errorf("Marshalling not working got %s instead of %s", ParseResponse(input), expectedOutput)
	}
}
