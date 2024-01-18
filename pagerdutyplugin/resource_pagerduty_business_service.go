package pagerduty

import (
	"context"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helperResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceBusinessService struct {
	client *pagerduty.Client
}

var (
	_ resource.ResourceWithConfigure   = (*resourceBusinessService)(nil)
	_ resource.ResourceWithImportState = (*resourceBusinessService)(nil)
)

func (r *resourceBusinessService) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_business_service"
}

func (r *resourceBusinessService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":     schema.StringAttribute{Required: true},
			"id":       schema.StringAttribute{Computed: true},
			"html_url": schema.StringAttribute{Computed: true},
			"self":     schema.StringAttribute{Computed: true},
			"summary":  schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Managed by Terraform"),
			},
			"type": schema.StringAttribute{
				Optional:           true,
				Computed:           true,
				Default:            stringdefault.StaticString("business_service"),
				DeprecationMessage: "This will change to a computed resource in the next major release.",
				Validators:         []validator.String{stringvalidator.OneOf("business_service")},
			},
			"point_of_contact": schema.StringAttribute{Optional: true},
			"team":             schema.StringAttribute{Optional: true},
		},
	}
}

func (r *resourceBusinessService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceBusinessServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	businessServicePlan := buildPagerdutyBusinessService(ctx, &plan)

	err := helperResource.RetryContext(ctx, 5*time.Minute, func() *helperResource.RetryError {
		log.Printf("[INFO] Creating PagerDuty business service %s", plan.Name)
		bs, err := r.client.CreateBusinessServiceWithContext(ctx, businessServicePlan)
		if err != nil {
			return helperResource.NonRetryableError(err)
		} else if bs != nil {
			businessServicePlan.ID = bs.ID
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling CreateBusinessServiceWithContext",
			err.Error(),
		)
		return
	}

	err = helperResource.RetryContext(ctx, 5*time.Minute, func() *helperResource.RetryError {
		businessService, err := r.client.GetBusinessServiceWithContext(ctx, businessServicePlan.ID)
		if err != nil {
			return helperResource.RetryableError(err)
		}
		plan = buildResourceBusinessServiceModel(businessService)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling GetBusinessServiceWithContext",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceBusinessService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceBusinessServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	log.Printf("[INFO] Reading PagerDuty business service %s", state.ID)
	if resp.Diagnostics.HasError() {
		return
	}
	err := helperResource.RetryContext(ctx, 5*time.Minute, func() *helperResource.RetryError {
		businessService, err := r.client.GetBusinessServiceWithContext(ctx, state.ID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) {
				return helperResource.NonRetryableError(err)
			}
			return helperResource.RetryableError(err)
		}
		state = buildResourceBusinessServiceModel(businessService)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling GetBusinessServiceWithContext",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceBusinessService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceBusinessServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	businessServicePlan := buildPagerdutyBusinessService(ctx, &plan)
	if businessServicePlan.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		businessServicePlan.ID = id
	}

	log.Printf("[DEBUG] poc: %v", businessServicePlan.PointOfContact)
	log.Printf("[DEBUG] point_of_contact: %v", plan.PointOfContact)
	log.Printf("[INFO] Updating PagerDuty business service %s", plan.ID)

	businessService, err := r.client.UpdateBusinessServiceWithContext(ctx, businessServicePlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling UpdateBusinessServiceWithContext",
			err.Error(),
		)
		return
	}
	plan = buildResourceBusinessServiceModel(businessService)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceBusinessService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty business service %s", id.String())

	err := r.client.DeleteBusinessServiceWithContext(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error calling DeleteBusinessServiceWithContext",
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceBusinessService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceBusinessService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceBusinessServiceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Summary        types.String `tfsdk:"summary"`
	Self           types.String `tfsdk:"self"`
	PointOfContact types.String `tfsdk:"point_of_contact"`
	HTMLUrl        types.String `tfsdk:"html_url"`
	Team           types.String `tfsdk:"team"`
}

func buildResourceBusinessServiceModel(src *pagerduty.BusinessService) resourceBusinessServiceModel {
	var model resourceBusinessServiceModel
	model.ID = types.StringValue(src.ID)
	model.Name = types.StringValue(src.Name)
	model.HTMLUrl = types.StringValue(src.HTMLUrl)
	model.Description = types.StringValue(src.Description)
	model.Type = types.StringValue(src.Type)
	model.PointOfContact = types.StringValue(src.PointOfContact)
	model.Summary = types.StringValue(src.Summary)
	model.Self = types.StringValue(src.Self)
	if src.Team != nil {
		model.Team = types.StringValue(src.Team.ID)
	}
	return model
}

func buildPagerdutyBusinessService(ctx context.Context, model *resourceBusinessServiceModel) *pagerduty.BusinessService {
	businessService := pagerduty.BusinessService{
		ID:             model.ID.ValueString(),
		Name:           model.Name.ValueString(),
		Description:    model.Description.ValueString(),
		Type:           model.Type.ValueString(),
		Summary:        model.Summary.ValueString(),
		Self:           model.Self.ValueString(),
		PointOfContact: model.PointOfContact.ValueString(),
		HTMLUrl:        model.HTMLUrl.ValueString(),
		Team:           &pagerduty.BusinessServiceTeam{ID: model.Team.ValueString()},
	}
	return &businessService
}
