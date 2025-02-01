variable "region" {
  description = "AWS region, can be either commercial (e.g. us-east-1) or Govcloud (e.g. us-gov-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "cluster_version" {
  description = "The Kubernetes version to use for the EKS cluster."
  type        = string
  default     = "1.29"
}

variable "machine_type" {
  description = "Machine type to use for EKS node group"
  type        = string
  default     = "t2.2xlarge"
}

variable "minimum_nodes" {
  description = "minimum number of nodes in EKS cluster"
  type        = number
  default     = 1
}

variable "maximum_nodes" {
  description = "maximum number of nodes in EKS cluster"
  type        = number
  default     = 3
}

variable "initial_nodes" {
  description = "Initial number of nodes in EKS cluster"
  type        = number
  default     = 1
}

variable "harness_node_port" {
  description = "Node port of ingress controller for harness"
  type = number
}

variable "airgap" {
  description = "Airgap installation"
  type = bool
  default = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}