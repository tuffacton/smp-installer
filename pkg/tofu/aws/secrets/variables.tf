variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "name of kubernetes cluster"
  type = string
}

variable "secrets_to_create" {
  description = "secret keys and values to create"
  type = map(string)
}

variable "namespace" {
  description = "namespace to deploy external operator"
  type = string
}