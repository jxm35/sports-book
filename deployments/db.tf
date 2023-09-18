#create a security group for RDS Database Instance
resource "aws_security_group" "rds_sg" {
  name = "rds_sg"
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

#create a RDS Database Instance
resource "aws_db_instance" "sports-book-db" {
  engine               = "mysql"
  identifier           = "sports-book"
  allocated_storage    =  10
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  username             = "admin"
  password             = data.aws_ssm_parameter.db_password.value
  parameter_group_name = "default.mysql5.7"
  vpc_security_group_ids = ["${aws_security_group.rds_sg.id}"]
#   deletion_protection = true
  skip_final_snapshot  = true
  publicly_accessible =  true
}