terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.11.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

resource "aws_sqs_queue" "data_queue" {
  name                       = "match-scraper-queue"
  visibility_timeout_seconds = 120

  tags = {
    Project   = "match-scraper"
    CreatedBy = "Terraform"
  }
}