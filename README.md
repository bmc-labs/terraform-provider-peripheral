<div align="center">

<img src="./assets/peripheral-banner-1024px.jpg" />
<br/>

# `terraform-provider-peripheral`

**Manage GitLab Runners in Docker via Terraform**

[![build badge](https://github.com/bmc-labs/terraform-provider-peripheral/actions/workflows/test.yml/badge.svg)](https://github.com/bmc-labs/terraform-provider-peripheral/actions/workflows/test.yml)
[![docs badge](https://img.shields.io/badge/docs-latest-7B42BC?logo=terraform)](https://registry.terraform.io/providers/bmc-labs/peripheral/latest/docs)

</div>

This is the official Terraform provider for peripheral. The provider allows you to manage Peripheral
resources using Terraform. The provider is currently in alpha and is under active development.

The main resource currently manageable using this provider are GitLab Runners, and it interacts with
[`runrs`](https://github.com/bmc-labs/runrs) to manage them. It therefore enables you to manage
GitLab Runners entirely through Terraform.

## Using the Provider

As a prerequisite, you'll need to be running `runrs` on your server running Docker or inside of your
Docker instance. Then, to use the provider, you need to add it to your Terraform configuration.

Your resulting Terraform code is going to look liket his:

```hcl
provider "peripheral" {
  url = "http://localhost:8080"
}

resource "peripheral_runner" "runner" {
  id = 42
  url = "https://gitlab.com"
  token = "glpat-1234567890abcdef"
  description = "my-runner"
  image = "alpine:latest"
  tag_list = "tag1,tag2"
  run_untagged = false
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your
machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary
in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Support

This is an open source project, so there isn't support per se. If you open an issue in the
repository, we'll try and help you, but no promises.

---

<div align="center">
© Copyright 2024 <b>bmc::labs</b> GmbH. All rights reserved.<br />
<em>solid engineering. sustainable code.</em>
</div>
