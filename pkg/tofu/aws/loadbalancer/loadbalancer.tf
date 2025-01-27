provider "aws" {
  region = var.region
}

module "alb" {
  source = "terraform-aws-modules/alb/aws"

  name    = "harness-smp-alb"
  vpc_id  = var.vpc
  subnets = var.subnets

  # Security Group
  security_group_ingress_rules = {
    all_http = {
      from_port   = 80
      to_port     = 80
      ip_protocol = "tcp"
      description = "HTTP web traffic"
      cidr_ipv4   = "0.0.0.0/0"
    }
    all_https = {
      from_port   = 443
      to_port     = 443
      ip_protocol = "tcp"
      description = "HTTPS web traffic"
      cidr_ipv4   = "0.0.0.0/0"
    }
  }
  security_group_egress_rules = {
    all = {
      ip_protocol = "-1"
      cidr_ipv4   = "10.0.0.0/16"
    }
  }

  listeners = {
    http-https-redirect = {
      port     = 80
      protocol = "HTTP"
      redirect = {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
    https = {
      port            = 443
      protocol        = "HTTPS"
      certificate_arn = "arn:aws:iam::123456789012:server-certificate/test_cert-123456789012"

      forward = {
        target_group_key = "smp-instance"
      }
    }
  }

  target_groups = {
    smp-instance = {
      name_prefix      = "harness-smp"
      protocol         = "HTTP"
      port             = var.harness_node_port
      target_type      = "instance"
      target_id        = "i-0f6d38a07d50d080f"
    }
  }
}