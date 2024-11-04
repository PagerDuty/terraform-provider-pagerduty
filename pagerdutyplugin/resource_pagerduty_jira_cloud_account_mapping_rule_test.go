package pagerduty

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyJiraCloudAccountsMappingRule_Basic(t *testing.T) {
	accountMappingID := os.Getenv("PAGERDUTY_ACC_JIRA_ACCOUNT_MAPPING_ID")
	if accountMappingID == "" {
		t.Skip("Missing env variable PAGERDUTY_ACC_JIRA_ACCOUNT_MAPPING_ID")
		return
	}

	rule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	ruleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyJiraCloudAccountsMappingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyJiraCloudAccountsMappingRuleConfig(service, user, rule, accountMappingID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyJiraCloudAccountsMappingRuleExists("pagerduty_jira_cloud_account_mapping_rule.foo"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "name", rule),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "account_mapping", accountMappingID),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "autocreate_jql_disabled_reason"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "autocreate_jql_disabled_until"),
					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.service"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.autocreate_jql", "priority = Highest"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.create_issue_on_incident_trigger", "true"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.0.source_incident_field", "incident_description"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.0.target_issue_field", "description"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.0.target_issue_field_name", "Description"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.0.type", "attribute"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.1.source_incident_field", "incident_number"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.1.target_issue_field", "customfield_10001"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.1.target_issue_field_name", "PagerDuty Incident Number"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.1.type", "attribute"),

					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.2.source_incident_field"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.2.target_issue_field", "security"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.2.target_issue_field_name", "Security Level"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.2.type", "jira_value"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.2.value", `{"displayName":"Sec Level 1","id":"10000"}`),

					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.3.source_incident_field"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.3.target_issue_field", "labels"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.3.target_issue_field_name", "Labels"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.3.type", "const"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.3.value", "pagerduty, incident"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.issue_type.id", "10001"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.issue_type.name", "Incident"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.0.jira_id", "1"),
					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.0.pagerduty_id"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.1.jira_id", "2"),
					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.1.pagerduty_id"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.2.jira_id", "3"),
					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.2.pagerduty_id"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.id", "10100"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.name", "IT Support"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.key", "ITS"),

					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.acknowledged.id", "2"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.acknowledged.name", "In Progress"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.resolved.id", "3"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.resolved.name", "Resolved"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.triggered.id", "1"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.triggered.name", "Open"),

					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.sync_notes_user"),
				),
			},

			{
				Config: testAccCheckPagerDutyJiraCloudAccountsMappingRuleConfigUpdated(service, ruleUpdated, accountMappingID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyJiraCloudAccountsMappingRuleExists("pagerduty_jira_cloud_account_mapping_rule.foo"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "name", ruleUpdated),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "account_mapping", accountMappingID),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "autocreate_jql_disabled_reason"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "autocreate_jql_disabled_until"),
					resource.TestCheckResourceAttrSet("pagerduty_jira_cloud_account_mapping_rule.foo", "config.service"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.autocreate_jql"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.create_issue_on_incident_trigger", "false"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.custom_fields.#", "0"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.issue_type.id", "10001"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.issue_type.name", "Incident"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.priorities.#", "0"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.id", "10100"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.name", "IT Support"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.project.key", "ITS"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.resolved"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.acknowledged"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.triggered.id", "1"),
					resource.TestCheckResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.status_mapping.triggered.name", "Open"),
					resource.TestCheckNoResourceAttr("pagerduty_jira_cloud_account_mapping_rule.foo", "config.jira.sync_notes_user"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyJiraCloudAccountsMappingRuleDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_jira_cloud_account_mapping_rule" {
			continue
		}

		ctx := context.Background()

		accountMappingID, ruleID, err := util.ResourcePagerDutyParseColonCompoundID(r.Primary.ID)
		if err != nil {
			return err
		}

		if _, err := testAccProvider.client.GetJiraCloudAccountsMappingRule(ctx, accountMappingID, ruleID); err == nil {
			return fmt.Errorf("Rule still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyJiraCloudAccountsMappingRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No rule ID is set")
		}

		accountMappingID, ruleID, err := util.ResourcePagerDutyParseColonCompoundID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := testAccProvider.client.GetJiraCloudAccountsMappingRule(ctx, accountMappingID, ruleID)
		if err != nil {
			return err
		}

		if found.ID != ruleID {
			return fmt.Errorf("Rule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyJiraCloudAccountsMappingRuleConfig(service, user, rule, accountMappingID string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

data "pagerduty_priority" "p1" {
	name = "P1"
}

data "pagerduty_priority" "p2" {
	name = "P2"
}

data "pagerduty_priority" "p3" {
	name = "P3"
}

resource "pagerduty_service" "foo" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_user" "foo" {
	name = "%s"
	email = "%[2]s@foo.test"
}

resource "pagerduty_jira_cloud_account_mapping_rule" "foo" {
	name = "%s"
	account_mapping = "%s"
	config {
		service = pagerduty_service.foo.id
		jira {
			autocreate_jql = "priority = Highest"
			create_issue_on_incident_trigger = true
			custom_fields {
				source_incident_field = "incident_description"
				target_issue_field = "description"
				target_issue_field_name = "Description"
				type = "attribute"
			}
			custom_fields {
				source_incident_field = "incident_number"
				target_issue_field = "customfield_10001"
				target_issue_field_name = "PagerDuty Incident Number"
				type = "attribute"
			}
			custom_fields {
				target_issue_field = "security"
				target_issue_field_name = "Security Level"
				type = "jira_value"
				value = jsonencode({
					displayName = "Sec Level 1"
					id = "10000"
				})
			}
			custom_fields {
				target_issue_field = "labels"
				target_issue_field_name = "Labels"
				type = "const"
				value = "pagerduty, incident"
			}
			issue_type {
				id = "10001"
				name = "Incident"
			}
			priorities {
				jira_id = "1"
				pagerduty_id = data.pagerduty_priority.p1.id
			}
			priorities {
				jira_id = "2"
				pagerduty_id = data.pagerduty_priority.p2.id
			}
			priorities {
				jira_id = "3"
				pagerduty_id = data.pagerduty_priority.p3.id
			}
			project {
				id = "10100"
				key = "ITS"
				name = "IT Support"
			}
			status_mapping {
				acknowledged {
					id = "2"
					name = "In Progress"
				}
				resolved {
					id = "3"
					name = "Resolved"
				}
				triggered {
					id = "1"
					name = "Open"
				}
			}
			sync_notes_user = pagerduty_user.foo.id
		}
	}
}`, service, user, rule, accountMappingID)
}

func testAccCheckPagerDutyJiraCloudAccountsMappingRuleConfigUpdated(service, rule, accountMappingID string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

data "pagerduty_priority" "p1" {
	name = "P1"
}

data "pagerduty_priority" "p2" {
	name = "P2"
}

data "pagerduty_priority" "p3" {
	name = "P3"
}

resource "pagerduty_service" "foo" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_jira_cloud_account_mapping_rule" "foo" {
	name = "%s"
	account_mapping = "%s"
	config {
		service = pagerduty_service.foo.id
		jira {
			issue_type {
				id = "10001"
				name = "Incident"
			}
			project {
				id = "10100"
				key = "ITS"
				name = "IT Support"
			}
			status_mapping {
				triggered {
					id = "1"
					name = "Open"
				}
			}
		}
	}
}`, service, rule, accountMappingID)
}
