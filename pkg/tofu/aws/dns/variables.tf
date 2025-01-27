variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "domain" {
  description = "Domain name to use with harness"
  type        = string
  default     = "abc.com"
}

variable "lbdns" {
  description = "Domain of ALB"
  type        = string
  default     = "abc.com"
}