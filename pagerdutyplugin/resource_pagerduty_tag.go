package pagerduty

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceTag struct {
	client *pagerduty.Client
}

var (
	_ resource.ResourceWithConfigure   = (*resourceTag)(nil)
	_ resource.ResourceWithImportState = (*resourceTag)(nil)
)

func (r *resourceTag) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceTag) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_tag"
}

func (r *resourceTag) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{Computed: true},
			"id":       schema.StringAttribute{Computed: true},
			"summary":  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *resourceTag) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceTagModel
	if d := req.Config.Get(ctx, &model); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	tagBody := buildTag(&model)
	log.Printf("[INFO] Creating PagerDuty tag %s", tagBody.Label)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		tag, err := r.client.CreateTagWithContext(ctx, tagBody)
		if err != nil {
			var apiErr pagerduty.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenTag(tag)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error calling CreateTagWithContext", err.Error())
	}
	resp.State.Set(ctx, &model)
}

func (r *resourceTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tagID types.String
	if d := req.State.GetAttribute(ctx, path.Root("id"), &tagID); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	log.Printf("[INFO] Reading PagerDuty tag %s", tagID)

	var model resourceTagModel
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		tag, err := r.client.GetTagWithContext(ctx, tagID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				log.Printf("[WARN] Removing %s because it's gone", tagID.String())
				resp.State.RemoveResource(ctx)
				return nil
			}
			return retry.RetryableError(err)
		}
		model = flattenTag(tag)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error calling GetTagWithContext", err.Error())
	}
	resp.State.Set(ctx, &model)
}

func (r *resourceTag) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *resourceTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model resourceTagModel
	if d := req.State.Get(ctx, &model); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	log.Printf("[INFO] Removing PagerDuty tag %s", model.ID)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.DeleteTagWithContext(ctx, model.ID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				resp.State.RemoveResource(ctx)
				return nil
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error calling DeleteTagWithContext", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceTag) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceTagModel struct {
	ID      types.String `tfsdk:"id"`
	HTMLURL types.String `tfsdk:"html_url"`
	Label   types.String `tfsdk:"label"`
	Summary types.String `tfsdk:"summary"`
}

func buildTag(model *resourceTagModel) *pagerduty.Tag {
	tag := &pagerduty.Tag{
		Label: model.Label.ValueString(),
	}
	tag.Type = "tag"
	return tag
}

func flattenTag(tag *pagerduty.Tag) resourceTagModel {
	model := resourceTagModel{
		ID:      types.StringValue(tag.ID),
		HTMLURL: types.StringValue(tag.HTMLURL),
		Label:   types.StringValue(tag.Label),
		Summary: types.StringValue(tag.Summary),
	}
	return model
}
