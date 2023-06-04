A registry credential extension to AWS ECR.

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=DRONE_ECRREGION=us-east-1 \
  --restart=always \
  --name=drone-ecr-registry-extension mudasirmirza/drone-ecr-registry-extension
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_REGISTRY_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_REGISTRY_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
```

---
**Note:** Make sure the pod / node has proper IAM permissions to pull ECR image(s)
