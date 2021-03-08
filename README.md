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

## References

[Go gRPC Quickstart](https://grpc.io/docs/languages/go/quickstart/) 

[Anatomy of modules in Go](https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16)

[Protocol Buffers Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial#the_protocol_buffer_api)