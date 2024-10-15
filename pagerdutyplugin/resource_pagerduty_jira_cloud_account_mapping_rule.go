package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceJiraCloudAccountMappingRule struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceJiraCloudAccountMappingRule)(nil)
	_ resource.ResourceWithImportState = (*resourceJiraCloudAccountMappingRule)(nil)
)

func (r *resourceJiraCloudAccountMappingRule) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_jira_cloud_account_mapping_rule"
}

func (r *resourceJiraCloudAccountMappingRule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthAtMost(100)},
			},
			"account_mapping": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"autocreate_jql_disabled_reason": schema.StringAttribute{
				Computed: true,
			},
			"autocreate_jql_disabled_until": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"config": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"service": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"jira": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"autocreate_jql": schema.StringAttribute{Optional: true},
							"create_issue_on_incident_trigger": schema.BoolAttribute{
								Optional: true,
								Computed: true,
								Default:  booldefault.StaticBool(false),
							},
							"sync_notes_user": schema.StringAttribute{Optional: true},
						},
						Blocks: map[string]schema.Block{
							"custom_fields": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"source_incident_field":   schema.StringAttribute{Optional: true},
										"target_issue_field":      schema.StringAttribute{Required: true},
										"target_issue_field_name": schema.StringAttribute{Required: true},
										"type":                    schema.StringAttribute{Required: true},
										"value":                   schema.StringAttribute{Optional: true},
									},
								},
							},
							"issue_type": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Required: true},
									"name": schema.StringAttribute{Required: true},
								},
							},
							"priorities": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"jira_id":      schema.StringAttribute{Required: true},
										"pagerduty_id": schema.StringAttribute{Required: true},
									},
								},
							},
							"project": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
									},
									"key": schema.StringAttribute{
										Required: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
									},
									"name": schema.StringAttribute{
										Required: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
									},
								},
							},
							"status_mapping": schema.SingleNestedBlock{
								Blocks: map[string]schema.Block{
									"acknowledged": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
												},
											},
											"name": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("id")),
												},
											},
										},
									},
									"resolved": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
												},
											},
											"name": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("id")),
												},
											},
										},
									},
									"triggered": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"id":   schema.StringAttribute{Required: true},
											"name": schema.StringAttribute{Required: true},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *resourceJiraCloudAccountMappingRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceJiraCloudAccountMappingRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountMappingID, plan := buildPagerdutyJiraCloudAccountsMappingRule(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Creating PagerDuty jira cloud account mapping rule %s for account mapping %s", plan.ID, accountMappingID)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateJiraCloudAccountsMappingRule(ctx, accountMappingID, plan)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		plan.ID = response.ID
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty jira cloud account mapping rule %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetJiraCloudAccountsMappingRule(ctx, r.client, accountMappingID, plan.ID, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty jira cloud account mapping rule %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceJiraCloudAccountMappingRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountMappingID, ruleID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}
	log.Printf("[INFO] Reading PagerDuty jira cloud account mapping rule %s for account mapping %s", ruleID, accountMappingID)

	state, err := requestGetJiraCloudAccountsMappingRule(ctx, r.client, accountMappingID, ruleID, false)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty jira cloud account mapping rule %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceJiraCloudAccountMappingRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceJiraCloudAccountMappingRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, plan := buildPagerdutyJiraCloudAccountsMappingRule(ctx, &model, &resp.Diagnostics)
	accountMappingID, ruleID, err := util.ResourcePagerDutyParseColonCompoundID(plan.ID)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}
	plan.ID = ruleID
	log.Printf("[INFO] Updating PagerDuty jira cloud account mapping rule %s for account mapping %s", ruleID, accountMappingID)

	jiraCloudAccountsMappingRule, err := r.client.UpdateJiraCloudAccountsMappingRule(ctx, accountMappingID, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty jira cloud account mapping rule %s for account mapping %s", ruleID, accountMappingID),
			err.Error(),
		)
		return
	}
	model = flattenJiraCloudAccountsMappingRule(jiraCloudAccountsMappingRule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceJiraCloudAccountMappingRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idValue types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountMappingID, ruleID, err := util.ResourcePagerDutyParseColonCompoundID(idValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}
	log.Printf("[INFO] Deleting PagerDuty jira cloud account mapping rule %s for account mapping %s", ruleID, accountMappingID)

	err = r.client.DeleteJiraCloudAccountsMappingRule(ctx, accountMappingID, ruleID)
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty jira cloud account mapping rule %s for account mapping %s", ruleID, accountMappingID),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceJiraCloudAccountMappingRule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceJiraCloudAccountMappingRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, err := util.ResourcePagerDutyParseColonCompoundID(req.ID); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceJiraCloudAccountMappingRuleModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	AccountsMapping             types.String `tfsdk:"account_mapping"`
	Config                      types.Object `tfsdk:"config"`
	AutocreateJQLDisabledReason types.String `tfsdk:"autocreate_jql_disabled_reason"`
	AutocreateJQLDisabledUntil  types.String `tfsdk:"autocreate_jql_disabled_until"`
}

