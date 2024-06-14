package pagerduty

import (
	"context"
	"fmt"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Testing Business Service Dependencies
func TestAccPagerDutyServiceDependency_BusinessBasic(t *testing.T) {
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	businessService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDependencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceDependencyExists("pagerduty_service_dependency.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.foo", "dependency.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.foo", "dependency.0.supporting_service.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.foo", "dependency.0.dependent_service.#", "1"),
				),
			},
			// Validating that externally removed business service dependencies are
			// detected and planned for re-creation
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyServiceDependency("pagerduty_service_dependency.foo", "pagerduty_business_service.foo", "pagerduty_service.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
			// Validating that externally removed dependent service are
			// detected and gracefully handled
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyedDependentService("pagerduty_service_dependency.foo", "pagerduty_business_service.foo", "pagerduty_service.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Testing Parallel creation of Business Service Dependencies
func TestAccPagerDutyServiceDependency_BusinessParallel(t *testing.T) {
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	businessService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	resCount := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDependencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyParallelConfig(service, businessService, username, email, escalationPolicy, resCount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceDependencyParallelExists("pagerduty_service_dependency.foo", resCount),
				),
			},
		},
	})
}

func testAccCheckPagerDutyBusinessServiceDependencyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Relationship ID is set")
		}
		businessService, _ := s.RootModule().Resources["pagerduty_business_service.foo"]

		ctx := context.Background()
		depResp, err := testAccProvider.client.ListBusinessServiceDependenciesWithContext(ctx, businessService.Primary.ID)
		if err != nil {
			return fmt.Errorf("Business Service not found: %v", err)
		}

		// loop serviceRelationships until relationship.IDs match
		var found *pagerduty.ServiceDependency
		for _, rel := range depResp.Relationships {
			if rel.ID == rs.Primary.ID {
				found = rel
				break
			}
		}
		if found == nil {
			return fmt.Errorf("Service Dependency not found: %v", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceDependencyParallelExists(n string, resCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := []*terraform.ResourceState{}
		for i := 0; i < resCount; i++ {
			resName := fmt.Sprintf("%s.%d", n, i)
			r, ok := s.RootModule().Resources[resName]
			if !ok {
				return fmt.Errorf("Not found: %s", resName)
			}
			rs = append(rs, r)
		}

		for _, r := range rs {
			if r.Primary.ID == "" {
				return fmt.Errorf("No Service Relationship ID is set")
			}
		}

		for i := 0; i < resCount; i++ {
			businessService, _ := s.RootModule().Resources["pagerduty_business_service.foo"]

			ctx := context.Background()

			depResp, err := testAccProvider.client.ListBusinessServiceDependenciesWithContext(ctx, businessService.Primary.ID)
			if err != nil {
				return fmt.Errorf("Business Service not found: %v", err)
			}

			// loop serviceRelationships until relationship.IDs match
			var found *pagerduty.ServiceDependency
			for _, rel := range depResp.Relationships {
				if rel.ID == rs[i].Primary.ID {
					found = rel
					break
				}
			}
			if found == nil {
				return fmt.Errorf("Service Dependency not found: %v", rs[i].Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceDependencyDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_service_dependency" {
			continue
		}
		businessService, _ := s.RootModule().Resources["pagerduty_business_service.foo"]

		// get business service
		ctx := context.Background()
		dependencies, err := testAccProvider.client.ListBusinessServiceDependenciesWithContext(ctx, businessService.Primary.ID)
		if err != nil {
			// if the business service doesn't exist, that's okay
			return nil
		}
		// get business service dependencies
		for _, rel := range dependencies.Relationships {
			if rel.ID == r.Primary.ID {
				return fmt.Errorf("supporting service relationship still exists")
			}
		}

	}
	return nil
}

func testAccCheckPagerDutyBusinessServiceDependencyParallelConfig(service, businessService, username, email, escalationPolicy string, resCount int) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%[1]s"
}

resource "pagerduty_user" "foo" {
	name        = "%[2]s"
	email       = "%[3]s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%[4]s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "supportBar" {
	count = %[6]d
	name = "%[5]s-${count.index}"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}

resource "pagerduty_service_dependency" "foo" {
	count = %[6]d
	dependency {
		dependent_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		supporting_service {
			id = pagerduty_service.supportBar[count.index].id
			type = "service"
		}
	}
}
`, businessService, username, email, escalationPolicy, service, resCount)
}

func testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}

resource "pagerduty_service_dependency" "foo" {
	dependency {
		dependent_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		supporting_service {
			id = pagerduty_service.foo.id
			type = "service"
		}
	}
}
`, businessService, username, email, escalationPolicy, service)
}

