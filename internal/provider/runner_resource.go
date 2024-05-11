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

// RunnerResourceModel describes the resource data model.
type RunnerResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Url         types.String `tfsdk:"url"`
	Token       types.String `tfsdk:"token"`
	Description types.String `tfsdk:"description"`
	Image       types.String `tfsdk:"image"`
	TagList     types.String `tfsdk:"tag_list"`
	RunUntagged types.Bool   `tfsdk:"run_untagged"`
}

// FromRunner converts a Runner to a RunnerResourceModel.
func FromRunner(runner *runrs.Runner) RunnerResourceModel {
	return RunnerResourceModel{
		Id:          types.Int64Value(runner.Id),
		Url:         types.StringValue(runner.Url),
		Token:       types.StringValue(runner.Token),
		Description: types.StringValue(runner.Description),
		Image:       types.StringValue(runner.Image),
		TagList:     types.StringValue(runner.TagList),
		RunUntagged: types.BoolValue(runner.RunUntagged),
	}
}

// ToRunner converts a RunnerResourceModel to a Runner.
func (m *RunnerResourceModel) ToRunner() runrs.Runner {
	return runrs.Runner{
		Id:          m.Id.ValueInt64(),
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
	_ resource.Resource                = &RunnerResource{}
	_ resource.ResourceWithConfigure   = &RunnerResource{}
	_ resource.ResourceWithImportState = &RunnerResource{}
)

// NewRunnerResource creates a new RunnerResource.
func NewRunnerResource() resource.Resource {
	return &RunnerResource{}
}

// RunnerResource defines the resource implementation.
type RunnerResource struct {
	client *runrs.ClientWithResponses
}

func (r *RunnerResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_runner"
}

func (r *RunnerResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runner resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Runner ID as provided by GitLab",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of GitLab instance for Runner",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token for Runner registration",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of Runner",
				Optional:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Docker image for Runner",
				Required:            true,
			},
			"tag_list": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of tags for Runner",
				Optional:            true,
			},
			"run_untagged": schema.BoolAttribute{
				MarkdownDescription: "Allow untagged jobs",
				Optional:            true,
			},
		},
	}
}

func (r *RunnerResource) Configure(
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

func (r *RunnerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data RunnerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.CreateWithResponse(ctx, data.ToRunner())
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
			fmt.Sprintf("Unable to create Runner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromRunner(clientResp.JSON201)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("created Runner with ID %d", data.Id.ValueInt64()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunnerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data RunnerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.ReadWithResponse(ctx, data.Id.ValueInt64())
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
			fmt.Sprintf("Unable to read Runner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromRunner(clientResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("read Runner with ID %d", data.Id.ValueInt64()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunnerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data RunnerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	runner := data.ToRunner()

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
			fmt.Sprintf("Unable to update Runner, got status: %s", clientResp.Status()),
		)
		return
	}

	data = FromRunner(clientResp.JSON200)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("updated Runner with ID %d", runner.Id))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunnerResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data RunnerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.DeleteWithResponse(ctx, data.Id.ValueInt64())
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
			fmt.Sprintf("Unable to delete Runner, got status: %s", clientResp.Status()),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("deleted Runner with ID %d", data.Id.ValueInt64()))
}

func (r *RunnerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
