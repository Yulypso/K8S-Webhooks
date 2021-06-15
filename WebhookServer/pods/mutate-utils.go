package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

/* Apply operations */
func getPatches(config Config, namespace Namespace, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {

	/* Add operation */
	for _, m := range config[namespace]["add"] { // For each map in operation "Add"
		for path, value := range m { // for each item within the current map
			operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
		}
	}

	/* Remove operation */
	for _, m := range config[namespace]["remove"] { // For each map in operation "Add"
		for _, value := range m { // for each item within the current map
			operations = append(operations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", value)))
		}
	}

	/* Replace operation */
	for _, m := range config[namespace]["replace"] { // For each map in operation "Add"
		for path, value := range m { // for each item within the current map
			operations = append(operations, admissioncontroller.ReplacePatchOperation(path, value))
		}
	}

	return operations
}

func verifyOperations(operations []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) []admissioncontroller.PatchOperation {

	/* parse pod yaml to JSONPath */
	podBytes, _ := r.Object.MarshalJSON()
	var podJson interface{}
	json.Unmarshal(podBytes, &podJson)

	fmt.Println(podJson)
	return operations
}
