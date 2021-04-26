# k8s-bmc

## Getting Started from Source

1. Configure your environment to communicate with a Kubernetes cluster (1.18+).
1. [Install Kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
1. [Install Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)
1. Clone this repository.
1. Run `make` to build the source.
1. Run `make install` to install the CRD on your cluster.
1. Configure your BMC credentials in your environment.
1. Run `make run` to run the controller locally.
1. Add your BMC credentials to a secret and wire that secret .
1. Run `make deploy`.

## Pulling the Image

The controller is available as a Docker image at: `docker.pkg.github.com/phoenixnap/k8s-bmc/bmc-server-controller:latest`.

## Testing and CI

1. Set `USE_EXISTING_CLUSTER=true` to execute a build and tests using an existing kubernetes cluster (as specified by the active context).
1. Set `KUBEBUILDER_ASSETS=<KUBEBUILDER_BINS_LOCATION>` to execute tests using the kubebuilder and kubernetes binaries at the specified location.

## Note to Maintainers

Becareful moving this repository. This project is written in Go and as such uses Git repo URLs as package identifiers. If the code URL is changed the code will need to be changed appropriately.

This is a `kubebuilder` project. Only minimal changes have been made to this codebase from the generated scaffolding so that maintainers can leverage as much off-the-shelf tooling and documentation as possible from the `kubebuilder` project. The bulk of the application code lives in the controller component at, `controllers/server_controller.go`. The API type definitions, defaulting and validating webhook logic live in the directory, `api/v1`.
