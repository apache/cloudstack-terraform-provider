---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_kubernetes_cluster_config"
sidebar_current: "docs-cloudstack-datasource-kubernetes-cluster-config"
description: |-
  Get the kubeconfig of a CloudStack Kubernetes cluster.
---

# cloudstack_kubernetes_cluster_config

Use this data source to retrieve the kubeconfig of a CloudStack Kubernetes (CKS) cluster, so that
the `kubernetes` and `helm` providers can be pointed at a cluster managed by Terraform.

The cluster must be running for its config to be available. CloudStack does not return a config
while a cluster is still being set up, or when the Kubernetes service plugin is disabled
(`cloud.kubernetes.service.enabled`).

~> **Note:** The kubeconfig contains the cluster administrator credentials, and every attribute of
this data source is stored in plain text in your Terraform state. Protect the state accordingly,
for example by using a remote backend with encryption at rest.

## Example Usage

```hcl
resource "cloudstack_kubernetes_cluster" "basic" {
  name               = "basic-cluster"
  zone               = "zone1"
  kubernetes_version = "1.25.0"
  service_offering   = "Medium Instance"
  size               = 3
  description        = "Basic Kubernetes cluster"
}

data "cloudstack_kubernetes_cluster_config" "basic" {
  cluster_id = cloudstack_kubernetes_cluster.basic.id
}
```

### Configuring the Kubernetes and Helm providers

```hcl
provider "kubernetes" {
  host                   = data.cloudstack_kubernetes_cluster_config.basic.endpoint
  cluster_ca_certificate = data.cloudstack_kubernetes_cluster_config.basic.cluster_ca_certificate
  client_certificate     = data.cloudstack_kubernetes_cluster_config.basic.client_certificate
  client_key             = data.cloudstack_kubernetes_cluster_config.basic.client_key
}

provider "helm" {
  kubernetes {
    host                   = data.cloudstack_kubernetes_cluster_config.basic.endpoint
    cluster_ca_certificate = data.cloudstack_kubernetes_cluster_config.basic.cluster_ca_certificate
    client_certificate     = data.cloudstack_kubernetes_cluster_config.basic.client_certificate
    client_key             = data.cloudstack_kubernetes_cluster_config.basic.client_key
  }
}
```

### Writing the kubeconfig to a file

```hcl
resource "local_sensitive_file" "kubeconfig" {
  filename = "${path.module}/kubeconfig"
  content  = data.cloudstack_kubernetes_cluster_config.basic.config_data
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kubernetes cluster to retrieve the config for.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes cluster.
* `name` - The name of the Kubernetes cluster.
* `config_data` - The raw kubeconfig of the Kubernetes cluster.
* `endpoint` - The URL of the Kubernetes API server.
* `cluster_ca_certificate` - The PEM encoded certificate authority of the Kubernetes API server.
* `client_certificate` - The PEM encoded client certificate used to authenticate against the
  Kubernetes API server.
* `client_key` - The PEM encoded client key used to authenticate against the Kubernetes API server.

The `endpoint`, `cluster_ca_certificate`, `client_certificate` and `client_key` attributes are
extracted from `config_data`, resolved through the kubeconfig's current context, and base64 decoded
where applicable. If a kubeconfig does not carry one of them, the corresponding attribute is empty
and `config_data` can be used directly instead:

```hcl
locals {
  kubeconfig = yamldecode(data.cloudstack_kubernetes_cluster_config.basic.config_data)
}
```
