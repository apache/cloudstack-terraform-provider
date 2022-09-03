---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_ssh_keypair"
sidebar_current: "docs-cloudstack-cloudstack_ssh_keypair"
description: |-
  Gets information about cloudstack ssh keypair.
---

# cloudstack_ssh_keypair

Use this datasource to get information about a ssh keypair for use in other resources.

### Example Usage

```hcl
  data "cloudstack_ssh_keypair" "ssh-keypair-data" {
	  filter {
	  name = "name" 
	  value = "myKey"
	}
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the keypair.
* `fingerprint` - Fingerprint of the public key.