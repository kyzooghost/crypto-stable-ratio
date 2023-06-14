// Single versioned S3 bucket

# Disable deletes on bucket
data "aws_iam_policy_document" "db" {
  statement {
    actions = [
      "s3:DeleteBucket",
      "s3:DeleteObject",
      "s3:DeleteObjectVersion",
    ]
    resources = [
      "${aws_s3_bucket.db.arn}",
      "${aws_s3_bucket.db.arn}/*",
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

resource "aws_s3_bucket" "db" {
  bucket = "crypto-stable-ratio-db"
  versioning {
      enabled = true
  }
  lifecycle {
      prevent_destroy = true
  }
}

resource "aws_s3_bucket_public_access_block" "db" {
  bucket = aws_s3_bucket.db.id
}

resource "aws_s3_bucket_policy" "backend" {
  bucket = aws_s3_bucket.db.id
  policy = data.aws_iam_policy_document.db.json
}