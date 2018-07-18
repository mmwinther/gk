# Kubectl and Gcloud wrapper

## Prereq

* kubectl installed
* gcloud installed and configured to communicate with your projects

## Dependencies

* gopkg.in/AlecAivazis/survey.v1 

## Run
```
#> gk
List through existing contexts in ~/.kube/config and picking one i.e "use-context"
#> gk -c
Listing availabe kubernetes clusteres in available gcloud projects then importing into kubeconfig 
#> gk -t 
Copy access-token to clipboard (when accessing kubernetes dashboard)
```