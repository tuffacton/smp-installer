
locals {
  secrets_from_k8s = flatten([
    for secret in var.secrets_in_k8s : [
      for key, value in data.kubernetes_secret.secrets[secret].data : {
        key = key
        value = value
        # service = secret
      }
    ]
  ])
  secrets_from_k8s_map = {
    for secret in local.secrets_from_k8s : secret.key => {
      # service = secret.service
      value = secret.value
    }
  }
  # secretkey_to_secretname = {
  #   for secret in var.secret_keys : secret => local.secrets_from_k8s_map[secret].service
  # }
  secrets_to_create = {
    for secret in var.secret_keys :
      secret => local.secrets_from_k8s_map[secret].value
  }
}


resource "aws_secretsmanager_secret" "secrets" {
  for_each = local.secrets_to_create
  name     = each.key
  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "secertValues" {
  for_each = local.secrets_to_create
  secret_id = aws_secretsmanager_secret.secrets[each.key].id
  secret_string = each.value
}

resource "kubernetes_manifest" "secret_store" {
  manifest = {
    apiVersion = "external-secrets.io/v1beta1"
    kind       = "SecretStore"
    metadata = {
      name      = "harness-secret-store"
      namespace = var.harness_namespace
    }

    spec = {
      provider = {
        aws = {
          region  = var.region
          service = "SecretsManager"
          auth = {
            jwt = {
              serviceAccountRef = {
                name      = var.service_account
                namespace = var.harness_namespace
              }
            }
          }
        }
      }
    }
  }
}