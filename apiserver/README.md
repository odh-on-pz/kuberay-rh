# KubeRay APIServer

The KubeRay APIServer provides gRPC and HTTP APIs to manage KubeRay resources.

## Note

```text
The KubeRay APIServer is an optional component. It provides a layer of simplified configuration for KubeRay resources. The KubeRay API server is used internally by some organizations to back user interfaces for KubeRay resource management.

The KubeRay APIServer is community-managed and is not officially endorsed by the Ray maintainers. At this time, the only officially supported methods for
managing KubeRay resources are:

- Direct management of KubeRay custom resources via kubectl, kustomize, and Kubernetes language clients.
- Helm charts.

KubeRay APIServer maintainer contacts (GitHub handles):
@Jeffwan @scarlet25151
```

## Installation

### Helm

Make sure the version of Helm is v3+. Currently, [existing CI tests](https://github.com/ray-project/kuberay/blob/master/.github/workflows/helm-lint.yaml) are based on Helm v3.4.1 and v3.9.4.

```sh
helm version
```

### Install KubeRay Operator

* Install a stable version via Helm repository (only supports KubeRay v0.4.0+)

  ```sh
    # Install the KubeRay helm repo
    helm repo add kuberay https://ray-project.github.io/kuberay-helm/

    # Install KubeRay Operator v1.0.0-rc.1.
    helm install kuberay-operator kuberay/kuberay-operator --version v1.0.0-rc.1

    # Check the KubeRay Operator Pod in `default` namespace
    kubectl get pods
    # NAME                                             READY   STATUS    RESTARTS   AGE
    # kuberay-operator-7456c6b69b-t6pt7                1/1     Running   0          172m
  ```

### Install KubeRay APIServer

```text
Please note that examples show here will only work with the nightly builds of the api-server. `v1.0.0-rc.1` does not yet contain critical fixes
to the api server that would allow Kuberay Serve endpoints to work properly
```

* Install a stable version via Helm repository (only supports KubeRay v0.4.0+)
  
  ```sh
  # Install the KubeRay helm repo
  helm repo add kuberay https://ray-project.github.io/kuberay-helm/

  # Install KubeRay APIServer.
  helm install kuberay-apiserver kuberay/kuberay-apiserver 

  # Check the KubeRay APIServer Pod in `default` namespace
  kubectl get pods
  # NAME                                             READY   STATUS    RESTARTS   AGE
  # kuberay-apiserver-857869f665-b94px               1/1     Running   0          86m
  # kuberay-operator-7456c6b69b-t6pt7                1/1     Running   0          172m
  ```

* Install the nightly version
  
  ```sh
  # Step1: Clone KubeRay repository

  # Step2: Move to `helm-chart/kuberay-apiserver`
  cd helm-chart/kuberay-apiserver

  # Step3: Install KubeRay APIServer
  helm install kuberay-apiserver .
  ```

* Install the current (working branch) version
  
  ```sh
  # Step1: Clone KubeRay repository

  # Step2: Change directory to the api server
  cd apiserver

  # Step3: Build docker image, create a local kind cluster and deploy api server (using helm)
  make docker-image cluster load-image deploy
  ```  

### List the chart

To list the deployments:

```sh
helm ls
# NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                           APP VERSION
# kuberay-apiserver       default         1               2023-09-25 10:42:34.267328 +0300 EEST   deployed        kuberay-apiserver-1.0.0-rc.1               
# kuberay-operator        default         1               2023-09-25 10:41:48.355831 +0300 EEST   deployed        kuberay-operator-1.0.0-rc.1                
```

### Uninstall the Chart

```sh
# Uninstall the `kuberay-apiserver` release
helm uninstall kuberay-apiserver

# The KubeRay APIServer Pod should be removed.
kubectl get pods
# No resources found in default namespace.
```

## Usage

After the deployment we may use the `{{baseUrl}}` to access the service. See [swagger support section](https://ray-project.github.io/kuberay/components/apiserver/#swagger-support) to get the complete definitions of APIs.

* (default) for nodeport access, use port `31888` for connection

* for ingress access, you will need to create your own ingress

The requests parameters detail can be seen in [KubeRay swagger](https://github.com/ray-project/kuberay/tree/master/proto/swagger), this document only presents basic examples.

### Setup a smoke test

The following steps allow you to validate that the KubeRay API Server components and KubeRay Operator integrate in your environment.

1. (Optional) You may use your local kind cluster or minikube

    ```bash
    cat <<EOF | kind create cluster --name ray-test --config -
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
    - role: control-plane
      image: kindest/node:v1.23.17@sha256:59c989ff8a517a93127d4a536e7014d28e235fb3529d9fba91b3951d461edfdb
      kubeadmConfigPatches:
        - |
          kind: InitConfiguration
          nodeRegistration:
            kubeletExtraArgs:
              node-labels: "ingress-ready=true"
      extraPortMappings:
      - containerPort: 30265
        hostPort: 8265
        listenAddress: "0.0.0.0"
        protocol: tcp
      - containerPort: 30001
        hostPort: 10001
        listenAddress: "0.0.0.0"
        protocol: tcp
      - containerPort: 8000
        hostPort: 8000
        listenAddress: "0.0.0.0"
      - containerPort: 31888
        hostPort: 31888
        listenAddress: "0.0.0.0"
      - containerPort: 31887 
        hostPort: 31887 
        listenAddress: "0.0.0.0"
    - role: worker
      image: kindest/node:v1.23.17@sha256:59c989ff8a517a93127d4a536e7014d28e235fb3529d9fba91b3951d461edfdb
    - role: worker
      image: kindest/node:v1.23.17@sha256:59c989ff8a517a93127d4a536e7014d28e235fb3529d9fba91b3951d461edfdb
    EOF
    ```

2. Deploy the KubeRay APIServer within the same cluster of KubeRay operator

    ```sh
    helm repo add kuberay https://ray-project.github.io/kuberay-helm/
    helm -n ray-system install kuberay-apiserver kuberay/kuberay-apiserver -n ray-system --create-namespace
    ```

3. The APIServer expose service using `NodePort` by default. You can test access by your host and port, the default port is set to `31888`. The examples below assume a kind (localhost) deployment. If Kuberay API server is deployed on another type of cluster, you'll need to adjust the hostname to match your environment.

    ```sh
    curl localhost:31888
    ...
    {
      "code": 5,
      "message": "Not Found"
    }
    ```

4. You can create `RayCluster`, `RayJobs` or `RayService` by dialing the endpoints. The following is a simple example for creating the `RayService` object, follow [swagger support](https://ray-project.github.io/kuberay/components/apiserver/#swagger-support) to get the complete definitions of APIs.

    Please note that:

    * The following examples use the `ray-system` namespace. If not already created by using the helm install steps above, you can create it prior to executing the curl examples by running `kubectl create namespace ray-system`
    * The examples assume that the cluster has at least 2 CPUs available and 4 GB of free memory. You can either increase the CPUs available to your cluster (docker settings) or reduce the CPU request in the `compute_templates` request.
    * If you are running the service and the kuberay operator on Apple Silicon Machine, you might want to use the `rayproject/ray:2.7.0-aarch64`image.

    ```sh
    # Create a template
    curl --silent -X POST 'localhost:31888/apis/v1alpha2/namespaces/ray-system/compute_templates' \
    --header 'Content-Type: application/json' \
    --data '{
      "name": "default-template",
      "namespace": "ray-system",
      "cpu": 2,
      "memory": 4
    }'

    # Create the "Fruit Stand" Ray Serve example (V1 Config spec)
    curl --silent -X POST 'localhost:31888/apis/v1alpha2/namespaces/ray-system/services' \
    --header 'Content-Type: application/json' \
    --data '{
      "name": "test-v1",
      "namespace": "ray-system",
      "user": "user",
      "serviceUnhealthySecondThreshold": 900,
      "deploymentUnhealthySecondThreshold": 300,
      "serveDeploymentGraphSpec": {
        "importPath": "fruit.deployment_graph",
        "runtimeEnv": "working_dir: \"https://github.com/ray-project/test_dag/archive/c620251044717ace0a4c19d766d43c5099af8a77.zip\"\n",
        "serveConfigs": [
          {
            "deploymentName": "OrangeStand",
            "replicas": 1,
            "userConfig": "price: 2",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "PearStand",
            "replicas": 1,
            "userConfig": "price: 1",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "FruitMarket",
            "replicas": 1,
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "DAGDriver",
            "replicas": 1,
            "routePrefix": "/",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          }
        ]
      },
      "clusterSpec": {
        "headGroupSpec": {
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "serviceType": "NodePort",
          "rayStartParams": {
            "dashboard-host": "0.0.0.0",
            "metrics-export-port": "8080"
          },
          "volumes": []
        },
        "workerGroupSpec": [
          {
            "groupName": "small-wg",
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0-py310",
            "replicas": 1,
            "minReplicas": 0,
            "maxReplicas": 5,
            "rayStartParams": {
              "node-ip-address": "$MY_POD_IP"
            }
          }
        ]
      }
    }'

    # Get the pods running in the namespace
    kubectl get pods -n ray-system
    # NAME                                             READY   STATUS    RESTARTS   AGE
    # kuberay-apiserver-5dd7c9b4c8-rnrcq               1/1     Running   0          17m
    # test-v1-raycluster-42vss-head-m4qj2              1/1     Running   0          91s
    # test-v1-raycluster-42vss-worker-small-wg-svr79   1/1     Running   0          90s

    # Delete the RayService in the namespace
    curl --silent -X 'DELETE' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services/test-v1' \
    -H 'accept: application/json'
    ```

## Swagger Support

Kuberay API server has support for Swagger UI. The swagger page can be reached at:

* [localhost:31888/swagger-ui](localhost:31888/swagger-ui) for local kind deployments
* [localhost:8888/swagger-ui](localhost:8888/swagger-ui) for instances started with `make run` (development machine builds)
* `<host name>:31888/swagger-ui` for nodeport deployments

## Full definition endpoints

### Compute Template

For the purpose to simplify the setting of resources, the Kuberay API server abstracts the resource of the pods template resource to the `compute template`. You can define the resources in the `compute template` and then choose the appropriate template for your `head` and `workergroup` when you are creating the objects of `RayCluster`, `RayJobs` or `RayService`.

The full definition of the compute template resource can be found in [config.proto](../proto/config.proto) or the Kuberay API server swagger doc.

#### Create compute templates in a given namespace

```text
POST {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/compute_templates
```

Examples:

* Request

  ```sh
  curl --silent -X 'POST' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/compute_templates' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
        "name": "default-template",
        "namespace": "ray-system",
        "cpu": 2,
        "memory": 4
      }'
  ```

* Response

  ```json
  {
    "name": "default-template",
    "namespace": "ray-system",
    "cpu": 2,
    "memory": 4
  }
  ```

#### List all compute templates in a given namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/compute_templates
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/compute_templates' \
    -H 'accept: application/json'
  ```

* Response

  ```json
  {
    "computeTemplates": [
      {
        "name": "default-template",
        "namespace": "ray-system",
        "cpu": 2,
        "memory": 4
      }
    ]
  }
  ```

#### List all compute templates in all namespaces

```text
GET {{baseUrl}}/apis/v1alpha2/compute_templates
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/compute_templates' \
  -H 'accept: application/json'
  ```

* Response

    ```json
    {
      "computeTemplates": [
        {
          "name": "default-template",
          "namespace": "ray-system",
          "cpu": 2,
          "memory": 4
        }
      ]
    }
    ```

#### Get compute template by name

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/compute_templates/<compute_template_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/compute_templates/default-template' \
  -H 'accept: application/json'
  ```

* Response:

  ```json
  {
    "name": "default-template",
    "namespace": "ray-system",
    "cpu": 2,
    "memory": 4
  }
  ```

#### Delete compute template by name

```text
DELETE {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/compute_templates/<compute_template_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'DELETE' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/compute_templates/default-template' \
  -H 'accept: application/json'
  ```

* Response

  ```json
  {}
  ```

### Clusters

#### Create cluster in a given namespace

```text
POST {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/clusters
```

Examples:

* Request

  ```sh
  curl --silent -X 'POST' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/clusters' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
    "name": "test-cluster",
    "namespace": "ray-system",
    "user": "3cpo",
    "version": "2.7.0",
    "environment": "DEV",
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        },
        "volumes": []
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0",
          "replicas": 1,
          "minReplicas": 1,
          "maxReplicas": 1,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    }
  }'
  ```

* Response

  ```json
  {
    "name": "test-cluster",
    "namespace": "ray-system",
    "user": "3cpo",
    "version": "2.7.0",
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0",
          "replicas": 1,
          "minReplicas": 1,
          "maxReplicas": 1,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    },
    "annotations": {
      "ray.io/creation-timestamp": "2023-09-25 10:48:35.766443417 +0000 UTC"
    },
    "createdAt": "2023-09-25T10:48:35Z",
    "events": [
      {
        "id": "test-cluster.178817bd10374138",
        "name": "test-cluster-test-cluster.178817bd10374138",
        "createdAt": "2023-09-25T08:42:40Z",
        "firstTimestamp": "2023-09-25T08:42:40Z",
        "lastTimestamp": "2023-09-25T08:42:40Z",
        "reason": "Created",
        "message": "Created service test-cluster-head-svc",
        "type": "Normal",
        "count": 1
      },
      {
        "id": "test-cluster.178817bd251e9c7c",
        "name": "test-cluster-test-cluster.178817bd251e9c7c",
        "createdAt": "2023-09-25T08:42:40Z",
        "firstTimestamp": "2023-09-25T08:42:40Z",
        "lastTimestamp": "2023-09-25T08:42:40Z",
        "reason": "Created",
        "message": "Created head pod test-cluster-head-rsbmm",
        "type": "Normal",
        "count": 1
      },
      {
        "id": "test-cluster.178817bd2b74493f",
        "name": "test-cluster-test-cluster.178817bd2b74493f",
        "createdAt": "2023-09-25T08:42:40Z",
        "firstTimestamp": "2023-09-25T08:42:40Z",
        "lastTimestamp": "2023-09-25T08:42:41Z",
        "reason": "Created",
        "message": "Created worker pod ",
        "type": "Normal",
        "count": 2
      }
    ]
  }
  ```

#### List all clusters in a given namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/clusters
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/clusters' \
    -H 'accept: application/json'
  ```

* Response

  ```json
  {
    "clusters": [
      {
        "name": "test-cluster",
        "namespace": "ray-system",
        "user": "3cpo",
        "version": "2.7.0",
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0",
              "replicas": 1,
              "minReplicas": 1,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "annotations": {
          "ray.io/creation-timestamp": "2023-09-25 10:48:35.766443417 +0000 UTC"
        },
        "createdAt": "2023-09-25T10:48:35Z",
        "clusterState": "ready",
        "events": [
          {
            "id": "test-cluster.178817bd10374138",
            "name": "test-cluster-test-cluster.178817bd10374138",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd251e9c7c",
            "name": "test-cluster-test-cluster.178817bd251e9c7c",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-rsbmm",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd2b74493f",
            "name": "test-cluster-test-cluster.178817bd2b74493f",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:41Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 2
          },
          {
            "id": "test-cluster.17881e9c2b82c449",
            "name": "test-cluster-test-cluster.17881e9c2b82c449",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c2e9cd4b8",
            "name": "test-cluster-test-cluster.17881e9c2e9cd4b8",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-nglmx",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c34460442",
            "name": "test-cluster-test-cluster.17881e9c34460442",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 1
          }
        ],
        "serviceEndpoint": {
          "dashboard": "31476",
          "head": "31850",
          "metrics": "32189",
          "redis": "30736"
        }
      }
    ]
  }
  ```

#### List all clusters in all namespaces

```text
GET {{baseUrl}}/apis/v1alpha2/clusters
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
    'http://localhost:31888/apis/v1alpha2/clusters' \
    -H 'accept: application/json'
  ```

* Response

  ```json
  {
    "clusters": [
      {
        "name": "test-cluster",
        "namespace": "ray-system",
        "user": "3cpo",
        "version": "2.7.0",
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0",
              "replicas": 1,
              "minReplicas": 1,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "annotations": {
          "ray.io/creation-timestamp": "2023-09-25 10:48:35.766443417 +0000 UTC"
        },
        "createdAt": "2023-09-25T10:48:35Z",
        "clusterState": "ready",
        "events": [
          {
            "id": "test-cluster.178817bd10374138",
            "name": "test-cluster-test-cluster.178817bd10374138",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd251e9c7c",
            "name": "test-cluster-test-cluster.178817bd251e9c7c",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-rsbmm",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd2b74493f",
            "name": "test-cluster-test-cluster.178817bd2b74493f",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:41Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 2
          },
          {
            "id": "test-cluster.17881e9c2b82c449",
            "name": "test-cluster-test-cluster.17881e9c2b82c449",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c2e9cd4b8",
            "name": "test-cluster-test-cluster.17881e9c2e9cd4b8",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-nglmx",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c34460442",
            "name": "test-cluster-test-cluster.17881e9c34460442",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 1
          }
        ],
        "serviceEndpoint": {
          "dashboard": "31476",
          "head": "31850",
          "metrics": "32189",
          "redis": "30736"
        }
      }
    ]
  }  
  ```

#### Get cluster by its name and namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/clusters/<cluster_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/clusters' \
    -H 'accept: application/json'
  ```

* Response

  ```json
  {
    "clusters": [
      {
        "name": "test-cluster",
        "namespace": "ray-system",
        "user": "3cpo",
        "version": "2.7.0",
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0",
              "replicas": 1,
              "minReplicas": 1,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "annotations": {
          "ray.io/creation-timestamp": "2023-09-25 10:48:35.766443417 +0000 UTC"
        },
        "createdAt": "2023-09-25T10:48:35Z",
        "clusterState": "ready",
        "events": [
          {
            "id": "test-cluster.178817bd10374138",
            "name": "test-cluster-test-cluster.178817bd10374138",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd251e9c7c",
            "name": "test-cluster-test-cluster.178817bd251e9c7c",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:40Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-rsbmm",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.178817bd2b74493f",
            "name": "test-cluster-test-cluster.178817bd2b74493f",
            "createdAt": "2023-09-25T08:42:40Z",
            "firstTimestamp": "2023-09-25T08:42:40Z",
            "lastTimestamp": "2023-09-25T08:42:41Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 2
          },
          {
            "id": "test-cluster.17881e9c2b82c449",
            "name": "test-cluster-test-cluster.17881e9c2b82c449",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created service test-cluster-head-svc",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c2e9cd4b8",
            "name": "test-cluster-test-cluster.17881e9c2e9cd4b8",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created head pod test-cluster-head-nglmx",
            "type": "Normal",
            "count": 1
          },
          {
            "id": "test-cluster.17881e9c34460442",
            "name": "test-cluster-test-cluster.17881e9c34460442",
            "createdAt": "2023-09-25T10:48:35Z",
            "firstTimestamp": "2023-09-25T10:48:35Z",
            "lastTimestamp": "2023-09-25T10:48:35Z",
            "reason": "Created",
            "message": "Created worker pod ",
            "type": "Normal",
            "count": 1
          }
        ],
        "serviceEndpoint": {
          "dashboard": "31476",
          "head": "31850",
          "metrics": "32189",
          "redis": "30736"
        }
      }
    ]
  }
  ```

#### Delete cluster by its name and namespace

```text
DELETE {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/clusters/<cluster_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'DELETE' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/clusters/test-cluster' \
  -H 'accept: application/json'
  ```

* Response:

  ```json
  {}
  ```

### RayJob

#### Create ray job in a given namespace

```text
POST {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/jobs
```

Examples:

* Request
  
  ```sh
  curl --silent -X 'POST' \
      'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/jobs' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "name": "rayjob-test",
      "namespace": "ray-system",
      "user": "3cp0",
      "entrypoint": "python -V",
      "clusterSpec": {
        "headGroupSpec": {
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0",
          "serviceType": "NodePort",
          "rayStartParams": {
            "dashboard-host": "0.0.0.0"
          }
        },
        "workerGroupSpec": [
          {
            "groupName": "small-wg",
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "replicas": 1,
            "minReplicas": 0,
            "maxReplicas": 1,
            "rayStartParams": {
              "metrics-export-port": "8080"
            }
          }
        ]
      }
    }
    '
  ```

* Response

  ```json
  {
    "name": "rayjob-test",
    "namespace": "ray-system",
    "user": "3cp0",
    "entrypoint": "python -V",
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0",
          "replicas": 1,
          "minReplicas": 1,
          "maxReplicas": 1,
          "rayStartParams": {
            "metrics-export-port": "8080"
          }
        }
      ]
    },
    "createdAt": "2023-09-25T11:36:02Z"
  }
  ```

The above example creates a new Ray cluster, executes a job on it and optionally deletes a cluster. As an alternative, the same command allows creating a new job on the existing cluster by referencing it in the payload.

Examples:

Start from creating Ray cluster (We assume here that the [template](test/cluster/template/simple) and [configmap](test/job/code.yaml) are already created).

* Request
  
```sh
curl -X POST 'localhost:31888/apis/v1alpha2/namespaces/default/clusters' \
--header 'Content-Type: application/json' \
--data '{
  "name": "job-test",
  "namespace": "default",
  "user": "boris",
  "version": "2.7.0",
  "environment": "DEV",
  "clusterSpec": {
    "headGroupSpec": {
      "computeTemplate": "default-template",
      "image": "rayproject/ray:2.7.0-py310",
      "serviceType": "NodePort",
      "rayStartParams": {
         "dashboard-host": "0.0.0.0",
         "metrics-export-port": "8080"
       },
       "volumes": [
         {
           "name": "code-sample",
           "mountPath": "/home/ray/samples",
           "volumeType": "CONFIGMAP",
           "source": "ray-job-code-sample",
           "items": {"sample_code.py" : "sample_code.py"}
         }
       ]
    },
    "workerGroupSpec": [
      {
        "groupName": "small-wg",
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0-py310",
        "replicas": 1,
        "minReplicas": 0,
        "maxReplicas": 5,
        "rayStartParams": {
           "node-ip-address": "$MY_POD_IP"
         },
        "volumes": [
          {
            "name": "code-sample",
            "mountPath": "/home/ray/samples",
            "volumeType": "CONFIGMAP",
            "source": "ray-job-code-sample",
            "items": {"sample_code.py" : "sample_code.py"}
          }
        ]
      }
    ]
  }
}'
```

* Response

```json
{
   "name":"job-test",
   "namespace":"default",
   "user":"boris",
   "version":"2.7.0",
   "clusterSpec":{
      "headGroupSpec":{
         "computeTemplate":"default-template",
         "image":"rayproject/ray:2.7.0-py310",
         "serviceType":"NodePort",
         "rayStartParams":{
            "dashboard-host":"0.0.0.0",
            "metrics-export-port":"8080"
         },
         "volumes":[
            {
               "mountPath":"/home/ray/samples",
               "volumeType":3,
               "name":"code-sample",
               "source":"ray-job-code-sample",
               "items":{
                  "sample_code.py":"sample_code.py"
               }               
            }
         ],
         "environment":{
            
         }
      },
      "workerGroupSpec":[
         {
            "groupName":"small-wg",
            "computeTemplate":"default-template",
            "image":"rayproject/ray:2.7.0-py310",
            "replicas":1,
            "minReplicas":5,
            "maxReplicas":1,
            "rayStartParams":{
               "node-ip-address":"$MY_POD_IP"
            },
            "volumes":[
               {
                  "mountPath":"/home/ray/samples",
                  "volumeType":3,
                  "name":"code-sample",
                  "source":"ray-job-code-sample",
                  "items":{
                     "sample_code.py":"sample_code.py"
                  }
               }
            ],
            "environment":{
               
            }
         }
      ]
   },
   "annotations":{
      "ray.io/creation-timestamp":"2023-10-18 08:47:48.058576 +0000 UTC"
   },
   "createdAt":"2023-10-18T08:47:48Z"
}
```

Once the cluster is created, we can create a job to run on it.

* Request
  
```sh
curl -X POST 'localhost:31888/apis/v1alpha2/namespaces/default/jobs' \
--header 'Content-Type: application/json' \
--data '{
  "name": "job-test",
  "namespace": "default",
  "user": "boris",
  "entrypoint": "python /home/ray/samples/sample_code.py",
   "runtimeEnv": "pip:\n  - requests==2.26.0\n  - pendulum==2.1.2\nenv_vars:\n  counter_name: test_counter\n",
  "clusterSelector": {
    "ray.io/cluster": "job-test"
  }
}'  
```

This request fails with the message:

```json
{
   "code":3,
   "message":"Failed to create a Ray Job: Create Job failed.: InvalidInputError: Failed to create a Ray Job: external Ray cluster requires Job submitter definition"
}
```

The reason for this failure is that in the case of the existing cluster the oprator can not figure out Python/Ray version of the cluster's nodes. To make it work, we need to add `jobSubmitter` component to the request to ensure that job submitter's Python/Ray version are the same as the ones of the cluster.

* Request
  
```sh
curl -X POST 'localhost:31888/apis/v1alpha2/namespaces/default/jobs' \
--header 'Content-Type: application/json' \
--data '{
  "name": "job-test",
  "namespace": "default",
  "user": "boris",
  "entrypoint": "python /home/ray/samples/sample_code.py",
   "runtimeEnv": "pip:\n  - requests==2.26.0\n  - pendulum==2.1.2\nenv_vars:\n  counter_name: test_counter\n",
  "clusterSelector": {
    "ray.io/cluster": "job-test"
  },
  "jobSubmitter": {
    "image": "rayproject/ray:2.7.0-py310"
  }
}'  
```

* Response

```json
{
   "name":"job-test",
   "namespace":"default",
   "user":"boris",
   "entrypoint":"python /home/ray/samples/sample_code.py",
   "runtimeEnv":"pip:\n  - requests==2.26.0\n  - pendulum==2.1.2\nenv_vars:\n  counter_name: test_counter\n",
   "clusterSelector":{
      "ray.io/cluster":"job-test"
   },
   "jobSubmitter":{
      "image":"rayproject/ray:2.7.0-py310"
   },
   "createdAt":"2023-10-18T10:19:49Z"
}
```

You should also see job submitter job completed, something like:

```text
job-test-2hhmf                   0/1     Completed   0          15s
```

To see job execution results run:

```sh
kubectl logs job-test-2hhmf 
```

And you should get something similar to:

```text
2023-10-18 03:19:51,524	INFO cli.py:36 -- Job submission server address: http://job-test-head-svc.default.svc.cluster.local:8265
2023-10-18 03:19:52,197	SUCC cli.py:60 -- -------------------------------------------
2023-10-18 03:19:52,197	SUCC cli.py:61 -- Job 'job-test-bbfqs' submitted successfully
2023-10-18 03:19:52,197	SUCC cli.py:62 -- -------------------------------------------
2023-10-18 03:19:52,197	INFO cli.py:274 -- Next steps
2023-10-18 03:19:52,197	INFO cli.py:275 -- Query the logs of the job:
2023-10-18 03:19:52,198	INFO cli.py:277 -- ray job logs job-test-bbfqs
2023-10-18 03:19:52,198	INFO cli.py:279 -- Query the status of the job:
2023-10-18 03:19:52,198	INFO cli.py:281 -- ray job status job-test-bbfqs
2023-10-18 03:19:52,198	INFO cli.py:283 -- Request the job to be stopped:
2023-10-18 03:19:52,198	INFO cli.py:285 -- ray job stop job-test-bbfqs
2023-10-18 03:19:52,203	INFO cli.py:292 -- Tailing logs until the job exits (disable with --no-wait):
2023-10-18 03:20:00,014	INFO worker.py:1329 -- Using address 10.244.0.10:6379 set in the environment variable RAY_ADDRESS
2023-10-18 03:20:00,014	INFO worker.py:1458 -- Connecting to existing Ray cluster at address: 10.244.0.10:6379...
2023-10-18 03:20:00,032	INFO worker.py:1633 -- Connected to Ray cluster. View the dashboard at 10.244.0.10:8265 
test_counter got 1
test_counter got 2
test_counter got 3
test_counter got 4
test_counter got 5
2023-10-18 03:20:03,304	SUCC cli.py:60 -- ------------------------------
2023-10-18 03:20:03,304	SUCC cli.py:61 -- Job 'job-test-bbfqs' succeeded
2023-10-18 03:20:03,304	SUCC cli.py:62 -- ------------------------------
```

#### List all jobs in a given namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/jobs
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/jobs' \
  -H 'accept: application/json'
  ```

* Response

  ```json
  {
    "jobs": [
      {
        "name": "rayjob-test",
        "namespace": "ray-system",
        "user": "3cp0",
        "entrypoint": "python -V",
        "jobId": "rayjob-test-drhlq",
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0",
              "replicas": 1,
              "minReplicas": 1,
              "maxReplicas": 1,
              "rayStartParams": {
                "metrics-export-port": "8080"
              }
            }
          ]
        },
        "createdAt": "2023-09-25T11:36:02Z",
        "jobStatus": "SUCCEEDED",
        "jobDeploymentStatus": "Running",
        "message": "Job finished successfully."
      }
    ]
  }
  ```

#### List all jobs in all namespaces

```text
GET {{baseUrl}}/apis/v1alpha2/jobs
```

Examples:

* Request:
  
  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/jobs' \
  -H 'accept: application/json' 
  ```

* Response

  ```json
  {
    "jobs": [
      {
        "name": "rayjob-test",
        "namespace": "ray-system",
        "user": "3cp0",
        "entrypoint": "python -V",
        "jobId": "rayjob-test-drhlq",
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0",
              "replicas": 1,
              "minReplicas": 1,
              "maxReplicas": 1,
              "rayStartParams": {
                "metrics-export-port": "8080"
              }
            }
          ]
        },
        "createdAt": "2023-09-25T11:36:02Z",
        "jobStatus": "SUCCEEDED",
        "jobDeploymentStatus": "Running",
        "message": "Job finished successfully."
      }
    ]
  }
  ```

#### Get job by its name and namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/jobs/<job_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/jobs/rayjob-test' \
  -H 'accept: application/json'
  ```

* Response:

  ```json
  {
    "name": "rayjob-test",
    "namespace": "ray-system",
    "user": "3cp0",
    "entrypoint": "python -V",
    "jobId": "rayjob-test-drhlq",
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0",
          "replicas": 1,
          "minReplicas": 1,
          "maxReplicas": 1,
          "rayStartParams": {
            "metrics-export-port": "8080"
          }
        }
      ]
    },
    "createdAt": "2023-09-25T11:36:02Z",
    "jobStatus": "SUCCEEDED",
    "jobDeploymentStatus": "Running",
    "message": "Job finished successfully."
  }
  ```

#### Delete job by its name and namespace

```text
DELETE {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/jobs/<job_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'DELETE' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/jobs/rayjob-test' \
  -H 'accept: application/json'
  ```

* Response
  
  ```json
  {}
  ```

### RayService

#### Create ray service in a given namespace

```text
POST {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/services
```

Examples:

* Request (V1)
  
  ```sh
  curl --silent -X 'POST' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
      "name": "test-v1",
      "namespace": "ray-system",
      "user": "user",
      "serviceUnhealthySecondThreshold": 900,
      "deploymentUnhealthySecondThreshold": 300,
      "serveDeploymentGraphSpec": {
        "importPath": "fruit.deployment_graph",
        "runtimeEnv": "working_dir: \"https://github.com/ray-project/test_dag/archive/c620251044717ace0a4c19d766d43c5099af8a77.zip\"\n",
        "serveConfigs": [
          {
            "deploymentName": "OrangeStand",
            "replicas": 1,
            "userConfig": "price: 2",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "PearStand",
            "replicas": 1,
            "userConfig": "price: 1",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "FruitMarket",
            "replicas": 1,
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          },
          {
            "deploymentName": "DAGDriver",
            "replicas": 1,
            "routePrefix": "/",
            "actorOptions": {
              "cpusPerActor": 0.1
            }
          }
        ]
      },
      "clusterSpec": {
        "headGroupSpec": {
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "serviceType": "NodePort",
          "rayStartParams": {
            "dashboard-host": "0.0.0.0",
            "metrics-export-port": "8080"
          },
          "volumes": []
        },
        "workerGroupSpec": [
          {
            "groupName": "small-wg",
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0-py310",
            "replicas": 1,
            "minReplicas": 0,
            "maxReplicas": 5,
            "rayStartParams": {
              "node-ip-address": "$MY_POD_IP"
            }
          }
        ]
      }
    }'
  ```

* Response (V1)

  ```json
  {
    "name": "test-v1",
    "namespace": "ray-system",
    "user": "user",
    "serveDeploymentGraphSpec": {
      "importPath": "fruit.deployment_graph",
      "runtimeEnv": "working_dir: \"https://github.com/ray-project/test_dag/archive/c620251044717ace0a4c19d766d43c5099af8a77.zip\"\n",
      "serveConfigs": [
        {
          "deploymentName": "OrangeStand",
          "replicas": 1,
          "userConfig": "price: 2",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "PearStand",
          "replicas": 1,
          "userConfig": "price: 1",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "FruitMarket",
          "replicas": 1,
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "DAGDriver",
          "replicas": 1,
          "routePrefix": "/",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        }
      ]
    },
    "serviceUnhealthySecondThreshold": 900,
    "deploymentUnhealthySecondThreshold": 300,
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0-py310",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "replicas": 1,
          "minReplicas": 5,
          "maxReplicas": 1,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    },
    "rayServiceStatus": {},
    "createdAt": "2023-09-25T11:42:11Z",
    "deleteAt": "1969-12-31T23:59:59Z"
  }
  ```

* Request (V2)  

  ```sh
  curl --silent -X 'POST' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '
  {
    "name": "test-v2",
    "namespace": "ray-system",
    "user": "user",
    "serviceUnhealthySecondThreshold": 900,
    "deploymentUnhealthySecondThreshold": 300,
    "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0-py310",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        },
        "volumes": []
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "replicas": 1,
          "minReplicas": 0,
          "maxReplicas": 5,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    }
  }'
  ```

* Response (V2)

  ```json
  {
    "name": "test-v2",
    "namespace": "ray-system",
    "user": "user",
    "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
    "serviceUnhealthySecondThreshold": 900,
    "deploymentUnhealthySecondThreshold": 300,
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0-py310",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "replicas": 1,
          "minReplicas": 5,
          "maxReplicas": 1,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    },
    "rayServiceStatus": {},
    "createdAt": "2023-09-25T11:44:41Z",
    "deleteAt": "1969-12-31T23:59:59Z"
  }
  ```

#### List all services in a given namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/services
```

