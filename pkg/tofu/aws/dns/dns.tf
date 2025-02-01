provider "aws" {
  region = var.region
}

module "zones" {
  count   = var.zone_id == "" ? 1 : 0
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

module "records" {
  count   = var.lbdns != "" ? 1 : 0
  source  = "terraform-aws-modules/route53/aws//modules/records"
  version = "~> 3.0"

  zone_name = var.zone_name == "" ? module.zones[0].route53_zone_name[var.domain] : var.zone_name
  zone_id = var.zone_id == "" ? module.zones[0].route53_zone_zone_id[var.domain] : var.zone_id
  private_zone = var.private_zone

  records = [
    {
      name = "harness"
      type = "A"
      alias = {
        name    = var.lbdns
        zone_id = var.lbzone
      }
    },
  ]

  depends_on = [module.zones]
}
