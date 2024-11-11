package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_alert_grouping_setting", &resource.Sweeper{
		Name: "pagerduty_alert_grouping_setting",
		F:    testSweepAlertGroupingSetting,
	})
}

func testSweepAlertGroupingSetting(_ string) error {
	ctx := context.Background()

	resp, err := testAccProvider.client.ListAlertGroupingSettings(ctx, pagerduty.ListAlertGroupingSettingsOptions{
		Limit: 100,
	})
	if err != nil {
		return err
	}

	for _, setting := range resp.AlertGroupingSettings {
		if strings.HasPrefix(setting.Name, "test") || strings.HasPrefix(setting.Name, "tf-") {
			log.Printf("Destroying alert grouping setting %s (%s)", setting.Name, setting.ID)
			if err := testAccProvider.client.DeleteAlertGroupingSetting(ctx, setting.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyAlertGroupingSetting_Basic(t *testing.T) {
	ref := fmt.Sprint("tf-", acctest.RandString(5))
	resourceRef := "pagerduty_alert_grouping_setting." + ref

	service := fmt.Sprint("tf-", acctest.RandString(5))
	serviceUpdated := fmt.Sprint("tf-", acctest.RandString(5))
	name := ref + "'s name"
	nameUpdated := ref + "'s name updated"

	configType := string(pagerduty.AlertGroupingSettingContentBasedType)
	config := pagerduty.AlertGroupingSettingConfigContentBased{
		TimeWindow: 0,
		Aggregate:  "all",
		Fields:     []string{"summary"},
	}
	configTypeUpdated := string(pagerduty.AlertGroupingSettingTimeType)
	configUpdated := pagerduty.AlertGroupingSettingConfigTime{Timeout: 60}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingConfig(ref, name, configType, service, config),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists(resourceRef),
					resource.TestCheckResourceAttr(resourceRef, "name", name),
					resource.TestCheckResourceAttrSet(resourceRef, "description"),
					resource.TestCheckResourceAttr(resourceRef, "type", configType),
					resource.TestCheckResourceAttr(resourceRef, "config.time_window", fmt.Sprint(300)),
					resource.TestCheckResourceAttr(resourceRef, "config.aggregate", config.Aggregate),
					resource.TestCheckResourceAttr(resourceRef, "config.fields.0", config.Fields[0]),
					resource.TestCheckResourceAttrSet(resourceRef, "services.0"),
				),
			},
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingConfig(ref, nameUpdated, configTypeUpdated, serviceUpdated, configUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists(resourceRef),
					resource.TestCheckResourceAttr(resourceRef, "name", nameUpdated),
					resource.TestCheckResourceAttrSet(resourceRef, "description"),
					resource.TestCheckResourceAttr(resourceRef, "type", configTypeUpdated),
					resource.TestCheckResourceAttr(resourceRef, "config.timeout", fmt.Sprint(configUpdated.Timeout)),
					resource.TestCheckResourceAttrSet(resourceRef, "services.0"),
				),
			},
		},
	})
}

func TestAccPagerDutyAlertGroupingSetting_AppendService(t *testing.T) {
	name := fmt.Sprint("tf-", acctest.RandString(5))
	service1 := fmt.Sprint("tf-", acctest.RandString(5))
	service2 := fmt.Sprint("tf-", acctest.RandString(5))
	service3 := fmt.Sprint("tf-", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingAppendServiceConfig(name, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting.foo"),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting.foo", "services.#", "2"),
				),
			},
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingAppendServiceConfigUpdated(name, service1, service2, service3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting.foo"),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting.foo", "services.#", "3"),
				),
			},
		},
	})
}

func TestAccPagerDutyAlertGroupingSetting_PopService(t *testing.T) {
	name := fmt.Sprint("tf-", acctest.RandString(5))
	service1 := fmt.Sprint("tf-", acctest.RandString(5))
	service2 := fmt.Sprint("tf-", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingAppendServiceConfig(name, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting.foo"),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting.foo", "services.#", "2"),
				),
			},
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingPopServiceConfigUpdated(name, service1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting.foo"),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting.foo", "services.#", "1"),
				),
			},
		},
	})
}

