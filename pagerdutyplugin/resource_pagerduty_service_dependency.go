package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceServiceDependency struct {
	client *pagerduty.Client
}

var (
	_ resource.ResourceWithConfigure   = (*resourceServiceDependency)(nil)
	_ resource.ResourceWithImportState = (*resourceServiceDependency)(nil)
)

func (r *resourceServiceDependency) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_dependency"
}

func (r *resourceServiceDependency) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	supportingServiceBlockObject := schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"business_service",
						"business_service_reference",
						"service",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	dependencyServiceBlockObject := schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"business_service",
						"business_service_reference",
						"service",
						"service_dependency",
						"technical_service_reference",
					),
				},
			},
		},
	}

	dependencyBlockObject := schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{Optional: true, Computed: true},
		},
		Blocks: map[string]schema.Block{
			"supporting_service": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: supportingServiceBlockObject,
			},
			"dependent_service": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: dependencyServiceBlockObject,
			},
		},
	}

	dependencyBlock := schema.ListNestedBlock{
		NestedObject: dependencyBlockObject,
		Validators: []validator.List{
			listvalidator.IsRequired(),
			listvalidator.SizeBetween(1, 1),
		},
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"dependency": dependencyBlock,
		},
	}
}

func (r *resourceServiceDependency) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceServiceDependencyModel

	if diags := req.Plan.Get(ctx, &model); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	serviceDependency, diags := buildServiceDependencyStruct(ctx, model)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	dependencies := &pagerduty.ListServiceDependencies{
		Relationships: []*pagerduty.ServiceDependency{serviceDependency},
	}

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		resourceServiceDependencyMu.Lock()
		list, err := r.client.AssociateServiceDependenciesWithContext(ctx, dependencies)
		resourceServiceDependencyMu.Unlock()
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		model = flattenServiceDependency(list.Relationships, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error associating service dependency", err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceServiceDependency) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model resourceServiceDependencyModel

	if diags := req.State.Get(ctx, &model); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	serviceDependency, diags := buildServiceDependencyStruct(ctx, model)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	log.Printf("Reading PagerDuty dependency %s", serviceDependency.ID)

	serviceDependency, err := r.requestGetServiceDependency(ctx, serviceDependency.ID, serviceDependency.DependentService.ID, serviceDependency.DependentService.Type)
	if serviceDependency == nil || util.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error listing service dependencies", err.Error())
		return
	}

	model = flattenServiceDependency([]*pagerduty.ServiceDependency{serviceDependency}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceServiceDependency) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Update for service dependency has no effect", "")
}

func (r *resourceServiceDependency) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model resourceServiceDependencyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dependencies []*resourceServiceDependencyItemModel
	if d := model.Dependency.ElementsAs(ctx, &dependencies, false); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	var dependents []types.Object
	if d := dependencies[0].DependentService.ElementsAs(ctx, &dependents, false); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	var dependent struct {
		ID   types.String `tfsdk:"id"`
		Type types.String `tfsdk:"type"`
	}
	if d := dependents[0].As(ctx, &dependent, basetypes.ObjectAsOptions{}); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	id := model.ID.ValueString()
	depID := dependent.ID.ValueString()
	rt := dependent.Type.ValueString()

	serviceDependency, err := r.requestGetServiceDependency(ctx, id, depID, rt)
	if serviceDependency == nil || util.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error listing service dependencies", err.Error())
		return
	}

	if serviceDependency.SupportingService != nil {
		serviceDependency.SupportingService.Type = convertServiceDependencyType(serviceDependency.SupportingService.Type)
	}
	if serviceDependency.DependentService != nil {
		serviceDependency.DependentService.Type = convertServiceDependencyType(serviceDependency.DependentService.Type)
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		_, err := r.client.DisassociateServiceDependenciesWithContext(ctx, &pagerduty.ListServiceDependencies{
			Relationships: []*pagerduty.ServiceDependency{serviceDependency},
		})
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty service dependency %s (%s) dependent of %s", id, rt, depID),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// requestGetServiceDependency requests the list of service dependencies
// according to its resource type, then searches and returns the
// ServiceDependency with an id equal to `id`, returns a nil ServiceDependency
// if it is not found.
func (r *resourceServiceDependency) requestGetServiceDependency(ctx context.Context, id, depID, rt string) (*pagerduty.ServiceDependency, error) {
	var found *pagerduty.ServiceDependency

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var list *pagerduty.ListServiceDependencies
		var err error

		switch rt {
		case "service", "technical_service", "technical_service_reference":
			list, err = r.client.ListTechnicalServiceDependenciesWithContext(ctx, depID)
		case "business_service", "business_service_reference":
			list, err = r.client.ListBusinessServiceDependenciesWithContext(ctx, depID)
		default:
			err = fmt.Errorf("RT not available: %v", rt)
			return retry.RetryableError(err)
		}
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		for _, rel := range list.Relationships {
			if rel.ID == id {
				found = rel
				break
			}
		}
		return nil
	})

	return found, err
}

