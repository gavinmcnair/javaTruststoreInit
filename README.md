# javaTruststoreInit

![GitHub](https://img.shields.io/github/license/gavinmcnair/javatruststoreinit)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gavinmcnair/javatruststoreinit)
![CircleCI](https://img.shields.io/circleci/build/github/gavinmcnair/javaTruststoreInit/main?token=aab7daba901f49034a2fb9f61895b61114b13de9)


## Problem statement

Do you have a Java application which uses a `p12` keys and certicates and a CA encoded into a  `JKS` file but you only have a standard pem encoded Keys, CA's and Certificate?

javaTruststoreInit is an `initContainer` which takes certificates from either local files or environment variables and writes out a Java Keystore (JKS) file to an emptyDir which can be shared with the main container

| Environment Variable  | Default  | Description  |
|---|---|---|
| PASSWORD  | password  | The password used for the keystore|
| FILE_MODE  | false | If to use the env vars or files  |
| KAFKA_KEY  |  NA | Public Key environment variable |
| KAFKA_CA  |  NA | CA Certificate environment variable |
| KAFKA_CERT  |  NA | Certificate environment variable  |
| KAFKA_KEY_FILE  |  NA |  Public Key file |
| KAFKA_CA_FILE  |  NA | CA Certificate file |
| KAFKA_CERT_FILE  | NA  | Certificate file  |
| OUTPUT_P12  | /var/run/secrets/truststore.p12  | The filename used to write the Private Key and Certificate into |
| OUTPUT_JKS  | /var/run/secrets/truststore.jks  | The filename used to write the CA into |

## How to use in Kubernetes

We can supply the PEM encoded `key` and `certificate` either within the environment variable or as files mounted upon the filesystem. Both of which can be sourced with secrets or configmaps as appropriate. When using files you need to set `FILE_MODE` to `true`

The init container will start and write the output file to the `OUTPUT_FILE` path.

This is then available to the target JVM.

### Example pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: KafkaClient
spec:
  initContainers:
    - name: javaTruststoreInit
      image: gavinmcnair/javatruststoreinit:v1.0.3
      env:
        - name: KEY
          value: "pem encoded key"
        - name: CERTIFICATE
          value: "pem encoded cert"
        - name: KAFKA_CA_FILE
          value: "pem encoded ca"
      volumeMounts:
        - mountPath: /var/run/secrets
          name: kafkasecrets
  containers:
    - name: kafkaclient
      image: kafkaclient:1.0.0
      env:
        - name: JAVA_P12_FILE
          value: "/var/run/secrets/truststore.p12"
        - name: JAVA_JKS_FILE
          value: "/var/run/secrets/truststore.jks"
        - name: JAVA_KEYSTORE_PASSWORD
          value: "password"
      volumeMounts:
        - mountPath: /var/run/secrets
          name: kafkasecrets
  volumes:
    - emptyDir: {}
      name: kafkasecrets

```

## Motivation

In the conventional way we need to use an insecure Java container which often contains an entire Linux operating system. 

This already large insecure container then has to execute multiple java keystore commands.

In comparison this container is a single binary build upon a scratch container. Its much smaller and has far less security implications.

It should be both quick and reliable.
