## 0.4.0 (Unreleased)

IMPROVEMENTS:

* Restore support for managing resource tags as CloudStack 4.11.3+ and 4.12+ support tags again [GH-65]
* Update license to Apache License, Version 2.0
* Remove obsolete vendor directory (no longer needed when using modules)

## 0.3.0 (May 29, 2019)

NOTES:

* While most resources of this provider should now work with CloudStack 4.12.0.0, there are a
  few resources (mainly the network related resources) that will not work properly yet. See this
  related CloudStack issue for more details: https://github.com/apache/cloudstack/issues/3321

IMPROVEMENTS:

* Add support to import resources when using projects ([#56](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/56))
* Updated the provider to work with CloudStack 4.12.0.0 ([#64](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/64))
* Add Terraform 0.12 support ([#64](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/64))

## 0.2.0 (January 03, 2019)

IMPROVEMENTS:

* Make the tests work with the CloudStack Simulator ([#46](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/46))
* Remove support for managing resource tags as CloudStack >= 4.9.3 does [not support tags](https://github.com/apache/cloudstack/issues/3002) properly ([#46](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/46))
* Add basic support for importing resources ([#47](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/47))
* `r/clouddstack_instance`: Allow user-data to be plain text or base64 encoded text ([#48](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/48))
* `r/cloudstack_loadbalancer_rule`: Add support for SSL offloading ([#49](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/49))
* `r/cloudstack_instance`: Support deploying a stopped VM ([#50](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/50))

## 0.1.5 (April 27, 2018)

IMPROVEMENTS:

* `r/cloudstack_secondary_ipaddress`: Read back the secondary IP address details after creation ([#34](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/34))

BUG FIXES:

* `r/cloudstack_instance`: Root volume size in returned in bytes instead of GiB ([#32](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/32))

## 0.1.4 (January 04, 2018)

IMPROVEMENTS:

* `r/cloudstack_instance`: Properly reflect changes to the `root_disk_size` ([#28](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/28))

BUG FIXES:

* Fix the check that determines if tags should be set ([#29](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/29))

## 0.1.3 (December 20, 2017)

BUG FIXES:

* Make sure tags work as expected on all resources ([#26](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/26))

## 0.1.2 (December 18, 2017)

FEATURES:

* **New Data Source:** `d/cloudstack_template` ([#21](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/21))

IMPROVEMENTS:

* `r/cloudstack_loadbalancer_rule`: Support setting a `protocol` ([#22](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/22))

BUG FIXES:

* `r/cloudstack_ipaddress`: Fix a panic when trying to disassociate public IP ([#18](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/18))
* Add tags in resources vpc, instance, ipadress, template and vpc ([#16](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/16))

## 0.1.1 (August 04, 2017)

IMPROVEMENTS:

* `r/cloudstack_customer_gateway`: Add support for using projects ([#3](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/3))

BUG FIXES:

* `r/cloudstack_instance`: Prevent a potential crash when deleting an instance ([#7](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/7))
* `r/cloudstack_security_group_rule`: Fix a panic when trying to read a deleted security group ([#2](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/2))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
