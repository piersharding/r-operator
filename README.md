## Kubernetes R Operator for ShinyApps

This operator manages R ShinyApp clusters.  The user provides a volume with the application in it, and can optionally provide a volume with additional libraries required.

### Requirements for running a ShinyApp

Typically, a ShinyApp would be based on the [Rocker project](https://www.rocker-project.org/images/) [shiny image](https://hub.docker.com/r/rocker/shiny) which may need customising - for instance, the following script builds a derivative image with [gdal support for rgdal](https://cran.r-project.org/web/packages/rgdal/index.html), and also creates an rlibs directory that can be unpacked into a volume (for example could be supplied over NFS, or host directory mount) for some further dependent R libraries.
```sh
#!/bin/sh
BASE=$(pwd)
DIR=${BASE}/rlibs
VOL="$DIR:/rlibs"
IMAGE=rocker/shiny:latest
TAG=piersharding/rgdal:latest

cat <<EOF | docker build -t ${TAG} -
FROM ${IMAGE}
RUN \
    apt update -y && \
    apt install -y libgdal-dev libproj-dev libssl-dev && \
    apt clean -y && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*
EOF

docker push ${TAG}

mkdir -p $DIR

cat <<EOF | docker run --rm -i -v ${VOL} ${TAG} bash

apt update && apt-get -y install libgdal-dev libproj-dev libssl-dev

R -e '.libPaths( c( .libPaths(), "/rlibs") ); setwd("/rlibs"); install.packages(c("dplyr", "readr", "tidyr", "lubridate", "scales", "tidyverse", "shinydashboard", "shinyBS", "shinyjs", "leaflet", "DT", "highcharter", "VennDiagram", "treemap", "circlize", "plotly", "sp", "rgdal", "png", "devtools"), lib="/rlibs/")'
EOF

cd ${DIR}
tar -czvf ${BASE}/rlibs.tar.gz *
```

### Prerequisites

* [Install Metacontroller](https://metacontroller.app/guide/install/)

MetaController can be installed with:
```sh
make metacontroller
```

### Install The Operator

r-operator can be installed with
```sh
make deploy # uninstall with 'make delete'
```

### Launch a ShinnyApp

```sh
cat <<EOF | kubectl apply -f -
---
apiVersion: piersharding.com/v1
kind: Rapp
metadata:
  name: app-1
spec:
  replicas: 3
  # default image
  # image: rocker/shiny:latest
  ingress: testapp.rapp.local 
  imagePullPolicy: IfNotPresent
  # add the following volumes and mounts to supply the application
  # and supporting libraries
  # volumes:
  #   - name: rscripts
  #     persistentVolumeClaim:
  #       claimName: app      
  #   - name: rlibs
  #     persistentVolumeClaim:
  #       claimName: rlibs      
  # volumeMounts:
  #   - mountPath: /rscripts
  #     readOnly: false
  #     name: rscripts
  #   - mountPath: /rlibs
  #     readOnly: false
  #     name: rlibs
EOF
```

Watch the cluster deploy:
```sh
$ kubectl get all,ingress,rapps
NAME                              READY   STATUS              RESTARTS   AGE
pod/rapp-app-1-59469cd577-2tbx2   0/1     ContainerCreating   0          2s
pod/rapp-app-1-59469cd577-4pcqs   0/1     ContainerCreating   0          2s
pod/rapp-app-1-59469cd577-lh28t   0/1     ContainerCreating   0          2s

NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP    33m
service/rapp-app-1   ClusterIP   10.99.65.165   <none>        8080/TCP   2s

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/rapp-app-1   0/3     3            0           2s

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/rapp-app-1-59469cd577   3         3         0       2s

NAME                            HOSTS                ADDRESS   PORTS   AGE
ingress.extensions/rapp-app-1   testapp.rapp.local             80      2s

NAME                          COMPONENTS   SUCCEEDED   AGE   STATE
rapp.piersharding.com/app-1   1            0           2s    Building
```

Get extended information with:
```sh
$ kubectl get rapps -o wide
NAME    COMPONENTS   SUCCEEDED   AGE   STATE     RESOURCES
app-1   1            1           57s   Running   Ingress: rapp-app-1 IP: 192.168.86.47, Hosts: http://testapp.rapp.local/ status: {"loadBalancer":{"ingress":[{"ip":"192.168.86.47"}]}} - Service: rapp-app-1 Type: ClusterIP, IP: 10.99.65.165, Ports: http/8080 status: {"loadBalancer":{}}
```
