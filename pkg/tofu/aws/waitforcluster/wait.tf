provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "kubernetes" {
  name = var.cluster_name
}

data "aws_instances" "kubernetes" {
  instance_tags = var.instance_tags
  instance_state_names = ["running"]
}

output "cluster_status" {
  value = data.aws_eks_cluster.kubernetes.status
}

output "instance_count" {
    value = length(data.aws_instances.kubernetes.ids)
}
