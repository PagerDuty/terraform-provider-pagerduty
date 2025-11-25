package pagerduty

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceUserContactMethod struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure      = (*resourceUserContactMethod)(nil)
	_ resource.ResourceWithImportState    = (*resourceUserContactMethod)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceUserContactMethod)(nil)
)

func (r *resourceUserContactMethod) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_user_contact_method"
}

func (r *resourceUserContactMethod) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"user_id":      schema.StringAttribute{Required: true},
			"label":        schema.StringAttribute{Required: true},
			"country_code": schema.Int64Attribute{Optional: true, Computed: true},
			"enabled":      schema.BoolAttribute{Computed: true},
			"blacklisted":  schema.BoolAttribute{Computed: true},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					util.ValidateContactAddress("type", "country_code"),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"email_contact_method",
						"phone_contact_method",
						"push_notification_contact_method",
						"sms_contact_method",
					),
				},
			},
			"send_short_email": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"device_type": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *resourceUserContactMethod) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg resourceUserContactMethodModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.Type.ValueString() == "push_notification_contact_method" {
		if !cfg.DeviceType.IsNull() && !cfg.DeviceType.IsUnknown() {
			if allowed, got := []string{"android", "ios"}, cfg.DeviceType.ValueString(); !slices.Contains(allowed, got) {
				resp.Diagnostics.AddAttributeError(path.Root("device_type"), "Invalid value", fmt.Sprintf("Attribute device_type value must be one of %q, got %q", allowed, got))
				return
			}
		}
	}

	a := cfg.Address.ValueString()
	t := cfg.Type.ValueString()
	var c int64
	if !cfg.CountryCode.IsNull() && !cfg.CountryCode.IsUnknown() {
		c = cfg.CountryCode.ValueInt64()
	}

	if t == "sms_contact_method" || t == "phone_contact_method" {
		maxLength := 40
		if len(a) > maxLength {
			resp.Diagnostics.AddAttributeError(path.Root("address"), "Invalid phone number", "phone numbers may not exceed 40 characters")
			return
		}

		for _, char := range a {
			isAllowedChar := char == ',' || char == '*' || char == '#'
			if _, err := strconv.ParseInt(string(char), 10, 64); err != nil && !isAllowedChar {
				resp.Diagnostics.AddAttributeError(path.Root("address"), "Invalid phone number", "phone numbers may only include digits from 0-9 and the symbols: comma (,), asterisk (*), and pound (#)")
				return
			}
		}

		isMexicoNumber := c == 52
		if t == "sms_contact_method" && isMexicoNumber && strings.HasPrefix(a, "1") {
			resp.Diagnostics.AddAttributeError(path.Root("address"), "Invalid Mexico SMS number", fmt.Sprintf("Mexico-based SMS numbers should be free of area code prefixes, so please remove the leading 1 in the number %q", a))
			return
		}

		isTrunkPrefixNotSupported := map[int64]string{
			33: "0",
			40: "0",
			44: "0",
			45: "0",
			49: "0",
			61: "0",
			66: "0",
			91: "0",
			1:  "1",
		}
		if prefix, ok := isTrunkPrefixNotSupported[c]; ok && strings.HasPrefix(a, prefix) {
			resp.Diagnostics.AddAttributeError(path.Root("address"), "Invalid phone number format", fmt.Sprintf("Trunk prefixes are not supported for following countries and regions: France, Romania, UK, Denmark, Germany, Australia, Thailand, India and North America, so must be formatted for international use without the leading %s", prefix))
			return
		}
	}
}

func (r *resourceUserContactMethod) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceUserContactMethodModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyContactMethod(&model)
	log.Printf("[INFO] Creating PagerDuty user contact method %s", plan.Label)

	response, err := r.client.CreateUserContactMethodWithContext(ctx, plan.UserID, plan.ContactMethod)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty user contact method %s", plan.Label),
			err.Error(),
		)
		return
	}

	model, err = requestGetUserContactMethod(ctx, r.client, plan.UserID, response.ID, true, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty user contact method %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceUserContactMethod) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	var userID types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("user_id"), &userID)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty user contact method %s", id)

	state, err := requestGetUserContactMethod(ctx, r.client, userID.ValueString(), id.ValueString(), false, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty user contact method %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceUserContactMethod) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceUserContactMethodModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyContactMethod(&model)
	if plan.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		plan.ID = id
	}
	log.Printf("[INFO] Updating PagerDuty user contact method %s", plan.ID)

	response, err := r.client.UpdateUserContactMethodWthContext(ctx, plan.UserID, plan.ContactMethod)
	processedResponse, err := r.processUpdateContactMethodResponse(ctx, plan.UserID, plan.ID, &plan.ContactMethod, response, err)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty user contact method %s", plan.ID),
			err.Error(),
		)
		return
	}
	_ = processedResponse

	model, err = requestGetUserContactMethod(ctx, r.client, plan.UserID, plan.ID, true, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty user contact method %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceUserContactMethod) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	var userID types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("user_id"), &userID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty user contact method %s", id)

	err := r.client.DeleteUserContactMethodWithContext(ctx, userID.ValueString(), id.ValueString())
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty user contact method %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceUserContactMethod) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceUserContactMethod) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ":")
	if len(ids) != 2 {
		resp.Diagnostics.AddError(
			"Error importing PagerDuty user contact method",
			"Expecting an ID formed as '<user_id>:<contact_method_id>'",
		)
	}
	uid, id := ids[0], ids[1]

	var d diag.Diagnostics
	d = resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))
	resp.Diagnostics.Append(d...)
	d = resp.State.SetAttribute(ctx, path.Root("user_id"), types.StringValue(uid))
	resp.Diagnostics.Append(d...)
}

