variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "name of kubernetes cluster"
  type = string
}

variable "override_file" {
  description = "Override file for harness installation"
  type        = string
}

variable "registry" {
  description = "Helm registry to download harness chart"
  type        = string
  default = "https://harness.github.io/helm-charts"
}

variable "namespace" {
  description = "namespace to install harness"
  type        = string
  default = "harness"
}

variable "chart_version" {
  description = "chart version to install harness"
  type        = string
}