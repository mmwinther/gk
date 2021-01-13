# Kubectl and Gcloud wrapper

## Prerequisites

* golang
* kubectl installed
* gcloud installed and configured to communicate with your projects

## Install

1. Clone this repo:
    ```shell
    git clone git@github.com:zaikinlv/gk.git
    ```
1. Change working directory into the repo
1. Run this command to install dependencies:
    ```shell
    go mod tidy
    ```
1. Build the executable:
    ```shell
    go build main.go
    ```
1. Create a link to the executable:
    ```shell
    ln -s /path/to/the/repo/gk/main /usr/local/bin/gk
    ```
1. Check that it works with `gk -h`

## Run
```
#> gk -v
Print current app version
#> gk
List through existing contexts in ~/.kube/config and picking one i.e "kubectl use-context <chosen_context>"
#> gk -c
Listing availabe kubernetes clusteres in available gcloud projects then importing into kubeconfig
#> gk -t
Copy access-token to clipboard (when accessing kubernetes dashboard)
#> gk -p
Set current gcp project
#> gk -i
Print current gcp project and current kubeconfig context
#> gk -n
Set "default" namespace for current kube context
#> gk -clean
Remove all clusters and contexts from .kube/config
```
