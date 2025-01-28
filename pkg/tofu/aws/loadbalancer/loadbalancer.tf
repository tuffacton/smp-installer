provider "aws" {
  region = var.region
}

data "aws_instances" "kubernetes" {

  # filter {
  #   name   = "instance.group-id"
  #   values = ["harness-smp-node-group"]
  # }
  instance_tags = {
    "Name": "harness-smp-node-group"
  }

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
    # all_https = {
    #   from_port   = 443
    #   to_port     = 443
    #   ip_protocol = "tcp"
    #   description = "HTTPS web traffic"
    #   cidr_ipv4   = "0.0.0.0/0"
    # }
  }
  security_group_egress_rules = {
    all = {
      ip_protocol = "-1"
      cidr_ipv4   = "10.0.0.0/16"
    }
  }

  listeners = {
    # http-https-redirect = {
    #   port     = 80
    #   protocol = "HTTP"
    #   redirect = {
    #     port        = "443"
    #     protocol    = "HTTPS"
    #     status_code = "HTTP_301"
    #   }
    # }
    http = {
      port            = 80
      protocol        = "HTTP"

      forward = {
        target_group_key = "smp-instance"
      }
    }
  }

  target_groups = {
    smp-instance = {
      name_prefix      = "smp"
      protocol         = "HTTP"
      port             = var.harness_node_port
      create_attachment = false
      target_type      = "instance"
      target_id        = "i-0f6d38a07d50d080f"
      health_check = {
          enabled = true
          healthy_threshold = 2
          interval = 30
          matcher = 404
          path = "/health"
        }
      
    }
  }
}

resource "aws_lb_target_group_attachment" "kubernetes" {
  # covert a list of instance objects to a map with instance ID as the key, and an instance
  # object as the value.
  for_each = toset(data.aws_instances.kubernetes.ids)

  target_group_arn = module.alb.target_groups["smp-instance"].arn
  target_id        = each.value
  port             = var.harness_node_port
}

output "dns_name" {
  value = module.alb.dns_name
}