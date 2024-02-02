package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

const (
	AppBaseUrl         = "https://app.pagerduty.com"
	StarWildcardConfig = "*"
)

func resourcePagerDutySlackConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutySlackConnectionCreate,
		Read:   resourcePagerDutySlackConnectionRead,
		Update: resourcePagerDutySlackConnectionUpdate,
		Delete: resourcePagerDutySlackConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutySlackConnectionImport,
		},
		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"service_reference",
					"team_reference",
				}),
			},
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SLACK_CONNECTION_WORKSPACE_ID", nil),
			},
			"notification_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"responder",
					"stakeholder",
				}),
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"events": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"priorities": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"urgency": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateDiagFunc: validateValueDiagFunc([]string{
								"high",
								"low",
							}),
						},
					},
				},
			},
		},
	}
}

func buildSlackConnectionStruct(d *schema.ResourceData) (*pagerduty.SlackConnection, error) {
	slackConn := pagerduty.SlackConnection{
		SourceID:         d.Get("source_id").(string),
		SourceName:       d.Get("source_name").(string),
		SourceType:       d.Get("source_type").(string),
		ChannelID:        d.Get("channel_id").(string),
		ChannelName:      d.Get("channel_name").(string),
		WorkspaceID:      d.Get("workspace_id").(string),
		NotificationType: d.Get("notification_type").(string),
		Config:           expandConnectionConfig(d.Get("config").(interface{})),
	}
	return &slackConn, nil
}

func resourcePagerDutySlackConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).SlackClient()
	if err != nil {
		return err
	}

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		slackConn, err := buildSlackConnectionStruct(d)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		log.Printf("[INFO] Creating PagerDuty slack connection for source %s and slack channel %s", slackConn.SourceID, slackConn.ChannelID)

		if slackConn, _, err = client.SlackConnections.Create(slackConn.WorkspaceID, slackConn); err != nil {
			return retry.RetryableError(err)
		} else if slackConn != nil {
			d.SetId(slackConn.ID)
			d.Set("workspace_id", slackConn.WorkspaceID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return resourcePagerDutySlackConnectionRead(d, meta)
}

func resourcePagerDutySlackConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).SlackClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty slack connection %s", d.Id())

	workspaceID := d.Get("workspace_id").(string)
	log.Printf("[DEBUG] Read Slack Connection: workspace_id %s", workspaceID)

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if slackConn, _, err := client.SlackConnections.Get(workspaceID, d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if slackConn != nil {
			d.Set("source_id", slackConn.SourceID)
			d.Set("source_name", slackConn.SourceName)
			d.Set("source_type", slackConn.SourceType)
			d.Set("channel_id", slackConn.ChannelID)
			d.Set("channel_name", slackConn.ChannelName)
			d.Set("notification_type", slackConn.NotificationType)
			d.Set("config", flattenConnectionConfig(slackConn.Config))
		}
		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return nil
}

func resourcePagerDutySlackConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).SlackClient()
	if err != nil {
		return err
	}

	slackConn, err := buildSlackConnectionStruct(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Updating PagerDuty slack connection %s", d.Id())

	if _, _, err := client.SlackConnections.Update(slackConn.WorkspaceID, d.Id(), slackConn); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutySlackConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).SlackClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty slack connection %s", d.Id())
	workspaceID := d.Get("workspace_id").(string)

	if _, err := client.SlackConnections.Delete(workspaceID, d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandConnectionConfig(v interface{}) pagerduty.ConnectionConfig {
	c := v.([]interface{})[0].(map[string]interface{})

	var config pagerduty.ConnectionConfig

	config = pagerduty.ConnectionConfig{
		Events:     expandConfigList(c["events"].([]interface{})),
		Priorities: expandStarWildcardConfig(expandConfigList(c["priorities"].([]interface{}))),
		Urgency:    nil,
	}
	if val, ok := c["urgency"]; ok {
		urgency := val.(string)
		if len(urgency) > 0 {
			config.Urgency = &urgency
		}
	}
	return config
}

func expandConfigList(v interface{}) []string {
	items := []string{}
	for _, i := range v.([]interface{}) {
		items = append(items, i.(string))
	}
	return items
}

// Expands the use of star wildcard ("*") configuration for an attribute to its
// matching expected value by PagerDuty's API, which is nil. This is necessary
// when the API accepts and interprets nil and empty configurations as valid
// settings. The state produced by this kind of config can be reverted to the API
// expected values with sibbling function `flattenStarWildcardConfig`.
func expandStarWildcardConfig(c []string) []string {
	if isUsingStarWildcardConfig := len(c) == 1 && c[0] == StarWildcardConfig; isUsingStarWildcardConfig {
		c = nil
	}
	return c
}

func flattenConnectionConfig(config pagerduty.ConnectionConfig) []map[string]interface{} {
	var configs []map[string]interface{}
	configMap := map[string]interface{}{
		"events":     flattenConfigList(config.Events),
		"priorities": flattenConfigList(flattenStarWildcardConfig(config.Priorities)),
	}
	if config.Urgency != nil {
		configMap["urgency"] = *config.Urgency
	}
	configs = append(configs, configMap)
	return configs
}

func flattenConfigList(list []string) interface{} {
	var items []interface{}

	for _, i := range list {
		items = append(items, i)
	}

	return items
}

// Flattens a `nil` configuration to its corresponding star wildcard ("*")
// configuration value for an attribute which is meant to be accepting this kind
// of configuration, with the only purpose to match the config stored in the
// Terraform's state.
func flattenStarWildcardConfig(c []string) []string {
	if hasStarWildcardConfigSet := c[:] == nil; hasStarWildcardConfigSet {
		c = append(c, StarWildcardConfig)
	}
	return c
}

func resourcePagerDutySlackConnectionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).SlackClient()
	if err != nil {
		return nil, err
	}

	ids := strings.Split(d.Id(), ".")

	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_slack_connection. Expecting an importation ID formed as '<workspace_id>.<slack_connection_id>'")
	}
	workspaceID, connectionID := ids[0], ids[1]

	_, _, err = client.SlackConnections.Get(workspaceID, connectionID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(connectionID)
	d.Set("workspace_id", workspaceID)

	return []*schema.ResourceData{d}, nil
}
