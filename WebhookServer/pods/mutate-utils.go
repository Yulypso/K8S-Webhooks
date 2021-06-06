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
	"strconv"
	"strings"

	"github.com/yalp/jsonpath"
	admission "k8s.io/api/admission/v1"
)

/** UTILS **/

//type Config map[string]Namespace
type Config map[Namespace]OperationType
type Namespace string
type OperationType map[string][]Operation
type Operation map[string]interface{}

/* Parse JSON Patches/*.json files */
func ParseConfig(config string) Config {
	jsonFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	var parsed Config
	json.Unmarshal([]byte(bytes), &parsed)
	return parsed
}

/* Annotate mutations */
func annotate(key string, value string, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	metadata := map[string]string{key: value}
	operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
	return operations
}

/* Apply operations */
func getPatches(config Config, namespace Namespace, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {

	/* Remove operation */
	for _, m := range config[namespace]["remove"] { // For each map in operation "Add"
		for _, value := range m { // for each item within the current map
			operations = append(operations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", value)))
		}
	}

	/* Add operation */
	for _, m := range config[namespace]["add"] { // For each map in operation "Add"
		for path, value := range m { // for each item within the current map
			operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
		}
	}
	return operations
}

/* dsl, pod, dsl (jsonpath)*/
func sanitizeOperation(currentOp int, podData interface{}, jp string, operations []admissioncontroller.PatchOperation, operationsLen int) (admissioncontroller.PatchOperation, string, []admissioncontroller.PatchOperation, int, error) {
	patchOps := operations[currentOp]
	var err error
	matched, _ := regexp.MatchString("^/([/A-Za-z0-9](\\[[\\*\\d]*\\])?(\\[[\\+\\d]*\\])?)+([A-Za-z0-9]|(\\[[\\*\\d]*\\])|(\\[[\\+\\d]*\\]))$", patchOps.Path)

	if matched && patchOps.Op == "add" {
		l := 1
		/* Add k new operations and remove operation containing [*] */
		if strings.Contains(jp, "[*]") {
			li := strings.Index(jp, "[*]")
			pattern := jp[:li]
			l = getSubPathLength(podData, pattern)

			path := JsonPathToPath(jp)
			for j := 0; j < l; j++ {
				if j == 0 {
					path = strings.Replace(path, "*", strconv.Itoa(j), 1)
					err = errors.New("removing [*] operation")
				} else {
					path = strings.Replace(path, strconv.Itoa(j-1), strconv.Itoa(j), 1)
				}
				fmt.Println(path)
				fmt.Println(patchOps.Value)
				operations = append(operations, admissioncontroller.AddPatchOperation(path, patchOps.Value))
				operationsLen++
			}
		} else if strings.Contains(patchOps.Path, "[+]") {
			/* Search for a free index if [+] */
			i := 0
			err = nil
			jp = strings.Replace(jp, "[+]", "["+strconv.Itoa(i)+"]", 1)
			for {
				fmt.Println(jp)
				_, err := jsonpath.Read(podData, jp)
				if err != nil {
					fmt.Println("-> found free id", i)
					patchOps.Path = strings.Replace(patchOps.Path, "[+]", "/"+strconv.Itoa(i), 1)
					return patchOps, jp, operations, operationsLen, nil
				} else {
					fmt.Println("-> id already exist")
					i++
					li := strings.LastIndex(jp, "["+strconv.Itoa(i-1)+"]")
					jp = jp[:li] + strings.Replace(jp[li:], "["+strconv.Itoa(i-1)+"]", "["+strconv.Itoa(i)+"]", 1)
				}
			}
		}
		return patchOps, jp, operations, operationsLen, err
	}
	return patchOps, "", operations, operationsLen, errors.New("error: path regex unmatch: " + patchOps.Path)
}

func removeInvalidOperation(index int, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	return append(operations[:index], operations[index+1:]...)
}

func getOperationPerType(t string, op []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	var operations []admissioncontroller.PatchOperation
	for _, o := range op {
		if o.Op == strings.ToLower(t) {
			operations = append(operations, o)
		}
	}
	return operations
}

func getSubPathLength(podData interface{}, pattern string) int {
	a, _ := jsonpath.Read(podData, "$.spec.containers") // ex: $.spec.container
	b, _ := json.Marshal(a)

	var res []interface{}
	json.Unmarshal(b, &res)
	return len(res)
}

/*
 * To JSONPath
 */
func pathToJsonPath(s string) string {
	jsonPath := "$"
	path := strings.Split(strings.TrimSpace(s), "/")

	for _, item := range path[1:] {
		if _, err := strconv.Atoi(item); err == nil { // if item looks like an integer
			fmt.Printf("%q looks like a number.\n", item)
			jsonPath += "[" + item + "]"
		} else {
			jsonPath += "." + item
		}
	}
	return jsonPath
}

func JsonPathToPath(s string) string {
	path := ""
	jsonPath := strings.Split(strings.TrimSpace(s), ".")

	for _, item := range jsonPath[1:] {
		if strings.Contains(item, "[*]") {
			path += "/" + strings.Replace(item, "[*]", "/*", -1)
		} else {
			path += "/" + item
		}
	}
	return path
}

/*func valueToJsonPath(s string) {
	jsonPath := "$"
}*/

func verifyAdd(op []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) []admissioncontroller.PatchOperation {
	fmt.Print("\nOperation list\n")
	fmt.Println(op)

	fmt.Print("\nOperation Add list\n")
	operations := getOperationPerType("add", op)
	fmt.Println(operations)
	operationsLen := len(operations)

	/* parse pod yaml to JSONPath */
	podBytes, _ := r.Object.MarshalJSON()
	var podData interface{}
	json.Unmarshal(podBytes, &podData)

	for i := 0; i < operationsLen; i++ {
		fmt.Printf("\n*** Operation [%d/%d] ***\n", i+1, operationsLen)
		fmt.Print("Op: ")
		fmt.Println(operations[i].Op)
		fmt.Print("Path: ")
		fmt.Println(operations[i].Path)
		fmt.Print("Value: ")
		fmt.Println(operations[i].Value)
		fmt.Print("* * * * * * * * * * * *\n")

		/* convert operation path to JSONPath */
		jp := pathToJsonPath(operations[i].Path)
		fmt.Println(jp)

		/* verify operation pattern */
		var err error
		operations[i], jp, operations, operationsLen, err = sanitizeOperation(i, podData, jp, operations, operationsLen)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("Removing operation: ")
			fmt.Println(operations[i])
			operations = removeInvalidOperation(i, operations)
			i--
			operationsLen--
		} else {
			/* try to read in pod JSONPath */
			_, err = jsonpath.Read(podData, jp)
			if err != nil {
				fmt.Println("field doesn't exist: good to add")
			} else {
				fmt.Println("field already exist: bad to add")
				fmt.Printf("Removing operation: ")
				fmt.Println(operations[i])
				operations = removeInvalidOperation(i, operations)
				i--
				operationsLen--
			}
		}
	}
	return operations
}

func verifyRemove(operations []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) []admissioncontroller.PatchOperation {
	return operations
}
