variable "key_name" {}

provider "aws" {
  region = "us-east-2"
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = var.key_name
  public_key = tls_private_key.example.public_key_openssh
}

resource "aws_security_group" "elk_sg" {
  name_prefix = "elk-sg"

  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 5601
    to_port   = 5601
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 9200
    to_port   = 9200
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 5044
    to_port   = 5044
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "elk_instance" {
  ami           = "ami-00eeedc4036573771"
  instance_type = "t2.micro"
  key_name      = aws_key_pair.generated_key.key_name
  security_groups = [aws_security_group.elk_sg.name]

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install docker -y
              sudo service docker start
              sudo usermod -a -G docker ec2-user
              sudo docker run -d -p 5601:5601 -p 9200:9200 -p 5044:5044 --name elk sebp/elk:781
              sudo docker run -d -p 80:8080 --name api -e ELK_HOST=http://localhost:9200 abhishek/webapp-golang-elk
              EOF

  tags = {
    Name = "ELK-with-API"
  }
}
