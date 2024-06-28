# Kuberay API Server User Guide

This guide covers the purpose, requirements, and deployment of the Kuberay API Server.

## Requirements

| Software | Version  |                                                                Link |
|:---------|:--------:|--------------------------------------------------------------------:|
| kubectl  | v1.18.3+ | [Download](https://kubernetes.io/docs/tasks/tools/install-kubectl/) |
| Go       |  v1.22.2 |                                  [Download](https://golang.org/dl/) |
| Docker   |  19.03+  |                        [Download](https://docs.docker.com/install/) |
| GNU Make |  3.81+   |                                                                     |
| curl     |  7.88+   |                                                                     |
| helm     | v3.12.1  |                      [Install](https://helm.sh/docs/intro/install/) |

### Optional Development Tools

These tools are downloaded and installed when they are needed. The directory of the download is `../bin`.
Typing `make dev-tools` will download and install all of them. The `make clean-dev-tools` command can be used to remove all the tools from the filesystem.

| Software      | Version  |                                                                    Link |
| :-------      | :------: | -----------------------------------------------------------------------:|
| kind          | v0.19.0  | [Install](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) |
| golangci-lint | v1.54.1  | [Install](https://golangci-lint.run/usage/install/)                     |
| kustomize     | v3.8.7   | [install](https://kubectl.docs.kubernetes.io/installation/kustomize/)   |
| gofumpt       | v0.3.1   | To install `go install mvdan.cc/gofumpt@v0.3.1`                         |
| goimports     | latest   | To install `go install golang.org/x/tools/cmd/goimports@latest`         |
| go-bindata    | v4.0.2   | To install `github.com/kevinburke/go-bindata/v4/...@v4.0.2`             |

## Purpose

The Kuberay API Server is designed to simplify the lifecycle management of Ray clusters for users who may not be well-versed in Kubernetes. It provides a RESTful web service to manage Ray cluster Kubernetes resources.

## Build and Deployment

The backend service can be deployed locally or within a Kubernetes cluster. The HTTP service listens on port 8888, the RPC port on 8887.

### Pre-requisites

Ensure that the admin Kubernetes configuration file is located at `~/.kube/config`. As a convenience, there are two makefile targets provided to help you manage a local kind cluster:

* `make cluster` -- creates a 3 node cluster (1 control plane 2 worker) named ray-api-server-cluster
* `make clean-cluster` -- deletes the cluster created with the `cluster` target

### Local Development

#### Build

```bash
#To build the executable
make build

#To start the executable build above
../bin/kuberay-apiserver
```

#### Test

```bash
#To run the unit tests
make test
```

#### Swagger UI updates

To update the swagger ui files deployed with the Kuberay API server, you'll need to:

* Manually run the [hack/update-swagger-ui.bash](hack/update-swagger-ui.bash) script. The script downloads the swagger ui release and copies the downloaded files
to the [../third_party/swagger-ui](../third_party/swagger-ui/) directory. It copies the [swagger-initializer.js](../third_party/swagger-ui/swagger-initializer.js)
to [swagger-initializer.js.backup](../third_party/swagger-ui/swagger-initializer.js.backup).
* Update the contents of the [swagger-initializer.js](../third_party/swagger-ui/swagger-initializer.js) to set the URLs for for the individual swagger docs. The content of the file is show below:

    ```javascript
    window.onload = function() {
    //<editor-fold desc="Changeable Configuration Block">

    // the following lines will be replaced by docker/configurator, when it runs in a docker-container
    window.ui = SwaggerUIBundle({
        spec: location.host,
        urls:  [{"url":"http://"+location.host+"/swagger/serve.swagger.json","name":"RayServe Service"},
                {"url":"http://"+location.host+"/swagger/error.swagger.json","name":"Errors API"},
                {"url":"http://"+location.host+"/swagger/job.swagger.json","name":"RayJob Service"},
                {"url":"http://"+location.host+"/swagger/config.swagger.json","name":"ComputeTemplate Service"},
                {"url":"http://"+location.host+"/swagger/cluster.swagger.json","name":"Cluster Service"}],
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
        ],
        plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
    });

    //</editor-fold>
    };
    ```

* Execute `make build-swagger` target to update the contents of the [datafile.go](pkg/swagger/datafile.go) file. This will package the content of the [swagger-ui](../third_party/swagger-ui/) directory for serving by the api server (see [func serveSwaggerUI(mux *http.ServeMux)](https://github.com/ray-project/kuberay/blob/f1067378bc99987f3eba1e5b12b4cc797465336d/apiserver/cmd/main.go#L149) in [main.go](cmd/main.go))

The swagger ui is available at the following URLs:

* [local service -- http://localhost:8888/swagger-ui/](http://localhost:8888/swagger-ui/)
* [kind cluster -- http://localhost:31888/swagger-ui/](http://localhost:31888/swagger-ui/)

Please note that for the local service the directory containing the `*.swagger.files` is specified using the `-localSwaggerPath` command line argument. The `make run` command sets the value correctly.

#### Start Local Service

This will start the api server on your development machine. The golang race detector is turned on when starting the api server this way. It will use Kubernetes configuration file located at `~/.kube/config`. The service will not start if you do not have a connection to a Kubernetes cluster.

```bash
make run
```

#### Access

Access the service at `localhost:8888` for http, and `locahost:8887` for the RPC port.

### Kubernetes Deployment

#### Build Image

```bash
#creates an image with the tag kuberay/apiserver:latest
make docker-image
```

#### Start Kubernetes Deployment

Note that you should make your KubeRay API server image available by either pushing it to an image registry, such as DockerHub or Quay, or by loading the image into the Kubernetes cluster. If you are using a Kind cluster for development, you can run `make load-image` to load the newly built API server image into the Kind cluster.  The operator image will also be needed to be loaded on your cluster. You can use `make operator-image` to build a fresh image from sources, and, if you are using a Kind cluster for development, you can run `make load-operator-image`.

```bash
#Optionally, to load the api server image into the local kind cluster created with make cluster
make load-image

#To use the  helm charts
make deploy

#To use the configuration from deploy/base
make install
```

#### Stop Kubernetes Deployment

```bash
#To use the helm charts
make undeploy

#To use the configuration
make uninstall
```

#### Local Kind Cluster Deployment

As a convenience for local development the following `make` targets are provided:

* `make cluster` -- creates a local kind cluster, using the configuration from `hack/kind-cluster-config.yaml`. It creates a port mapping allowing for the service running in the kind cluster to be accessed on  `localhost:31888` for HTTP and `localhost:31887` for RPC.
* `make clean-cluster` -- deletes the local kind cluster created with `make cluster`
* `load-image` -- loads the docker image defined by the `IMG` make variable into the kind cluster. The default value for variable is: `kuberay/apiserver:latest`. The name of the image can be changed by using `make load-image -e IMG=<your image name and tag>`
* `operator-image` -- Build the operator image to be loaded in your kind cluster. The tag for the operator image is `kuberay/operator:latest`. This step is optional.
* `load-operator-image` -- Load the operator image to the kind cluster created with `create-kind-cluster`. The tag for the operator image is `kuberay/operator:latest`, and the tag can be overridden using `make load-operator-image -E OPERATOR_IMAGE_TAG=<operator tag>`. To use the nightly operator tag, set `OPERATOR_IMAGE_TAG` to `nightly`.
* `deploy-operator` -- Deploy operator into your cluster.  The tag for the operator image is `kuberay/operator:latest`.
* `undeploy-operator` -- Undeploy operator from your cluster

When developing and testing with kind you might want to execute these targets together:

```bash
#To create a new API server image and to deploy it on a new cluster
make docker-image cluster load-image deploy

#To create a new API server image, operator image and deploy them on a new cluster
make docker-image operator-image cluster load-image load-operator-image deploy deploy-operator
```

#### Access API Server in the Cluster

Access the service at `localhost:31888` for http and `locahost:31887` for the RPC port.
