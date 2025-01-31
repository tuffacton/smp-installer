provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "kubernetes" {
  name = var.cluster_name
}

data "aws_eks_cluster_auth" "kubernetes" {
  name = var.cluster_name
}


provider "kubernetes" {
  host                   = data.aws_eks_cluster.kubernetes.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.kubernetes.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.kubernetes.token
}

provider "helm" {
  debug = var.debug
  kubernetes {
    host                   = data.aws_eks_cluster.kubernetes.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.kubernetes.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.kubernetes.token
  }
}


resource "helm_release" "harness_smp" {
  name = "harness"

  repository = var.repository
  chart      = "harness"

  namespace = var.namespace

  values = [
    "${file(var.override_file)}"
  ]

  upgrade_install  = true
  create_namespace = true

  version = var.chart_version
  wait    = false
  max_history = 10
}