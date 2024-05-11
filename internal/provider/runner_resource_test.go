// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(42),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"peripheral_runner.test",
						"id",
						"42",
					),
					resource.TestCheckResourceAttr(
						"peripheral_runner.test",
						"defaulted",
						"example value when not configured",
					),
					resource.TestCheckResourceAttr("peripheral_runner.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "peripheral_runner.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig(1337),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"peripheral_runner.test",
						"id",
						"1337",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(id int64) string {
	return fmt.Sprintf(`
resource "peripheral_runner" "test" {
  id           = %d
  url          = "https://gitlab.com"
  token        = "glpat-1234567890abcdef"
  description  = "my-runner"
  image        = "alpine:latest"
  tag_list     = "tag1,tag2"
  run_untagged = false
}
`, id)
}
