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
  kubernetes {
    host                   = data.aws_eks_cluster.kubernetes.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.kubernetes.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.kubernetes.token
  }
}

