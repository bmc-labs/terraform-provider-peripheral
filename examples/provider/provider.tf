provider "peripheral" {
  endpoint = "http://0.0.0.0:3000"
  token    = var.peripheral_token
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
