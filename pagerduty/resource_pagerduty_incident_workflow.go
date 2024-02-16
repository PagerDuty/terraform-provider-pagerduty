package pagerduty

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// This integer controls the level of inline_steps_inputs recursion allowed in the Incident Workflow schema.
// A value of 1 indicates that the top level workflow steps may specify inline_steps_inputs of an action whose
// configuration requires it, but a step that exists within an inline_steps_input cannot itself specify an
// inline_steps_input. For example:
//
//	resource "pagerduty_incident_workflow" "test_workflow" {
//	  name = "Test Terraform Incident Workflow"
//	  step {
//	    name   = "Step 1"
//	    action = "pagerduty.com:incident-workflows:action:1"
//	    inline_steps_input { # Allowed iff inlineStepSchemaLevel>0
//	      name  = "Actions"
//	      step {
//	        name = "Inline Step 1"
//	        action = "pagerduty.com:incident-workflows:action:1"
//	        inline_steps_input { # Allowed iff inlineStepSchemaLevel>1
//	          ...
//	            inline_steps_input { # Allowed iff inlineStepSchemaLevel>2
//	          ...
//	        }
//	      }
//	    }
//	  }
//	}
//
// Each time the value is incremented by 1, it will allow another level of "inline_steps_input".
const inlineStepSchemaLevel = 1

func resourcePagerDutyIncidentWorkflow() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentWorkflowRead,
		UpdateContext: resourcePagerDutyIncidentWorkflowUpdate,
		DeleteContext: resourcePagerDutyIncidentWorkflowDelete,
		CreateContext: resourcePagerDutyIncidentWorkflowCreate,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyIncidentWorkflowImport,
		},
		CustomizeDiff: customizeIncidentWorkflowDiff(),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"team": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"step": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"input": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"generated": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"inline_steps_input": inlineStepSchema(inlineStepSchemaLevel),
					},
				},
			},
		},
	}
}

