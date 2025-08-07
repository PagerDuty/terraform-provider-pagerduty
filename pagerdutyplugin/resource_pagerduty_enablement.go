package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceEnablement struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceEnablement)(nil)
	_ resource.ResourceWithImportState = (*resourceEnablement)(nil)
)

func (r *resourceEnablement) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_enablement"
}

func (r *resourceEnablement) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"entity_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("service", "event_orchestration"),
				},
			},
			"entity_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"feature": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("aiops"),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

func (r *resourceEnablement) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceEnablementModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enablement := buildEnablement(&model)
	log.Printf("[INFO] Creating PagerDuty enablement %s for %s entity with ID %s with enabled=%t", enablement.Feature, enablement.EntityType, enablement.EntityID, enablement.Enabled)

	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var enablementResult *pagerduty.Enablement
		var err error

		switch enablement.EntityType {
		case "service":
			enablementResult, err = r.client.UpdateServiceEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, enablement.Enabled)
		case "event_orchestration":
			enablementResult, err = r.client.UpdateEventOrchestrationEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, enablement.Enabled)
		default:
			return retry.NonRetryableError(fmt.Errorf("unsupported entity type: %s", enablement.EntityType))
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		// Update the model with the actual enabled state from API response
		model.Enabled = types.BoolValue(enablementResult.Enabled)
		model.ID = flattenEnablementID(enablement.EntityType, enablement.EntityID, enablement.Feature)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceEnablement) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceEnablementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty enablement %s", state.ID)

	isFound := r.requestGetEnablement(ctx, &state, &resp.Diagnostics)
	if !isFound {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceEnablement) requestGetEnablement(ctx context.Context, model *resourceEnablementModel, diags *diag.Diagnostics) bool {
	enablement := buildEnablement(model)

	// First validate that the entity exists
	isFound := r.isFoundEntity(ctx, enablement.EntityType, enablement.EntityID, diags)
	if !isFound {
		return false
	}

	// Then check the actual enablement status
	enablementExists := false
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var enablements []pagerduty.Enablement
		var err error

		switch enablement.EntityType {
		case "service":
			enablements, err = r.client.ListServiceEnablementsWithContext(ctx, enablement.EntityID)
		case "event_orchestration":
			enablements, err = r.client.ListEventOrchestrationEnablementsWithContext(ctx, enablement.EntityID)
		default:
			return retry.NonRetryableError(fmt.Errorf("unsupported entity type: %s", enablement.EntityType))
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}

		// Find the specific feature in the enablements list
		for _, e := range enablements {
			if e.Feature == enablement.Feature {
				enablementExists = true
				// Update the model with the current state
				model.Enabled = types.BoolValue(e.Enabled)
				break
			}
		}
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID),
			err.Error(),
		)
		return false
	}

	return enablementExists
}

func (r *resourceEnablement) isFoundEntity(ctx context.Context, entityType, entityID string, diags *diag.Diagnostics) bool {
	isFound := false

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var err error

		switch entityType {
		case "service":
			opts := pagerduty.GetServiceOptions{}
			_, err = r.client.GetServiceWithContext(ctx, entityID, &opts)
		case "event_orchestration":
			opts := &pagerduty.GetOrchestrationOptions{}
			_, err = r.client.GetOrchestrationWithContext(ctx, entityID, opts)
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}

		isFound = true
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error finding %s entity with ID %s", entityType, entityID),
			err.Error(),
		)
		isFound = false
	}

	return isFound
}

