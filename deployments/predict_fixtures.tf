resource "null_resource" "predictor_lambda_dependencies" {
  provisioner "local-exec" {
    command = "cd ${path.module}/../cmd/predict_fixtures && make build"
  }

  triggers = {
    build_number = "${timestamp()}"
  }
}

data "archive_file" "predictor_zip" {
  type        = "zip"
  source_file = "${path.module}/../bin/predict_fixtures"
  output_path = "${path.module}/.predict_fixtures.zip"

  depends_on = [
    resource.null_resource.predictor_lambda_dependencies
  ]
}

resource "aws_lambda_function" "predictor_lambda" {
  function_name    = "fixture-predictor"
#   s3_bucket        = aws_s3_bucket.reddit-api-binary-bucket.id
#   s3_key           = aws_s3_object.loader_upload.key
  runtime          = "go1.x"
  handler          = "predict_fixtures"
  role             = aws_iam_role.predictor_lambda_role.arn
  memory_size      = 128
  timeout          = 100
  filename         = data.archive_file.predictor_zip.output_path
  source_code_hash = data.archive_file.predictor_zip.output_base64sha256

  environment  {
    variables = {
#       POST_TABLE_NAME = aws_dynamodb_table.posts_table.name
#       COMMENT_TABLE_NAME = aws_dynamodb_table.comments_table.name
        DB_HOST = aws_db_instance.sports-book-db.address
        DB_NAME = "sports-book"
        DB_PASSWORD = data.aws_ssm_parameter.db_password.value
        DB_USER = aws_db_instance.sports-book-db.username
        DISCORD_WEBHOOK_URL = data.aws_ssm_parameter.discord_url.value
        env = "live"
        ODDS_API_KEY = data.aws_ssm_parameter.odds_api_key.value
    }
  }
}

resource "aws_lambda_event_source_mapping" "predict_fixtures" {
  event_source_arn = aws_sqs_queue.data_queue.arn
  function_name = aws_lambda_function.predictor_lambda.function_name
  batch_size = 10
}

# resource "aws_s3_object" "loader_upload" {
#   bucket = aws_s3_bucket.reddit-api-binary-bucket.id
#   key    = "loader.zip"
#   source = data.archive_file.loader_zip.output_path
#   etag   = data.archive_file.loader_zip.output_base64sha256
# }