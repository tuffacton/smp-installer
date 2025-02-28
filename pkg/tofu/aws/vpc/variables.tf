variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "harness_node_port" {
  description = "Node port of ingress controller for harness"
  type        = number
}

variable "airgap" {
  description = "Airgap installation"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

variable "vpc" {
  description = "vpc to create or use"
  type = object({
    id   = string,
    cidr = string,
    availability_zones = list(string),
    subnets = object({
      public = list(object({
        id = string,
        cidr = string,
      })),
      private = list(object({
        id = string,
        cidr = string,
      }))
    })
  })
}

