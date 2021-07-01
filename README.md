# K8S-Webhooks

## Author

[![Linkedin: Thierry Khamphousone](https://img.shields.io/badge/-Thierry_Khamphousone-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/tkhamphousone/)](https://www.linkedin.com/in/tkhamphousone)

---

<br/>

## Setup

```bash
$ git clone https://github.com/Yulypso/K8S-Webhooks.git
$ cd K8S-Webhooks
```

---

<br/>

## Start Webhooks 
(Namespace: webhookserver-ns)

> Generate CA certificate/CA private key 
> Generate server certificate/server private key signed by CA private key

```bash
$ ./generate-keys.sh
```

> Generate Kubernetes objects for webhooks within webhookserver-ns namespace
```bash
$ ./generate-K8S-objects.sh
```

<br/>

## Start test pod 
(Namespace: admissionwebhook-ns)

> Generate pod within admissionwebhook-ns namespace which calls the created webhook
```bash
$ kubectl apply -f ./TestDeployments/pod-1.yml
```

---

<br/>

## Reset

> Reset K8S webhook server. 
```bash
$ ./reset.sh
```

> Clean certificates
```bash
$ ./reset.sh certificates
```

> Build and Push to Docker Hub (for development uses)
```bash
$ ./reset.sh docker
```

> Reset all K8S Cluster
```bash
$ ./reset cluster
```

---

<br/>

## Certificates

> Verify CSR content
```bash
$ openssl req -text -noout -verify -in Certificates/webhookservertls.csr
```

> Verify CERT content
```bash
$ openssl x509 -text -noout -in Certificates/webhookservertls.cert
```

---

<br/>

## Questions and Answers

What happens if you have both a MutatingWebhook and a ValidatingWebhook, which is applied first during a deployment? 
- [Admission control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
> "Admission webhooks are HTTP callbacks that receive admission requests and do something with them. You can define two types of admission webhooks, validating admission webhook and mutating admission webhook. Mutating admission webhooks are invoked first, and can modify objects sent to the API server to enforce custom defaults. After all object modifications are complete, and after the incoming object is validated by the API server, validating admission webhooks are invoked and can reject requests to enforce custom policies."

What happens when you have two MutatingWebhooks, which of the two is run first? How is it going ?
> The metadata.name field is used to define the order of application of MutatingWebhooks and ValidatingWebhooks. Within the prototype, a-mutatingwebhook is applied before b-mutatingwebhook 
> 
> Therefore, we can have in this order of execution, a MutatingWebhook [a] which can do [0-n] operations and have a second MutatingWebhook [b] which can also do [0-n] operations.

Should we specify runAsUser, runAsGroup, fsGroup field in addition to the field runAsNonRoot: true ? 
> It is not recommanded because some images such as the jenkins/jenkins official server image runs as group:user == jenkins:jenkins. That is why we should not specify any of those fields in order to make sure the server will work properly in this case.

---

<br/>

## JSONPath

Filter must be written with: []
- OK
> $.spec.containers[?(@[name] == 'node-app0')]
- KO
> $.spec.containers[?(@.name == 'node-app0')]


Keys containing '.' must be surrounded by []
- OK
> $.metadata.annotations.[artemis.site]
- KO
> $.metadata.annotations.artemis.site

---

<br/>

## Default DSL Config (common to all deployed pods)
```json
{
    "add":[
        {
            "path":"$.spec.securityContext.runAsNonRoot",
            "value":true
        }
    ],
    "remove":[
        {
            "path":"$.spec.volumes[*].hostPath"
        },
        {
            "path":"$.spec.securityContext.runAsNonRoot",
            "value":false
        }
    ],
    "mandatorydata":[
        {
            "path":"$.spec.securityContext.runAsNonRoot",
            "value":true
        }
    ],
    "forbiddendata":[
        {
            "path":"$.spec.volumes[*].hostPath"
        },
        {
            "path": "$.spec.securityContext.runAsUser",
            "value": 0
        },
        {
            "path":"$.spec.securityContext.runAsNonRoot",
            "value":false
        }
    ]
}
```

---

<br/>

## HTTP Requests (import from Postman directory)

PATCH: DSL Config (example)

```json
> Method: PUT  
> Endpoint: localhost:31000/namespace/<NAMESPACE>  
> Body: 
{
    "add": [
        {
            "path": "$.spec.containers[*].securityContext.allowPrivilegeEscalation",
            "value": false
        },

        {
            "path": "$.metadata.annotations.[co.elastic.logs/multiline.pattern]",
            "value": "^\d{4}-\d{2}-\d{2}"
        },
        {
            "path": "$.spec.containers",
            "value": {
                "env": [
                    {
                        "name": "ELASTOMCAT_HOST",
                        "value": "http://srvelasprod.technique.artemis:9200"
                    },
                    {
                        "name": "ELASTOMCAT_USERNAME",
                        "value": "UpXo3on-wowrT8g"
                    },
                    {
                        "name": "ELASTOMACT_USERPWD",
                        "value": "Artemis2019****"
                    }
                ],
                "image": "tomcat:8.0-alpine",
                "imagePullPolicy": "Always",
                "name": "tomcatelas--1185509365",
                "ports": [
                    {
                        "containerPort": 8080,
                        "name": "http-ext",
                        "protocol": "TCP"
                    }
                ]
            }
        }
    ],
    "replace": [
        {
            "path": "$.spec.containers[?(@[name] == 'node-app0')]",
            "value": {
                "env": [
                    {
                        "name": "elas_host",
                        "value": "http://srvelasprod.technique.artemis:9200"
                    },
                    {
                        "name": "username",
                        "value": "anonymous"
                    },
                    {
                        "name": "password",
                        "value": "na"
                    }
                ],
                "image": "tomcat:8.5-alpine",
                "imagePullPolicy": "Always",
                "name": "containerexistant",
                "ports": [
                    {
                        "containerPort": 8080,
                        "name": "http-ext",
                        "protocol": "TCP"
                    }
                ]
            }
        }
    ],
    "remove": [
        {
            "path": "$.spec.securityContext.runAsUser"
        }
    ],
    "mandatorydata": [
        {
            "path": "$.metadata.labels.a4c_nodeid"
        },
        {
            "path": "$.metadata.annotations.[artemis.site]",
            "value": "prod"
        }
    ],
    "forbiddendata": [
        {
            "path": "$.spec.securityContext.runAsUser",
            "value": 0
        }
    ]
}
```

RESET: DSL Config to the initial default.json (empty config => {})
``` json
> Method: DELETE
> Endpoint: localhost:31000/reset  
> Body:   
```


CLEAR: DSL Config by namespace
```json
> Method: DELETE 
> Endpoint: localhost:31000/namespace/<NAMESPACE>  
> Body:   
```