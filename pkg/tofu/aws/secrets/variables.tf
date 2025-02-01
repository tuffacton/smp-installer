variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "name of kubernetes cluster"
  type = string
}

variable "secrets_in_k8s" {
  description = "secret names to find from kubernetes"
  type = list(string)
}

variable "secret_keys" {
  description = "secret names to create in secrets manager"
  type = list(string)
}

variable "namespace" {
  description = "namespace to deploy external operator"
  type = string
}

variable "harness_namespace" {
  description = "namespace where harness is deployed"
  type = string
}

variable "service_account" {
  description = "service account to use for operator"
  type = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}