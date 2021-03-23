# image-cloner
A k8s controller which watches the applications and caches the images by re-uploading to our own registry and re-configures the applications to use these copies.

### Prerequisits
1. A cluster is required, it can be a managed cluster like GKE etc. or a local `kind` or `minikube` cluster.  
  The makefile contains the way to create the `kind` cluster with the following command.
  If you don't have `kind` installed:
  ```shell
  $ make install-kind
  ```
  ```shell
  $ make cluster
  ```
2. The operator need a secret which holds the registry creds to the private docker registry where all the image will be cached.  
This has to be created upfront before deploying the operator to your cluster. Below are the steps
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: private-registry-creds
  namespace: system
type: Opaque
stringData:
  DOCKER_REGISTRY_URI: <registry_name/repo_name>  # example "index.docker.io/ishankhare07"
  DOCKER_REGISTRY_USERNAME: <username>
  DOCKER_REGISTRY_PASSWORD: <password>
```
  And then create this secret in the cluster:
```shell
$ kubectl apply -f secret.yaml
```

### Use your registry to host operator
The controller image is currently available as a public image on dockerhub at `ishankhare07/image-cloner`, in case you want to host your own image  
build and push to your registry in the following way:
```shell
$ docker login <your_registry>
...
...
...

$ IMG="<your_registry/username>" make build
$ IMG="<your_registry/username>" make push
$ IMG="<your_registry/username>" make deploy
```

### Use existing registry image
If you are okay with using the current image on the public dockerhub, just go ahead and deploy:
```shell
$ make deploy
```

This will apply the following to your cluster:
1. Create namespace `system`
2. Create `clusterrole` for the controllers
3. Create `clusterrolebindings`
4. Create `deployment` which will actually run the controller.

### Running locally
When running locally, create a file called `config/config.env` which contains the private registry cred as follows:
```
DOCKER_REGISTRY_URI="<registry_name/repo_name>"
DOCKER_REGISTRY_USERNAME="<username>"
DOCKER_REGISTRY_PASSWORD="<password>"
```

This allows easily testing the whole flow, next just run:
```shell
$ make cluster
$ make run
```
