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

type Config map[string]interface{}

/* Parse JSON Patches/*.json files */
func ParseConfig(configName string) Config {
	configPath := "../../Patches/"
	jsonFile, err := os.Open(configPath + configName)
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
	fmt.Println("Log: Add /metadata/annotations...")
	metadata := map[string]string{key: value}
	operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
	return operations
}

/* Apply operations */
func getPatches(config Config, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	// Operate
	for path, value := range config { // for each patch within the config file
		operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
	}

	// Add annotation
	operations = annotate("mutate-CREATE", "mutateRunAsRoot", operations)
	return operations
}
