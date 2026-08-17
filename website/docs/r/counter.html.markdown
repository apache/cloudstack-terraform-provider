---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_counter"
sidebar_current: "docs-cloudstack-counter"
description: |-
  Creates a counter for autoscale policies.
---

# cloudstack_counter

Creates a counter that can be used in autoscale conditions to monitor performance metrics.

## Example Usage

```hcl
resource "cloudstack_counter" "cpu_counter" {
  name   = "cpu-counter"
  source = "cpu"
  value  = "cpuused"
}

resource "cloudstack_counter" "memory_counter" {
  name   = "memory-counter" 
  source = "memory"
  value  = "memoryused"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the counter.

* `source` - (Required) The source of the counter (e.g., "cpu", "memory", "network").

* `value` - (Required) The specific metric value to monitor (e.g., "cpuused", "memoryused").

## Attributes Reference

The following attributes are exported:

* `id` - The counter ID.
