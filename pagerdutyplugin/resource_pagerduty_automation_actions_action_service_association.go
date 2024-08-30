package pagerduty

import (
	"context"
	"errors"
	"fmt"
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

type resourceAutomationActionsActionServiceAssociation struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceAutomationActionsActionServiceAssociation)(nil)
	_ resource.ResourceWithImportState = (*resourceAutomationActionsActionServiceAssociation)(nil)
)

func (r *resourceAutomationActionsActionServiceAssociation) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_action_service_association"
}

func (r *resourceAutomationActionsActionServiceAssociation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"action_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"service_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *resourceAutomationActionsActionServiceAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var actionID, serviceID types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("action_id"), &actionID)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("service_id"), &serviceID)...)

	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Creating PagerDuty automation action service association %s:%s", actionID, serviceID)
	id := fmt.Sprintf("%s:%s", actionID.ValueString(), serviceID.ValueString())

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		o := pagerduty.AssociateAutomationActionServiceOptions{
			Service: pagerduty.APIReference{
				ID:   serviceID.ValueString(),
				Type: "service_reference",
			},
		}
		if _, err := r.client.AssociateAutomationActionServiceWithContext(ctx, actionID.ValueString(), o); err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty automation action service association %s", id),
			err.Error(),
		)
		return
	}

	model, err := requestGetAutomationActionsActionServiceAssociation(ctx, r.client, id, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action service association %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsActionServiceAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty automation action service association %s", id)

	state, err := requestGetAutomationActionsActionServiceAssociation(ctx, r.client, id.ValueString(), false)
	if err != nil {
		if errors.Is(err, errAutomationActionServiceNotAssociated) || util.IsNotFoundError(err) {
			log.Printf("[WARN] Removing automation action service association %s: %s", id, err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action service association %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsActionServiceAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceAutomationActionsActionServiceAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty automation action service association %s", id)

	actionID, serviceID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		return
	}

	err = r.client.DisassociateAutomationActionServiceWithContext(ctx, actionID, serviceID)
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty automation action service association %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAutomationActionsActionServiceAssociation) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAutomationActionsActionServiceAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAutomationActionsActionServiceAssociationModel struct {
	ID        types.String `tfsdk:"id"`
	ActionID  types.String `tfsdk:"action_id"`
	ServiceID types.String `tfsdk:"service_id"`
}

var errAutomationActionServiceNotAssociated = errors.New("service is not associated to this action")

func requestGetAutomationActionsActionServiceAssociation(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool) (resourceAutomationActionsActionServiceAssociationModel, error) {
	var model resourceAutomationActionsActionServiceAssociationModel

	actionID, serviceID, err := util.ResourcePagerDutyParseColonCompoundID(id)
	if err != nil {
		return model, err
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := client.GetAutomationActionServiceWithContext(ctx, actionID, serviceID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		if response.Service.ID != serviceID {
			return retry.NonRetryableError(errAutomationActionServiceNotAssociated)
		}
		model.ID = types.StringValue(id)
		model.ActionID = types.StringValue(actionID)
		model.ServiceID = types.StringValue(serviceID)
		return nil
	})

	return model, err
}
