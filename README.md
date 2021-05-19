# K8S-Webhooks

## Author

[![Linkedin: Thierry Khamphousone](https://img.shields.io/badge/-Thierry_Khamphousone-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/tkhamphousone/)](https://www.linkedin.com/in/tkhamphousone)

---

<br/>

## First implementation

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