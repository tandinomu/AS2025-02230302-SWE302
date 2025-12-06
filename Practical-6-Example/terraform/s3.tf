# KMS key for customer-managed encryption
resource "aws_kms_key" "s3_encryption" {
  description             = "KMS key for S3 bucket encryption"
  deletion_window_in_days = 10
  enable_key_rotation     = true

  tags = {
    Name        = "S3 Encryption Key"
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_kms_alias" "s3_encryption" {
  name          = "alias/${var.project_name}-s3-encryption-${var.environment}"
  target_key_id = aws_kms_key.s3_encryption.key_id
}

# S3 bucket for deployment (website hosting)
resource "aws_s3_bucket" "deployment" {
  bucket = "${var.project_name}-deployment-${var.environment}"

  tags = {
    Name        = "Deployment Bucket"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Configure deployment bucket for static website hosting
resource "aws_s3_bucket_website_configuration" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "404.html"
  }
}

# Enable encryption for deployment bucket with customer-managed KMS key
resource "aws_s3_bucket_server_side_encryption_configuration" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_encryption.arn
    }
  }
}

# Enable versioning for deployment bucket
resource "aws_s3_bucket_versioning" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Make deployment bucket public (for static website access)
resource "aws_s3_bucket_public_access_block" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# Bucket policy to allow public read access
resource "aws_s3_bucket_policy" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.deployment.arn}/*"
      }
    ]
  })
}

# S3 bucket for access logs
resource "aws_s3_bucket" "logs" {
  bucket = "${var.project_name}-logs-${var.environment}"

  tags = {
    Name        = "Access Logs Bucket"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Enable encryption for logs bucket with customer-managed KMS key
resource "aws_s3_bucket_server_side_encryption_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_encryption.arn
    }
  }
}

# Enable versioning for logs bucket
resource "aws_s3_bucket_versioning" "logs" {
  bucket = aws_s3_bucket.logs.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Block public access for logs bucket (security best practice)
resource "aws_s3_bucket_public_access_block" "logs" {
  bucket = aws_s3_bucket.logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Enable access logging for deployment bucket
resource "aws_s3_bucket_logging" "deployment" {
  bucket = aws_s3_bucket.deployment.id

  target_bucket = aws_s3_bucket.logs.id
  target_prefix = "deployment-logs/"
}
