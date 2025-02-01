provider "tls" {

}

resource "tls_private_key" "tls_self_signed_key" {
  algorithm = "RSA"
  lifecycle {
    create_before_destroy = true
  }
}

locals {
  domain_name = var.domain_name != "" ? var.domain_name : module.alb.dns_name
}

resource "tls_self_signed_cert" "tls_self_signed_cert" {
  private_key_pem = tls_private_key.tls_self_signed_key.private_key_pem

  # Certificate expires after 1 year.
  validity_period_hours = 24 * 30 * 12

  # Generate a new certificate if Terraform is run within three
  # hours of the certificate's expiration time.
  early_renewal_hours = 24 * 30 * 12

  # Reasonable set of uses for a server SSL certificate.
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  dns_names = [local.domain_name]

  subject {
    common_name = local.domain_name
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [tls_private_key.tls_self_signed_key]
}

resource "aws_acm_certificate" "cert" {
  count            = (var.tls_enabled && var.existing_certificate_id == "") ? 1 : 0
  private_key      = tls_private_key.tls_self_signed_key.private_key_pem
  certificate_body = tls_self_signed_cert.tls_self_signed_cert.cert_pem
  lifecycle {
    create_before_destroy = true

  }

  tags = var.tags

  depends_on = [tls_private_key.tls_self_signed_key, tls_self_signed_cert.tls_self_signed_cert]
}
