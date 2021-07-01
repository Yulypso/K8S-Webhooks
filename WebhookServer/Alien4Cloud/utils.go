package Alien4Cloud

import (
	"strings"

	"github.com/spyzhov/ajson"
)

func CheckPathPattern(path string, podNodes []*ajson.Node) []string {
	var frontPathSplitted []string
	var midPath string
	var jsonPathSplitted []string

	frontPath := path[:strings.Index(path, ".[")]
	midPath = path[strings.Index(path, "[") : strings.Index(path, "]")+1]

	frontPathSplitted = strings.SplitN(strings.TrimSpace(frontPath), ".", -1)
	midPath = strings.Replace(midPath, "/", "~1", -1)

	jsonPathSplitted = append(jsonPathSplitted, frontPathSplitted...)
	jsonPathSplitted = append(jsonPathSplitted, midPath)

	return jsonPathSplitted
}
