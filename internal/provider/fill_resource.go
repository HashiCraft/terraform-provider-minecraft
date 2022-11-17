package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = fillResourceType{}
var _ tfsdk.Resource = fillResource{}
var _ tfsdk.ResourceWithImportState = fillResource{}

type fillResourceType struct{}

func (t fillResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A Minecraft block",

		Attributes: map[string]tfsdk.Attribute{
			"material": {
				MarkdownDescription: "The material of the block",
				Required:            true,
				Type:                types.StringType,
			},
			"start": {
				MarkdownDescription: "The start position of the block",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"x": {
						MarkdownDescription: "X coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"y": {
						MarkdownDescription: "Y coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"z": {
						MarkdownDescription: "Z coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
				}),
			},
			"end": {
				MarkdownDescription: "The end position of the block",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"x": {
						MarkdownDescription: "X coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"y": {
						MarkdownDescription: "Y coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"z": {
						MarkdownDescription: "Z coordinate of the block",
						Type:                types.NumberType,
						Required:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
				}),
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the block",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t fillResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return fillResource{
		provider: provider,
	}, diags
}

type fillResourceData struct {
	Id       types.String `tfsdk:"id"`
	Material string       `tfsdk:"material"`
	Start    struct {
		X int `tfsdk:"x"`
		Y int `tfsdk:"y"`
		Z int `tfsdk:"z"`
	} `tfsdk:"start"`
	End struct {
		X int `tfsdk:"x"`
		Y int `tfsdk:"y"`
		Z int `tfsdk:"z"`
	} `tfsdk:"end"`
}

type fillResource struct {
	provider provider
}

func (r fillResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data fillResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.provider.GetClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	err = client.FillBlock(ctx, data.Material, data.Start.X, data.Start.Y, data.Start.Z, data.End.X, data.End.Y, data.End.Z)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create block, got error: %s", err))
		return
	}

	data.Id = types.String{Value: fmt.Sprintf("block-%d-%d-%d", data.Start.X, data.Start.Y, data.Start.Z)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// It is impossible to read a block without an entity, so we do nothing for now.
func (r fillResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data fillResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r fillResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data fillResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.provider.GetClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	err = client.FillBlock(ctx, data.Material, data.Start.X, data.Start.Y, data.Start.Z, data.End.X, data.End.Y, data.End.Z)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update block, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r fillResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data fillResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.provider.GetClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	err = client.FillBlock(ctx, "minecraft:air", data.Start.X, data.Start.Y, data.Start.Z, data.End.X, data.End.Y, data.End.Z)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete block, got error: %s", err))
		return
	}
}

func (r fillResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
