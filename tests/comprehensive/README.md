# Comprehensive Provider Test Suite

Tests the CloudStack Terraform Provider capabilities and identifies limitations.

## Test Coverage

1. SSH Keypair
2. VPC and Networking
3. Isolated Network with Source NAT
4. IP Address allocation
5. Security Groups and Firewall rules
6. Network ACL
7. VM Deployment (2 VMs with cloud-init, SSH keys, affinity groups, tags)
8. Affinity Groups
9. Volumes
10. Volume Attachments
11. Port Forwarding
12. Load Balancer
13. Static NAT
14. Secondary IP

## Usage

```bash
terraform init
terraform plan
terraform apply
terraform destroy
```

## Configuration

Copy terraform.tfvars.example to terraform.tfvars and configure:

- api_url - CloudStack API endpoint
- api_key - API key
- secret_key - API secret
- zone - Zone name
- network_offering - Network offering name
- service_offering - VM service offering
- template - Template name

## Known Limitations

Some resources have stub implementations or missing update functions. See provider documentation for details.
