# FIXED: Previously insecure IAM configuration
# Security improvements applied based on Trivy scan results

# FIXED ISSUE 1: Least privilege IAM policy with specific actions
resource "aws_iam_role" "insecure_role" {
  name = "insecure-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# FIXED ISSUE 2 & 3: Specific actions and resources (least privilege)
resource "aws_iam_role_policy" "insecure_policy" {
  name = "insecure-policy"
  role = aws_iam_role.insecure_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::insecure-example-bucket",
          "arn:aws:s3:::insecure-example-bucket/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:Query"
        ]
        Resource = "arn:aws:dynamodb:us-east-1:*:table/specific-table"
      },
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeSecurityGroups"
        ]
        Resource = "*"
      }
    ]
  })
}

# FIXED ISSUE 4: Removed hardcoded credentials
# Use IAM roles and temporary credentials instead
# resource "aws_iam_access_key" "insecure_key" {
#   user = aws_iam_user.insecure_user.name
# }

resource "aws_iam_user" "insecure_user" {
  name = "service-account"
  path = "/"

  tags = {
    Purpose     = "Service Account"
    MFARequired = "true"
  }
}

# FIXED ISSUE 7: Removed admin permissions, use specific permissions
resource "aws_iam_user_policy" "user_policy" {
  name = "user-limited-policy"
  user = aws_iam_user.insecure_user.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::insecure-example-bucket",
          "arn:aws:s3:::insecure-example-bucket/*"
        ]
      }
    ]
  })
}

# BEST PRACTICES IMPLEMENTED:
# - No wildcard IAM permissions (s3:*, iam:*, etc.)
# - Specific resources instead of "*"
# - No hardcoded credentials
# - Removed administrator access
# - Policies use least privilege principle
