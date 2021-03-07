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
    $ protoc -I=$SRC_DIR --go_out=$DST_DIR ../posting.proto
    ```
  
* make sure to not include the full package name in DST_DIR, or else it will result in additional nested dirs.
