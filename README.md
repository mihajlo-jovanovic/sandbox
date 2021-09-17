[![codecov](https://codecov.io/gh/linus18/sandbox/branch/master/graph/badge.svg?token=UOQB1OKEWY)](https://codecov.io/gh/linus18/sandbox)
# Sandbox

My personal Golang playground to try out language features such as channels as well as concepts such as gRPC, async logging, etc.

## Packages

### Driver

This is a simple test driver client utility. It works by sending a stream of requests to a server endpoint and logs responses (along with timestamp) to a file.

- To run:

    ```console
    $ cd driver && go build
    $ ./driver
    ```
  
### tcp_server

This is the server endpoint that listens on a tcp port 8888 for a stream of string requests and simply echoes then back. Requests are 
separated by end-of-line \n.

- To run:

    ```console
    $ cd tcp_server && go build
    $ ./tcp_server
    ```
  
### posting_api_grpc

Simplest possible posting API

- To compile:

    ```console
    $ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative posting_api_grpc/posting.proto
    ```
  
* This will generate both the types and service interface definition
* Make sure to be in project root when running this command, note the relative path to proto file

## Running in Minikube

- Integrate minikube with local docker:

  ```console
  eval $(minikube -p minikube docker-env)
  ```

## Secret Management

Posting server assumes there is an instance of Hashi Vault running and accessible at vault:8200, with the database
password inserted at the following path and accessible by its service account. It used an init container to fetch
the secret and write it to a file from which the server reads it. See references for more info on how to set this up.
I will automate this when I get some time, but for now this needs to be done manually.

```
     ROOT_DB_PASSWD:                vault:secret/data/dev/postgres#ROOT_DB_PASSWD
     VAULT_ADDR:                    https://vault:8200
```

## References

[Go gRPC Quickstart](https://grpc.io/docs/languages/go/quickstart/) 

[Anatomy of modules in Go](https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16)

[Protocol Buffers Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial#the_protocol_buffer_api)

[Create the smallest and secured golang docker image based on scratch](https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324)

[Basic Postgres database in Kubernetes](https://itnext.io/basic-postgres-database-in-kubernetes-23c7834d91ef)

[Dynamic secrets on kubernetes pods using vault](https://gmaliar.medium.com/dynamic-secrets-on-kubernetes-pods-using-vault-35d9094d169)