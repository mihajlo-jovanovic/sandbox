## How to run Kafka locally

1. Start by getting a local k8s cluster

```
minikube start --kubernetes-version=v1.19.14 --driver=hyperkit --container-runtime=docker
```

Note we are using hyperkit on a Mac as well as docker runtime; this seems more CPU efficient than Docker Desktop, which 
I used previously. Version of k8s is required for Confluent helm charts to work correctly, see https://docs.confluent.io/5.1.0/installation/installing_cp/cp-helm-charts/docs/index.html

Confluent Operator 1.7.0 did not work for me; it seems to hang after starting the first Zookeeper node.

2. Check out Helm charts repo from GitHub; then in the root dir:

```
helm install --set cp-schema-registry.enabled=false,cp-kafka-rest.enabled=false,cp-kafka-connect.enabled=false,cp-zookeeper.servers=3,cp-kafka.brokers=3,cp-ksql-server.enabled=false,cp-control-center.enabled=false --generate-name .
```

This will start only Zookeeper & Kafka, with NodePort enabled for access outside the k8s cluster (remember to use 
correct `node` ip)

```bash
$(minikube ip):31090
```

## Client libs for Go

Note we are currently using sarama from Shopify, which is a pure Go lib. Apparently there is also a Confluent Kafka lib for Go, which is a
wrapper around `librdkafka`, a well-known and mature lib written in C.

## References

[Helm Charts from Confluent](https://github.com/confluentinc/cp-helm-charts)

[Confluent Operator](https://docs.confluent.io/operator/1.7.0/overview.html#operator-supported-environments)

[Minikube on MacOs without Docker Desktop](https://itnext.io/goodbye-docker-desktop-hello-minikube-3649f2a1c469)

[Minikube Drivers](https://minikube.sigs.k8s.io/docs/drivers/)