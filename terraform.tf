variable "key_name" {}

variable "git_access_token" {
  description = "Access token for the private Git repository"
  type        = string
}

provider "github" {
    token = var.token
}

provider "aws" {
  region = "us-east-1"
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = var.key_name
  public_key = tls_private_key.example.public_key_openssh
}

resource "aws_security_group" "elk_sg" {
  name_prefix = "elk-sg"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 5601
    to_port     = 5601
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 9200
    to_port     = 9200
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 5044
    to_port     = 5044
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name = "ELK Docker Server SG"
  }
}

# Launch ec2 instance with the below configuration
resource "aws_instance" "elk_instance" {
  ami             = "ami-0557a15b87f6559cf"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.generated_key.key_name
  security_groups = [aws_security_group.elk_sg.name]

  tags = {
    Name = "ELK API Server"
  }
}

resource "null_resource" "setup" {

  # ssh into the ec2 instance 
  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = tls_private_key.example.private_key_pem
    host        = aws_instance.elk_instance.public_ip
  }
  
 provisioner "file" {
    source      = "setup_server.sh"
    destination = "/home/ubuntu/setup_server.sh"
  }

  provisioner "file" {
    source      = ".env.sample"
    destination = "/home/ubuntu/.env"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo chmod +x /home/ubuntu/setup_server.sh",
      "/home/ubuntu/setup_server.sh",
      "git clone https://oauth2:var.git_access_token@github.com/PSNAppz/Fold-ELK.git && cd Fold-ELK",
      "docker compose up -d --build"
    ]
  }

  # wait for ec2 to be created
  depends_on = [aws_instance.elk_instance]

}

# print the url of the container
output "container_url" {
  value = join("", ["http://", aws_instance.elk_instance.public_dns])
}
