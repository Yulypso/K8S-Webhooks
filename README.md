# K8S-Webhooks

## Author

[Thierry Khamphousone](https://www.linkedin.com/in/tkhamphousone/)

---

<br/>

## Setup

```bash
$ git clone https://github.com/Yulypso/K8S-Webhooks.git
$ cd K8S-Webhooks
```


## Start Webhooks

```bash
# generate CA certificate/private key 
# generate certificate/private key signed by CA private key

$ /bin/bash generate-keys.sh

# generate Kubernetes objects for webhooks

$ /bin/bash generate-K8S-objects.sh
```

## Reset

```bash
# Reset the workspace, clearning all created objects. 
$ /bin/bash reset.sh
```