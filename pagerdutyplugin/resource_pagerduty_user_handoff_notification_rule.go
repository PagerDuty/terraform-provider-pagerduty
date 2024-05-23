package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helperResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceUserHandoffNotificationRule struct {
	client *pagerduty.Client
}

var (
	_ resource.ResourceWithConfigure   = (*resourceUserHandoffNotificationRule)(nil)
	_ resource.ResourceWithImportState = (*resourceUserHandoffNotificationRule)(nil)
)

func (r *resourceUserHandoffNotificationRule) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_user_handoff_notification_rule"
}

func (r *resourceUserHandoffNotificationRule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	contactMethodBlockObject := schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of contact method to notify for. Possible values are 'email_contact_method', 'email_contact_method_reference', 'phone_contact_method', 'phone_contact_method_reference', 'push_notification_contact_method', 'push_notification_contact_method_reference', 'sms_contact_method', 'sms_contact_method_reference'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"email_contact_method",
						"email_contact_method_reference",
						"phone_contact_method",
						"phone_contact_method_reference",
						"push_notification_contact_method",
						"push_notification_contact_method_reference",
						"sms_contact_method",
						"sms_contact_method_reference",
					),
				},
			},
		},
	}

	contactMethodBlock := schema.ListNestedBlock{
		NestedObject: contactMethodBlockObject,
		Description:  "The contact method to notify for the user handoff notification rule.",
		Validators: []validator.List{
			listvalidator.IsRequired(),
			listvalidator.SizeBetween(1, 1),
		},
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user_id": schema.StringAttribute{
				Required: true,
			},
			"handoff_type": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The type of handoff to notify for. Possible values are 'both', 'oncall', 'offcall'.",
				Validators: []validator.String{
					stringvalidator.OneOf("both", "oncall", "offcall"),
				},
			},
			"notify_advance_in_minutes": schema.Int64Attribute{
				Required:    true,
				Description: "The number of minutes before the handoff to notify the user. Must be greater than or equal to 0.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
		},

		Blocks: map[string]schema.Block{
			"contact_method": contactMethodBlock,
		},
	}
}

func (r *resourceUserHandoffNotificationRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceUserHandoffNotificationRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	userHandoffNotificationRule, diags := buildPagerdutyUserHandoffNotificationRule(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	log.Printf("[INFO] Creating PagerDuty User Handoff Notification Rule %s", plan.ID)

	retryErr := helperResource.RetryContext(ctx, 2*time.Minute, func() *helperResource.RetryError {
		rule, err := r.client.CreateUserOncallHandoffNotificationRuleWithContext(ctx, plan.UserID.ValueString(), *userHandoffNotificationRule)
		if util.IsNotFoundError(err) {
			return helperResource.RetryableError(err)
		}
		if err != nil {
			return helperResource.NonRetryableError(err)
		} else if rule != nil {
			userHandoffNotificationRule.ID = rule.ID
		}
		return nil
	})
	if retryErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating User Handoff Notification Rule %s", plan.ID),
			retryErr.Error(),
		)
		return
	}

	plan = requestGetUserHandoffNotificationRule(ctx, r.client, plan.UserID.ValueString(), userHandoffNotificationRule.ID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceUserHandoffNotificationRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceUserHandoffNotificationRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty User Handoff Notification Rule %s", state.ID)

	var diags diag.Diagnostics
	state = requestGetUserHandoffNotificationRule(ctx, r.client, state.UserID.ValueString(), state.ID.ValueString(), &diags)
	if diags.HasError() {
		for _, d := range diags.Errors() {
			if d.Summary() == "resource not found." {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceUserHandoffNotificationRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceUserHandoffNotificationRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userHandoffNotificationRulePayload, diags := buildPagerdutyUserHandoffNotificationRule(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if userHandoffNotificationRulePayload.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		userHandoffNotificationRulePayload.ID = id
	}
	log.Printf("[INFO] Updating PagerDuty User Handoff Notification Rule %s", userHandoffNotificationRulePayload.ID)

	userHandoffNotificationRuleUpdated, err := r.client.UpdateUserOncallHandoffNotificationRuleWithContext(ctx, plan.UserID.ValueString(), *userHandoffNotificationRulePayload)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating User Handoff Notification Rule %s", userHandoffNotificationRulePayload.ID),
			err.Error(),
		)
		return
	}
	model, d := flattenUserHandoffNotificationRule(plan.UserID.ValueString(), userHandoffNotificationRuleUpdated)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *resourceUserHandoffNotificationRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		id     types.String
		userID types.String
	)

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("user_id"), &userID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty User Handoff Notification Rule %s", id.String())

	err := r.client.DeleteUserOncallHandoffNotificationRuleWithContext(ctx, userID.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting User Handoff Notification Rule %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceUserHandoffNotificationRule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceUserHandoffNotificationRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ".")
	if len(ids) != 2 {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_user_handoff_notification_rule",
			"Expecting an importation ID formed as '<user_id>.<user_handoff_notification_rule_id>'",
		)
		return
	}

	userID, ruleID := ids[0], ids[1]

	userHandoffNotificationRule, err := r.client.GetUserOncallHandoffNotificationRuleWithContext(ctx, userID, ruleID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading User Handoff Notification Rule %s", userID),
			err.Error(),
		)
		return
	}

	if userHandoffNotificationRule == nil || userHandoffNotificationRule.ID != ruleID {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Error importing pagerduty_user_handoff_notification_rule", "User Handoff Notification Rule not found")
		return
	}

	state, diags := flattenUserHandoffNotificationRule(userID, userHandoffNotificationRule)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