func requestGetJiraCloudAccountsMappingRule(ctx context.Context, client *pagerduty.Client, accountMappingID, ruleID string, retryNotFound bool) (resourceJiraCloudAccountMappingRuleModel, error) {
	var model resourceJiraCloudAccountMappingRuleModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		jiraCloudAccountsMappingRule, err := client.GetJiraCloudAccountsMappingRule(ctx, accountMappingID, ruleID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenJiraCloudAccountsMappingRule(jiraCloudAccountsMappingRule)
		return nil
	})

	return model, err
}

func buildPagerdutyJiraCloudAccountsMappingRule(ctx context.Context, model *resourceJiraCloudAccountMappingRuleModel, diags *diag.Diagnostics) (string, pagerduty.JiraCloudAccountsMappingRule) {
	return model.AccountsMapping.ValueString(), pagerduty.JiraCloudAccountsMappingRule{
		ID:     model.ID.ValueString(),
		Name:   model.Name.ValueString(),
		Config: buildPagerdutyJiraCloudAccountsMappingRuleConfig(ctx, model, diags),
	}
}

func buildPagerdutyJiraCloudAccountsMappingRuleConfig(ctx context.Context, model *resourceJiraCloudAccountMappingRuleModel, diags *diag.Diagnostics) pagerduty.JiraCloudAccountsMappingRuleConfig {
	var target struct {
		Jira    types.Object `tfsdk:"jira"`
		Service types.String `tfsdk:"service"`
	}

	d := model.Config.As(ctx, &target, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return pagerduty.JiraCloudAccountsMappingRuleConfig{}
	}

	return pagerduty.JiraCloudAccountsMappingRuleConfig{
		Jira: buildJiraCloudSettings(ctx, target.Jira, diags),
		Service: pagerduty.APIObject{
			ID:   target.Service.ValueString(),
			Type: "service_reference",
		},
	}
}

func buildJiraCloudSettings(ctx context.Context, obj types.Object, diags *diag.Diagnostics) pagerduty.JiraCloudSettings {
	var target struct {
		AutocreateJQL                types.String `tfsdk:"autocreate_jql"`
		CreateIssueOnIncidentTrigger types.Bool   `tfsdk:"create_issue_on_incident_trigger"`
		CustomFields                 types.List   `tfsdk:"custom_fields"`
		IssueType                    types.Object `tfsdk:"issue_type"`
		Priorities                   types.List   `tfsdk:"priorities"`
		Project                      types.Object `tfsdk:"project"`
		StatusMapping                types.Object `tfsdk:"status_mapping"`
		SyncNotesUser                types.String `tfsdk:"sync_notes_user"`
	}
	d := obj.As(ctx, &target, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return pagerduty.JiraCloudSettings{}
	}

	jira := pagerduty.JiraCloudSettings{
		AutocreateJQL:                target.AutocreateJQL.ValueStringPointer(),
		CreateIssueOnIncidentTrigger: target.CreateIssueOnIncidentTrigger.ValueBool(),
		CustomFields:                 buildPagerdutyJiraCloudCustomFields(ctx, target.CustomFields, diags),
		IssueType:                    buildJiraCloudRef(ctx, target.IssueType, diags),
		Priorities:                   buildJiraCloudPriorities(ctx, target.Priorities, diags),
		Project:                      buildJiraCloudProject(ctx, target.Project, diags),
		StatusMapping:                buildJiraCloudStatusMapping(ctx, target.StatusMapping, diags),
		SyncNotesUser:                nil,
	}
	if !target.SyncNotesUser.IsNull() && !target.SyncNotesUser.IsUnknown() {
		jira.SyncNotesUser = &pagerduty.UserJiraCloud{
			APIObject: pagerduty.APIObject{
				ID:   target.SyncNotesUser.ValueString(),
				Type: "user_reference",
			},
		}
	}

	return jira
}

