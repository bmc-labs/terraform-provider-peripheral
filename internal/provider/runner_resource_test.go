// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = "peripheral_gitlab_runner"
const resourceName = "test_runner"
const resourceCoordinate = resourceType + "." + resourceName

func testRunnerResourceConfig(description string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
		  id           = "42"
		  url          = "https://gitlab.com"
		  token        = "glpat-1234567890abcdef"
		  description  = "%s"
		  image        = "alpine:latest"
		  tag_list     = "tag1,tag2"
		  run_untagged = false
		}`,
		resourceType,
		resourceName,
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
						resourceCoordinate,
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"url",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"token",
						"glpat-1234567890abcdef",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"description",
						"my-runner",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"image",
						"alpine:latest",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"tag_list",
						"tag1,tag2",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"run_untagged",
						"false",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceCoordinate,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testRunnerResourceConfig("updated-runner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"url",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"token",
						"glpat-1234567890abcdef",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"description",
						"updated-runner",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"image",
						"alpine:latest",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"tag_list",
						"tag1,tag2",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"run_untagged",
						"false",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
