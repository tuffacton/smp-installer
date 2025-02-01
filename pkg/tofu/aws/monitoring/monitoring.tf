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

resource "helm_release" "prometheus" {
  count = var.prometheus_version != "" ? 1 : 0
  name  = "harness-prometheus"

  repository = var.prometheus_repository
  chart      = var.prometheus_chart_name

  namespace = var.namespace

  values = [
    "${file(var.prometheus_override_file)}"
  ]

  upgrade_install  = true
  create_namespace = true

  version     = var.prometheus_version
  wait        = false
  max_history = 10
}

resource "helm_release" "grafana" {
  count = var.grafana_version != "" ? 1 : 0
  name  = "harness-grafana"

  repository = var.grafana_repository
  chart      = var.grafana_chart_name

  namespace = var.namespace

  values = [
    "${file(var.grafana_override_file)}"
  ]

  upgrade_install  = true
  create_namespace = true

  version     = var.grafana_version
  wait        = false
  max_history = 10
}
