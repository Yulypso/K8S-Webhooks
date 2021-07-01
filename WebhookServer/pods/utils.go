package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/spyzhov/ajson"
)

type Config map[Namespace]OperationType
type Namespace string
type OperationType map[string][]Operation
type Operation map[string]interface{}

/* Convert []byte to Config */
func Byte2Config(bytes []byte) Config {
	var config Config
	err := json.Unmarshal([]byte(bytes), &config)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}
	return config
}

/* Convert Config to []byte */
func Config2Byte(config Config) []byte {
	bytes, err := json.Marshal(config)
	if err != nil {
		log.Printf("cannot marshal data: %v", err)
	}
	return bytes
}

/* Convert []byte to Config */
func Byte2OpType(bytes []byte) OperationType {
	var opType OperationType
	err := json.Unmarshal([]byte(bytes), &opType)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}
	return opType
}

/* Convert Config to []byte */
func OpType2Byte(opType OperationType) []byte {
	bytes, err := json.Marshal(opType)
	if err != nil {
		log.Printf("cannot marshal data: %v", err)
	}
	return bytes
}

/* Retrieves JSON/YAML pod bytes */
func ReadFile(file string) []byte {
	var mutex sync.Mutex
	mutex.Lock()
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	mutex.Unlock()
	if err != nil {
		log.Printf("cannot read data: %v", err)
	}
	return bytes
}

/* Convert []byte to Node */
func Byte2Node(bytes []byte) *ajson.Node {
	node, err := ajson.Unmarshal(bytes)
	if err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}
	return node
}

/* Convert Node to []byte */
func Node2Byte(node *ajson.Node) []byte {
	bytes, err := ajson.Marshal(node)
	if err != nil {
		log.Printf("cannot marshal data: %v", err)
	}
	return bytes
}

/* Append element at a specific index */
func appendAtIndex(array []admissioncontroller.PatchOperation, val admissioncontroller.PatchOperation, i int) ([]admissioncontroller.PatchOperation, error) {
	length := len(array)
	if i > length {
		return array, errors.New("error: index is out of slice range")
	}

	array = append(array, array[length-1])
	copy(array[i+1:], array[i:])
	array[i] = val

	return array, nil
}

func JSONPath2XPath(jpo admissioncontroller.PatchOperation, podNodes []*ajson.Node, re1 *regexp.Regexp, re2 *regexp.Regexp) (string, error) {
	path := ""
	jsonPathSplitted := strings.SplitN(strings.TrimSpace(jpo.Path), ".", -1)

	for _, item := range jsonPathSplitted[1:] {
		if strings.Contains(item, "?(@") {
			path += "/" + item[:strings.Index(item, "[?")] + "/" + fmt.Sprintf("%v", (podNodes[0].Index()))
		} else {
			path += "/" + item
		}
	}

	if len(podNodes) > 0 {
		if podNodes[0].IsArray() {
			path += "/" + fmt.Sprintf("%v", (len(podNodes[0].MustArray())))
		}

		path = re1.ReplaceAllString(path, `/`)
		path = re2.ReplaceAllString(path, ``)

		return path, nil
	}
	path = re1.ReplaceAllString(path, `/`)
	path = re2.ReplaceAllString(path, ``)
	return path, errors.New("error: no path found")
}

func recursiveCheckTypeArray(podNode *ajson.Node) []interface{} {
	var items []interface{}

	for _, item := range podNode.MustArray() {
		switch item.Type() {
		case ajson.String:
			items = append(items, item.MustString())
		case ajson.Numeric:
			items = append(items, item.MustNumeric())
		case ajson.Bool:
			items = append(items, item.MustBool())
		case ajson.Array:
			items = append(items, recursiveCheckTypeArray(item))
		case ajson.Object:
			keys := item.Keys()
			for _, k := range keys {
				items = append(items, recursiveCheckTypeObject(item, k))
			}
		}
	}
	return items
}

func recursiveCheckTypeObject(podNode *ajson.Node, key string) map[string]interface{} {
	item := make(map[string]interface{})

	switch podNode.MustObject()[key].Type() {
	case ajson.String:
		item[key] = podNode.MustObject()[key].MustString()
	case ajson.Numeric:
		item[key] = podNode.MustObject()[key].MustNumeric()
	case ajson.Bool:
		item[key] = podNode.MustObject()[key].MustBool()
	case ajson.Array:
		item[key] = recursiveCheckTypeArray(podNode.MustObject()[key])
	case ajson.Object:
		subItem := make(map[string]interface{})
		keys := podNode.MustObject()[key].Keys()
		for _, k := range keys {
			subItem[k] = recursiveCheckTypeObject(podNode.MustObject()[key], k)[k]
		}
		item[key] = subItem
	}
	return item
}

func PatchJPOperations(jpOperations []admissioncontroller.PatchOperation, podNode *ajson.Node) []admissioncontroller.PatchOperation {
	var patchedJpoperations []admissioncontroller.PatchOperation
	var err error
	nbOp := len(jpOperations)

	for i := 0; i < nbOp; i++ {
		if strings.Contains(jpOperations[i].Path, "[*]") {
			split := strings.SplitN(jpOperations[i].Path, "[*]", -1)
			podNodes, _ := podNode.JSONPath(split[0])

			if len(podNodes) > 0 {
				if podNodes[0].IsArray() {
					ignored := strings.Replace(jpOperations[i].Path, split[0]+"[*]", "", 1)
					for j := 0; j < len(podNodes[0].MustArray()); j++ {
						newPath := split[0] + "[" + fmt.Sprint(j) + "]" + ignored
						switch jpOperations[i].Op {
						case "add":
							jpOperations, err = appendAtIndex(jpOperations, admissioncontroller.AddPatchOperation(fmt.Sprintf("%v", newPath), jpOperations[i].Value), i+1)
							if err != nil {
								log.Println(err)
							}
							nbOp++
						case "remove":
							jpOperations, err = appendAtIndex(jpOperations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", newPath)), i+1)
							if err != nil {
								log.Println(err)
							}
							nbOp++
						case "replace":
							jpOperations, err = appendAtIndex(jpOperations, admissioncontroller.ReplacePatchOperation(fmt.Sprintf("%v", newPath), jpOperations[i].Value), i+1)
							if err != nil {
								log.Println(err)
							}
							nbOp++
						case "mandatorydata":
							jpOperations, err = appendAtIndex(jpOperations, admissioncontroller.MandatoryDataCheckOperation(fmt.Sprintf("%v", newPath), jpOperations[i].Value), i+1)
							if err != nil {
								log.Println(err)
							}
							nbOp++
						case "forbiddendata":
							jpOperations, err = appendAtIndex(jpOperations, admissioncontroller.ForbiddenDataCheckOperation(fmt.Sprintf("%v", newPath), jpOperations[i].Value), i+1)
							if err != nil {
								log.Println(err)
							}
							nbOp++
						default:
							log.Printf("- error: Operation: undefined")
						}
					}
				}
			}
		} else {
			patchedJpoperations = append(patchedJpoperations, jpOperations[i])
		}
	}
	return patchedJpoperations
}
