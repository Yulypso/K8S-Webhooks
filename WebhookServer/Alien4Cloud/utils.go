package Alien4Cloud

import (
	"strings"

	"github.com/spyzhov/ajson"
)

func CheckPathPattern(path string, podNodes []*ajson.Node) []string {
	var frontPathSplitted []string
	var midPath string
	var jsonPathSplitted []string

	if strings.Contains(path, ".[") && strings.Contains(path, "]") {
		frontPath := path[:strings.Index(path, ".[")]
		midPath = path[strings.Index(path, "[") : strings.Index(path, "]")+1]

		frontPathSplitted = strings.SplitN(strings.TrimSpace(frontPath), ".", -1)
		midPath = strings.Replace(midPath, "/", "~1", -1)

		jsonPathSplitted = append(jsonPathSplitted, frontPathSplitted...)
		jsonPathSplitted = append(jsonPathSplitted, midPath)
	} else {
		jsonPathSplitted = strings.SplitN(strings.TrimSpace(path), ".", -1)
	}
	return jsonPathSplitted
}

func CheckValuePattern() {

}
