// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	runrs "terraform-provider-peripheral/internal/clients"
	"time"

	"github.com/deepmap/oapi-codegen/v2/pkg/securityprovider"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure peripheralProvider satisfies various provider interfaces.
var _ provider.Provider = &peripheralProvider{}
var _ provider.ProviderWithFunctions = &peripheralProvider{}

// peripheralProvider defines the provider implementation.
type peripheralProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// peripheralProviderModel describes the provider data model.
type peripheralProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *peripheralProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "peripheral"
	resp.Version = p.version
}

func (p *peripheralProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "URL for the peripheral API.",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Access token for peripheral.",
				Required:            true,
			},
		},
	}
}

func (p *peripheralProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data peripheralProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "peripheral",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	encodedToken, err := jwtToken.SignedString([]byte(data.Token.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Token Encoding Error",
			fmt.Sprintf("Failed to encode JWT token: %s", err),
		)
		return
	}

	bearerToken, err := securityprovider.NewSecurityProviderBearerToken(encodedToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Auth Setup Error",
			fmt.Sprintf("Failed to set up auth: %s", err),
		)
		return
	}

	client, err := runrs.NewClientWithResponses(
		data.Endpoint.ValueString(),
		runrs.WithRequestEditorFn(bearerToken.Intercept),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Setup Error",
			fmt.Sprintf("Failed to set up client: %s", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *peripheralProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGitLabRunnerResource,
	}
}

func (p *peripheralProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *peripheralProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &peripheralProvider{
			version: version,
		}
	}
}
