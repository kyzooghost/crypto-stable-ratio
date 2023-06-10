// Single versioned S3 bucket

resource "aws_dynamodb_table" "db" {
  name           = "Crypto-Stable-Ratio-Table"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "timestamp"

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "crypto-stable-ratio"
    type = "N"
  }
}