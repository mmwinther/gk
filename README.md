# Kubectl and Gcloud wrapper

## Prereq

* kubectl installed
* gcloud installed and configured to communicate with your projects

## External dependencies

* gopkg.in/AlecAivazis/survey.v1
* github.com/fatih/color

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
