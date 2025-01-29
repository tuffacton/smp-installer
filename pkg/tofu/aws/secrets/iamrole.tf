resource "aws_iam_role" "secret_viewer_role" {
  name = "harness-secret-viewer-role"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Version" : "2012-10-17",
        "Statement" : {
          "Effect" : "Allow",
          "Action" : "secretsmanager:GetSecretValue",
          "Resource" : "*"
        }
      },
    ]
  })
}
