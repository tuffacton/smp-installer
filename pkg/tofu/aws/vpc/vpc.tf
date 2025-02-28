provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

locals {
  existing_azs = slice(data.aws_availability_zones.available.names, 0, 2)
  provided_azs = var.vpc.availability_zones
  public_subnet_cidr = [ for subnet in var.vpc.subnets.public : subnet.cidr if subnet.id != "" ]
  private_subnet_cidr = [ for subnet in var.vpc.subnets.private : subnet.cidr if subnet.id != "" ]
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.16.0"

  name = "harness-smp-vpc"

  cidr = var.vpc.cidr == "" ? "10.0.0.0/16" : var.vpc.cidr
  azs  = length(local.provided_azs) > 0 ? local.provided_azs : local.existing_azs

  private_subnets = length(local.private_subnet_cidr) > 0 ? local.private_subnet_cidr : ["10.0.16.0/20", "10.0.32.0/20"]
  public_subnets  = length(local.public_subnet_cidr) > 0 ? local.public_subnet_cidr : ["10.0.48.0/20", "10.0.64.0/20"]

  enable_nat_gateway   = var.airgap ? false : true
  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  tags = var.tags
}

output "vpc" {
  value = module.vpc.vpc_id
}

output "subnets" {
  value = module.vpc.public_subnets
}

output "private_subnets" {
  value = module.vpc.private_subnets
}

output "availability_zones" {
  value = module.vpc.azs
}