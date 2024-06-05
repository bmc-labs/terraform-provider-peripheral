terraform {
  required_providers {
    peripheral = {
      source  = "peripheral-cloud/peripheral"
      version = "0.1.0"
    }
  }
}

# provide the token from the environment by 
# setting the "TF_VAR_token" envvar, or just
# "token" in Terraform Cloud or Enterprise
variable "token" {
  type      = string
  sensitive = true
}

provider "peripheral" {
  endpoint = "http://0.0.0.0:3000"
  token    = var.token
}

resource "peripheral_gitlab_runner" "gitlab_runner" {
  id  = "42"
  url = "https://gitlab.com"
  # this token is just for testing; when setting it
  # in production, use secure secret management
  token        = "glpat-1234567890abcdef"
  description  = "my-runner"
  image        = "alpine:latest"
  tag_list     = "tag1,tag2"
  run_untagged = false
}

output "gitlab_runner" {
  value = peripheral_gitlab_runner.gitlab_runner
}
