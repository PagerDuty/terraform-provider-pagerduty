package pagerduty

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyIncidentWorkflow() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentWorkflowRead,
		UpdateContext: resourcePagerDutyIncidentWorkflowUpdate,
		DeleteContext: resourcePagerDutyIncidentWorkflowDelete,
		CreateContext: resourcePagerDutyIncidentWorkflowCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	inputCountRegex := regexp.MustCompile(`^(step\.(\d+)\.input)\.#$`)

	nameDoesNotExistAlready := func(inputs []interface{}, name string) bool {
		for _, v := range inputs {
			m := v.(map[string]interface{})
			if n, ok := m["name"]; ok && n == name {
				return false
			}
		}
		return true
	}

	addMissingGeneratedInputs := func(withGenerated []interface{}, maybeWithoutGenerated []interface{}) (interface{}, []string) {
		result := maybeWithoutGenerated
		addedNames := make([]string, 0)

		for _, v := range withGenerated {
			m := v.(map[string]interface{})
			if b, ok := m["generated"]; ok && b.(bool) && nameDoesNotExistAlready(maybeWithoutGenerated, m["name"].(string)) {
				result = append(result, m)
				addedNames = append(addedNames, m["name"].(string))
			}
		}

		return result, addedNames
	}

	updateStepInput := func(step interface{}, index int64, input interface{}) {
		step.([]interface{})[index].(map[string]interface{})["input"] = input
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
					inputsFromState, inputsAfterDiff := d.GetChange(inputKey)

					replacementInput, addedNames := addMissingGeneratedInputs(inputsFromState.([]interface{}), inputsAfterDiff.([]interface{}))
					stepIndex, _ := strconv.ParseInt(indexMatch[2], 10, 32)
					updateStepInput(newStep, stepIndex, replacementInput)
					log.Printf("[INFO] Updating diff for step %d to include generated inputs %v.", stepIndex, addedNames)
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

	iw, err := buildIncidentWorkflowStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident workflow %s.", iw.Name)

	createdWorkflow, _, err := client.IncidentWorkflows.CreateContext(ctx, iw)
	if err != nil {
		return diag.FromErr(err)
	}

	stepIdMapping := map[int]string{}
	for i, s := range createdWorkflow.Steps {
		stepIdMapping[i] = s.ID
	}

	nonGeneratedInputNames := createNonGeneratedInputNamesFromWorkflowStepsWithoutIDs(iw, stepIdMapping)
	err = flattenIncidentWorkflow(d, createdWorkflow, true, nonGeneratedInputNames)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty incident workflow %s", d.Id())
	err := fetchIncidentWorkflow(ctx, d, meta, handleNotFoundError)
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

	nonGeneratedInputNames := createNonGeneratedInputNamesFromWorkflowStepsInResourceData(d)

	iw, err := buildIncidentWorkflowStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident workflow %s", d.Id())

	updatedWorkflow, _, err := client.IncidentWorkflows.UpdateContext(ctx, d.Id(), iw)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentWorkflow(d, updatedWorkflow, true, nonGeneratedInputNames)
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

func fetchIncidentWorkflow(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	nonGeneratedInputNames := createNonGeneratedInputNamesFromWorkflowStepsInResourceData(d)

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		iw, _, err := client.IncidentWorkflows.GetContext(ctx, d.Id())
		if err != nil {
			log.Printf("[WARN] Incident workflow read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenIncidentWorkflow(d, iw, true, nonGeneratedInputNames); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}

func createNonGeneratedInputNamesFromWorkflowStepsInResourceData(d *schema.ResourceData) map[string][]string {
	nonGeneratedInputNames := map[string][]string{}

	if _, ok := d.GetOk("name"); ok {
		if s, ok := d.GetOk("step"); ok {
			steps := s.([]interface{})

			for _, v := range steps {
				stepData := v.(map[string]interface{})
				if id, idOk := stepData["id"]; idOk {
					nonGeneratedInputNamesForStep := make([]string, 0)
					if sdis, inputOk := stepData["input"]; inputOk {
						for _, sdi := range sdis.([]interface{}) {
							inputData := sdi.(map[string]interface{})
							if b, ok := inputData["generated"]; !ok || !b.(bool) {
								nonGeneratedInputNamesForStep = append(nonGeneratedInputNamesForStep, inputData["name"].(string))
							}
						}
					}
					nonGeneratedInputNames[id.(string)] = nonGeneratedInputNamesForStep
				}
			}
		}
	}
	return nonGeneratedInputNames
}

func createNonGeneratedInputNamesFromWorkflowStepsWithoutIDs(iw *pagerduty.IncidentWorkflow, stepIdMapping map[int]string) map[string][]string {
	nonGeneratedInputNames := map[string][]string{}

	if iw != nil {
		for i, step := range iw.Steps {
			nonGeneratedInputNamesForStep := make([]string, 0)
			if step.Configuration != nil {
				for _, input := range step.Configuration.Inputs {
					nonGeneratedInputNamesForStep = append(nonGeneratedInputNamesForStep, input.Name)
				}
			}
			nonGeneratedInputNames[stepIdMapping[i]] = nonGeneratedInputNamesForStep
		}
	}
	return nonGeneratedInputNames
}

func flattenIncidentWorkflow(d *schema.ResourceData, iw *pagerduty.IncidentWorkflow, includeSteps bool, nonGeneratedInputNames map[string][]string) error {
	d.SetId(iw.ID)
	d.Set("name", iw.Name)
	if iw.Description != nil {
		d.Set("description", *(iw.Description))
	}
	if iw.Team != nil {
		d.Set("team", iw.Team.ID)
	}

	if includeSteps {
		steps := flattenIncidentWorkflowSteps(iw, nonGeneratedInputNames)
		d.Set("step", steps)
	}

	return nil
}

func flattenIncidentWorkflowSteps(iw *pagerduty.IncidentWorkflow, nonGeneratedInputNames map[string][]string) []map[string]interface{} {
	newSteps := make([]map[string]interface{}, len(iw.Steps))
	for i, s := range iw.Steps {
		nonGeneratedInputNamesForStep, ok := nonGeneratedInputNames[s.ID]
		if !ok {
			nonGeneratedInputNamesForStep = make([]string, 0)
		}

		m := make(map[string]interface{})
		m["id"] = s.ID
		m["name"] = s.Name
		m["action"] = s.Configuration.ActionID
		m["input"] = flattenIncidentWorkflowStepInput(s.Configuration.Inputs, nonGeneratedInputNamesForStep)

		newSteps[i] = m
	}
	return newSteps
}

func flattenIncidentWorkflowStepInput(inputs []*pagerduty.IncidentWorkflowActionInput, nonGeneratedInputNames []string) *[]interface{} {
	newInputs := make([]interface{}, len(inputs))

	for i, v := range inputs {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["value"] = v.Value

		if !isInputInNonGeneratedInputNames(v, nonGeneratedInputNames) {
			m["generated"] = true
		}

		newInputs[i] = m
	}
	return &newInputs
}

func isInputInNonGeneratedInputNames(i *pagerduty.IncidentWorkflowActionInput, names []string) bool {
	for _, in := range names {
		if i.Name == in {
			return true
		}
	}
	return false
}

func buildIncidentWorkflowStruct(d *schema.ResourceData) (*pagerduty.IncidentWorkflow, error) {
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

	if steps, ok := d.GetOk("step"); ok {
		iw.Steps = buildIncidentWorkflowStepsStruct(steps)
	}

	return &iw, nil

}

func buildIncidentWorkflowStepsStruct(s interface{}) []*pagerduty.IncidentWorkflowStep {
	steps := s.([]interface{})
	newSteps := make([]*pagerduty.IncidentWorkflowStep, len(steps))

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

		step.Configuration.Inputs = buildIncidentWorkflowInputsStruct(stepData["input"])

		newSteps[i] = &step
	}
	return newSteps
}

func buildIncidentWorkflowInputsStruct(in interface{}) []*pagerduty.IncidentWorkflowActionInput {
	inputs := in.([]interface{})
	newInputs := make([]*pagerduty.IncidentWorkflowActionInput, len(inputs))

	for i, v := range inputs {
		inputData := v.(map[string]interface{})
		input := pagerduty.IncidentWorkflowActionInput{
			Name:  inputData["name"].(string),
			Value: inputData["value"].(string),
		}

		newInputs[i] = &input
	}
	return newInputs
}
