# Getting Started
## Prerequisites
[Install Bazel](https://docs.bazel.build/versions/master/install.html)

[Install GoLang](https://golang.org/doc/install)

[Install Docker](https://docs.docker.com/get-docker/)

[Install protoc, the protobuf compiler](https://developers.google.com/protocol-buffers/docs/downloads)

[Install Google Cloud SDK/Command Line Tools](https://cloud.google.com/sdk/docs/install)

## Initialize the proto submodule
```bash
git submodule update --init
cd proto
git checkout main
```

## Download the Service Account Credentials
These are necessary to connect to cloud services.

* [Go to service account page](https://console.cloud.google.com/iam-admin/serviceaccounts?project=food-prod&folder=&organizationId=&supportedpurview=project)
* For "default-service-account@food-prod.iam.gserviceaccount.com", click the 3 dots under "Actions" on the far right.
* Click "Create Key"
* In .bash_profile, `export GOOGLE_APPLICATION_CREDENTIALS=path/to/service_account_keys`.
    This will let you connect to cloud services like databases.
* For Container Registry access, follow the instructions below
    ([details here](https://cloud.google.com/container-registry/docs/advanced-authentication))
    *   ```bash
        sudo usermod -a -G docker ${USER}
        gcloud auth login
        gcloud auth activate-service-account default-service-account@food-prod.iam.gserviceaccount.com \
            --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
        gcloud auth configure-docker
        ```

Google cloud APIs will look for environment variable $GOOGLE_APPLICATION_CREDENTIALS
when connecting to cloud services. [More info here](https://cloud.google.com/docs/authentication/production)
We don't need to download service account keys on Compute Engine servers,
since they will automatically use the compute engine service account.

#### Keys are secret. Do not commit to git.

## Deploy and test server locally
Execute following commands from repo root.
```bash
export FDA_ROOT="path/to/fda-backend"
go get ...
./scripts/generate_buildfiles_and_build.sh
bazel run //server

# In another tab
bazel run //client
```

## Deploy remote server
#### NOTE:
* `./scripts/ON_VM_deploy_food_grpc_server.sh` must reside on the
target machine at location `./scripts/ON_VM_deploy_food_grpc_server.sh` to successfully deploy.
See `./scripts/deploy_server.sh` for details.
* [Uber H3 Library](https://github.com/uber/h3) requires core C libraries.
  If your glibc version is > 2.27, the docker image you build on your machine
  will not run on Compute Engine VMs. This happens for Ubuntu > 18.10,
  and possible other OS.
  A temporary hack is to spin up a high CPU VM, clone this repo on the instance,
  and build/deploy on the VM.
  Remember to stop the VM when done building to save $$$!


#### Steps
Execute following commands from repo root.
```bash
# Define gRPC service in Cloud Endpoints
./scripts/deploy_endpoints_service.sh
# Deploy server to Compute Engine instance.
./scripts/deploy_server.sh
```

## Test remote server
```
cd $FDA_ROOT
./scripts/test_server.sh
```

# Info

## Service Account
https://cloud.google.com/docs/authentication/production#command-line

## Containers
To run container locally: https://cloud.google.com/run/docs/testing/local

### Authentication for Container Registry
https://cloud.google.com/container-registry/docs/pushing-and-pulling
https://cloud.google.com/container-registry/docs/advanced-authentication

### Docker Error
https://stackoverflow.com/questions/21871479/docker-cant-connect-to-docker-daemon


## Bazel/Gazelle
[Bazel](https://docs.bazel.build/versions/master/build-ref.html)
is a build tool that generates a dependency graph using the BUILD file at the
root of each bazel package directory.

A [Gazelle](https://github.com/bazelbuild/bazel-gazelle) is any of many antelope
species in the genus Gazella

### Adding dependencies to WORKSPACE
`bazel run //:gazelle update-repos REPO_NAME`

Example:
`bazel run //:gazelle update-repos github.com/google/uuid`

#### Resources to get started with bazel and go
https://hardyantz.medium.com/getting-started-monorepo-golang-application-with-bazel-370ed1069b4f
https://filipnikolovski.com/posts/managing-go-monorepo-with-bazel/

### Proto dependencies reference
https://github.com/bazelbuild/rules_go/blob/0.19.0/go/workspace.rst#proto-dependencies

## gRPC Resources
[Simple Example](https://tutorialedge.net/golang/go-grpc-beginners-tutorial/)
[Setting HTTP headers](http://www.inanzzz.com/index.php/post/7l4u/sending-and-receiving-grpc-client-server-headers-in-golang)
[gRPC web Example](https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/main.go) 

## SSL
[Enabling SSL for Cloud Endpoints ESP](https://cloud.google.com/endpoints/docs/grpc/enabling-ssl)
[SSL Certificate for Load Balancer](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs)

### SSL Errors when connecting to server with self signed cert
```
ERR_CERT_AUTHORITY_INVALID
ERR_CERT_COMMON_NAME_INVALID
```
**TRY THIS:** Start Chrome with the following flag:
`google-chrome-stable --ignore-certificate-errors`

**ELSE TRY THIS:**

https://superuser.com/a/1036062

OR

Attempt to connect to insecure address. In the top left next to URL,
you see a red triangle that says "Not Secure." Click on it.
Export the certificate.
chrome://settings/certificates -> Authorities tab -> Import the certificate

### Load Balancer
An external HTTPS load balancer is configured in Google Cloud.
In order to use gRPC with the load balancer,
[it must communicate with backend servers over HTTP/2](https://cloud.google.com/load-balancing/docs/https#using_grpc_with_your_applications)

## Getting Started additional resources
[Getting Started Github Repo](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/endpoints/getting-started-grpc)

[Cloud Endpoints gRPC Overview/Tutorial](https://cloud.google.com/endpoints/docs/grpc/about-grpc)
