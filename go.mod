module K8S-Webhooks

replace github.com/go-logr/logr v0.1.0 => github.com/go-logr/logr v0.2.0

go 1.16

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/structured-merge-diff/v4 v4.1.1 // indirect
)
