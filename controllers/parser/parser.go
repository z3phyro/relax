package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"stoicdynamics.com/relax/types"
)

var (
	Variables map[string]string = make(map[string]string)
	Requests  []types.Request   = []types.Request{}
)

func parseVariables(text string) {
	Variables = make(map[string]string)

	pattern := `@(\w+)\s=(.*)\n`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) == 3 {
			Variables[match[1]] = strings.Trim(match[2], " \n")
		}
	}
}

func parseRequests(text string, fileName string) {
	Requests = []types.Request{}

	re := regexp.MustCompile(`###(.*)\n`)

	rawRequests := re.Split(text, -1)[1:]

	for _, rawRequest := range rawRequests {
		for key, value := range Variables {
			rawRequest = strings.ReplaceAll(rawRequest, fmt.Sprintf("{{%s}}", key), value)
		}

		requestMatch := regexp.MustCompile(`(POST|GET|PUT|PATCH|DELETE)\s(\S+)\s(\S+)\n((.+: .+\n)*)(\n)?((.*\n?)*)`).FindStringSubmatch(rawRequest)

		if len(requestMatch) < 5 {
			continue
		}

		name := requestMatch[1]

		Requests = append(Requests, types.Request{
			Raw:      rawRequest,
			Name:     name,
			Verb:     types.Method(string(requestMatch[1])),
			Headers:  requestMatch[4],
			Body:     requestMatch[7],
			Url:      requestMatch[2],
			FileName: fileName,
		})
	}
}

func ParseRequestText(text string, filename string) {
	parseVariables(text)
	parseRequests(text, filename)
}

func IsFileDirectory(fullPath string) bool {
	if fileInfo, err := os.Stat(fullPath); err == nil {
		return fileInfo.IsDir()
	}

	return false
}

func OpenFile(fullPath string) string {
	if IsFileDirectory(fullPath) {
		return ""
	}

	content, err := os.ReadFile(fullPath)
	defer os.Stdin.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func ParseResponse(text string) string {
	if text[0] == '{' {
		var jsonData interface{}
		if err := json.Unmarshal([]byte(text), &jsonData); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
		}

		if value, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
			return string(value)
		} else {
			fmt.Printf("There was an erro %s", err)
		}

	}
	return text
}
