# [1.0.0](https://github.com/Longsight/cloudstack-terraform-provider/compare/...v1.0.0) (2025-03-28)


### Bug Fixes

* hmm ([b9f0467](https://github.com/Longsight/cloudstack-terraform-provider/commit/b9f04674f4a4d9417bac0fd6ba41674e2859e23f))
* hmm ([5394dc0](https://github.com/Longsight/cloudstack-terraform-provider/commit/5394dc087f6985fe57196c1f5afbca12940d9e81))
* hmm ([2abb925](https://github.com/Longsight/cloudstack-terraform-provider/commit/2abb9257b579b4254933fe557208fd7e611aeea0))
* hmm ([98e51b3](https://github.com/Longsight/cloudstack-terraform-provider/commit/98e51b3c5705d730fa1b2a8a976c9039828701dc))
* hmm ([2e8507d](https://github.com/Longsight/cloudstack-terraform-provider/commit/2e8507d69a400712fd986c06c16f2b2d69340743))
* hmm ([f6b4efa](https://github.com/Longsight/cloudstack-terraform-provider/commit/f6b4efa871557fa0e5eb4f632ea7b1cf3cb94dd2))
* hmm ([baf1ed4](https://github.com/Longsight/cloudstack-terraform-provider/commit/baf1ed4a26f49ed745cc86177e2a4f40c9e92b7f))
* hmm ([f839c0e](https://github.com/Longsight/cloudstack-terraform-provider/commit/f839c0e4b50e5c6a7711aa02a8bb0c95e7162071))
* hmm ([d90170d](https://github.com/Longsight/cloudstack-terraform-provider/commit/d90170ddd3c99f998b00913a4c98e564db7e988a))
* hmm ([97cffef](https://github.com/Longsight/cloudstack-terraform-provider/commit/97cffef1359816691654ebe61956929084ff4a52))
* hmm ([7825192](https://github.com/Longsight/cloudstack-terraform-provider/commit/7825192c4e520405235e7767b3f68cf397381c42))
* hmm ([2034b5a](https://github.com/Longsight/cloudstack-terraform-provider/commit/2034b5a4ab6a1cdd606e86a688bcea737edf373e))
* hmm ([4b20551](https://github.com/Longsight/cloudstack-terraform-provider/commit/4b20551364d3f997d4748124ebc376181604fef3))
* hmm ([b90a757](https://github.com/Longsight/cloudstack-terraform-provider/commit/b90a757e8aa4c6f3819e4753c298fc4a0ccf1876))
* hmm ([1656e5b](https://github.com/Longsight/cloudstack-terraform-provider/commit/1656e5bb12dae0479e218275c2850a0f7e54dbcb))
* hmm ([6d1d333](https://github.com/Longsight/cloudstack-terraform-provider/commit/6d1d3331d727ebf9bcec33797effd1ae1825d1ac))
* hmm ([65a7822](https://github.com/Longsight/cloudstack-terraform-provider/commit/65a7822e3d64227f5f80af397608aa38b9c43964))
* hmm ([12e0ba7](https://github.com/Longsight/cloudstack-terraform-provider/commit/12e0ba7c0f2e1a9e2ff147a9a96f06cd082692af))
* hmm ([8f639a8](https://github.com/Longsight/cloudstack-terraform-provider/commit/8f639a8c619439eb3346d42fd4d107014255eecf))
* hmm ([ca30601](https://github.com/Longsight/cloudstack-terraform-provider/commit/ca30601f11039ff4633e2b80310505077bf2face))
* remove user-data conditional ([#121](https://github.com/Longsight/cloudstack-terraform-provider/issues/121)) ([02058b4](https://github.com/Longsight/cloudstack-terraform-provider/commit/02058b4f545131d52a42dbe6d5795d3f06d9a54c))
* set vlan in private network ([1edf0cf](https://github.com/Longsight/cloudstack-terraform-provider/commit/1edf0cf312fb44b243a4429edea812a615fba971))
* tests indentation ([2033b3f](https://github.com/Longsight/cloudstack-terraform-provider/commit/2033b3fa8cfedf27795b38112c82cdb1d5d9479a))


### Features

* migrate to terraform plugin framework ([#113](https://github.com/Longsight/cloudstack-terraform-provider/issues/113)) ([e6bcc74](https://github.com/Longsight/cloudstack-terraform-provider/commit/e6bcc7489e84db74e30117146b6dec8f76461eeb))


### Reverts

* Revert "docs: consistent use of array configuration syntax" ([34bddeb](https://github.com/Longsight/cloudstack-terraform-provider/commit/34bddebca9468e9cf19dc305b126cf8ee60d49a5))

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
