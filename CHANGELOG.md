## 0.1.2 (Unreleased)
## 0.1.1 (August 04, 2017)

IMPROVEMENTS:

* `cloudstack_customer_gateway`: Add support for using projects ([#3](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/3))

BUG FIXES:

* `cloudstack_instance`: Prevent a potential crash when deleting an instance ([#7](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/7))
* `cloudstack_security_group_rule`: Fix a panic when trying to read a deleted security group ([#2](https://github.com/terraform-providers/terraform-provider-cloudstack/issues/2))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
