data "aws_ssm_parameter" "db_password" {
  name        = "/sports-book/db/password"
}

data "aws_ssm_parameter" "discord_url" {
  name        = "/sports-book/discord-url"
}

data "aws_ssm_parameter" "odds_api_key" {
  name        = "/sports-book/odds-api-key"
}