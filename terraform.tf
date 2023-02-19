variable "key_name" {}

provider "aws" {
  region = "us-east-2"
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
}

# Launch ec2 instance with the below configuration
resource "aws_instance" "elk_instance" {
  ami             = "ami-00eeedc4036573771"
  instance_type   = "t3.micro"
  key_name        = aws_key_pair.generated_key.key_name
  security_groups = [aws_security_group.elk_sg.name]

  tags = {
    Name = "ELK API Server"
  }
}

# an empty resource block
resource "null_resource" "name" {

  # ssh into the ec2 instance 
  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = tls_private_key.example.private_key_pem
    host        = aws_instance.elk_instance.public_ip
  }

  provisioner "file" {
    source      = "docker-compose.yml"
    destination = "/home/ubuntu/docker-compose.yml"
  }

  # copy the build_docker_image.sh from your computer to the ec2 instance 
  provisioner "file" {
    source      = "build_docker_image.sh"
    destination = "/home/ubuntu/build_docker_image.sh"
  }

  # set permissions and run the build_docker_image.sh file
  provisioner "remote-exec" {
    inline = [
      "sudo chmod +x /home/ubuntu/build_docker_image.sh",
      "/home/ubuntu/build_docker_image.sh"
    ]
  }

  # wait for ec2 to be created
  depends_on = [aws_instance.elk_instance]

}

# print the url of the container
output "container_url" {
  value = join("", ["http://", aws_instance.elk_instance.public_dns])
}
