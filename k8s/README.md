## Running in Minikube

- Integrate minikube with local docker:

  ```console
  eval $(minikube -p minikube docker-env)
  ```
  
- Apply k8s resources:

  ```console
  kubectl apply -f posting-processor.yaml
  ```