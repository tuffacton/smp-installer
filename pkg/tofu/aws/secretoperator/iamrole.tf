provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "kubernetes" {
  name = var.cluster_name
}

data "aws_eks_cluster_auth" "kubernetes" {
  name = var.cluster_name
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.kubernetes.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.kubernetes.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.kubernetes.token
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.kubernetes.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.kubernetes.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.kubernetes.token
  }
}

resource "aws_iam_role" "secret_viewer_role" {
  name = "harness-secret-viewer-role"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : "sts:AssumeRoleWithWebIdentity",
        "Principal" : {
          "Federated" : var.oidc_provider_arn
        },
        "Condition" : {
          "StringEquals" : {
            "${var.oidc_provider_url}:aud" : [
              "sts.amazonaws.com"
            ]
          }
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_policy" "secret_viewer_policy" {
  name = "harness-secret-viewer-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
        ]
        Resource = "*"
      },
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "secret_viewer_policy_attachment" {
  role       = aws_iam_role.secret_viewer_role.name
  policy_arn = aws_iam_policy.secret_viewer_policy.arn
}