func inlineStepSchema(level int) *schema.Schema {
	if level > 1 {
		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"step": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
								"action": {
									Type:     schema.TypeString,
									Required: true,
								},
								"input": {
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"name": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:     schema.TypeString,
												Required: true,
											},
											"generated": {
												Type:     schema.TypeBool,
												Computed: true,
											},
										},
									},
								},
								"inline_steps_input": inlineStepSchema(level - 1),
							},
						},
					},
				},
			},
		}
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"step": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"action": {
								Type:     schema.TypeString,
								Required: true,
							},
							"input": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"value": {
											Type:     schema.TypeString,
											Required: true,
										},
										"generated": {
											Type:     schema.TypeBool,
											Computed: true,
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

// It is allowed for an incident workflow to return more inputs for a step than are present in the
// Terraform configuration. This can happen when there are inputs which have a default value.
// These inputs get persisted in the state with `generated` set to `true`
// but then should not be removed by a `terraform apply`
func customizeIncidentWorkflowDiff() schema.CustomizeDiffFunc {
	// Regex targets all changed keys that have affected either an input list's length, a specific input's name, or a
	// specific input's value. The diff function is built to handle all input changes.
	inputCountRegex := regexp.MustCompile(`^(step\.\d+\.?.*\.input)\.(#|\d+\.name|\d+.value)$`)

	excludeInput := func(inputs []interface{}, inputName string) []interface{} {
		result := make([]interface{}, 0)
		for _, input := range inputs {
			if !(input.(map[string]interface{})["name"] == inputName) {
				result = append(result, input)
			}
		}
		return result
	}

	determineNewInputs := func(oldInputs []interface{}, newInputs []interface{}) interface{} {
		result := make([]interface{}, 0)

		// Statically set generated to false for all newInputs because we know the input was found in the config so it
		// is by definition not generated. Terraform may think it is generated if the ordering of the inputs changes
		// around, so it helpfully assigns generated since it was not *precisely* the input specified in the config
		// since the ordering is different. This static assignment overrides that so ordering issues alone do not
		// generate a plan diff.
		for _, newInputValue := range newInputs {
			newInputValue.(map[string]interface{})["generated"] = false
		}

		// 1. Iterate through old versions from tfstate.
		for _, oldInputValue := range oldInputs {
			oldInput := oldInputValue.(map[string]interface{})

			// 2. For each old input, search for the same input within all new inputs.
			var input map[string]interface{}
			for _, newInputValue := range newInputs {
				newInput := newInputValue.(map[string]interface{})
				if oldInput["name"].(string) == newInput["name"].(string) {
					// 3. If there is a new version of the same input, use that and then remove it from newInputs.
					// This keeps the ordering of inputs the same to prevent the diff from thinking a change occurred
					// when it is just a different ordering of inputs (ordering does not matter).
					// Excluding the new input lets us track which new inputs have been handled or not.
					input = newInput
					newInputs = excludeInput(newInputs, input["name"].(string))
					break
				}
			}

			// 4. If there was no matching new input to this old input and the old input is generated, use the old
			// input. This lets the diff maintain the cached tfstate generated inputs as they will always exist unless
			// overridden with a non-generated input from the config.
			if input == nil && oldInput["generated"].(bool) {
				input = oldInput
			}

			// 5. If an input was found, append it to the result.
			if input != nil {
				result = append(result, input)
			}
		}

		// 6. Any remaining new inputs that were not already matched to old inputs should now be added to the result.
		// This lets new inputs that do not have default values (which are already cached as generated=true) maintain
		// their state in the diff, while not disrupting ordering of unchanged inputs.
		result = append(result, newInputs...)

		return result
	}

	var updateStepInput func(ctx interface{}, path string, input interface{})
	updateStepInput = func(ctx interface{}, path string, input interface{}) {
		if path == "input" {
			ctx.(map[string]interface{})["input"] = input
			return
		}
		parts := strings.Split(path, ".")
		if len(parts) <= 1 {
			log.Printf("[WARN] Unexpected input key not terminated by `.input`")
			return
		}

		top := parts[0]
		newPath := strings.Join(parts[1:], ".")

		idx, err := strconv.ParseInt(top, 10, 32)
		if err == nil {
			updateStepInput(ctx.([]interface{})[idx], newPath, input)
		} else if top == "step" || top == "inline_steps_input" {
			updateStepInput(ctx.(map[string]interface{})[top], newPath, input)
		} else {
			log.Printf("[WARN] Unexpected input key part: %s - expected integer,\"step\",\"inline_steps_input\"", top)
		}
	}

	return func(ctx context.Context, d *schema.ResourceDiff, _ interface{}) error {
		id := d.Id()
		if id != "" {
			keys := d.GetChangedKeysPrefix("step")

			_, newStep := d.GetChange("step")
			needToSetNew := false

			for _, key := range keys {
				indexMatch := inputCountRegex.FindStringSubmatch(key)
				if len(indexMatch) == 3 {
					inputKey := indexMatch[1]
					oldInputs, newInputs := d.GetChange(inputKey)
					replacementInputs := determineNewInputs(oldInputs.([]interface{}), newInputs.([]interface{}))

					log.Printf("[INFO] Updating diff for input key %s to include generated inputs.", inputKey)
					// Initiate the path after `step.` as we already know we are in a step context from the regex
					updateStepInput(newStep, inputKey[5:], replacementInputs)
					needToSetNew = true
				}
			}

			if needToSetNew {
				err := d.SetNew("step", newStep)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func resourcePagerDutyIncidentWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iw, specifiedSteps, err := buildIncidentWorkflowStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident workflow %s.", iw.Name)

	createdWorkflow, _, err := client.IncidentWorkflows.CreateContext(ctx, iw)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentWorkflow(d, createdWorkflow, true, specifiedSteps, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty incident workflow %s", d.Id())
	err := fetchIncidentWorkflow(ctx, d, meta, handleNotFoundError, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iw, specifiedSteps, err := buildIncidentWorkflowStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident workflow %s", d.Id())

	updatedWorkflow, _, err := client.IncidentWorkflows.UpdateContext(ctx, d.Id(), iw)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentWorkflow(d, updatedWorkflow, true, specifiedSteps, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.IncidentWorkflows.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := fetchIncidentWorkflow(ctx, d, m, handleNotFoundError, true)
	return []*schema.ResourceData{d}, err
}

func fetchIncidentWorkflow(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error, isImport bool) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	_, specifiedSteps, err := buildIncidentWorkflowStruct(d)
	if err != nil {
		return err
	}

	return retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		iw, _, err := client.IncidentWorkflows.GetContext(ctx, d.Id())
		if err != nil {
			log.Printf("[WARN] Incident workflow read error")
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenIncidentWorkflow(d, iw, true, specifiedSteps, isImport); err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

func flattenIncidentWorkflow(
	d *schema.ResourceData,
	iw *pagerduty.IncidentWorkflow,
	includeSteps bool,
	specifiedSteps []*SpecifiedStep,
	isImport bool,
) error {
	d.SetId(iw.ID)
	d.Set("name", iw.Name)
	if iw.Description != nil {
		d.Set("description", *(iw.Description))
	}
	if iw.Team != nil {
		d.Set("team", iw.Team.ID)
	}

	if includeSteps {
		steps := flattenIncidentWorkflowSteps(iw, specifiedSteps, isImport)
		d.Set("step", steps)
	}

	return nil
}

func flattenIncidentWorkflowSteps(iw *pagerduty.IncidentWorkflow, specifiedSteps []*SpecifiedStep, isImport bool) []map[string]interface{} {
	newSteps := make([]map[string]interface{}, len(iw.Steps))
	for i, s := range iw.Steps {
		m := make(map[string]interface{})
		m["id"] = s.ID
		m["name"] = s.Name
		m["action"] = s.Configuration.ActionID

		var inputNames []string
		inlineInputs := make(map[string][]*SpecifiedStep)
		if !isImport {
			specifiedStep := *specifiedSteps[i]
			inputNames = specifiedStep.SpecifiedInputNames
			inlineInputs = specifiedStep.SpecifiedInlineInputs
		}

		m["input"] = flattenIncidentWorkflowStepInput(s.Configuration.Inputs, inputNames, isImport)
		m["inline_steps_input"] = flattenIncidentWorkflowStepInlineStepsInput(
			s.Configuration.InlineStepsInputs,
			inlineInputs,
			isImport,
		)

		newSteps[i] = m
	}

	return newSteps
}

func flattenIncidentWorkflowStepInput(inputs []*pagerduty.IncidentWorkflowActionInput, specifiedInputNames []string, isImport bool) *[]interface{} {
	newInputs := make([]interface{}, len(inputs))

	for i, v := range inputs {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["value"] = v.Value

		if !isImport && !isInputInNonGeneratedInputNames(v, specifiedInputNames) {
			m["generated"] = true
		}

		newInputs[i] = m
	}
	return &newInputs
}

func flattenIncidentWorkflowStepInlineStepsInput(
	inlineStepsInputs []*pagerduty.IncidentWorkflowActionInlineStepsInput,
	specifiedInlineInputs map[string][]*SpecifiedStep,
	isImport bool,
) *[]interface{} {
	newInlineStepsInputs := make([]interface{}, len(inlineStepsInputs))

	for i, v := range inlineStepsInputs {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["step"] = flattenIncidentWorkflowStepInlineStepsInputSteps(v.Value.Steps, specifiedInlineInputs[v.Name], isImport)

		newInlineStepsInputs[i] = m
	}
	return &newInlineStepsInputs
}

func flattenIncidentWorkflowStepInlineStepsInputSteps(
	inlineSteps []*pagerduty.IncidentWorkflowActionInlineStep,
	specifiedSteps []*SpecifiedStep,
	isImport bool,
) *[]interface{} {
	newInlineSteps := make([]interface{}, len(inlineSteps))

	for i, v := range inlineSteps {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["action"] = v.Configuration.ActionID

		var inputNames []string
		inlineInputs := make(map[string][]*SpecifiedStep)
		if !isImport {
			specifiedStep := *specifiedSteps[i]
			inputNames = specifiedStep.SpecifiedInputNames
			inlineInputs = specifiedStep.SpecifiedInlineInputs
		}

		m["input"] = flattenIncidentWorkflowStepInput(v.Configuration.Inputs, inputNames, isImport)
		if v.Configuration.InlineStepsInputs != nil && len(v.Configuration.InlineStepsInputs) > 0 {
			// We should prefer to not set inline_steps_input if the array is empty. This doubles as a schema edge guard
			// and prevents an invalid set if we try to set inline_steps_input to an empty array where the schema
			// disallows setting any value whatsoever.
			m["inline_steps_input"] = flattenIncidentWorkflowStepInlineStepsInput(
				v.Configuration.InlineStepsInputs,
				inlineInputs,
				isImport,
			)
		}

		newInlineSteps[i] = m
	}
	return &newInlineSteps
}

func isInputInNonGeneratedInputNames(i *pagerduty.IncidentWorkflowActionInput, names []string) bool {
	for _, in := range names {
		if i.Name == in {
			return true
		}
	}
	return false
}

// Tracks specified inputs recursively to identify which are generated or not
type SpecifiedStep struct {
	SpecifiedInputNames   []string
	SpecifiedInlineInputs map[string][]*SpecifiedStep
}

func buildIncidentWorkflowStruct(d *schema.ResourceData) (
	*pagerduty.IncidentWorkflow,
	[]*SpecifiedStep,
	error,
) {
	iw := pagerduty.IncidentWorkflow{
		Name: d.Get("name").(string),
	}
	if desc, ok := d.GetOk("description"); ok {
		str := desc.(string)
		iw.Description = &str
	}
	if team, ok := d.GetOk("team"); ok {
		iw.Team = &pagerduty.TeamReference{
			ID: team.(string),
		}
	}

	specifiedSteps := make([]*SpecifiedStep, 0)
	if steps, ok := d.GetOk("step"); ok {
		iw.Steps, specifiedSteps = buildIncidentWorkflowStepsStruct(steps)
	}

	return &iw, specifiedSteps, nil
}

func buildIncidentWorkflowStepsStruct(s interface{}) (
	[]*pagerduty.IncidentWorkflowStep,
	[]*SpecifiedStep,
) {
	steps := s.([]interface{})
	newSteps := make([]*pagerduty.IncidentWorkflowStep, len(steps))
	specifiedSteps := make([]*SpecifiedStep, len(steps))

	for i, v := range steps {
		stepData := v.(map[string]interface{})
		step := pagerduty.IncidentWorkflowStep{
			Name: stepData["name"].(string),
			Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
				ActionID: stepData["action"].(string),
			},
		}
		if id, ok := stepData["id"]; ok {
			step.ID = id.(string)
		}

		specifiedStep := SpecifiedStep{
			SpecifiedInputNames:   make([]string, 0),
			SpecifiedInlineInputs: map[string][]*SpecifiedStep{},
		}
		step.Configuration.Inputs,
			specifiedStep.SpecifiedInputNames = buildIncidentWorkflowInputsStruct(stepData["input"])
		step.Configuration.InlineStepsInputs,
			specifiedStep.SpecifiedInlineInputs = buildIncidentWorkflowInlineStepsInputsStruct(stepData["inline_steps_input"])

		newSteps[i] = &step
		specifiedSteps[i] = &specifiedStep
	}
	return newSteps, specifiedSteps
}

func buildIncidentWorkflowInputsStruct(in interface{}) (
	[]*pagerduty.IncidentWorkflowActionInput,
	[]string,
) {
	inputs := in.([]interface{})
	newInputs := make([]*pagerduty.IncidentWorkflowActionInput, len(inputs))
	specifiedInputNames := make([]string, 0)

	for i, v := range inputs {
		inputData := v.(map[string]interface{})
		input := pagerduty.IncidentWorkflowActionInput{
			Name:  inputData["name"].(string),
			Value: inputData["value"].(string),
		}

		generated := inputData["generated"].(bool)
		if !generated {
			specifiedInputNames = append(specifiedInputNames, input.Name)
		}
		newInputs[i] = &input
	}
	return newInputs, specifiedInputNames
}

func buildIncidentWorkflowInlineStepsInputsStruct(in interface{}) (
	[]*pagerduty.IncidentWorkflowActionInlineStepsInput,
	map[string][]*SpecifiedStep,
) {
	specifiedInlineInputs := map[string][]*SpecifiedStep{}
	if in == nil {
		// We need to catch the case where the schema stops allowing inline_steps_input where this will be nil so that
		// we can return an empty list.
		return make([]*pagerduty.IncidentWorkflowActionInlineStepsInput, 0), specifiedInlineInputs
	}
	inputs := in.([]interface{})
	newInputs := make([]*pagerduty.IncidentWorkflowActionInlineStepsInput, len(inputs))

	for i, v := range inputs {
		inputData := v.(map[string]interface{})
		input := pagerduty.IncidentWorkflowActionInlineStepsInput{
			Name: inputData["name"].(string),
		}
		inputValue := pagerduty.IncidentWorkflowActionInlineStepsInputValue{}
		steps, specifiedInlineSteps := buildIncidentWorkflowActionInlineStepsInputSteps(inputData["step"])
		inputValue.Steps = steps
		input.Value = &inputValue
		specifiedInlineInputs[input.Name] = specifiedInlineSteps
		newInputs[i] = &input
	}
	return newInputs, specifiedInlineInputs
}

func buildIncidentWorkflowActionInlineStepsInputSteps(in interface{}) (
	[]*pagerduty.IncidentWorkflowActionInlineStep,
	[]*SpecifiedStep,
) {
	inlineSteps := in.([]interface{})
	newInlineSteps := make([]*pagerduty.IncidentWorkflowActionInlineStep, len(inlineSteps))
	specifiedInlineSteps := make([]*SpecifiedStep, len(inlineSteps))

	for i, v := range inlineSteps {
		inlineStepData := v.(map[string]interface{})
		inlineStep := pagerduty.IncidentWorkflowActionInlineStep{
			Name: inlineStepData["name"].(string),
			Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
				ActionID: inlineStepData["action"].(string),
			},
		}

		specifiedInlineStep := SpecifiedStep{}
		inlineStep.Configuration.Inputs,
			specifiedInlineStep.SpecifiedInputNames = buildIncidentWorkflowInputsStruct(inlineStepData["input"])
		inlineStep.Configuration.InlineStepsInputs,
			specifiedInlineStep.SpecifiedInlineInputs = buildIncidentWorkflowInlineStepsInputsStruct(inlineStepData["inline_steps_input"])

		newInlineSteps[i] = &inlineStep
		specifiedInlineSteps[i] = &specifiedInlineStep
	}
	return newInlineSteps, specifiedInlineSteps
}
