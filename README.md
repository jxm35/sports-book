# Sports-Book
This repo contains a predictive model for footabll matches implemented in golang & python. It leverages historical match data, player statistics, and other relevant factors to generate predictions the morning of each match's kickoff.

### Key Features
* <strong>Web Scraping</strong>: Automated web scrapers collect fixture and result data from reliable sources daily.
* <strong>Predictive Model</strong>: Predictive model (developed by myself) to analyze historical data to forecast match outcomes and identify value bets.
* <strong>Discord Notifications</strong>: Users receive notifications via Discord regarding potential value bets for upcoming matches.
* <strong>Backtesting Platform</strong>: A backtesting platform allows users to test and refine prediction strategies using historical data.

### Architecture
* Web Scrapers: Python scripts running on AWS CloudWatch Events collect fixture and result data, sending it to an AWS SQS queue.
* Main Application: A Go-based application processes incoming data, updates the database, and creates predictions and notifications. This is triggered by new entries to the SQS queue.
* Database: An SQL database stores historical match data, team statistics, and prediction results.
* Notification System: Discord integration sends notifications to users with details on potential value bets.
* Backtesting Platform: A separate module allows for historical strategy testing and refinement. This saves graphs and statistics to a data folder to help visualise results of the tests.

### Core Tools & Technologies
* Go
* Python
* Terraform
* Docker
* MySql
