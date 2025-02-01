variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
}

variable "cluster_name" {
  description = "name of kubernetes cluster"
  type        = string
}

variable "namespace" {
  description = "namespace to deploy prometheus and grafana"
  type        = string
}

variable "prometheus_repository" {
  description = "Helm repository to download prometheus chart"
  type        = string
}

variable "prometheus_chart_name" {
    description = "name of prometheus chart"
    type        = string
}

variable "prometheus_version" {
  description = "version of prometheus to install"
  type        = string
}

variable "prometheus_override_file" {
    description = "Override file for prometheus installation"
    type        = string
}

variable "grafana_repository" {
  description = "Helm repository to download grafana chart"
  type        = string
}

variable "grafana_chart_name" {
    description = "name of grafana chart"
    type        = string
}

variable "grafana_version" {
  description = "version of grafana to install"
  type        = string
}

variable "grafana_override_file" {
    description = "Override file for grafana installation"
    type        = string
}

variable "debug" {
  description = "value of debug"
  type        = bool
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}