Examples

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services' \
  -H 'accept: application/json'
  ```

* Response:

  ```json
  {
    "services": [
      {
        "name": "test-v2",
        "namespace": "ray-system",
        "user": "user",
        "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
        "serviceUnhealthySecondThreshold": 900,
        "deploymentUnhealthySecondThreshold": 300,
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0-py310",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0-py310",
              "replicas": 1,
              "minReplicas": 5,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "rayServiceStatus": {
          "rayServiceEvents": [
            {
              "id": "test-v2.178821ac4b15c743",
              "name": "test-v2-test-v2.178821ac4b15c743",
              "createdAt": "2023-09-25T11:44:43Z",
              "firstTimestamp": "2023-09-25T11:44:43Z",
              "lastTimestamp": "2023-09-25T11:46:15Z",
              "reason": "ServiceUnhealthy",
              "message": "The service is in an unhealthy state. Controller will perform a round of actions in 10s.",
              "type": "Normal",
              "count": 11
            }
          ]
        },
        "createdAt": "2023-09-25T11:44:41Z",
        "deleteAt": "1969-12-31T23:59:59Z"
      },
      {
        "name": "test-v2",
        "namespace": "ray-system",
        "user": "user",
        "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
        "serviceUnhealthySecondThreshold": 900,
        "deploymentUnhealthySecondThreshold": 300,
        "clusterSpec": {
          "headGroupSpec": {
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "replicas": 1,
              "minReplicas": 5,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "rayServiceStatus": {
          "rayServiceEvents": [
            {
              "id": "test-v2.178821ac4b15c743",
              "name": "test-v2-test-v2.178821ac4b15c743",
              "createdAt": "2023-09-25T11:44:43Z",
              "firstTimestamp": "2023-09-25T11:44:43Z",
              "lastTimestamp": "2023-09-25T11:46:15Z",
              "reason": "ServiceUnhealthy",
              "message": "The service is in an unhealthy state. Controller will perform a round of actions in 10s.",
              "type": "Normal",
              "count": 11
            }
          ]
        },
        "createdAt": "2023-09-25T11:44:41Z",
        "deleteAt": "1969-12-31T23:59:59Z"
      }
    ]
  }
  ```

#### List all services in all namespaces

```text
GET {{baseUrl}}/apis/v1alpha2/services
```

Examples:

* Request

  ```sh
  curl --silent -X 'GET' \
  'http://localhost:31888/apis/v1alpha2/services' \
  -H 'accept: application/json'
  ```

* Response:

  ```json
  {
    "services": [
      {
        "name": "test-v2",
        "namespace": "ray-system",
        "user": "user",
        "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
        "serviceUnhealthySecondThreshold": 900,
        "deploymentUnhealthySecondThreshold": 300,
        "clusterSpec": {
          "headGroupSpec": {
            "computeTemplate": "default-template",
            "image": "rayproject/ray:2.7.0-py310",
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "computeTemplate": "default-template",
              "image": "rayproject/ray:2.7.0-py310",
              "replicas": 1,
              "minReplicas": 5,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "rayServiceStatus": {
          "rayServiceEvents": [
            {
              "id": "test-v2.178821ac4b15c743",
              "name": "test-v2-test-v2.178821ac4b15c743",
              "createdAt": "2023-09-25T11:44:43Z",
              "firstTimestamp": "2023-09-25T11:44:43Z",
              "lastTimestamp": "2023-09-25T11:47:55Z",
              "reason": "ServiceUnhealthy",
              "message": "The service is in an unhealthy state. Controller will perform a round of actions in 10s.",
              "type": "Normal",
              "count": 21
            }
          ]
        },
        "createdAt": "2023-09-25T11:44:41Z",
        "deleteAt": "1969-12-31T23:59:59Z"
      },
      {
        "name": "test-v2",
        "namespace": "ray-system",
        "user": "user",
        "serveConfigV2": "applications:\n  - name: fruit_app\n    import_path: fruit.deployment_graph\n    route_prefix: /fruit\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: MangoStand\n        num_replicas: 1\n        user_config:\n          price: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: OrangeStand\n        num_replicas: 1\n        user_config:\n          price: 2\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: PearStand\n        num_replicas: 1\n        user_config:\n          price: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: FruitMarket\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: DAGDriver\n        num_replicas: 1\n        ray_actor_options:\n          num_cpus: 0.1\n  - name: math_app\n    import_path: conditional_dag.serve_dag\n    route_prefix: /calc\n    runtime_env:\n      working_dir: \"https://github.com/ray-project/test_dag/archive/41d09119cbdf8450599f993f51318e9e27c59098.zip\"\n    deployments:\n      - name: Adder\n        num_replicas: 1\n        user_config:\n          increment: 3\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Multiplier\n        num_replicas: 1\n        user_config:\n          factor: 5\n        ray_actor_options:\n          num_cpus: 0.1\n      - name: Router\n        num_replicas: 1\n      - name: create_order\n        num_replicas: 1\n      - name: DAGDriver\n        num_replicas: 1\n",
        "serviceUnhealthySecondThreshold": 900,
        "deploymentUnhealthySecondThreshold": 300,
        "clusterSpec": {
          "headGroupSpec": {
            "serviceType": "NodePort",
            "rayStartParams": {
              "dashboard-host": "0.0.0.0",
              "metrics-export-port": "8080"
            }
          },
          "workerGroupSpec": [
            {
              "groupName": "small-wg",
              "replicas": 1,
              "minReplicas": 5,
              "maxReplicas": 1,
              "rayStartParams": {
                "node-ip-address": "$MY_POD_IP"
              }
            }
          ]
        },
        "rayServiceStatus": {
          "rayServiceEvents": [
            {
              "id": "test-v2.178821ac4b15c743",
              "name": "test-v2-test-v2.178821ac4b15c743",
              "createdAt": "2023-09-25T11:44:43Z",
              "firstTimestamp": "2023-09-25T11:44:43Z",
              "lastTimestamp": "2023-09-25T11:47:55Z",
              "reason": "ServiceUnhealthy",
              "message": "The service is in an unhealthy state. Controller will perform a round of actions in 10s.",
              "type": "Normal",
              "count": 21
            }
          ]
        },
        "createdAt": "2023-09-25T11:44:41Z",
        "deleteAt": "1969-12-31T23:59:59Z"
      }
    ]
  }
  ```

#### Get service by its name and namespace

```text
GET {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/services/<service_name>
```

Examples:

* Request:

  ```sh
  curl --silent -X 'GET' \
    'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services/test-v1' \
    -H 'accept: application/json'  
  ```

* Response:

  ```json
  {
    "name": "test-v1",
    "namespace": "ray-system",
    "user": "user",
    "serveDeploymentGraphSpec": {
      "importPath": "fruit.deployment_graph",
      "runtimeEnv": "working_dir: \"https://github.com/ray-project/test_dag/archive/c620251044717ace0a4c19d766d43c5099af8a77.zip\"\n",
      "serveConfigs": [
        {
          "deploymentName": "OrangeStand",
          "replicas": 1,
          "userConfig": "price: 2",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "PearStand",
          "replicas": 1,
          "userConfig": "price: 1",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "FruitMarket",
          "replicas": 1,
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        },
        {
          "deploymentName": "DAGDriver",
          "replicas": 1,
          "routePrefix": "/",
          "actorOptions": {
            "cpusPerActor": 0.1
          }
        }
      ]
    },
    "serviceUnhealthySecondThreshold": 900,
    "deploymentUnhealthySecondThreshold": 300,
    "clusterSpec": {
      "headGroupSpec": {
        "computeTemplate": "default-template",
        "image": "rayproject/ray:2.7.0-py310",
        "serviceType": "NodePort",
        "rayStartParams": {
          "dashboard-host": "0.0.0.0",
          "metrics-export-port": "8080"
        }
      },
      "workerGroupSpec": [
        {
          "groupName": "small-wg",
          "computeTemplate": "default-template",
          "image": "rayproject/ray:2.7.0-py310",
          "replicas": 1,
          "minReplicas": 5,
          "maxReplicas": 1,
          "rayStartParams": {
            "node-ip-address": "$MY_POD_IP"
          }
        }
      ]
    },
    "rayServiceStatus": {
      "rayServiceEvents": [
        {
          "id": "test-v1.1788218987842e1a",
          "name": "test-v1-test-v1.1788218987842e1a",
          "createdAt": "2023-09-25T11:42:14Z",
          "firstTimestamp": "2023-09-25T11:42:14Z",
          "lastTimestamp": "2023-09-25T11:42:15Z",
          "reason": "ServiceUnhealthy",
          "message": "The service is in an unhealthy state. Controller will perform a round of actions in 10s.",
          "type": "Normal",
          "count": 2
        },
        {
          "id": "test-v1.1788218a0a86434b",
          "name": "test-v1-test-v1.1788218a0a86434b",
          "createdAt": "2023-09-25T11:42:16Z",
          "firstTimestamp": "2023-09-25T11:42:16Z",
          "lastTimestamp": "2023-09-25T11:42:24Z",
          "reason": "WaitForServeDeploymentReady",
          "message": "Fail to create / update Serve deployments. If you observe this error consistently, please check \"Issue 5: Fail to create / update Serve applications.\" in https://github.com/ray-project/kuberay/blob/master/docs/guidance/rayservice-troubleshooting.md for more details. err: Put \"http://test-v1-raycluster-7grg7-head-svc.ray-system.svc.cluster.local:52365/api/serve/deployments/\": dial tcp 10.96.69.78:52365: connect: connection refused",
          "type": "Normal",
          "count": 6
        },
        {
          "id": "test-v1.1788218cf7437d6f",
          "name": "test-v1-test-v1.1788218cf7437d6f",
          "createdAt": "2023-09-25T11:42:29Z",
          "firstTimestamp": "2023-09-25T11:42:29Z",
          "lastTimestamp": "2023-09-25T11:42:29Z",
          "reason": "SubmittedServeDeployment",
          "message": "Controller sent API request to update Serve deployments on cluster test-v1-raycluster-7grg7",
          "type": "Normal",
          "count": 1
        },
        {
          "id": "test-v1.1788218cf9823bdf",
          "name": "test-v1-test-v1.1788218cf9823bdf",
          "createdAt": "2023-09-25T11:42:29Z",
          "firstTimestamp": "2023-09-25T11:42:29Z",
          "lastTimestamp": "2023-09-25T11:42:35Z",
          "reason": "ServiceNotReady",
          "message": "The service is not ready yet. Controller will perform a round of actions in 2s.",
          "type": "Normal",
          "count": 4
        },
        {
          "id": "test-v1.1788218ee2f173bc",
          "name": "test-v1-test-v1.1788218ee2f173bc",
          "createdAt": "2023-09-25T11:42:37Z",
          "firstTimestamp": "2023-09-25T11:42:37Z",
          "lastTimestamp": "2023-09-25T11:47:16Z",
          "reason": "Running",
          "message": "The Serve applicaton is now running and healthy.",
          "type": "Normal",
          "count": 140
        },
        {
          "id": "test-v1-raycluster-7grg7.1788218976cd5268",
          "name": "test-v1-test-v1-raycluster-7grg7.1788218976cd5268",
          "createdAt": "2023-09-25T11:42:13Z",
          "firstTimestamp": "2023-09-25T11:42:13Z",
          "lastTimestamp": "2023-09-25T11:42:13Z",
          "reason": "Created",
          "message": "Created service test-v1-raycluster-7grg7-head-svc",
          "type": "Normal",
          "count": 1
        },
        {
          "id": "test-v1-raycluster-7grg7.178821897b8cddee",
          "name": "test-v1-test-v1-raycluster-7grg7.178821897b8cddee",
          "createdAt": "2023-09-25T11:42:14Z",
          "firstTimestamp": "2023-09-25T11:42:14Z",
          "lastTimestamp": "2023-09-25T11:42:14Z",
          "reason": "Created",
          "message": "Created head pod test-v1-raycluster-7grg7-head-jfkq5",
          "type": "Normal",
          "count": 1
        },
        {
          "id": "test-v1-raycluster-7grg7.178821898168eedb",
          "name": "test-v1-test-v1-raycluster-7grg7.178821898168eedb",
          "createdAt": "2023-09-25T11:42:14Z",
          "firstTimestamp": "2023-09-25T11:42:14Z",
          "lastTimestamp": "2023-09-25T11:42:14Z",
          "reason": "Created",
          "message": "Created worker pod ",
          "type": "Normal",
          "count": 1
        }
      ],
      "rayClusterName": "test-v1-raycluster-7grg7",
      "serveApplicationStatus": [
        {
          "name": "default",
          "status": "RUNNING",
          "serveDeploymentStatus": [
            {
              "deploymentName": "DAGDriver",
              "status": "HEALTHY"
            },
            {
              "deploymentName": "FruitMarket",
              "status": "HEALTHY"
            },
            {
              "deploymentName": "MangoStand",
              "status": "HEALTHY"
            },
            {
              "deploymentName": "OrangeStand",
              "status": "HEALTHY"
            },
            {
              "deploymentName": "PearStand",
              "status": "HEALTHY"
            }
          ]
        }
      ]
    },
    "createdAt": "2023-09-25T11:42:11Z",
    "deleteAt": "1969-12-31T23:59:59Z"
  }
  ```

#### Delete service by its name and namespace

```text
DELETE {{baseUrl}}/apis/v1alpha2/namespaces/<namespace>/services/<service_name>
```

Examples:

* Request

  ```sh
  curl --silent -X 'DELETE' \
  'http://localhost:31888/apis/v1alpha2/namespaces/ray-system/services/test-v1' \
  -H 'accept: application/json'  
  ```

* Response

  ```json
  {}
  ```
