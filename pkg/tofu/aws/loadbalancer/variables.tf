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

variable "vpc" {
    description = "vpc to create loadbalancer in"
    type = string
}

variable "subnets" {
  description = "subnet names for availabilty"
  type = list(string)
}

variable "harness_node_port" {
  description = "Node port of ingress controller for harness"
  type = number
}