func (r *resourceEnablement) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceEnablementModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enablement := buildEnablement(&model)
	log.Printf("[INFO] Updating PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID)

	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var enablementResult *pagerduty.Enablement
		var err error

		switch enablement.EntityType {
		case "service":
			enablementResult, err = r.client.UpdateServiceEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, enablement.Enabled)
		case "event_orchestration":
			enablementResult, err = r.client.UpdateEventOrchestrationEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, enablement.Enabled)
		default:
			return retry.NonRetryableError(fmt.Errorf("unsupported entity type: %s", enablement.EntityType))
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		// Update the model with the actual enabled state from API response
		model.Enabled = types.BoolValue(enablementResult.Enabled)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceEnablement) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model resourceEnablementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enablement := buildEnablement(&model)
	log.Printf("[INFO] Deleting PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID)

	// Disable the enablement by setting enabled=false
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var err error

		switch enablement.EntityType {
		case "service":
			_, err = r.client.UpdateServiceEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, false)
		case "event_orchestration":
			_, err = r.client.UpdateEventOrchestrationEnablementWithContext(ctx, enablement.EntityID, enablement.Feature, false)
		default:
			return retry.NonRetryableError(fmt.Errorf("unsupported entity type: %s", enablement.EntityType))
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				// If the entity or enablement is not found, consider the deletion successful
				return nil
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error disabling PagerDuty enablement %s for %s entity with ID %s", enablement.Feature, enablement.EntityType, enablement.EntityID),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceEnablement) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceEnablement) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ".")
	if len(ids) != 3 {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_enablement",
			"Expecting an importation ID formed as '<entity_type>.<entity_id>.<feature>'",
		)
		return
	}
	entityType, entityID, feature := ids[0], ids[1], ids[2]

	// Validate entity_type
	if entityType != "service" && entityType != "event_orchestration" {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_enablement",
			fmt.Sprintf("Invalid entity_type '%s'. Must be 'service' or 'event_orchestration'", entityType),
		)
		return
	}

	// Validate feature
	if feature != "aiops" {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_enablement",
			fmt.Sprintf("Invalid feature '%s'. Currently supported features: 'aiops'", feature),
		)
		return
	}

	// Validate that the entity exists
	var diags diag.Diagnostics
	isFound := r.isFoundEntity(ctx, entityType, entityID, &diags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !isFound {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_enablement",
			fmt.Sprintf("%s entity with ID %s not found", entityType, entityID),
		)
		return
	}

	// Validate that the enablement exists for the entity/feature combination
	isFound = r.validateEnablementExists(ctx, entityType, entityID, feature, &diags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !isFound {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_enablement",
			fmt.Sprintf("Enablement %s not found for %s entity with ID %s", feature, entityType, entityID),
		)
		return
	}

	// Get the actual enabled state from the API
	var enabledState bool
	var err error
	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var enablements []pagerduty.Enablement
		var err error

		switch entityType {
		case "service":
			enablements, err = r.client.ListServiceEnablementsWithContext(ctx, entityID)
		case "event_orchestration":
			enablements, err = r.client.ListEventOrchestrationEnablementsWithContext(ctx, entityID)
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		// Find the specific feature in the enablements list
		for _, e := range enablements {
			if e.Feature == feature {
				enabledState = e.Enabled
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading enablement state for %s entity with ID %s", entityType, entityID),
			err.Error(),
		)
		return
	}

	state := flattenEnablement(entityType, entityID, feature, enabledState)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceEnablement) validateEnablementExists(ctx context.Context, entityType, entityID, feature string, diags *diag.Diagnostics) bool {
	enablementExists := false
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var enablements []pagerduty.Enablement
		var err error

		switch entityType {
		case "service":
			enablements, err = r.client.ListServiceEnablementsWithContext(ctx, entityID)
		case "event_orchestration":
			enablements, err = r.client.ListEventOrchestrationEnablementsWithContext(ctx, entityID)
		default:
			return retry.NonRetryableError(fmt.Errorf("unsupported entity type: %s", entityType))
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}

		// Find the specific feature in the enablements list
		for _, e := range enablements {
			if e.Feature == feature {
				enablementExists = true
				break
			}
		}
		return nil
	})

	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error validating enablement %s for %s entity with ID %s", feature, entityType, entityID),
			err.Error(),
		)
		return false
	}

	return enablementExists
}

type resourceEnablementModel struct {
	ID         types.String `tfsdk:"id"`
	EntityID   types.String `tfsdk:"entity_id"`
	EntityType types.String `tfsdk:"entity_type"`
	Feature    types.String `tfsdk:"feature"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

type enablement struct {
	ID         string
	EntityID   string
	EntityType string
	Feature    string
	Enabled    bool
}

func buildEnablement(model *resourceEnablementModel) enablement {
	return enablement{
		ID:         model.ID.ValueString(),
		EntityID:   model.EntityID.ValueString(),
		EntityType: model.EntityType.ValueString(),
		Feature:    model.Feature.ValueString(),
		Enabled:    model.Enabled.ValueBool(),
	}
}

func flattenEnablement(entityType, entityID, feature string, enabled bool) resourceEnablementModel {
	model := resourceEnablementModel{
		ID:         flattenEnablementID(entityType, entityID, feature),
		EntityID:   types.StringValue(entityID),
		EntityType: types.StringValue(entityType),
		Feature:    types.StringValue(feature),
		Enabled:    types.BoolValue(enabled),
	}
	return model
}

func flattenEnablementID(entityType, entityID, feature string) types.String {
	return types.StringValue(fmt.Sprintf("%s.%s.%s", entityType, entityID, feature))
}