func buildPagerdutyJiraCloudCustomFields(ctx context.Context, list types.List, diags *diag.Diagnostics) []pagerduty.JiraCloudCustomField {
	var target []struct {
		SourceIncidentField  types.String `tfsdk:"source_incident_field"`
		TargetIssueField     types.String `tfsdk:"target_issue_field"`
		TargetIssueFieldName types.String `tfsdk:"target_issue_field_name"`
		Type                 types.String `tfsdk:"type"`
		Value                types.String `tfsdk:"value"`
	}

	d := list.ElementsAs(ctx, &target, true)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	var customFields []pagerduty.JiraCloudCustomField
	for _, cf := range target {
		field := pagerduty.JiraCloudCustomField{
			SourceIncidentField:  nil,
			TargetIssueField:     cf.TargetIssueField.ValueString(),
			TargetIssueFieldName: cf.TargetIssueFieldName.ValueString(),
			Type:                 cf.Type.ValueString(),
			Value:                cf.Value.ValueString(),
		}
		if !cf.SourceIncidentField.IsNull() && !cf.SourceIncidentField.IsUnknown() {
			field.SourceIncidentField = cf.SourceIncidentField.ValueStringPointer()
		}
		customFields = append(customFields, field)
	}

	return customFields
}

func buildJiraCloudRef(ctx context.Context, obj types.Object, diags *diag.Diagnostics) pagerduty.JiraCloudReference {
	var target struct {
		ID   string `tfsdk:"id"`
		Name string `tfsdk:"name"`
	}

	d := obj.As(ctx, &target, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diags.Append(d...); d.HasError() {
		return pagerduty.JiraCloudReference{}
	}

	return pagerduty.JiraCloudReference{ID: target.ID, Name: target.Name}
}

func buildJiraCloudPriorities(ctx context.Context, list types.List, diags *diag.Diagnostics) []pagerduty.JiraCloudPriority {
	var target []struct {
		JiraID      types.String `tfsdk:"jira_id"`
		PagerDutyID types.String `tfsdk:"pagerduty_id"`
	}

	d := list.ElementsAs(ctx, &target, true)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	priorities := make([]pagerduty.JiraCloudPriority, 0, len(list.Elements()))
	for _, p := range target {
		priorities = append(priorities, pagerduty.JiraCloudPriority{
			JiraID:      p.JiraID.ValueString(),
			PagerDutyID: p.PagerDutyID.ValueString(),
		})
	}

	return priorities
}

func buildJiraCloudProject(ctx context.Context, obj types.Object, diags *diag.Diagnostics) pagerduty.JiraCloudReference {
	var target struct {
		ID   string `tfsdk:"id"`
		Name string `tfsdk:"name"`
		Key  string `tfsdk:"key"`
	}

	d := obj.As(ctx, &target, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diags.Append(d...); d.HasError() {
		return pagerduty.JiraCloudReference{}
	}

	return pagerduty.JiraCloudReference{ID: target.ID, Name: target.Name, Key: target.Key}
}

func buildJiraCloudStatusMapping(ctx context.Context, obj types.Object, diags *diag.Diagnostics) pagerduty.JiraCloudStatusMapping {
	var target struct {
		Acknowledged types.Object `tfsdk:"acknowledged"`
		Resolved     types.Object `tfsdk:"resolved"`
		Triggered    types.Object `tfsdk:"triggered"`
	}

	d := obj.As(ctx, &target, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diags.Append(d...); d.HasError() {
		return pagerduty.JiraCloudStatusMapping{}
	}

	var acknowledged *pagerduty.JiraCloudReference
	if v := buildJiraCloudRef(ctx, target.Acknowledged, diags); v.ID != "" {
		acknowledged = &v
	}

	var resolved *pagerduty.JiraCloudReference
	if v := buildJiraCloudRef(ctx, target.Resolved, diags); v.ID != "" {
		resolved = &v
	}

	var triggered *pagerduty.JiraCloudReference
	if v := buildJiraCloudRef(ctx, target.Triggered, diags); v.ID != "" {
		triggered = &v
	}

	return pagerduty.JiraCloudStatusMapping{
		Acknowledged: acknowledged,
		Resolved:     resolved,
		Triggered:    triggered,
	}
}

func flattenJiraCloudAccountsMappingRule(response *pagerduty.JiraCloudAccountsMappingRule) resourceJiraCloudAccountMappingRuleModel {
	id := types.StringNull()
	if response.AccountsMapping != nil {
		id = types.StringValue(response.AccountsMapping.ID + ":" + response.ID)
	}

	autocreateJQLDisabledReason := types.StringNull()
	if response.AutocreateJqlDisabledReason != "" {
		autocreateJQLDisabledReason = types.StringValue(response.AutocreateJqlDisabledReason)
	}

	autocreateJQLDisabledUntil := types.StringNull()
	if response.AutocreateJqlDisabledUntil != "" {
		autocreateJQLDisabledUntil = types.StringValue(response.AutocreateJqlDisabledUntil)
	}

	model := resourceJiraCloudAccountMappingRuleModel{
		ID:                          id,
		Config:                      flattenJiraCloudAccountsMappingRuleConfig(response),
		Name:                        types.StringValue(response.Name),
		AccountsMapping:             types.StringValue(response.AccountsMapping.ID),
		AutocreateJQLDisabledReason: autocreateJQLDisabledReason,
		AutocreateJQLDisabledUntil:  autocreateJQLDisabledUntil,
	}

	return model
}

func flattenJiraCloudAccountsMappingRuleConfig(response *pagerduty.JiraCloudAccountsMappingRule) types.Object {
	var configJiraObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"autocreate_jql":                   types.StringType,
			"create_issue_on_incident_trigger": types.BoolType,
			"custom_fields":                    types.ListType{ElemType: jiraCloudCustomFieldObjectType},
			"issue_type":                       jiraCloudRefObjectType,
			"priorities":                       types.ListType{ElemType: jiraCloudPriorityObjectType},
			"project": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   types.StringType,
					"name": types.StringType,
					"key":  types.StringType,
				},
			},
			"status_mapping":  jiraCloudStatusMappingObjectType,
			"sync_notes_user": types.StringType,
		},
	}

	var configObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service": types.StringType,
			"jira":    configJiraObjectType,
		},
	}

	autocreateJQL := types.StringNull()
	if response.Config.Jira.AutocreateJQL != nil {
		autocreateJQL = types.StringValue(*response.Config.Jira.AutocreateJQL)
	}

	syncNotesUser := types.StringNull()
	if response.Config.Jira.SyncNotesUser != nil {
		syncNotesUser = types.StringValue(response.Config.Jira.SyncNotesUser.ID)
	}

	return types.ObjectValueMust(configObjectType.AttrTypes, map[string]attr.Value{
		"service": types.StringValue(response.Config.Service.ID),
		"jira": types.ObjectValueMust(configJiraObjectType.AttrTypes, map[string]attr.Value{
			"autocreate_jql":                   autocreateJQL,
			"create_issue_on_incident_trigger": types.BoolValue(response.Config.Jira.CreateIssueOnIncidentTrigger),
			"sync_notes_user":                  syncNotesUser,
			"custom_fields":                    flattenJiraCloudCustomFields(response),
			"issue_type":                       flattenJiraCloudRef(response.Config.Jira.IssueType),
			"priorities":                       flattenJiraCloudPriorities(response),
			"project":                          flattenJiraCloudRefKeyed(response.Config.Jira.Project),
			"status_mapping":                   flattenJiraCloudStatusMapping(response),
		}),
	})
}

