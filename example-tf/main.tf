terraform {
  required_version = "1.0.11"

  required_providers {
    sshtunnel = {
      source = "github.com/maclarensg/sshtunnel"
      version = ">= 0.1.0"
    }
    mysql = {
      source  = "winebarrel/mysql"
      version = "~> 1.10.2"
    }
  }
}

provider "sshtunnel" {
  aws_profile = "<AWS PROFILE>"
  aws_region = "<AWS REGION>"
  jumphost = "<JUMPHOST>"
  jumphost_port = "<Port No. of the Jumphost>"
  target_host = "<TARGET HOST>"
  user = "<USER LOGIN of the JUMPHOST>"
  target_port = "<Port No. of the DB>"
}

data "sshtunnel_port" "listen" {  
}

output "endpoint" {
  value = "localhost:${data.sshtunnel_port.listen.port}"
}

provider "mysql" {    
  endpoint = "localhost:${data.sshtunnel_port.listen.port}"
  username = "<ADMIN USERID>"
  password = "<ADMIN PASSWORD>"
}

resource "mysql_user" "gavin_yap_tst" {
  user               = "gavin_yap_test"
  plaintext_password = "IamGr00t"
  host               = "%"
}
