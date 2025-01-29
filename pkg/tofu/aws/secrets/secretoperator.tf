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

resource "helm_release" "harness_smp" {
  name = "harness"

  repository = "https://charts.external-secrets.io"
  chart      = "external-secrets"

  namespace = var.namespace

  upgrade_install  = true
  create_namespace = true

  wait        = false
  max_history = 10
}

resource "kubernetes_service_account" "secretmanager" {
  metadata {
    name      = "harness-external-secret"
    namespace = var.namespace
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.secret_viewer_role.arn
    }
  }
}

resource "kubernetes_manifest" "secret_store" {
  manifest = {
    apiVersion = "external-secrets.io/v1beta1"
    kind       = "SecretStore"
    metadata = {
      name      = "harness-secret-store"
      namespace = var.namespace
    }

    spec = {
      provider = {
        aws = {
          region  = var.region
          service = "SecretsManager"
          auth = {
            jwt = {
              serviceAccountRef = {
                name      = kubernetes_service_account.secretmanager.metadata[0].name
                namespace = kubernetes_service_account.secretmanager.metadata[0].namespace
              }
            }
          }
        }
      }
    }
  }
}
