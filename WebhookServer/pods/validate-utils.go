package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/spyzhov/ajson"
	admission "k8s.io/api/admission/v1"
)

func GetJsonPathCommonVerifications(opType OperationType, jpVerifications []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	/* MandatoryData verification */
	for _, m := range opType["mandatorydata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpVerifications = append(jpVerifications, admissioncontroller.MandatoryDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* ForbiddenData verification */
	for _, m := range opType["forbiddendata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpVerifications = append(jpVerifications, admissioncontroller.ForbiddenDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}
	return jpVerifications
}

func GetJsonPathVerifications(config Config, namespace Namespace, jpVerifications []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	/* MandatoryData verification */
	for _, m := range config[namespace]["mandatorydata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpVerifications = append(jpVerifications, admissioncontroller.MandatoryDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* ForbiddenData verification */
	for _, m := range config[namespace]["forbiddendata"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: MandatoryDataCheck: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpVerifications = append(jpVerifications, admissioncontroller.ForbiddenDataCheckOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}
	return jpVerifications
}

func VerifyValidation(jpVerifications []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) error {
	podBytes, _ := r.Object.MarshalJSON()
	podNode := Byte2Node(podBytes)

	jpVerifications = PatchJPOperations(jpVerifications, podNode)

	var re1 = regexp.MustCompile(`\[`)
	var re2 = regexp.MustCompile(`\]`)

	for _, jpo := range jpVerifications {
		podNodes, _ := podNode.JSONPath(jpo.Path)
		_, err := JSONPath2XPath(jpo, podNodes, re1, re2)

		switch jpo.Op {
		case "mandatorydata":
			if err != nil { // no path found
				fmt.Println("error - mandatory data: missing required data:", jpo.Path)
				return fmt.Errorf("error - mandatory data: missing required data: %v", jpo.Path)
			} else { //path found
				if jpo.Value != nil { // check if the value exists
					switch podNodes[0].Type() {
					case ajson.Bool:
						if jpo.Value != podNodes[0].MustBool() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - mandatory data: missing required data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Numeric:
						if jpo.Value != podNodes[0].MustNumeric() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - mandatory data: missing required data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.String:
						if jpo.Value != podNodes[0].MustString() {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - mandatory data: missing required data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Array:
						var items []interface{}
						items = recursiveCheckTypeArray(podNodes[0])
						if !reflect.DeepEqual(jpo.Value, items) {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - mandatory data: missing required data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Object:
						keys := podNodes[0].Keys()
						item := make(map[string]interface{})
						for _, k := range keys {
							item[k] = recursiveCheckTypeObject(podNodes[0], k)[k]
						}
						if !reflect.DeepEqual(jpo.Value, item) {
							fmt.Println("error - mandatory data: missing required data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - mandatory data: missing required data: %v at: %v", jpo.Value, jpo.Path)
						}
					}
				}
			}
		case "forbiddendata":
			if err == nil { //path found
				if jpo.Value != nil { //check if the value exists
					switch podNodes[0].Type() {
					case ajson.Bool:
						if jpo.Value == podNodes[0].MustBool() {
							fmt.Println("error - forbidden data: found forbidden data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - forbidden data: found forbidden data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Numeric:
						if jpo.Value == podNodes[0].MustNumeric() {
							fmt.Println("error - forbidden data: found forbidden data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - forbidden data: found forbidden data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.String:
						if jpo.Value == podNodes[0].MustString() {
							fmt.Println("error - forbidden data: found forbidden data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - forbidden data: found forbidden data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Array:
						var items []interface{}
						items = recursiveCheckTypeArray(podNodes[0])
						if reflect.DeepEqual(jpo.Value, items) {
							fmt.Println("error - forbidden data: found forbidden data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - forbidden data: found forbidden data: %v at: %v", jpo.Value, jpo.Path)
						}
					case ajson.Object:
						keys := podNodes[0].Keys()
						item := make(map[string]interface{})
						for _, k := range keys {
							item[k] = recursiveCheckTypeObject(podNodes[0], k)[k]
						}
						if reflect.DeepEqual(jpo.Value, item) {
							fmt.Println("error - forbidden data: found forbidden data:", jpo.Value, "at:", jpo.Path)
							return fmt.Errorf("error - forbidden data: found forbidden data: %v at: %v", jpo.Value, jpo.Path)
						}
					}
				} else {
					fmt.Println("error - forbidden data: found forbidden data:", jpo.Path)
					return fmt.Errorf("error - forbidden data: found forbidden data: %v", jpo.Path)
				}
			}
		default:
			log.Printf("- error: Operation: undefined")
		}
	}
	return nil
}
