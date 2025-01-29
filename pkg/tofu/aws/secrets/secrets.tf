provider "aws" {
  region = var.region
}


resource "aws_secretsmanager_secret" "secrets" {
  for_each = var.secrets_to_create
  name     = each.key
}

resource "aws_secretsmanager_secret_version" "secertValues" {
  for_each = var.secrets_to_create
  secret_id = aws_secretsmanager_secret.secrets[each.key].id
  secret_string = each.value
}