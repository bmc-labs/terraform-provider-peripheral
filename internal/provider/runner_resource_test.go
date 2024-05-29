// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testRunnerResourceConfig(description string) string {
	return fmt.Sprintf(`
		resource "peripheral_gitlab_runner" "test_runner" {
		  id           = "42"
		  url          = "https://gitlab.com"
		  token        = "glpat-1234567890abcdef"
		  description  = "%s"
		  image        = "alpine:latest"
		  tag_list     = "tag1,tag2"
		  run_untagged = false
		}`,
		description,
	)
}

func TestAccRunnerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testRunnerResourceConfig("my-runner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"url",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"token",
						"glpat-1234567890abcdef",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"description",
						"my-runner",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"image",
						"alpine:latest",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"tag_list",
						"tag1,tag2",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"run_untagged",
						"false",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "peripheral_runner.test_runner",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testRunnerResourceConfig("updated-runner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"url",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"token",
						"glpat-1234567890abcdef",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"description",
						"updated-runner",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"image",
						"alpine:latest",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"tag_list",
						"tag1,tag2",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test_runner",
						"run_untagged",
						"false",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
