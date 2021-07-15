module K8S-Webhooks

replace github.com/go-logr/logr v0.1.0 => github.com/go-logr/logr v0.2.0

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	github.com/spyzhov/ajson v0.4.2
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/klog/v2 v2.9.0
)
