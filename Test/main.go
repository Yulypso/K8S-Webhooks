package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

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
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return config
}

/* Convert Config to []byte */
func Config2Byte(config Config) []byte {
	bytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("cannot marshal data: %v", err)
	}
	return bytes
}

/* Retrieves JSON/YAML pod bytes */
func ReadFile(pod string) []byte {
	jsonFile, err := os.Open(pod)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("cannot read data: %v", err)
	}
	return bytes
}

/* Convert []byte to Node */
func Byte2Node(bytes []byte) *ajson.Node {
	node, err := ajson.Unmarshal(bytes)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return node
}

/* Convert Node to []byte */
func Node2Byte(node *ajson.Node) []byte {
	bytes, err := ajson.Marshal(node)
	if err != nil {
		log.Fatalf("cannot marshal data: %v", err)
	}
	return bytes
}

func JSONPath2XPath(jpo PatchOperation, podNodes []*ajson.Node) (string, error) {
	path := ""
	jsonPathSplitted := strings.Split(strings.TrimSpace(jpo.Path), ".")
	fmt.Println(jsonPathSplitted)

	for _, item := range jsonPathSplitted[1:] {
		if strings.Contains(item, "?(@") {
			path += "/" + item[:strings.Index(item, "[?")] + "/" + fmt.Sprintf("%v", (podNodes[0].Index()))
		} else {
			path += "/" + item
		}
	}

	// TODO if: contains [*] => Recursive
	if len(podNodes) > 0 {
		if podNodes[0].IsArray() {
			path += "/" + fmt.Sprintf("%v", (len(podNodes[0].MustArray())))
		}
		return path, nil
	}
	return "", errors.New("error: path undefined")
}

/* 1 */
func getJsonPathOperations(config Config, namespace Namespace, jpOperations []PatchOperation) []PatchOperation {

	/* Add operation */
	for _, m := range config[namespace]["add"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: AddOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, AddPatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* Remove operation */
	for _, m := range config[namespace]["remove"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: RemoveOperation: No path specified")
			continue
		}
		jpOperations = append(jpOperations, RemovePatchOperation(fmt.Sprintf("%v", m["path"])))
	}

	/* Replace operation */
	for _, m := range config[namespace]["replace"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: ReplaceOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, ReplacePatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* MandatoryData check */
	for _, m := range config[namespace]["mandatorydata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, MandatoryDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* ForbiddenData check */
	for _, m := range config[namespace]["forbiddendata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, ForbiddenDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}
	return jpOperations
}

func recursiveCheckType(podNodes *ajson.Node) []interface{} {
	var items []interface{}
	for _, item := range podNodes.MustArray() {
		switch item.Type() {
		case ajson.String:
			items = append(items, item.MustString())
		case ajson.Numeric:
			items = append(items, item.MustNumeric())
		case ajson.Bool:
			items = append(items, item.MustBool())
		case ajson.Array:
			items = append(items, recursiveCheckType(item))
		}
	}
	return items
}

/* 2-3 returns sanitized jpOperations */ //TODO : replace filename to r.*admission.AdmissionRequest
func verifyDeployment(jpOperations []PatchOperation, filename string) []PatchOperation {
	var operations []PatchOperation

	podBytes := ReadFile(filename) // podBytes, _ := r.Object.MarshalJSON() au lieu de filename
	podNode := Byte2Node(podBytes)

	for _, jpo := range jpOperations {
		podNodes, _ := podNode.JSONPath(jpo.Path) // read JSONPath
		xPath, err := JSONPath2XPath(jpo, podNodes)

		switch jpo.Op {
		case "add":
			if err == nil {
				operations = append(operations, AddPatchOperation(xPath, jpo.Value))
			}
		case "remove":
			if len(podNodes) > 0 && err == nil {
				operations = append(operations, RemovePatchOperation(xPath))
			}
		case "replace":
			if len(podNodes) > 0 && err == nil {
				operations = append(operations, ReplacePatchOperation(xPath, jpo.Value))
			}
		case "mandatorydata":
			fmt.Println("mandatorydata")
			if err != nil { // no path found
				fmt.Println("error - mandatory data: missing required data:", jpo.Path)
				return operations[:0]
			} else { //path found
				if jpo.Value != nil { // check the value if it exists
					switch podNodes[0].Type() {
					case ajson.Bool:
						if jpo.Value != podNodes[0].MustBool() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return operations[:0]
						}
					case ajson.Numeric:
						if jpo.Value != podNodes[0].MustNumeric() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return operations[:0]
						}
					case ajson.String:
						if jpo.Value != podNodes[0].MustString() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return operations[:0]
						}
					case ajson.Array:
						var items []interface{}
						items = recursiveCheckType(podNodes[0])

						if !reflect.DeepEqual(jpo.Value, items) {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return operations[:0]
						}
						/* TODO: case JSON Object */
						/*case ajson.Object:
						fmt.Println("object")
						fmt.Println(podNodes[0].SetObject(podNodes[0].MustObject()), "***", jpo.Value)

						if jpo.Value != podNodes[0].SetObject(podNodes[0].MustObject()) {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return operations[:0]
						}*/
					}
				}
			}
		case "forbiddendata":
			if err == nil {
				fmt.Println("error - forbidden data: found forbidden data:", jpo.Path)
				return operations[:0]
			}
		default:
			log.Printf("- error: Operation: undefined")
		}
	}
	return operations
}

func main() {

	/* 1 - Get PATCHES -> jpOperations (JSON PATH) */
	config := Byte2Config(ReadFile("default.json"))
	var jpOperations []PatchOperation
	jpOperations = getJsonPathOperations(config, "admissionwebhook-ns", jpOperations)

	/* 2 - validate data [mandatorydata/forbiddendata] */
	/* 3 - jpOperations -> Verify Pod if field already exist or not */
	operations := verifyDeployment(jpOperations, "pod.json")

	/* 4 - Create operation [add/remove/replace] -> add to operations (XPATH) */

	/* 5 - Apply operations */

	fmt.Println("----")
	fmt.Println(operations)
	/* Get Config from default.json */
	/*configBytes := ReadFile("default.json")
	configNode := Byte2Node(configBytes)
	configNodes := configNode.MustKey("admissionwebhook-ns").MustObject()*/
	//configNodes, _ := configNode.JSONPath("$.admissionwebhook-ns.add")

	/* Get Node from pod.json */
	/*podBytes := ReadFile("pod.json")
	podNode := Byte2Node(podBytes)
	podNodes, _ := podNode.JSONPath("$.spec.volumes[*].hostPath")
	fmt.Println(podNodes)*/

	/*podBytes := ReadFile("pod.json")
	podNode := Byte2Node(podBytes)
	podNodes, _ := podNode.JSONPath("$.spec.containers[?(@[name] == 'node-app1')]")
	fmt.Println(podNodes[0].Index())*/

	/* Operate */
	/*for i := 0; i < len(n); i++ {
		fmt.Println(n[i])
		//fmt.Println(o["path"])
		//fmt.Println(o["value"])
		//podNodes[0].AppendObject()
	}*/

	/*for _, configNode := range configNodes {
		configNode.SetArray(configNode.MustArray())
	}
	res := Node2Byte(configNode)
	fmt.Printf("%s", res)*/
}
