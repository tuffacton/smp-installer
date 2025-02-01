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

variable "oidc_provider_arn" {
  description = "OIDC provider for EKS cluster"
  type = string
}

variable "oidc_provider_url" {
  description = "OIDC issuer URL for EKS cluster"
  type = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}