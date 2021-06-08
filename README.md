# K8S-Webhooks

## Author

[![Linkedin: Thierry Khamphousone](https://img.shields.io/badge/-Thierry_Khamphousone-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/tkhamphousone/)](https://www.linkedin.com/in/tkhamphousone)

---

<br/>

## Forth implementation (JSONPath first implementation)


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

