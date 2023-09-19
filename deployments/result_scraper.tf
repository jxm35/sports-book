resource "null_resource" "result_install_python_dependencies" {
  provisioner "local-exec" {
    command = "bash ${path.module}/scripts/create_pkg.sh"

    environment = {
      source_code_path = "${path.module}/../python/scrape_results"
      filename = "scrape_results.py"
      function_name = "result-scraper"
      dir_name = "scrape_results_dist_pkg/"
      path_module = path.module
      runtime = "python3.9"
      path_cwd = path.cwd
    }
  }

    triggers = {
#       build_number = "${timestamp()}"
    }
}

data "archive_file" "result_python_lambda_package" {
  depends_on = [null_resource.result_install_python_dependencies]
  source_dir = "${path.cwd}/scrape_results_dist_pkg/"
  output_path = "scrape_results.zip"
  type = "zip"
}

resource "aws_lambda_function" "scrape-results-lambda-func" {
        function_name = "result-scraper"
        filename      = "scrape_results.zip"
        source_code_hash = data.archive_file.result_python_lambda_package.output_base64sha256
        role          = aws_iam_role.scraper_lambda_role.arn
        runtime       = "python3.9"
        handler       = "scrape_results.handle_results"
        timeout       = 30

        environment  {
            variables = {
              SQS_QUEUE_URL = aws_sqs_queue.data_queue.id
            }
          }
}

resource "aws_cloudwatch_event_rule" "every-day-10-55" {
  name                  = "run-lambda-function-results"
  description           = "Schedule lambda function"
  schedule_expression   = "cron(55 22 * * ? *)"
}

resource "aws_cloudwatch_event_target" "run-result-scraper" {
  target_id = "lambda-function-target"
  rule      = aws_cloudwatch_event_rule.every-day-10-55.name
  arn       = aws_lambda_function.scrape-results-lambda-func.arn
}

resource "aws_lambda_permission" "allow_result_cloudwatch" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.scrape-results-lambda-func.function_name
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.every-day-10-55.arn
}
