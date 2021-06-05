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
func sanitizeOperation(patchOps admissioncontroller.PatchOperation, podData interface{}, jp string) (admissioncontroller.PatchOperation, string, error) {
	matched, _ := regexp.MatchString("^/([/A-Za-z0-9](\\[[\\*\\d]*\\])?(\\[[\\+\\d]*\\])?)+([A-Za-z0-9]|(\\[[\\*\\d]*\\])|(\\[[\\+\\d]*\\]))$", patchOps.Path)

	if matched {
		/* Search for a free index if [+] */
		if patchOps.Op == "add" && strings.Contains(patchOps.Path, "[+]") {
			i := 0
			jp = strings.Replace(jp, "[+]", "["+strconv.Itoa(i)+"]", 1)
			for {
				fmt.Println(jp)
				_, err := jsonpath.Read(podData, jp)
				if err != nil {
					fmt.Println("-> found free id", i)
					//TODO: Also replace [*] and add new operation
					patchOps.Path = strings.Replace(patchOps.Path, "[+]", "/"+strconv.Itoa(i), 1)
					return patchOps, jp, nil
				} else {
					fmt.Println("-> id already exist")
					i++
					li := strings.LastIndex(jp, "["+strconv.Itoa(i-1)+"]")
					jp = jp[:li] + strings.Replace(jp[li:], "["+strconv.Itoa(i-1)+"]", "["+strconv.Itoa(i)+"]", 1)
					//jp = strings.Replace(jp, "["+strconv.Itoa(i-1)+"]", "["+strconv.Itoa(i)+"]", 1)
				}
			}
		}
		return patchOps, jp, nil
	}
	return patchOps, "", errors.New("error: path regex unmatch: " + patchOps.Path)
}

func removeIncorrectOperation(index int, operation []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	return append(operation[:index], operation[index+1:]...)
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

/*func valueToJsonPath(s string) {
	jsonPath := "$"
}*/

func verifyAdd(op []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) []admissioncontroller.PatchOperation {
	fmt.Print("\nOperation list\n")
	fmt.Println(op)

	fmt.Print("\nOperation Add list\n")
	operations := getOperationPerType("add", op)
	fmt.Println(operations)
	operationLen := len(operations)
	for i := 0; i < operationLen; i++ {
		fmt.Printf("\n*** Operation [%d/%d] ***\n", i+1, operationLen)
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

		/* parse pod yaml to JSONPath */
		podBytes, _ := r.Object.MarshalJSON()
		var podData interface{}
		json.Unmarshal(podBytes, &podData)

		/* verify operation pattern */
		var err error
		operations[i], jp, err = sanitizeOperation(operations[i], podData, jp)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("Removing operation: ")
			fmt.Println(operations[i])
			operations = removeIncorrectOperation(i, operations)
			i--
			operationLen--
		} else {
			/* try to read in pod JSONPath */
			_, err = jsonpath.Read(podData, jp)
			if err != nil {
				fmt.Println("field doesn't exist: good to add")
			} else {
				fmt.Println("field already exist: bad to add")
				fmt.Printf("Removing operation: ")
				fmt.Println(operations[i])
				operations = removeIncorrectOperation(i, operations)
				i--
				operationLen--
			}
		}
	}
	return operations
}

func verifyRemove(operations []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) []admissioncontroller.PatchOperation {
	return operations
}
