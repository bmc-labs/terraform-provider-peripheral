// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// REST API Client generation
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
