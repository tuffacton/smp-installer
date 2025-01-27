provider "aws" {
  region = var.region
}

module "zones" {
  source  = "terraform-aws-modules/route53/aws//modules/zones"
  version = "~> 3.0"

  zones = {
    "${var.domain}" = {
      comment = "${var.domain} harness"
      tags = {
        env = "production"
      }
    }
  }

  tags = {
    ManagedBy = "Terraform"
  }
}

module "records" {
  source  = "terraform-aws-modules/route53/aws//modules/records"
  version = "~> 3.0"

  zone_name = keys(module.zones.route53_zone_zone_id)[0]

  records = [
    {
      name    = "harness"
      type    = "A"
      alias   = {
        name    = var.lbdns
        zone_id = module.zones.route53_zone_zone_id
      }
    },
  ]

  depends_on = [module.zones]
}