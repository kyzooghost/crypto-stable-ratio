# MANUALLY CREATE THIS BUCKET, BEFORE STARTING ANY TERRAFORM WORKFLOW - Terraform seems to hang for 20 minutes when creating buckets
# Then we can modify the bucket using Terraform
resource "aws_s3_bucket" "backend" {
    bucket = var.BACKEND_BUCKET_NAME
    versioning {
        enabled = true
    }
    lifecycle {
        prevent_destroy = true
    }
}

# Disable deletes on backend bucket
data "aws_iam_policy_document" "backend" {
  statement {
    actions = [
      "s3:DeleteBucket",
      "s3:DeleteObject",
      "s3:DeleteObjectVersion",
    ]
    resources = [
      "${aws_s3_bucket.backend.arn}",
      "${aws_s3_bucket.backend.arn}/*",
    ]
    principals {
      type = "AWS"
      identifiers = [
        "*",
      ]
    }
    effect = "Deny"
  }
}

resource "aws_s3_bucket_policy" "backend" {
  bucket = aws_s3_bucket.backend.id
  policy = data.aws_iam_policy_document.backend.json
}

# TODO Add bucket policy to enable only the CICD pipeline (the CodeBuild instance?) to edit the bucket