type resourceUserHandoffNotificationRuleContactMethodModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type resourceUserHandoffNotificationRuleModel struct {
	ID                     types.String `tfsdk:"id"`
	UserID                 types.String `tfsdk:"user_id"`
	HandoffType            types.String `tfsdk:"handoff_type"`
	NotifyAdvanceInMinutes types.Int64  `tfsdk:"notify_advance_in_minutes"`
	ContactMethod          types.List   `tfsdk:"contact_method"`
}

var resourceContactMethodObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
	},
}

func requestGetUserHandoffNotificationRule(ctx context.Context, client *pagerduty.Client, userID, ruleID string, diags *diag.Diagnostics) resourceUserHandoffNotificationRuleModel {
	var userHandoffNotificationRule *pagerduty.OncallHandoffNotificationRule

	retryErr := helperResource.RetryContext(ctx, 2*time.Minute, func() *helperResource.RetryError {
		var err error
		userHandoffNotificationRule, err = client.GetUserOncallHandoffNotificationRuleWithContext(ctx, userID, ruleID)
		if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
			return helperResource.NonRetryableError(err)
		}
		if err != nil {
			return helperResource.RetryableError(err)
		}
		return nil
	})

	var model resourceUserHandoffNotificationRuleModel
	if util.IsNotFoundError(retryErr) {
		log.Printf("User Handoff Notification Rule %s not found. Removing from state", ruleID)
		diags.AddError("resource not found.", "")
		return model
	}
	if retryErr != nil {
		diags.AddError(
			fmt.Sprintf("Error reading User Handoff Notification Rule %s", userID),
			retryErr.Error(),
		)
		return model
	}
	model, d := flattenUserHandoffNotificationRule(userID, userHandoffNotificationRule)
	if d.HasError() {
		diags.Append(d...)
	}

	return model
}

func buildPagerdutyUserHandoffNotificationRule(ctx context.Context, plan *resourceUserHandoffNotificationRuleModel) (*pagerduty.OncallHandoffNotificationRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	var contactMethodPlan []*resourceUserHandoffNotificationRuleContactMethodModel
	if diags = plan.ContactMethod.ElementsAs(ctx, &contactMethodPlan, false); diags.HasError() {
		return nil, diags
	}

	if len(contactMethodPlan) < 1 {
		diags.AddError("contact_method is required", "")
		return nil, diags
	}
	if diags.HasError() {
		return nil, diags
	}

	userHandoffNotificationRule := pagerduty.OncallHandoffNotificationRule{
		ID:                     plan.ID.ValueString(),
		NotifyAdvanceInMinutes: int(plan.NotifyAdvanceInMinutes.ValueInt64()),
		HandoffType:            plan.HandoffType.ValueString(),
		ContactMethod: &pagerduty.ContactMethod{
			ID:   convertContactMethodDependencyType(contactMethodPlan[0].ID.ValueString()),
			Type: contactMethodPlan[0].Type.ValueString(),
		},
	}

	return &userHandoffNotificationRule, diags
}

func flattenUserHandoffNotificationRule(userID string, src *pagerduty.OncallHandoffNotificationRule) (model resourceUserHandoffNotificationRuleModel, diags diag.Diagnostics) {
	model = resourceUserHandoffNotificationRuleModel{
		ID:                     types.StringValue(src.ID),
		UserID:                 types.StringValue(userID),
		HandoffType:            types.StringValue(src.HandoffType),
		NotifyAdvanceInMinutes: types.Int64Value(int64(src.NotifyAdvanceInMinutes)),
	}
	if src.ContactMethod != nil {
		contactMethodRef, d := types.ObjectValue(resourceContactMethodObjectType.AttrTypes, map[string]attr.Value{
			"id":   types.StringValue(src.ContactMethod.ID),
			"type": types.StringValue(src.ContactMethod.Type),
		})

		contactMethodList, d := types.ListValue(resourceContactMethodObjectType, []attr.Value{contactMethodRef})
		if diags.Append(d...); diags.HasError() {
			return model, diags
		}

		model.ContactMethod = contactMethodList
	}
	return model, diags
}

// convertContactMethodDependencyType is needed because the PagerDuty API
// returns without '*_reference' values in the response but uses the other kind
// of values in requests
func convertContactMethodDependencyType(s string) string {
	switch s {
	case "email_contact_method":
		s = "email_contact_method_reference"
	case "phone_contact_method":
		s = "phone_contact_method_reference"
	case "push_notification_contact_method":
		s = "push_notification_contact_method_reference"
	case "sms_contact_method":
		s = "sms_contact_method_reference"
	}
	return s
}
