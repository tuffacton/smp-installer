# resource "kubernetes_manifest" "cluster_external_secret_crd" {
#     manifest = yamldecode(file("cluster-external-secret-crd.yaml"))
# }

# resource "kubernetes_manifest" "cluster_secret_store_crd" {
#     manifest = yamldecode(file("cluster-secret-store-crd.yaml"))
# }

# resource "kubernetes_manifest" "external_secret_crd" {
#     manifest = yamldecode(file("external-secret-crd.yaml"))
# }

# resource "kubernetes_manifest" "secret_store_crd" {
#     manifest = yamldecode(file("secret-store-crd.yaml"))
# }

resource "helm_release" "external_secrets" {
  name = "external-secret"

  repository = "https://charts.external-secrets.io"
  chart      = "external-secrets"
  version    = "0.13.0"

  namespace = var.namespace

  upgrade_install  = true
  create_namespace = true

  wait        = true
  max_history = 10
  skip_crds = false
}


resource "kubernetes_service_account" "secretmanager" {
  metadata {
    name      = "harness-external-secret"
    namespace = var.harness_namespace
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.secret_viewer_role.arn
    }
  }
}

output "service_account" {
  value = kubernetes_service_account.secretmanager.metadata[0].name
}


