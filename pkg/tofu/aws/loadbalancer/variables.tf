variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "availability_zones" {
    description = "number of availability zones to create loadbalancer"
    type = number
    default = 2
}