func TestAccPagerDutyAlertGroupingSetting_ContentBased_WithTimeWindow(t *testing.T) {
	ref := fmt.Sprint("tf-", acctest.RandString(5))
	resourceRef := "pagerduty_alert_grouping_setting." + ref
	name := ref + "'s name"

	configType := string(pagerduty.AlertGroupingSettingContentBasedType)
	config := pagerduty.AlertGroupingSettingConfigContentBased{
		TimeWindow: 600,
		Aggregate:  "all",
		Fields:     []string{"summary"},
	}

	service := fmt.Sprint("tf-", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingConfig(ref, name, configType, service, config),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists(resourceRef),
					resource.TestCheckResourceAttr(resourceRef, "name", name),
					resource.TestCheckResourceAttrSet(resourceRef, "description"),
					resource.TestCheckResourceAttr(resourceRef, "type", configType),
					resource.TestCheckResourceAttr(resourceRef, "config.time_window", fmt.Sprint(config.TimeWindow)),
					resource.TestCheckResourceAttr(resourceRef, "config.aggregate", config.Aggregate),
					resource.TestCheckResourceAttr(resourceRef, "config.fields.0", config.Fields[0]),
					resource.TestCheckResourceAttrSet(resourceRef, "services.0"),
				),
			},
		},
	})
}

func TestAccPagerDutyAlertGroupingSetting_Time_WithTimeoutZero(t *testing.T) {
	ref := fmt.Sprint("tf-", acctest.RandString(5))
	name := ref

	configType := string(pagerduty.AlertGroupingSettingTimeType)
	config := pagerduty.AlertGroupingSettingConfigTime{Timeout: 0}

	service := fmt.Sprint("tf-", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAlertGroupingSettingConfig(ref, name, configType, service, config),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting."+ref),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting."+ref, "name", name),
					resource.TestCheckResourceAttrSet("pagerduty_alert_grouping_setting."+ref, "description"),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting."+ref, "type", configType),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting."+ref, "config.timeout", fmt.Sprint(config.Timeout)),
					resource.TestCheckResourceAttrSet("pagerduty_alert_grouping_setting."+ref, "services.0"),
				),
			},
		},
	})
}

