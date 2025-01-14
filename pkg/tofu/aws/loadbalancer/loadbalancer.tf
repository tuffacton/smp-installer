provider "aws" {
  region = var.region
}

resource "aws_eip" "lb_eip" {
  count = var.availability_zones
  domain   = "vpc"
}

output "lb_eip_public_ip" {
  value       = aws_eip.lb_eip[0].public_ip
  description = "The public IP address of the reserved Elastic IP for nginx Load Balancer."
}