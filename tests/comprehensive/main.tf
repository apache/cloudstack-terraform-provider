# Test Infrastructure for CloudStack Terraform Provider Verification
# This tests all claimed provider capabilities and limitations

provider "cloudstack" {
  # Using environment variables: CLOUDSTACK_API_URL, CLOUDSTACK_API_KEY, CLOUDSTACK_SECRET_KEY
}

# Generate unique test identifier
locals {
  test_prefix = "tf-test"
}

# ============================================
# TEST 1: SSH Keypair (Tests: missing Update function)
# ============================================

resource "cloudstack_ssh_keypair" "test" {
  name       = "${local.test_prefix}-keypair"
  public_key = file("~/.ssh/id_rsa.pub")
}

# ============================================
# TEST 2: VPC and Networking
# ============================================

resource "cloudstack_vpc" "test" {
  name         = "${local.test_prefix}-vpc"
  cidr         = "10.10.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = var.zone
}

# ============================================
# TEST 3: Network (Isolated Network with Source NAT)
# ============================================

resource "cloudstack_network" "test" {
  name             = "${local.test_prefix}-network"
  cidr             = "10.10.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id           = cloudstack_vpc.test.id
  zone = var.zone
}

# ============================================
# TEST 4: IP Address (Tests: missing Update function)
# ============================================

resource "cloudstack_ipaddress" "test" {
  vpc_id = cloudstack_vpc.test.id
  zone = var.zone
}

# ============================================
# TEST 5: Security Group and Firewall Rules
# ============================================

resource "cloudstack_security_group" "test" {
  name        = "${local.test_prefix}-sg"
  description = "Test security group for provider verification"
}

resource "cloudstack_firewall" "test" {
  ip_address_id = cloudstack_ipaddress.test.id

  rule {
    cidr_list = ["0.0.0.0/0"]
    protocol  = "tcp"
    ports     = ["22", "80", "443"]
  }

  rule {
    cidr_list = ["0.0.0.0/0"]
    protocol  = "icmp"
    icmp_code = "-1"
    icmp_type = "-1"
  }

  managed = true # Workaround for firewall rule issues
}

# ============================================
# TEST 6: Network ACL (Tests: missing Update function)
# ============================================

resource "cloudstack_network_acl" "test" {
  name        = "${local.test_prefix}-acl"
  description = "Test ACL for provider verification"
  vpc_id      = cloudstack_vpc.test.id
}

resource "cloudstack_network_acl_rule" "test" {
  acl_id = cloudstack_network_acl.test.id

  rule {
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = 22
    traffic_type = "ingress"
    rule_number  = 100
  }

  rule {
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = 80
    traffic_type = "ingress"
    rule_number  = 101
  }
}

# ============================================
# TEST 7: VM Deployment - Frontend (Web Server)
# ============================================

resource "cloudstack_instance" "frontend" {
  name             = "${local.test_prefix}-frontend"
  service_offering = var.service_offering # Smallest offering
  template         = var.template
  zone = var.zone
  network_id       = cloudstack_network.test.id

  # Small disk override
  root_disk_size = 10

  # SSH keypair
  keypair = cloudstack_ssh_keypair.test.name

  # User data for simple web server
  user_data = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl start nginx
    systemctl enable nginx
    echo "<h1>Frontend VM - Test Infrastructure</h1>" > /var/www/html/index.html
  EOF
  )

  tags = {
    role      = "frontend"
    test_run  = "tf-test"
    component = "web"
  }
}

# ============================================
# TEST 8: VM Deployment - Backend (Database/API)
# ============================================

