variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "name of kubernetes cluster"
  type = string
}

variable "namespace" {
  description = "namespace to deploy external operator"
  type = string
}

variable "harness_namespace" {
  description = "namespace where harness is deployed"
  type = string
}