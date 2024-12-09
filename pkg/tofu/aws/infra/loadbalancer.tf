resource "aws_eip" "lb_eip" {
  count = length(module.vpc.public_subnets)
  domain   = "vpc"
}
