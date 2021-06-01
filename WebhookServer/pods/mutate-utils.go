package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
			fmt.Println("Remove: ", value)
			operations = append(operations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", value)))
		}
	}

	/* Add operation */
	for _, m := range config[namespace]["add"] { // For each map in operation "Add"
		for path, value := range m { // for each item within the current map
			fmt.Println("Add: ", path, value)
			operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
		}
	}
	return operations
}
