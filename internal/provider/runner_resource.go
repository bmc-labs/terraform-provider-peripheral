// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	uuidpkg "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	runrs "terraform-provider-peripheral/internal/clients"
)

// GitLabRunnerResourceModel describes the resource data model.
type GitLabRunnerResourceModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Id              types.Int32  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Url             types.String `tfsdk:"url"`
	Token           types.String `tfsdk:"token"`
	TokenObtainedAt types.String `tfsdk:"token_obtained_at"`
	DockerImage     types.String `tfsdk:"docker_image"`
}

// FromGitLabRunner converts a GitLabRunner to a GitLabRunnerResourceModel.
func FromGitLabRunner(runner *runrs.GitLabRunner) GitLabRunnerResourceModel {
	ts, _ := runner.TokenObtainedAt.MarshalText()

	return GitLabRunnerResourceModel{
		Uuid:            types.StringValue(runner.Uuid.String()),
		Id:              types.Int32Value(runner.Id),
		Name:            types.StringValue(*runner.Name),
		Url:             types.StringValue(runner.Url),
		Token:           types.StringValue(runner.Token),
		TokenObtainedAt: types.StringValue(string(ts)),
		DockerImage:     types.StringValue(runner.DockerImage),
	}
}

// ToGitLabRunner converts a GitLabRunnerResourceModel to a GitLabRunner.
func (m *GitLabRunnerResourceModel) ToGitLabRunner() runrs.GitLabRunner {
	var uuid *uuidpkg.UUID
	if uuidVal, err := uuidpkg.Parse(m.Uuid.ValueString()); err == nil {
		uuid = &uuidVal
	}

	name := m.Name.ValueString()

	var tokenObtainedAt *time.Time
	if t, err := time.Parse(time.RFC3339, m.TokenObtainedAt.ValueString()); err == nil {
		tokenObtainedAt = &t
	}

	return runrs.GitLabRunner{
		Uuid:            uuid,
		Id:              m.Id.ValueInt32(),
		Name:            &name,
		Url:             m.Url.ValueString(),
		Token:           m.Token.ValueString(),
		TokenObtainedAt: tokenObtainedAt,
		DockerImage:     m.DockerImage.ValueString(),
	}
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &GitLabRunnerResource{}
	_ resource.ResourceWithConfigure   = &GitLabRunnerResource{}
	_ resource.ResourceWithImportState = &GitLabRunnerResource{}
)

// NewGitLabRunnerResource creates a new GitLabRunnerResource.
func NewGitLabRunnerResource() resource.Resource {
	return &GitLabRunnerResource{}
}

// GitLabRunnerResource defines the resource implementation.
type GitLabRunnerResource struct {
	client *runrs.ClientWithResponses
}

func (r *GitLabRunnerResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_gitlab_runner"
}

func (r *GitLabRunnerResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GitLabRunner resource",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				MarkdownDescription: "UUID of GitLabRunner",
				Computed:            true,
			},
			"id": schema.Int32Attribute{
				MarkdownDescription: "GitLab Runner instance ID as provided by GitLab",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Description of GitLabRunner",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of GitLab instance for GitLabRunner",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`https?://.+/`),
						"URL must start with 'http[s]://' and end with '/'",
					),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token for GitLabRunner registration",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^glrt-\w{20}$`),
						"GitLab Runner Token must start with 'glrt-' "+
							"followed by 20 alphanumeric characters",
					),
				},
			},
			"token_obtained_at": schema.StringAttribute{
				MarkdownDescription: "Time when GitLabRunner token was obtained",
				Computed:            true,
			},
			"docker_image": schema.StringAttribute{
				MarkdownDescription: "Docker image for GitLabRunner",
				Required:            true,
			},
		},
	}
}

func (r *GitLabRunnerResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*runrs.ClientWithResponses)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *runrs.Client, got: %T. Report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	r.client = client
}

func (r *GitLabRunnerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data GitLabRunnerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.CreateWithResponse(ctx, data.ToGitLabRunner())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if err := apiResp.GetError(); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf(
				"Unable to create GitLabRunner: %s (%s)",
				err.Msg,
				apiResp.Status(),
			),
		)
		return
	}

	data = FromGitLabRunner(apiResp.JSON201)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("created GitLabRunner with ID %d", data.Id.ValueInt32()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitLabRunnerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data GitLabRunnerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.ReadWithResponse(ctx, uuidpkg.MustParse(data.Uuid.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if err := apiResp.GetError(); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf(
				"Unable to create GitLabRunner: %s (%s)",
				err.Msg,
				apiResp.Status(),
			),
		)
		return
	}

	data = FromGitLabRunner(apiResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("read GitLabRunner with UUID %s", data.Uuid.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitLabRunnerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data GitLabRunnerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read runner UUID from Terraform state
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("uuid"), &data.Uuid)...)
	if resp.Diagnostics.HasError() {
		return
	}

	runner := data.ToGitLabRunner()

	apiResp, err := r.client.UpdateWithResponse(ctx, *runner.Uuid, runner)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if err := apiResp.GetError(); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf(
				"Unable to update GitLabRunner: %s (%s)",
				err.Msg,
				apiResp.Status(),
			),
		)
		return
	}

	resp.Diagnostics.AddWarning(
		"DEBUG PRINT",
		fmt.Sprintf(
			"%+v",
			apiResp.JSON200,
		),
	)

	data = FromGitLabRunner(apiResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("updated GitLabRunner with UUID %s", data.Uuid.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitLabRunnerResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data GitLabRunnerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.DeleteWithResponse(ctx, uuidpkg.MustParse(data.Uuid.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if err := apiResp.GetError(); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf(
				"Unable to create GitLabRunner: %s (%s)",
				err.Msg,
				apiResp.Status(),
			),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("deleted GitLabRunner with UUID %s", data.Uuid.ValueString()))
}

func (r *GitLabRunnerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
