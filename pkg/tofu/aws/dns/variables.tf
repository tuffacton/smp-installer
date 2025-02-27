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

variable "vpc" {
  description = "vpc to create loadbalancer in"
  type        = string
}

variable "lbdns" {
  description = "Domain of ALB"
  type        = string
  default     = "abc.com"
}

variable "lbzone" {
  description = "Zone id of ALB"
  type        = string
}

variable "zone_id" {
  description = "Route53 zone id"
  type        = string
}

variable "private_zone" {
  description = "Route53 zone is private"
  type        = bool
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}