# CloudStack Terraform Provider Test Suites

This directory contains test configurations for validating the CloudStack Terraform Provider.

## Test Suites

### comprehensive/
Full provider feature test covering VMs, networks, volumes, firewall rules, load balancers, and more.

### networking/
Focused test for networking configuration with egress rules.

## Prerequisites

1. CloudStack API credentials:
   ```bash
   export CLOUDSTACK_API_URL="your-api-url"
   export CLOUDSTACK_API_KEY="your-api-key"
   export CLOUDSTACK_SECRET_KEY="your-secret-key"
   ```

2. Terraform 1.0+

3. Built provider (for local testing):
   ```bash
   make install
   ```

## Running Tests

```bash
cd tests/comprehensive  # or tests/networking
terraform init
terraform plan
terraform apply
```

## Cleanup

```bash
terraform destroy
```
