# [0.6.0](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.5...v0.6.0) (2025-04-02)


### Features

* update to context-aware funcs ([5da876e](https://github.com/Longsight/cloudstack-terraform-provider/commit/5da876ecfdb9368c2d72fa1fd3a9b906a7cbc014))

## [0.5.5](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.4...v0.5.5) (2025-04-01)


### Bug Fixes

* revert ([3ab3b00](https://github.com/Longsight/cloudstack-terraform-provider/commit/3ab3b00943a9c68248a391574df62741a1dd36b8))
* revert ([6ed6a6e](https://github.com/Longsight/cloudstack-terraform-provider/commit/6ed6a6e4318bedf9ae2e5efe820b74c60445e794))


### Reverts

* Revert "revert to earlier version" ([6cddb58](https://github.com/Longsight/cloudstack-terraform-provider/commit/6cddb583a4ea5803738d261f8b2bad1678f0b82a))
* Revert "revert to earlier version" ([342c770](https://github.com/Longsight/cloudstack-terraform-provider/commit/342c770d2995a8c003a2f518011334eaea222e94))

## [0.5.4](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.3...v0.5.4) (2025-04-01)


### Bug Fixes

* revert to ubuntu-20.04 to get libc6-2.31 ([6baa8bb](https://github.com/Longsight/cloudstack-terraform-provider/commit/6baa8bbab77ba0f4d6a294259c0205cb9dffa340))
* revert to ubuntu-20.04 to get libc6-2.31 ([050ca09](https://github.com/Longsight/cloudstack-terraform-provider/commit/050ca098bfdcd9fc0d67267c1d96fb5bf413b4e6))

## [0.5.3](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.2...v0.5.3) (2025-04-01)


### Bug Fixes

* build on 24.04 ([3271dde](https://github.com/Longsight/cloudstack-terraform-provider/commit/3271dde05fe01afc4bd609cf83ab4ec72b410c7a))
* build on 24.04 ([18cde5c](https://github.com/Longsight/cloudstack-terraform-provider/commit/18cde5c89940e89da611b710009ed44454300fbe))

## [0.5.2](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.1...v0.5.2) (2025-04-01)


### Bug Fixes

* revert to earlier version of libs ([2b0bf51](https://github.com/Longsight/cloudstack-terraform-provider/commit/2b0bf51a7c02050edbd44eec0d89e11f4adc3865))

## [0.5.1](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.0...v0.5.1) (2025-04-01)


### Bug Fixes

* trigger release ([e780bfc](https://github.com/Longsight/cloudstack-terraform-provider/commit/e780bfcccb1848eb9df4602576da0fa1ca60182c))
* trigger release ([189d271](https://github.com/Longsight/cloudstack-terraform-provider/commit/189d2712c362934b73898cf7320c6acb6ed11413))

## [0.5.1](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.0...v0.5.1) (2025-03-31)


### Bug Fixes

* trigger release ([e780bfc](https://github.com/Longsight/cloudstack-terraform-provider/commit/e780bfcccb1848eb9df4602576da0fa1ca60182c))
* trigger release ([189d271](https://github.com/Longsight/cloudstack-terraform-provider/commit/189d2712c362934b73898cf7320c6acb6ed11413))

## [0.5.1](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.0...v0.5.1) (2025-03-31)


### Bug Fixes

* trigger release ([e780bfc](https://github.com/Longsight/cloudstack-terraform-provider/commit/e780bfcccb1848eb9df4602576da0fa1ca60182c))
* trigger release ([189d271](https://github.com/Longsight/cloudstack-terraform-provider/commit/189d2712c362934b73898cf7320c6acb6ed11413))

## [0.5.1](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.0...v0.5.1) (2025-03-31)


### Bug Fixes

* trigger release ([e780bfc](https://github.com/Longsight/cloudstack-terraform-provider/commit/e780bfcccb1848eb9df4602576da0fa1ca60182c))
* trigger release ([189d271](https://github.com/Longsight/cloudstack-terraform-provider/commit/189d2712c362934b73898cf7320c6acb6ed11413))

## [0.5.1](https://github.com/Longsight/cloudstack-terraform-provider/compare/v0.5.0...v0.5.1) (2025-03-28)


### Bug Fixes

* trigger release ([e780bfc](https://github.com/Longsight/cloudstack-terraform-provider/commit/e780bfcccb1848eb9df4602576da0fa1ca60182c))
* trigger release ([189d271](https://github.com/Longsight/cloudstack-terraform-provider/commit/189d2712c362934b73898cf7320c6acb6ed11413))

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
