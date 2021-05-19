# K8S-Webhooks

## Author

[![Linkedin: Thierry Khamphousone](https://img.shields.io/badge/-Thierry_Khamphousone-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/tkhamphousone/)](https://www.linkedin.com/in/tkhamphousone)

---

<br/>

## Second implementation (1 MutatingWebhook + 1 ValidatingWebhook)

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

> Reset the workspace, clean K8S objects. 
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

> Both
```bash
$ ./reset certificates docker
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

## Questions & Answers

What happens if you have both a MutatingWebhook and a ValidatingWebhook, which is applied first during a deployment?
- [Admission control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
> "Admission webhooks are HTTP callbacks that receive admission requests and do something with them. You can define two types of admission webhooks, validating admission webhook and mutating admission webhook. Mutating admission webhooks are invoked first, and can modify objects sent to the API server to enforce custom defaults. After all object modifications are complete, and after the incoming object is validated by the API server, validating admission webhooks are invoked and can reject requests to enforce custom policies."