func flattenJiraCloudCustomFields(response *pagerduty.JiraCloudAccountsMappingRule) types.List {
	elements := make([]attr.Value, 0)

	for _, cf := range response.Config.Jira.CustomFields {
		values := map[string]attr.Value{
			"source_incident_field":   types.StringNull(),
			"target_issue_field":      types.StringValue(cf.TargetIssueField),
			"target_issue_field_name": types.StringValue(cf.TargetIssueFieldName),
			"type":                    types.StringValue(cf.Type),
			"value":                   types.StringNull(),
		}
		if cf.SourceIncidentField != nil {
			values["source_incident_field"] = types.StringValue(*cf.SourceIncidentField)
		}
		if !util.IsNilFunc(cf.Value) {
			s, ok := cf.Value.(string)
			if !ok {
				buf, _ := json.Marshal(cf.Value)
				s = string(buf)
			}
			values["value"] = types.StringValue(s)
		}
		elements = append(elements, types.ObjectValueMust(jiraCloudCustomFieldObjectType.AttrTypes, values))
	}

	return types.ListValueMust(jiraCloudCustomFieldObjectType, elements)
}

func flattenJiraCloudRef(v pagerduty.JiraCloudReference) types.Object {
	return types.ObjectValueMust(jiraCloudRefObjectType.AttrTypes, map[string]attr.Value{
		"id":   types.StringValue(v.ID),
		"name": types.StringValue(v.Name),
	})
}