func TestAccPagerDutyAlertGroupingSetting_serviceNotExist(t *testing.T) {
	ref := fmt.Sprint("tf-", acctest.RandString(5))
	service := fmt.Sprint("tf-", acctest.RandString(5))
	name := fmt.Sprintf("%s grouping", service)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAlertGroupingSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyAlertGroupingSettingServiceNotExist(ref, service, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAlertGroupingSettingExists("pagerduty_alert_grouping_setting."+ref),
					resource.TestCheckResourceAttr("pagerduty_alert_grouping_setting."+ref, "name", name),
					resource.TestCheckResourceAttrSet("pagerduty_alert_grouping_setting."+ref, "description"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAlertGroupingSettingDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_alert_grouping_setting" {
			continue
		}

		ctx := context.Background()

		if _, err := testAccProvider.client.GetAlertGroupingSetting(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("Alert grouping setting still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyAlertGroupingSettingExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alert grouping setting ID is set")
		}

		found, err := testAccProvider.client.GetAlertGroupingSetting(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Alert grouping setting not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func helperConfigPagerDutyAlertGroupingSettingConfig(config interface{}) string {
	switch c := config.(type) {
	case pagerduty.AlertGroupingSettingConfigContentBased:
		timeWindowStr := ""
		if c.TimeWindow != 0 {
			timeWindowStr = fmt.Sprintf("time_window = %d", c.TimeWindow)
		}
		return fmt.Sprintf(`{
			%s
			aggregate = "%s"
			fields = ["%s"]
		}`, timeWindowStr, c.Aggregate, strings.Join(c.Fields, `","`))
	case pagerduty.AlertGroupingSettingConfigIntelligent:
		timeWindowStr := ""
		if c.TimeWindow != 0 {
			timeWindowStr = fmt.Sprintf("time_window = %d", c.TimeWindow)
		}
		return fmt.Sprintf("{\n\t%s\n}", timeWindowStr)
	case pagerduty.AlertGroupingSettingConfigTime:
		timeoutStr := ""
		if c.Timeout != 0 {
			timeoutStr = fmt.Sprintf("timeout = %d", c.Timeout)
		}
		return fmt.Sprintf("{\n%s\n}", timeoutStr)
	}
	return "{}"
}

func testAccCheckPagerDutyAlertGroupingSettingConfig(ref, name, cfgType, serviceName string, cfg interface{}) string {
	config := helperConfigPagerDutyAlertGroupingSettingConfig(cfg)
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "%[4]s" {
	name = "%[4]s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "%[1]s" {
  name = "%[2]s"
  type = "%[3]s"
  services = [pagerduty_service.%[4]s.id]
  config %[5]s
}`, ref, name, cfgType, serviceName, config)
}

func testAccCheckPagerDutyAlertGroupingSettingConfigUpdated(ref, name, cfgType, serviceName string, cfg interface{}) string {
	config := helperConfigPagerDutyAlertGroupingSettingConfig(cfg)
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "%[4]s" {
	name = "%[4]s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}
resource "pagerduty_service" "%[4]s-copy" {
	name = "Copy of %[4]s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "%[1]s" {
  name = "%[2]s"
  type = "%[3]s"
  services = [pagerduty_service.%[4]s.id, pagerduty_service.%[4]s-copy.id]
  config %[5]s
}`, ref, name, cfgType, serviceName, config)
}

func testAccCheckPagerDutyAlertGroupingSettingAppendServiceConfig(name, service1, service2 string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "foo" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}
resource "pagerduty_service" "bar" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "foo" {
	name = "%s"
	type = "content_based"
	config {
		time_window = 1440
		aggregate = "all"
		fields = ["summary"]
	}
	services = [pagerduty_service.foo.id, pagerduty_service.bar.id]
}`, service1, service2, name)
}

func testAccCheckPagerDutyAlertGroupingSettingAppendServiceConfigUpdated(name, service1, service2, service3 string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "foo" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}
resource "pagerduty_service" "bar" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}
resource "pagerduty_service" "qux" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "foo" {
	name = "%s"
	type = "content_based"
	config {
		time_window = 1440
		aggregate = "all"
		fields = ["summary"]
	}
	services = [pagerduty_service.foo.id, pagerduty_service.qux.id, pagerduty_service.bar.id]
}`, service1, service2, service3, name)
}

func testAccCheckPagerDutyAlertGroupingSettingPopServiceConfigUpdated(name, service1 string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
	name = "Default"
}

resource "pagerduty_service" "foo" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.default.id
}

resource "pagerduty_alert_grouping_setting" "foo" {
	name = "%s"
	type = "content_based"
	config {
		time_window = 1440
		aggregate = "all"
		fields = ["summary"]
	}
	services = [pagerduty_service.foo.id]
}`, service1, name)
}

func testAccPagerDutyAlertGroupingSettingServiceNotExist(ref, service, name string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "default" {
  name = "Default"
}

resource "pagerduty_service" "foo" {
  name                    = "%[2]s"
  auto_resolve_timeout    = "null"
  acknowledgement_timeout = 1800
  escalation_policy       = data.pagerduty_escalation_policy.default.id
  alert_creation          = "create_alerts_and_incidents"

  auto_pause_notifications_parameters {
    enabled = true
    timeout = 120
  }

  incident_urgency_rule {
    type    = "constant"
    urgency = "high"
  }
}

resource "pagerduty_alert_grouping_setting" "%[1]s" {
  name     = "%[3]s"
  services = [pagerduty_service.foo.id]
  type     = "intelligent"
  config {
    time_window = 300
  }
}`, ref, service, name)
}
