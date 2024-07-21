// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = "peripheral_gitlab_runner"
const resourceName = "test_runner"
const resourceCoordinate = resourceType + "." + resourceName

const initialRunnerName = "initial-runner"
const updatedRunnerName = "updated-runner"

func testRunnerResourceConfig(runnerName string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
		  id           = 42
		  name         = "%s"
		  url          = "https://gitlab.com/"
		  token        = "glrt-0123456789-abcdefXYZ"
		  docker_image = "alpine:latest"
		}`,
		resourceType,
		resourceName,
		runnerName,
	)
}

func TestAccRunnerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testRunnerResourceConfig(initialRunnerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"name",
						initialRunnerName,
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"url",
						"https://gitlab.com/",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"token",
						"glrt-0123456789-abcdefXYZ",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"docker_image",
						"alpine:latest",
					),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceCoordinate,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					resource := s.RootModule().Resources[resourceCoordinate]
					return resource.Primary.Attributes["uuid"], nil
				},
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testRunnerResourceConfig(updatedRunnerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"name",
						updatedRunnerName,
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"url",
						"https://gitlab.com/",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"token",
						"glrt-0123456789-abcdefXYZ",
					),
					resource.TestCheckResourceAttr(
						resourceCoordinate,
						"docker_image",
						"alpine:latest",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
