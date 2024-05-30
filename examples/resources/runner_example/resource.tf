terraform {
  required_providers {
    peripheral = {
      source  = "peripheral-cloud/peripheral"
      version = "0.1.0"
    }
  }
}

provider "peripheral" {
  endpoint = "http://0.0.0.0:3000"
}

resource "peripheral_gitlab_runner" "gitlab_runner" {
  id           = "42"
  url          = "https://gitlab.com"
  token        = "glpat-1234567890abcdef"
  description  = "my-runner"
  image        = "alpine:latest"
  tag_list     = "tag1,tag2"
  run_untagged = false
}

output "gitlab_runner" {
  value = peripheral_gitlab_runner.gitlab_runner
}
