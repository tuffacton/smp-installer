provider "aws" {
  region = var.region
}

module "zones" {
  count = var.zone_id != "" ? 0 : 1
  source  = "terraform-aws-modules/route53/aws//modules/zones"
  version = "~> 3.0"

  zones = {
    "${var.domain}" = {
      comment = "${var.domain} harness"
      vpc = var.private_zone ? {
        vpc_id = var.vpc
      } : null
      tags = {
        env = "production"
      }
    }
  }

  tags = var.tags
}

resource "aws_route53_record" "this" {
  count = var.lbdns != "" ? 1 : 0

  zone_id = var.zone_id != "" ? var.zone_id : module.zones[0].route53_zone_zone_id[var.domain]
  name    = var.domain
  type    = "A"

  alias {
    name                   = var.lbdns
    zone_id                = var.lbzone
    evaluate_target_health = false
  }

  depends_on = [module.zones]
}
