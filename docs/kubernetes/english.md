Contents
=============

    * [Deploy RadonDB MySQL cluster on Kubernetes](#deploy-radondb-mysql-cluster on-kubernetes-)
       * [Introduction](#Introduction)
       * [Deployment preparation](#deployment preparation)
       * [Deployment step](#deployment step)
          * [Step 1: Clone code](#step-1-clone code)
          * [Step 2: Deploy Operator](#step-2-deployment-operator)
          * [Step 3: Deploy RadonDB MySQL Cluster](#step-3-deployment-radondb-mysql-cluster)
       * [Deployment Verification](#Deployment Verification)
          * [Verify RadonDB MySQL Operator](#verify-radondb-mysql-operator)
          * [Verify RadonDB MySQL Cluster](#verify-radondb-mysql-cluster)
       * [Connect RadonDB MySQL](#connect-radondb-mysql)
          * [Same NameSpace access](#same-namespace-access)
          * [Cross-namespace access](#cross-namespace-access)
       * [Uninstall](#uninstall)
          * [Uninstall Operator](#uninstall-Operator)
          * [Uninstall RadonDB MySQL](#uninstall-RadonDB-MySQL)
          * [Uninstall custom resources](#uninstall custom resources)
       * [Configuration](#configuration)
          * [Container Configuration](#Container Configuration)
          * [Node configuration](#node configuration)
          * [Persistence Configuration](#Persistence Configuration)
       * [reference](#reference)

# Deploy RadonDB MySQL cluster (Operator) on Kubernetes

## Introduction

RadonDB MySQL is an open source, highly available, cloud-native cluster solution based on MySQL. It supports a high-availability architecture with one master and multiple slaves, and has a full set of management functions such as security, automatic backup, monitoring and alarming, and automatic expansion. It has been used on a large scale in production environments, including banks, insurance companies, traditional large enterprises, etc.

RadonDB MySQL supports installation, deployment and management on Kubernetes, automating tasks associated with running a RadonDB MySQL cluster.

This tutorial mainly demonstrates how to deploy RadonDB MySQL cluster (Operator) on Kubernetes.

## Deployment preparation

* Kubernetes cluster is ready for use.

## Deployment steps

### Step 1: Add helm repository

```
helm repo add radondb https://radondb.github.io/radondb-mysql-kubernetes/
```

Verify the warehouse and you can see the chart named `radondb/mysql-operator`.
```
helm search repo
NAME CHART VERSION APP VERSION DESCRIPTION
radondb/mysql-operator 0.1.1 latest Open Source, High Availability Cluster, based on MySQL
```

### Step 2: Deploy Operator



The following specifies the release name as `demo` and creates a [Deployment](#7-deployments) named `demo-mysql-operator`.

```
helm install demo radondb/mysql-operator
```

> Note: In this step, the [CRD](#8-CRD) required by the cluster will be created at the same time by default.

### Step 3: Deploy RadonDB MySQL Cluster

Execute the following command to create an instance for CRD `mysqlclusters.mysql.radondb.com` with default parameters, that is, create a RadonDB MySQL cluster. You can refer to [Configuration](#configuration) to customize cluster deployment parameters.

```kubectl
kubectl apply -f https://raw.githubusercontent.com/radondb/radondb-mysql-kubernetes/main/config/samples/mysql_v1alpha1_mysqlcluster.yaml
```

## Deployment verification

### Verify RadonDB MySQL Operator

Check the Deployment and corresponding monitoring service of `demo`. If the following information is displayed, the deployment is successful.

```shell
kubectl get deployment,svc
NAME READY UP-TO-DATE AVAILABLE AGE
demo-mysql-operator 1/1 1 1 7h50m


NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
service/mysql-operator-metrics ClusterIP 10.96.142.22 <none> 8443/TCP 8h
```

### Verify RadonDB MySQL Cluster

Execute the following command and you will see the following CRD.

```shell
kubectl get crd | grep mysql.radondb.com
backups.mysql.radondb.com 2021-11-02T07:00:01Z
mysqlclusters.mysql.radondb.com 2021-11-02T07:00:01Z
mysqlusers.mysql.radondb.com 2021-11-02T07:00:01Z
```

Taking the default deployment as an example, execute the following command to view the three-node RadonDB MySQL cluster named `sample-mysql` and the services used to access the nodes.

```shell
kubectl get statefulset,svc
NAME READY AGE
sample-mysql 3/3 7h33m

NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
service/sample-follower ClusterIP 10.96.131.84 <none> 3306/TCP 7h37m
service/sample-leader ClusterIP 10.96.111.214 <none> 3306/TCP 7h37m
service/sample-mysql ClusterIP None <none> 3306/TCP 7h37m
```

## Connect to RadonDB MySQL

You need to prepare a client to connect to MySQL.

### Same as NameSpace access

When the client is in the same NameSpace as the RadonDB MySQL cluster, the leader/follower service name can be used instead of the specific IP and port.

* Connect to the main node (read-write node).

     ```shell
     $ mysql -h <leader service name> -u <username> -p
     ```

    The user name is `radondb_usr`, and the release name is `sample`. The example of connecting to the master node is as follows:

     ```shell
     $ mysql -h sample-leader -u radondb_usr -p
     ```

* Connect to the slave node (read-only node).

     ```shell
     $ mysql -h <follower service name> -u <username> -p
     ```

    The user name is `radondb_usr`, the release name is `sample`, and the example of connecting the slave node is as follows:

     ```shell
     $ mysql -h sample-follower -u radondb_usr -p
     ```

### Cross NameSpace access

When the client is not in the same NameSpace as the RadonDB MySQL cluster, it can connect to the corresponding node through podIP or service ClusterIP.

1. Query the pod list and service list to obtain the pod name or corresponding service name of the node to be connected.

     ```shell
     $ kubectl get pod,svc
     ```

2. View the detailed information of the pod/service and obtain the corresponding IP.

     ```shell
     $ kubectl describe pod <pod name>
     $ kubectl describe svc <service name>
     ```

     > Note: The pod IP will be updated after the pod is restarted. You need to obtain the pod IP again. It is recommended to use the ClusterIP of the service to connect to the node.

3. Connect the nodes.

     ```shell
     $ mysql -h <pod IP/service ClusterIP> -u <username> -p
     ```

     The user name is `radondb_usr`, the Cluster IP is `10.10.128.136`, and the connection example is as follows:

     ```shell
     $ mysql -h 10.10.128.136 -u radondb_usr -p
     ```

## uninstall

### Uninstall Operator

Uninstall the RadonDB MySQL Operator named `demo` released in the current namespace.

```shell
helm delete demo
```

### Uninstall RadonDB MySQL

Uninstall the release named `sample` RadonDB MySQL cluster.

```shell
kubectl delete mysqlclusters.mysql.radondb.com sample
```

### Uninstall custom resources

```shell
kubectl delete customresourcedefinitions.apiextensions.k8s.io mysqlclusters.mysql.radondb.com
kubectl delete customresourcedefinitions.apiextensions.k8s.io mysqlu