resource "cloudstack_instance" "backend" {
  name             = "${local.test_prefix}-backend"
  service_offering = var.service_offering # Smallest offering
  template         = var.template
  zone = var.zone
  network_id       = cloudstack_network.test.id

  # Small disk override
  root_disk_size = 10

  # SSH keypair
  keypair = cloudstack_ssh_keypair.test.name

  # User data for simple backend
  user_data = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y python3 python3-pip
    cat > /tmp/simple_api.py << 'PYEOF'
    from http.server import HTTPServer, SimpleHTTPRequestHandler
    import json
    
    class APIHandler(SimpleHTTPRequestHandler):
        def do_GET(self):
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {"status": "ok", "role": "backend", "message": "Backend API running"}
            self.wfile.write(json.dumps(response).encode())
    
    if __name__ == "__main__":
        server = HTTPServer(('', 8080), APIHandler)
        print("Backend API running on port 8080")
        server.serve_forever()
    PYEOF
    nohup python3 /tmp/simple_api.py > /tmp/api.log 2>&1 &
  EOF
  )

  tags = {
    role      = "backend"
    test_run  = "tf-test"
    component = "api"
  }
}

# ============================================
# TEST 9: Affinity Group (Tests: missing Update function)
# DISABLED: Domain admin cannot create affinity groups
# ============================================

# resource "cloudstack_affinity_group" "test" {
#   name = "${local.test_prefix}-affinity"
#   type = "host anti-affinity"
# }

# ============================================
# TEST 10: Volume (Tests: missing Update function)
# ============================================

resource "cloudstack_volume" "test" {
  name             = "${local.test_prefix}-data-volume"
  disk_offering_id = var.disk_offering_id
  zone_id          = var.zone_id
}

# ============================================
# TEST 11: Volume Attachment
# ============================================

resource "cloudstack_attach_volume" "test" {
  volume_id          = cloudstack_volume.test.id
  virtual_machine_id = cloudstack_instance.backend.id
}

# ============================================
# TEST 12: Port Forwarding
# ============================================

resource "cloudstack_port_forward" "frontend_http" {
  ip_address_id = cloudstack_ipaddress.test.id

  forward {
    protocol           = "tcp"
    private_port       = 80
    public_port        = 80
    virtual_machine_id = cloudstack_instance.frontend.id
  }
}

resource "cloudstack_port_forward" "frontend_https" {
  ip_address_id = cloudstack_ipaddress.test.id

  forward {
    protocol           = "tcp"
    private_port       = 443
    public_port        = 443
    virtual_machine_id = cloudstack_instance.frontend.id
  }
}

resource "cloudstack_port_forward" "backend_api" {
  ip_address_id = cloudstack_ipaddress.test.id

  forward {
    protocol           = "tcp"
    private_port       = 8080
    public_port        = 8080
    virtual_machine_id = cloudstack_instance.backend.id
  }
}

resource "cloudstack_port_forward" "ssh" {
  ip_address_id = cloudstack_ipaddress.test.id

  forward {
    protocol           = "tcp"
    private_port       = 22
    public_port        = 22
    virtual_machine_id = cloudstack_instance.frontend.id
  }
}

# ============================================
# TEST 13: Load Balancer Rule
# ============================================

resource "cloudstack_loadbalancer_rule" "test" {
  name          = "${local.test_prefix}-lb"
  description   = "Test load balancer rule"
  ip_address_id = cloudstack_ipaddress.test.id
  algorithm     = "roundrobin"
  private_port  = 80
  public_port   = 8888
  member_ids    = [cloudstack_instance.frontend.id]
}

# ============================================
# TEST 14: Static NAT (Tests: missing Update function)
# ============================================

resource "cloudstack_static_nat" "test" {
  ip_address_id      = cloudstack_ipaddress.test.id
  virtual_machine_id = cloudstack_instance.frontend.id
}

# ============================================
# TEST 15: Secondary IP Address (Tests: missing Update function)
# ============================================

resource "cloudstack_secondary_ipaddress" "test" {
  virtual_machine_id = cloudstack_instance.backend.id
}

# ============================================
# Outputs for Verification
# ============================================

output "frontend_ip" {
  value       = cloudstack_ipaddress.test.ip_address
  description = "Public IP address for accessing frontend"
}

output "frontend_vm_id" {
  value = cloudstack_instance.frontend.id
}

output "backend_vm_id" {
  value = cloudstack_instance.backend.id
}

output "test_urls" {
  value = {
    frontend = "http://${cloudstack_ipaddress.test.ip_address}"
    backend  = "http://${cloudstack_ipaddress.test.ip_address}:8080"
  }
  description = "URLs to verify deployment"
}