func testAccExternallyDestroyServiceDependency(resName, depName, suppName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Dependency ID is set for %q", resName)
		}

		dep, ok := s.RootModule().Resources[depName]
		if !ok {
			return fmt.Errorf("Not found: %s", depName)
		}
		if dep.Primary.ID == "" {
			return fmt.Errorf("No Dependent Business Service ID is set for %q", depName)
		}
		depServiceType := dep.Primary.Attributes["type"]

		supp, ok := s.RootModule().Resources[suppName]
		if !ok {
			return fmt.Errorf("Not found: %s", suppName)
		}
		if supp.Primary.ID == "" {
			return fmt.Errorf("No Supporting Service ID is set for %q", suppName)
		}
		suppServiceType := supp.Primary.Attributes["type"]

		var r []*pagerduty.ServiceDependency
		r = append(r, &pagerduty.ServiceDependency{
			ID: rs.Primary.ID,
			DependentService: &pagerduty.ServiceObj{
				ID:   dep.Primary.ID,
				Type: depServiceType,
			},
			SupportingService: &pagerduty.ServiceObj{
				ID:   supp.Primary.ID,
				Type: suppServiceType,
			},
		})

		ctx := context.Background()
		input := pagerduty.ListServiceDependencies{
			Relationships: r,
		}
		_, err := testAccProvider.client.DisassociateServiceDependenciesWithContext(ctx, &input)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccExternallyDestroyedDependentService(resName, depName, suppName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Dependency ID is set for %q", resName)
		}

		dep, ok := s.RootModule().Resources[depName]
		if !ok {
			return fmt.Errorf("Not found: %s", depName)
		}
		if dep.Primary.ID == "" {
			return fmt.Errorf("No Dependent Business Service ID is set for %q", depName)
		}
		depServiceType := dep.Primary.Attributes["type"]

		supp, ok := s.RootModule().Resources[suppName]
		if !ok {
			return fmt.Errorf("Not found: %s", suppName)
		}
		if supp.Primary.ID == "" {
			return fmt.Errorf("No Supporting Service ID is set for %q", suppName)
		}

		ctx := context.Background()
		if depServiceType == "business_service" {
			err := testAccProvider.client.DeleteBusinessServiceWithContext(ctx, dep.Primary.ID)
			if err != nil {
				return err
			}
		} else {
			err := testAccProvider.client.DeleteServiceWithContext(ctx, dep.Primary.ID)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Testing Technical Service Dependencies
func TestAccPagerDutyServiceDependency_TechnicalBasic(t *testing.T) {
	dependentService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	supportingService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTechnicalServiceDependencyDestroy("pagerduty_service.supportBar"),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTechnicalServiceDependencyConfig(dependentService, supportingService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTechnicalServiceDependencyExists("pagerduty_service_dependency.bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.bar", "dependency.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.bar", "dependency.0.supporting_service.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_dependency.bar", "dependency.0.dependent_service.#", "1"),
				),
			},
			// Validating that externally removed technical service dependencies are
			// detected and planned for re-creation
			{
				Config: testAccCheckPagerDutyTechnicalServiceDependencyConfig(dependentService, supportingService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyServiceDependency("pagerduty_service_dependency.bar", "pagerduty_service.dependBar", "pagerduty_service.supportBar"),
				),
				ExpectNonEmptyPlan: true,
			},
			// Validating that externally removed dependent service are
			// detected and gracefully handled
			{
				Config: testAccCheckPagerDutyTechnicalServiceDependencyConfig(dependentService, supportingService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyedDependentService("pagerduty_service_dependency.bar", "pagerduty_service.dependBar", "pagerduty_service.supportBar"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Testing Parallel creation of Technical Service Dependencies
func TestAccPagerDutyServiceDependency_TechnicalParallel(t *testing.T) {
	dependentService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	supportingService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	resCount := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTechnicalServiceDependencyParallelDestroy("pagerduty_service.supportBar", resCount),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTechnicalServiceDependencyParallelConfig(dependentService, supportingService, username, email, escalationPolicy, resCount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTechnicalServiceDependencyParallelExists("pagerduty_service_dependency.bar", resCount),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTechnicalServiceDependencyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Relationship ID is set")
		}
		supportService, _ := s.RootModule().Resources["pagerduty_service.supportBar"]

		ctx := context.Background()
		depResp, err := testAccProvider.client.ListTechnicalServiceDependenciesWithContext(ctx, supportService.Primary.ID)
		if err != nil {
			return fmt.Errorf("Technical Service not found: %v", err)
		}
		var foundRel *pagerduty.ServiceDependency

		// loop serviceRelationships until relationship.IDs match
		for _, rel := range depResp.Relationships {
			if rel.ID == rs.Primary.ID {
				foundRel = rel
				break
			}
		}
		if foundRel == nil {
			return fmt.Errorf("Service Dependency not found: %v", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyTechnicalServiceDependencyParallelExists(n string, resCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := []*terraform.ResourceState{}
		for i := 0; i < resCount; i++ {
			resName := fmt.Sprintf("%s.%d", n, i)
			r, ok := s.RootModule().Resources[resName]
			if !ok {
				return fmt.Errorf("Not found: %s", resName)
			}
			rs = append(rs, r)
		}

		for _, r := range rs {
			if r.Primary.ID == "" {
				return fmt.Errorf("No Service Relationship ID is set")
			}
		}

		for i := 0; i < resCount; i++ {
			resName := fmt.Sprintf("pagerduty_service.supportBar.%d", i)
			supportService, _ := s.RootModule().Resources[resName]

			ctx := context.Background()
			depResp, err := testAccProvider.client.ListTechnicalServiceDependenciesWithContext(ctx, supportService.Primary.ID)
			if err != nil {
				return fmt.Errorf("Technical Service not found: %v", err)
			}
			var foundRel *pagerduty.ServiceDependency

			// loop serviceRelationships until relationship.IDs match
			for _, rel := range depResp.Relationships {
				if rel.ID == rs[i].Primary.ID {
					foundRel = rel
					break
				}
			}
			if foundRel == nil {
				return fmt.Errorf("Service Dependency not found: %v", rs[i].Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckPagerDutyTechnicalServiceDependencyParallelDestroy(n string, resCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for i := 0; i < resCount; i++ {
			if err := testAccCheckPagerDutyTechnicalServiceDependencyDestroy(fmt.Sprintf("%s.%d", n, i))(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckPagerDutyTechnicalServiceDependencyDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type != "pagerduty_service_dependency" {
				continue
			}
			supportService, _ := s.RootModule().Resources[n]

			// get service dependencies
			ctx := context.Background()
			dependencies, err := testAccProvider.client.ListTechnicalServiceDependenciesWithContext(ctx, supportService.Primary.ID)
			if err != nil {
				// if the dependency doesn't exist, that's okay
				return nil
			}
			// find desired dependency
			for _, rel := range dependencies.Relationships {
				if rel.ID == r.Primary.ID {
					return fmt.Errorf("supporting service relationship still exists")
				}
			}

		}
		return nil
	}
}

func testAccCheckPagerDutyTechnicalServiceDependencyConfig(dependentService, supportingService, username, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "bar" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}
resource "pagerduty_escalation_policy" "bar" {
	name        = "%s"
	description = "bar-desc"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.bar.id
		}
	}
}
resource "pagerduty_service" "supportBar" {
	name = "%s"
	description             = "supportBarDesc"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "dependBar" {
	name = "%s"
	description             = "dependBarDesc"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service_dependency" "bar" {
	dependency {
		dependent_service {
			id = pagerduty_service.dependBar.id
			type = "service"
		}
		supporting_service {
			id = pagerduty_service.supportBar.id
			type = "service"
		}
	}
}
`, username, email, escalationPolicy, supportingService, dependentService)
}

func testAccCheckPagerDutyTechnicalServiceDependencyParallelConfig(dependentService, supportingService, username, email, escalationPolicy string, resCount int) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "bar" {
	name        = "%[1]s"
	email       = "%[2]s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}
resource "pagerduty_escalation_policy" "bar" {
	name        = "%[3]s"
	description = "bar-desc"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.bar.id
		}
	}
}
resource "pagerduty_service" "supportBar" {
	count                   = %[6]d
	name                    = "%[4]s-${count.index}"
	description             = "supportBarDesc"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "dependBar" {
	name                    = "%[5]s"
	description             = "dependBarDesc"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service_dependency" "bar" {
	count = %[6]d
	dependency {
		dependent_service {
			id   = pagerduty_service.dependBar.id
			type = "service"
		}
		supporting_service {
			id   = pagerduty_service.supportBar[count.index].id
			type = "service"
		}
	}
}
`, username, email, escalationPolicy, supportingService, dependentService, resCount)
}
