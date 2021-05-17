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

```bash
# generate CA certificate/private key 
# generate certificate/private key signed by CA private key

$ ./generate-keys.sh

# generate Kubernetes objects for webhooks

$ ./generate-K8S-objects.sh
```

## Reset

```bash
# Reset the workspace, clean K8S objects. 
$ ./reset.sh

# Clean certificates
$ ./reset.sh certificates

# Build and Push to Docker Hub (for development uses)
$ ./reset.sh docker
```

---

<br/>

## Certificates

```bash
# verify CSR content
$ openssl req -text -noout -verify -in Certificates/webhookservertls.csr

# verify CERT content
$ openssl x509 -text -noout -in Certificates/webhookservertls.cert
```