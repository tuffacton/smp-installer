provider "aws" {
  region = var.region
}

resource "random_string" "suffix" {
  length  = 8
  special = false
}

locals {
  sanitized_version = replace(var.cluster_version, ".", "-")
  cluster_name = "harness-smp-eks-${local.sanitized_version}-${random_string.suffix.result}"
}

data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.31.1"

  cluster_name    = local.cluster_name
  cluster_version = var.cluster_version

  cluster_endpoint_public_access           = true
  enable_cluster_creator_admin_permissions = true

  cluster_addons = {
    aws-ebs-csi-driver = {
      service_account_role_arn = module.irsa-ebs-csi.iam_role_arn
    }
    snapshot-controller = {}
    coredns                = {}
    eks-pod-identity-agent = {}
    kube-proxy             = {}
    vpc-cni                = {}
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  eks_managed_node_group_defaults = {
    ami_type = "AL2_x86_64"
  }

  eks_managed_node_groups = {
    one = {
      name = "harness-smp-node-group"
      
      # Change to an ec2 size you would like to utilize
      instance_types = [var.machine_type]

      min_size     = var.minimum_nodes
      max_size     = var.maximum_nodes
      desired_size = var.initial_nodes
      # ensure this aligns with the var.region
      availability_zones = slice(data.aws_availability_zones.available.names, 0, 2)
    }
  }

  node_security_group_additional_rules = {
    allow-ingress-private-subnet = {
      from_port = var.harness_node_port
      cidr_blocks = module.vpc.private_subnets_cidr_blocks
      to_port = var.harness_node_port
      type = "ingress"
      protocol = "tcp"
    }
    allow-ingress-public-subnet = {
      from_port = var.harness_node_port
      cidr_blocks = module.vpc.public_subnets_cidr_blocks
      to_port = var.harness_node_port
      type = "ingress"
      protocol = "tcp"
    }
  }
}

data "aws_iam_policy" "ebs_csi_policy" {
  arn = contains(["us-gov-west-1", "us-gov-east-1"], var.region) ? "arn:aws-us-gov:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy" : "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
}

module "irsa-ebs-csi" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "5.39.0"

  create_role                   = true
  role_name                     = "AmazonEKSTFEBSCSIRole-${module.eks.cluster_name}"
  provider_url                  = module.eks.oidc_provider
  role_policy_arns              = [data.aws_iam_policy.ebs_csi_policy.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:kube-system:ebs-csi-controller-sa"]
}

# data "aws_eks_cluster_auth" "smp" {
#   name = local.cluster_name
# }

# output "eksout" {
#   value = module.eks
# }

# output authout {
#   value = data.aws_eks_cluster_auth.smp.token
#   sensitive = true
# }

output clustername {
  value = local.cluster_name
}