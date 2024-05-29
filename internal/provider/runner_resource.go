// Copyright (c) bmc::labs GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	runrs "terraform-provider-peripheral/internal/clients"
)

// GitLabRunnerResourceModel describes the resource data model.
type GitLabRunnerResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Url         types.String `tfsdk:"url"`
	Token       types.String `tfsdk:"token"`
	Description types.String `tfsdk:"description"`
	Image       types.String `tfsdk:"image"`
	TagList     types.String `tfsdk:"tag_list"`
	RunUntagged types.Bool   `tfsdk:"run_untagged"`
}

// FromGitLabRunner converts a GitLabRunner to a GitLabRunnerResourceModel.
func FromGitLabRunner(runner *runrs.GitLabRunner) GitLabRunnerResourceModel {
	return GitLabRunnerResourceModel{
		Id:          types.StringValue(runner.Id),
		Url:         types.StringValue(runner.Url),
		Token:       types.StringValue(runner.Token),
		Description: types.StringValue(runner.Description),
		Image:       types.StringValue(runner.Image),
		TagList:     types.StringValue(runner.TagList),
		RunUntagged: types.BoolValue(runner.RunUntagged),
	}
}

// ToGitLabRunner converts a GitLabRunnerResourceModel to a GitLabRunner.
func (m *GitLabRunnerResourceModel) ToGitLabRunner() runrs.GitLabRunner {
	return runrs.GitLabRunner{
		Id:          m.Id.ValueString(),
		Url:         m.Url.ValueString(),
		Token:       m.Token.ValueString(),
		Description: m.Description.ValueString(),
		Image:       m.Image.ValueString(),
		TagList:     m.TagList.ValueString(),
		RunUntagged: m.RunUntagged.ValueBool(),
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
			"id": schema.StringAttribute{
				MarkdownDescription: "GitLabRunner ID as provided by GitLab",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of GitLab instance for GitLabRunner",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token for GitLabRunner registration",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of GitLabRunner",
				Optional:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Docker image for GitLabRunner",
				Required:            true,
			},
			"tag_list": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of tags for GitLabRunner",
				Optional:            true,
			},
			"run_untagged": schema.BoolAttribute{
				MarkdownDescription: "Allow untagged jobs",
				Optional:            true,
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

	clientResp, err := r.client.CreateWithResponse(ctx, data.ToGitLabRunner())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if clientResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create GitLabRunner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromGitLabRunner(clientResp.JSON201)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("created GitLabRunner with ID %s", data.Id.ValueString()))

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

	clientResp, err := r.client.ReadWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read GitLabRunner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromGitLabRunner(clientResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("read GitLabRunner with ID %s", data.Id.ValueString()))

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

	runner := data.ToGitLabRunner()

	clientResp, err := r.client.UpdateWithResponse(ctx, runner.Id, runner)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update GitLabRunner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromGitLabRunner(clientResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("updated GitLabRunner with ID %s", runner.Id))

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

	clientResp, err := r.client.DeleteWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to talk to client, got error: %s", err),
		)
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete GitLabRunner, got status: %s", clientResp.Status()),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("deleted GitLabRunner with ID %s", data.Id.ValueString()))
}

func (r *GitLabRunnerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
