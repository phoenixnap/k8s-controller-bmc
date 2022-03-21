<h1 align="center">
  <br>
  <a href="https://phoenixnap.com/bare-metal-cloud"><img src="https://user-images.githubusercontent.com/78744488/109779287-16da8600-7c06-11eb-81a1-97bf44983d33.png" alt="phoenixnap Bare Metal Cloud" width="300"></a>
  <br>
  Bare Metal Cloud Controller for Kubernetes
  <br>
</h1>

<p align="center">
The Bare Metal Cloud Controller for Kubernetes allows developers to define, deploy, and manage Bare Metal Cloud servers in a unified way directly from a Kubernetes cluster. 
</p>

<p align="center">
  <a href="https://phoenixnap.com/bare-metal-cloud">Bare Metal Cloud</a> •
  <a href="https://developers.phoenixnap.com/">Developers Portal</a> •
  <a href="https://developers.phoenixnap.com/apis">API Documentation</a> •
  <a href="http://phoenixnap.com/kb">Knowledge Base</a> •
  <a href="https://developers.phoenixnap.com/support">Support</a>
</p>

## Creating a Bare Metal Cloud Account
You need to have a Bare Metal Cloud account to use this Kubernetes controller.  

1. Go to the [Bare Metal Cloud signup page](https://support.phoenixnap.com/wap-jpost3/bmcSignup).
2. Follow the prompts to set up your account.
3. Use your credentials to [log in to the Bare Metal Cloud portal](https://bmc.phoenixnap.com).

:arrow_forward: **Video tutorial:** [How to Create a Bare Metal Cloud Account](https://www.youtube.com/watch?v=RLRQOisEB-k)
<br>
:arrow_forward: **Video tutorial:** [Introduction to Bare Metal Cloud](https://www.youtube.com/watch?v=8TLsqgLDMN4)

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

The controller is available as a Docker image here: [docker.pkg.github.com/phoenixnap/k8s-bmc/bmc-server-controller:latest](docker.pkg.github.com/phoenixnap/k8s-bmc/bmc-server-controller:latest).

## Testing and CI

1. Set `USE_EXISTING_CLUSTER=true` to execute a build and test using an existing Kubernetes cluster, as specified by the active context.
1. Set `KUBEBUILDER_ASSETS=<KUBEBUILDER_BINS_LOCATION>` to execute tests using the Kubebuilder and Kubernetes binaries at the specified location.

## Retrieving BMC Credentials

1. [Log in to the Bare Metal Cloud portal](https://bmc.phoenixnap.com). 
2. On the left side menu, click on API Credentials. 
3. Click the Create Credentials button. 
4. Fill in the Name and Description fields, select the permissions scope and click Create. 
5. In the table, click on Actions and select View Credentials from the dropdown.  

:bulb: **Tutorial:** [How to create and manage BMC credentials](https://developers.phoenixnap.com/resources)

## Note to Maintainers

Be careful moving this repository. This project is written in Go and as such uses Git repo URLs as package identifiers. If the code URL is changed the code will need to be changed appropriately.

This is a `kubebuilder` project. Only minimal changes have been made to this codebase from the generated scaffolding so that maintainers can leverage as much off-the-shelf tooling and documentation as possible from the `kubebuilder` project. The bulk of the application code lives in the controller component at, `controllers/server_controller.go`. The API type definitions, defaulting and validating webhook logic live in the directory, `api/v1`.

## Bare Metal Cloud Community
Become part of the Bare Metal Cloud community to get updates on new features, help us improve the platform, and engage with developers and other users. 

-   Follow [@phoenixNAP on Twitter](https://twitter.com/phoenixnap)
-   Join the [official Slack channel](https://phoenixnap.slack.com)
-   Sign up for our [Developers Monthly newsletter](https://phoenixnap.com/developers-monthly-newsletter)

### Bare Metal Cloud Resources
-	[Product page](https://phoenixnap.com/bare-metal-cloud)
-	[Instance pricing](https://phoenixnap.com/bare-metal-cloud/instances)
-	[YouTube tutorials](https://www.youtube.com/watch?v=8TLsqgLDMN4&list=PLWcrQnFWd54WwkHM0oPpR1BrAhxlsy1Rc&ab_channel=PhoenixNAPGlobalITServices)
-	[Developers Portal](https://developers.phoenixnap.com)
-	[Knowledge Base](https://phoenixnap.com/kb)
-	[Blog](https:/phoenixnap.com/blog)

### Documentation
-	[API documentation](https://developers.phoenixnap.com/apis)

### Contact phoenixNAP
Get in touch with us if you have questions or need help with Bare Metal Cloud. 

<p align="left">
  <a href="https://twitter.com/phoenixNAP">Twitter</a> •
  <a href="https://www.facebook.com/phoenixnap">Facebook</a> •
  <a href="https://www.linkedin.com/company/phoenix-nap">LinkedIn</a> •
  <a href="https://www.instagram.com/phoenixnap">Instagram</a> •
  <a href="https://www.youtube.com/user/PhoenixNAPdatacenter">YouTube</a> •
  <a href="https://developers.phoenixnap.com/support">Email</a> 
</p>

<p align="center">
  <br>
  <a href="https://phoenixnap.com/bare-metal-cloud"><img src="https://user-images.githubusercontent.com/81640346/115243282-0c773b80-a123-11eb-9de7-59e3934a5712.jpg" alt="phoenixnap Bare Metal Cloud"></a>
</p>