type resourceUserContactMethodModel struct {
	ID             types.String `tfsdk:"id"`
	UserID         types.String `tfsdk:"user_id"`
	Address        types.String `tfsdk:"address"`
	Blacklisted    types.Bool   `tfsdk:"blacklisted"`
	CountryCode    types.Int64  `tfsdk:"country_code"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Label          types.String `tfsdk:"label"`
	SendShortEmail types.Bool   `tfsdk:"send_short_email"`
	Type           types.String `tfsdk:"type"`
	DeviceType     types.String `tfsdk:"device_type"`
}

type ContactMethod struct {
	pagerduty.ContactMethod
	UserID string
}

func requestGetUserContactMethod(ctx context.Context, client *pagerduty.Client, userID, id string, retryNotFound bool, diags *diag.Diagnostics) (resourceUserContactMethodModel, error) {
	var model resourceUserContactMethodModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		contactMethod, err := client.GetUserContactMethodWithContext(ctx, userID, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenUserContactMethod(contactMethod, userID)
		return nil
	})

	return model, err
}

func buildPagerdutyContactMethod(model *resourceUserContactMethodModel) ContactMethod {
	contactMethod := pagerduty.ContactMethod{
		Label:          model.Label.ValueString(),
		Address:        model.Address.ValueString(),
		SendShortEmail: model.SendShortEmail.ValueBool(),
		Enabled:        model.Enabled.ValueBool(),
	}
	contactMethod.Type = model.Type.ValueString()

	if !model.CountryCode.IsNull() && !model.CountryCode.IsUnknown() {
		contactMethod.CountryCode = int(model.CountryCode.ValueInt64())
	}
	if !model.DeviceType.IsNull() && !model.DeviceType.IsUnknown() {
		contactMethod.DeviceType = model.DeviceType.ValueString()
	}

	return ContactMethod{ContactMethod: contactMethod, UserID: model.UserID.ValueString()}
}

func (r *resourceUserContactMethod) updateContactMethodCall(ctx context.Context, userID, contactMethodID string, contactMethod *pagerduty.ContactMethod) (*pagerduty.ContactMethod, error) {
	return r.client.UpdateUserContactMethodWthContext(ctx, userID, *contactMethod)
}

func (r *resourceUserContactMethod) processUpdateContactMethodResponse(ctx context.Context, userID, contactMethodID string, contactMethod *pagerduty.ContactMethod, response *pagerduty.ContactMethod, err error) (*pagerduty.ContactMethod, error) {
	if err != nil {
		log.Println("[CG]", err.Error())
		log.Printf("[CG] %#v", err)
		isUniqueContactError := strings.Contains(err.Error(), "User Contact method must be unique")
		if !isUniqueContactError {
			return nil, err
		}

		existingContact, err := r.findExistingContactMethod(ctx, userID, contactMethod)
		if err != nil {
			return nil, err
		}

		err = r.client.DeleteUserContactMethodWithContext(ctx, userID, existingContact.ID)
		if err != nil {
			return nil, err
		}

		_, err = r.updateContactMethodCall(ctx, userID, contactMethodID, contactMethod)
		if err != nil {
			return nil, err
		}

		existingContact, err = r.findExistingContactMethod(ctx, userID, contactMethod)
		if err != nil {
			return nil, err
		}
		return existingContact, nil
	}

	return response, nil
}

func (r *resourceUserContactMethod) findExistingContactMethod(ctx context.Context, userID string, contactMethod *pagerduty.ContactMethod) (*pagerduty.ContactMethod, error) {
	contactMethods, err := r.client.ListUserContactMethodsWithContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[User Contact method must be unique] but failed to fetch existing ones: %w", err)
	}

	for _, contact := range contactMethods.ContactMethods {
		if r.isSameContactMethod(&contact, contactMethod) {
			return r.client.GetUserContactMethodWithContext(ctx, userID, contact.ID)
		}
	}

	return nil, fmt.Errorf("[User Contact method must be unique]")
}

func (r *resourceUserContactMethod) isSameContactMethod(existingContact, newContact *pagerduty.ContactMethod) bool {
	return existingContact.Type == newContact.Type &&
		existingContact.Address == newContact.Address &&
		existingContact.CountryCode == newContact.CountryCode
}

func flattenUserContactMethod(response *pagerduty.ContactMethod, userID string) resourceUserContactMethodModel {
	model := resourceUserContactMethodModel{
		ID:             types.StringValue(response.ID),
		Address:        types.StringValue(response.Address),
		Blacklisted:    types.BoolValue(response.Blacklisted),
		CountryCode:    types.Int64Value(int64(response.CountryCode)),
		Enabled:        types.BoolValue(response.Enabled),
		Label:          types.StringValue(response.Label),
		SendShortEmail: types.BoolValue(response.SendShortEmail),
		Type:           types.StringValue(response.Type),
		UserID:         types.StringValue(userID),
	}

	if response.DeviceType == "" {
		model.DeviceType = types.StringNull()
	} else {
		model.DeviceType = types.StringValue(response.DeviceType)
	}

	return model
}
