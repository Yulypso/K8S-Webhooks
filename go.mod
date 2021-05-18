module K8S-Webhooks

replace github.com/go-logr/logr v0.1.0 => github.com/go-logr/logr v0.2.0

go 1.16

require (
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v0.0.0-20150909031657-73d445a93680 // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/klog/v2 v2.8.0 // indirect
)
