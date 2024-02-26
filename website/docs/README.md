# CloudStack Provider

The CloudStack provider is used to interact with the many resources
supported by CloudStack. The provider needs to be configured with a
URL pointing to a running CloudStack API and the proper credentials
before it can be used.

In order to provide the required configuration options you can either
supply values for the `api_url`, `api_key` and `secret_key` fields, or
for the `config` and `profile` fields. A combination of both is not
allowed and will not work.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the CloudStack Provider
provider "cloudstack" {
  api_url    = "${var.cloudstack_api_url}"
  api_key    = "${var.cloudstack_api_key}"
  secret_key = "${var.cloudstack_secret_key}"
}

# Create a web server
resource "cloudstack_instance" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `api_url` - (Optional) This is the CloudStack API URL. It can also be sourced
  from the `CLOUDSTACK_API_URL` environment variable.

* `api_key` - (Optional) This is the CloudStack API key. It can also be sourced
  from the `CLOUDSTACK_API_KEY` environment variable.

* `secret_key` - (Optional) This is the CloudStack secret key. It can also be
  sourced from the `CLOUDSTACK_SECRET_KEY` environment variable.

* `config` - (Optional) The path to a `CloudMonkey` config file. If set the API
  URL, key and secret will be retrieved from this file.

* `profile` - (Optional) Used together with the `config` option. Specifies which
  `CloudMonkey` profile in the config file to use.

* `http_get_only` - (Optional) Some cloud providers only allow HTTP GET calls to
  their CloudStack API. If using such a provider, you need to set this to `true`
  in order for the provider to only make GET calls and no POST calls. It can also
  be sourced from the `CLOUDSTACK_HTTP_GET_ONLY` environment variable.

* `timeout` - (Optional) A value in seconds. This is the time allowed for Cloudstack
  to complete each asynchronous job triggered. If unset, this can be sourced from the
  `CLOUDSTACK_TIMEOUT` environment variable. Otherwise, this will default to 300
  seconds.

## Data Sources

- [instance](./d/instance.html.markdown)
- [ipaddress](./d/ipaddress.html.markdown)
- [network_offering](./d/network_offering.html.markdown)
- [service_offering](./d/service_offering.html.markdown)
- [ssh_keypair](./d/ssh_keypair.html.markdown)
- [template](./d/template.html.markdown)
- [user](./d/user.html.markdown)
- [volume](./d/volume.html.markdown)
- [vpc](./d/vpc.html.markdown)
- [vpn_connection](./d/vpn_connection.html.markdown)
- [zone](./d/zone.html.markdown)
## Resources

- [account](./r/account.html.markdown)
- [affinity_group](./r/affinity_group.html.markdown)
- [autoscale_vm_profile](./r/autoscale_vm_profile.html.markdown)
- [disk](./r/disk.html.markdown)
- [disk_offering](./r/disk_offering.html.markdown)
- [domain](./r/domain.html.markdown)
- [egress_firewall](./r/egress_firewall.html.markdown)
- [firewall](./r/firewall.html.markdown)
- [instance](./r/instance.html.markdown)
- [ipaddress](./r/ipaddress.html.markdown)
- [kubernetes_cluster](./r/kubernetes_cluster.html.markdown)
- [kubernetes_version](./r/kubernetes_version.html.markdown)
- [loadbalancer_rule](./r/loadbalancer_rule.html.markdown)
- [network](./r/network.html.markdown)
- [network_acl](./r/network_acl.html.markdown)
- [network_acl_rule](./r/network_acl_rule.html.markdown)
- [network_offering](./r/network_offering.html.markdown)
- [nic](./r/nic.html.markdown)
- [port_forward](./r/port_forward.html.markdown)
- [private_gateway](./r/private_gateway.html.markdown)
- [secondary_ipaddress](./r/secondary_ipaddress.html.markdown)
- [security_group](./r/security_group.html.markdown)
- [security_group_rule](./r/security_group_rule.html.markdown)
- [service_offering](./r/service_offering.html.markdown)
- [ssh_keypair](./r/ssh_keypair.html.markdown)
- [static_nat](./r/static_nat.html.markdown)
- [static_route](./r/static_route.html.markdown)
- [template](./r/template.html.markdown)
- [user](./r/user.html.markdown)
- [volume](./r/volume.html.markdown)
- [vpc](./r/vpc.html.markdown)
- [vpn_connection](./r/vpn_connection.html.markdown)
- [vpn_customer_gateway](./r/vpn_customer_gateway.html.markdown)
- [vpn_gateway](./r/vpn_gateway.html.markdown)
- [zone](./r/zone.html.markdown)
