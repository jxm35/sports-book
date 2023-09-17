data "archive_file" "python_lambda_package" {
  type = "zip"
  source_file = "${path.module}/../python/scrape_fixtures/scrape_fixtures.py"
  output_path = "scrape_fixtures.zip"
}

resource "aws_lambda_function" "scrape-fixtures-lambda-func" {
        function_name = "fixture-scraper"
        filename      = "scrape_fixtures.zip"
        source_code_hash = data.archive_file.python_lambda_package.output_base64sha256
        role          = aws_iam_role.scraper_lambda_role.arn
        runtime       = "python3.9"
        handler       = "lambda_function.lambda_handler"
        timeout       = 100

        environment  {
            variables = {
              SQS_QUEUE_URL = aws_sqs_queue.data_queue.id
            }
          }
}

resource "aws_cloudwatch_event_rule" "every-day-8" {
  name                  = "run-lambda-function"
  description           = "Schedule lambda function"
  schedule_expression   = "cron(0 8 * * ? *)"
}

resource "aws_cloudwatch_event_target" "run-fixture-scraper" {
  target_id = "lambda-function-target"
  rule      = aws_cloudwatch_event_rule.every-day-8.name
  arn       = aws_lambda_function.scrape-fixtures-lambda-func.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.scrape-fixtures-lambda-func.function_name
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.every-day-8.arn
}
