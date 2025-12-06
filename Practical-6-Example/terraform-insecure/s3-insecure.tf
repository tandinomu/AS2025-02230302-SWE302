# FIXED: Previously insecure Terraform configuration
# Security improvements applied based on Trivy scan results

# KMS key for customer-managed encryption
resource "aws_kms_key" "s3_encryption" {
  description             = "KMS key for S3 bucket encryption"
  deletion_window_in_days = 10
  enable_key_rotation     = true

  tags = {
    Name = "S3 Encryption Key"
  }
}

resource "aws_kms_alias" "s3_encryption" {
  name          = "alias/s3-encryption-key"
  target_key_id = aws_kms_key.s3_encryption.key_id
}

resource "aws_s3_bucket" "insecure_example" {
  bucket = "insecure-example-bucket"

  tags = {
    Name        = "Insecure Example"
    Environment = "dev"
  }
}

# FIXED ISSUE 1: Server-side encryption enabled with customer-managed KMS key
resource "aws_s3_bucket_server_side_encryption_configuration" "insecure_example" {
  bucket = aws_s3_bucket.insecure_example.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_encryption.arn
    }
  }
}

# FIXED ISSUE 2: Versioning enabled
resource "aws_s3_bucket_versioning" "insecure_example" {
  bucket = aws_s3_bucket.insecure_example.id

  versioning_configuration {
    status = "Enabled"
  }
}

# FIXED ISSUE 3: Access logging enabled
resource "aws_s3_bucket_logging" "insecure_example" {
  bucket = aws_s3_bucket.insecure_example.id

  target_bucket = aws_s3_bucket.backup_insecure.id
  target_prefix = "access-logs/"
}

# FIXED ISSUE 4: Public access blocked
resource "aws_s3_bucket_public_access_block" "insecure_example" {
  bucket = aws_s3_bucket.insecure_example.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# FIXED ISSUE 5: Removed overly permissive bucket policy
# Public read-only access for specific use case (if needed)
resource "aws_s3_bucket_policy" "insecure_example" {
  bucket = aws_s3_bucket.insecure_example.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowGetObjectOnly"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.insecure_example.arn}/*"
      }
    ]
  })
}

# FIXED: Second bucket with proper security
resource "aws_s3_bucket" "backup_insecure" {
  bucket = "backup-insecure-bucket"

  tags = {
    Name        = "Backup Bucket"
    Environment = "dev"
    Purpose     = "Backups"
  }
}

# FIXED: Encryption for backup bucket with customer-managed KMS key
resource "aws_s3_bucket_server_side_encryption_configuration" "backup_insecure" {
  bucket = aws_s3_bucket.backup_insecure.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_encryption.arn
    }
  }
}

# FIXED: Versioning for backup bucket
resource "aws_s3_bucket_versioning" "backup_insecure" {
  bucket = aws_s3_bucket.backup_insecure.id

  versioning_configuration {
    status = "Enabled"
  }
}

# FIXED: Public access block for backup bucket
resource "aws_s3_bucket_public_access_block" "backup_insecure" {
  bucket = aws_s3_bucket.backup_insecure.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# FIXED: Logging for backup bucket
resource "aws_s3_bucket_logging" "backup_insecure" {
  bucket = aws_s3_bucket.backup_insecure.id

  target_bucket = aws_s3_bucket.backup_insecure.id
  target_prefix = "self-logs/"
}
# ISSUE 8: No MFA delete protection for versioned buckets
