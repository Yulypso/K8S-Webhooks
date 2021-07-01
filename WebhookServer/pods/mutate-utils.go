package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"
	"regexp"

	admission "k8s.io/api/admission/v1"
)

func GetJSONPathCommonOperations(opType OperationType, jpOperations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	/* Add operation */
	for _, m := range opType["add"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: AddOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.AddPatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* Remove operation */
	for _, m := range opType["remove"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: RemoveOperation: No path specified")
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", m["path"])))
	}

	/* Replace operation */
	for _, m := range opType["replace"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: ReplaceOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.ReplacePatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}
	return jpOperations
}

func GetJsonPathOperations(config Config, namespace Namespace, jpOperations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	/* Add operation */
	for _, m := range config[namespace]["add"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: AddOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.AddPatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}

	/* Remove operation */
	for _, m := range config[namespace]["remove"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: RemoveOperation: No path specified")
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.RemovePatchOperation(fmt.Sprintf("%v", m["path"])))
	}

	/* Replace operation */
	for _, m := range config[namespace]["replace"] {
		if m["path"] == nil || m["path"] == "" { //TODO: Compiled Regex
			log.Printf("- error: ReplaceOperation: no path specified for \"value\": %v\n", m["value"])
			continue
		}
		jpOperations = append(jpOperations, admissioncontroller.ReplacePatchOperation(fmt.Sprintf("%v", m["path"]), m["value"]))
	}
	return jpOperations
}

func VerifyMutation(jpOperations []admissioncontroller.PatchOperation, r *admission.AdmissionRequest) ([]admissioncontroller.PatchOperation, error) {
	var operations []admissioncontroller.PatchOperation

	podBytes, _ := r.Object.MarshalJSON()
	podNode := Byte2Node(podBytes)

	jpOperations = PatchJPOperations(jpOperations, podNode)

	var re1 = regexp.MustCompile(`\[`)
	var re2 = regexp.MustCompile(`\]`)

	for _, jpo := range jpOperations {
		podNodes, _ := podNode.JSONPath(jpo.Path)
		xPath, err := JSONPath2XPath(jpo, podNodes, re1, re2)

		switch jpo.Op {
		case "add":
			operations = append(operations, admissioncontroller.AddPatchOperation(xPath, jpo.Value))
		case "remove":
			if len(podNodes) > 0 && err == nil {
				operations = append(operations, admissioncontroller.RemovePatchOperation(xPath))
			} else {
				log.Println("remove:", err, ":", xPath)
			}
		case "replace":
			if len(podNodes) > 0 && err == nil {
				operations = append(operations, admissioncontroller.ReplacePatchOperation(xPath, jpo.Value))
			} else {
				log.Println("replace:", err, ":", xPath, ",", jpo.Value)
			}
		default:
			log.Printf("- error: Operation: undefined")
		}
	}
	return operations, nil
}
