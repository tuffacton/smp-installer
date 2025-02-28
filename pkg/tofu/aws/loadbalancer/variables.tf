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

variable "tls_enabled" {
  description = "Enable TLS for the load balancer"
  type = bool
}

variable "existing_certificate_id" {
  description = "Existing certificate ARN to use for the load balancer"
  type = string
}

variable "domain_name" {
  description = "Domain name for the certificate"
  type = string
}

variable "subject_alternative_names" {
  description = "Subject alternative names for the certificate"
  type = list(string)
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

variable "instance_tags" {
  description = "Tags to filter kubernetes cluster nodes"
  type        = map(string)
  default     = {}
}

variable "kubernetes_cluster_security_group" {
  description = "Security group for the kubernetes cluster"
  type        = string
}