func flattenJiraCloudRefKeyed(v pagerduty.JiraCloudReference) types.Object {
	return types.ObjectValueMust(jiraCloudRefKeyedObjectType.AttrTypes, map[string]attr.Value{
		"id":   types.StringValue(v.ID),
		"name": types.StringValue(v.Name),
		"key":  types.StringValue(v.Key),
	})
}

func flattenJiraCloudPriorities(response *pagerduty.JiraCloudAccountsMappingRule) types.List {
	elements := make([]attr.Value, 0, len(response.Config.Jira.Priorities))
	for _, p := range response.Config.Jira.Priorities {
		elements = append(elements, types.ObjectValueMust(jiraCloudPriorityObjectType.AttrTypes, map[string]attr.Value{
			"jira_id":      types.StringValue(p.JiraID),
			"pagerduty_id": types.StringValue(p.PagerDutyID),
		}))
	}
	return types.ListValueMust(jiraCloudPriorityObjectType, elements)
}

func flattenJiraCloudStatusMapping(response *pagerduty.JiraCloudAccountsMappingRule) types.Object {
	acknowledged := types.ObjectNull(jiraCloudRefObjectType.AttrTypes)
	if response.Config.Jira.StatusMapping.Acknowledged != nil {
		acknowledged = flattenJiraCloudRef(*response.Config.Jira.StatusMapping.Acknowledged)
	}

	resolved := types.ObjectNull(jiraCloudRefObjectType.AttrTypes)
	if response.Config.Jira.StatusMapping.Resolved != nil {
		resolved = flattenJiraCloudRef(*response.Config.Jira.StatusMapping.Resolved)
	}

	triggered := types.ObjectNull(jiraCloudRefObjectType.AttrTypes)
	if response.Config.Jira.StatusMapping.Triggered != nil {
		triggered = flattenJiraCloudRef(*response.Config.Jira.StatusMapping.Triggered)
	}

	return types.ObjectValueMust(jiraCloudStatusMappingObjectType.AttrTypes, map[string]attr.Value{
		"acknowledged": acknowledged,
		"resolved":     resolved,
		"triggered":    triggered,
	})
}

var jiraCloudCustomFieldObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"source_incident_field":   types.StringType,
		"target_issue_field":      types.StringType,
		"target_issue_field_name": types.StringType,
		"type":                    types.StringType,
		"value":                   types.StringType,
	},
}

var jiraCloudRefObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	},
}

var jiraCloudRefKeyedObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"key":  types.StringType,
	},
}

var jiraCloudPriorityObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"jira_id":      types.StringType,
		"pagerduty_id": types.StringType,
	},
}

var jiraCloudStatusMappingObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"acknowledged": jiraCloudRefObjectType,
		"resolved":     jiraCloudRefObjectType,
		"triggered":    jiraCloudRefObjectType,
	},
}
