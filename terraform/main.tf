terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
  required_version = ">= 1.2.0"
  backend "s3" {
    bucket = "crypto-stable-ratio-terraform-backend"
    key    = "terraform.tfstate"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
  profile = "so"
}

module "backend" {
  source = "./modules/backend"

  BACKEND_BUCKET_NAME = var.BACKEND_BUCKET_NAME
}

module "db" {
  source = "./modules/db"
}

module "api" {
  source = "./modules/api"
}