func (r *resourceServiceDependency) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceServiceDependency) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ".")
	if len(ids) != 3 {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_service_dependency",
			"Expecting an importation ID formed as '<supporting_service_id>.<supporting_service_type>.<service_dependency_id>'",
		)
		return
	}
	supID, supRt, id := ids[0], ids[1], ids[2]
	serviceDependency, err := r.requestGetServiceDependency(ctx, id, supID, supRt)
	if serviceDependency == nil || util.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error listing service dependencies", err.Error())
		return
	}

	model := flattenServiceDependency([]*pagerduty.ServiceDependency{serviceDependency}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

var supportingServiceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	},
}

var dependentServiceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	},
}

var serviceDependencyObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"supporting_service": types.ListType{
			ElemType: supportingServiceObjectType,
		},
		"dependent_service": types.ListType{
			ElemType: supportingServiceObjectType,
		},
	},
}

type resourceServiceDependencyItemModel struct {
	SupportingService types.List   `tfsdk:"supporting_service"`
	DependentService  types.List   `tfsdk:"dependent_service"`
	Type              types.String `tfsdk:"type"`
}

type resourceServiceDependencyModel struct {
	ID         types.String `tfsdk:"id"`
	Dependency types.List   `tfsdk:"dependency"`
}

var resourceServiceDependencyMu sync.Mutex

func buildServiceDependencyStruct(ctx context.Context, model resourceServiceDependencyModel) (*pagerduty.ServiceDependency, diag.Diagnostics) {
	var diags diag.Diagnostics

	var dependency []*resourceServiceDependencyItemModel
	if d := model.Dependency.ElementsAs(ctx, &dependency, false); d.HasError() {
		return nil, d
	}

	// These branches should not happen because of schema Validation
	if len(dependency) < 1 {
		diags.AddError("dependency length < 1", "")
		return nil, diags
	}
	if len(dependency[0].SupportingService.Elements()) < 1 {
		diags.AddError("supporting service not found for dependency", "")
	}
	if len(dependency[0].DependentService.Elements()) < 1 {
		diags.AddError("dependent service not found for dependency", "")
	}
	if diags.HasError() {
		return nil, diags
	}
	// ^These branches should not happen because of schema Validation

	ss, d := buildServiceObj(ctx, dependency[0].SupportingService.Elements()[0])
	if d.HasError() {
		diags.Append(d...)
		return nil, diags
	}
	ds, d := buildServiceObj(ctx, dependency[0].DependentService.Elements()[0])
	if d.HasError() {
		diags.Append(d...)
		return nil, diags
	}

	serviceDependency := &pagerduty.ServiceDependency{
		ID:                model.ID.ValueString(),
		Type:              dependency[0].Type.ValueString(),
		SupportingService: ss,
		DependentService:  ds,
	}

	return serviceDependency, diags
}

func buildServiceObj(ctx context.Context, model attr.Value) (*pagerduty.ServiceObj, diag.Diagnostics) {
	var diags diag.Diagnostics
	obj, ok := model.(types.Object)
	if !ok {
		diags.AddError("Not ok", "")
		return nil, diags
	}
	var serviceRef struct {
		ID   string `tfsdk:"id"`
		Type string `tfsdk:"type"`
	}
	obj.As(ctx, &serviceRef, basetypes.ObjectAsOptions{})
	serviceObj := pagerduty.ServiceObj(serviceRef)
	return &serviceObj, diags
}

func flattenServiceReference(objType types.ObjectType, src *pagerduty.ServiceObj) (list types.List, diags diag.Diagnostics) {
	if src == nil {
		diags.AddError("service reference is null", "")
		return
	}

	serviceRef, d := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
		"id":   types.StringValue(src.ID),
		"type": types.StringValue(convertServiceDependencyType(src.Type)),
	})
	if diags.Append(d...); diags.HasError() {
		return
	}

	list, d = types.ListValue(supportingServiceObjectType, []attr.Value{serviceRef})
	diags.Append(d...)
	return
}

func flattenServiceDependency(list []*pagerduty.ServiceDependency, diags *diag.Diagnostics) (model resourceServiceDependencyModel) {
	if len(list) < 1 {
		diags.AddError("Pagerduty did not responded with any dependency", "")
		return
	}
	item := list[0]

	supportingService, d := flattenServiceReference(supportingServiceObjectType, item.SupportingService)
	if diags.Append(d...); d.HasError() {
		return
	}

	dependentService, d := flattenServiceReference(dependentServiceObjectType, item.DependentService)
	if diags.Append(d...); d.HasError() {
		return
	}

	dependency, d := types.ObjectValue(
		serviceDependencyObjectType.AttrTypes,
		map[string]attr.Value{
			"type":               types.StringValue(item.Type),
			"supporting_service": supportingService,
			"dependent_service":  dependentService,
		},
	)
	if diags.Append(d...); d.HasError() {
		return model
	}

	dependencyList, d := types.ListValue(serviceDependencyObjectType, []attr.Value{dependency})
	if diags.Append(d...); d.HasError() {
		return model
	}

	model.ID = types.StringValue(item.ID)
	model.Dependency = dependencyList
	return model
}

// convertServiceDependencyType is needed because the PagerDuty API returns
// '*_reference' values in the response but uses the other kind of values in
// requests
func convertServiceDependencyType(s string) string {
	switch s {
	case "business_service_reference":
		s = "business_service"
	case "technical_service_reference":
		s = "service"
	}
	return s
}
