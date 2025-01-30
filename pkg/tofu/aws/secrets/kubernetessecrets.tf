data "kubernetes_secret" "secrets" {
  for_each = toset(var.secrets_in_k8s)
  metadata {
    name      = each.key
    namespace = var.harness_namespace
  }
}