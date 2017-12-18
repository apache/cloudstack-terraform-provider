## 0.1.2 (December 18, 2017)

FEATURES:

* **New Data Source:** `cloudstack_template` ([#21](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/21))

IMPROVEMENTS:

* `cloudstack_loadbalancer_rule`: Support setting a `protocol` ([#22](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/22))

BUG FIXES:

* `cloudstack_ipaddress`: Fix a panic when trying to disassociate public IP ([#18](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/18))
* Add tags in resources vpc, instance, ipadress, template and vpc ([#16](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/16))

## 0.1.1 (August 04, 2017)

IMPROVEMENTS:

* `cloudstack_customer_gateway`: Add support for using projects ([#3](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/3))

BUG FIXES:

* `cloudstack_instance`: Prevent a potential crash when deleting an instance ([#7](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/7))
* `cloudstack_security_group_rule`: Fix a panic when trying to read a deleted security group ([#2](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/2))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
