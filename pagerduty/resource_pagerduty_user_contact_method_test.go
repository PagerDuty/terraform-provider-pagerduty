package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPagerDutyUserContactMethodEmail_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	usernameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	emailUpdated := fmt.Sprintf("%s@foo.test", usernameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserContactMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserContactMethodEmailConfig(username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserContactMethodEmailConfigUpdated(usernameUpdated, emailUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
		},
	})
}

func TestAccPagerDutyUserContactMethodPhone_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	usernameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	emailUpdated := fmt.Sprintf("%s@foo.test", usernameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserContactMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserContactMethodPhoneConfig(username, email, "4153013250"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
			{
				Config:      testAccCheckPagerDutyUserContactMethodPhoneConfig(username, email, "04153013250"),
				ExpectError: regexp.MustCompile("phone numbers starting with a 0 are not supported"),
			},
			{
				Config: testAccCheckPagerDutyUserContactMethodPhoneConfig(usernameUpdated, emailUpdated, "8019351337"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
		},
	})
}

func TestAccPagerDutyUserContactMethodPhone_EnforceUpdateIfAlreadyExist(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	phoneNumber := "4153013250"
	newPhoneNumber := "4153013251"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserContactMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserContactMethodPhoneConfig(username, email, phoneNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
					testAccAddPhoneContactOutsideTerraform("pagerduty_user_contact_method.foo", newPhoneNumber),
				),
			},
			{
				Config: testAccCheckPagerDutyUserContactMethodPhoneConfig(username, email, newPhoneNumber),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_user_contact_method.foo", "label", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user_contact_method.foo", "address", newPhoneNumber),
				),
			},
		},
	})
}

func TestAccPagerDutyUserContactMethodSMS_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	usernameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	emailUpdated := fmt.Sprintf("%s@foo.test", usernameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserContactMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserContactMethodSMSConfig(username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserContactMethodSMSConfigUpdated(usernameUpdated, emailUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserContactMethodExists("pagerduty_user_contact_method.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyUserContactMethodDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_user_contact_method" {
			continue
		}

		if _, _, err := client.Users.GetContactMethod(r.Primary.Attributes["user_id"], r.Primary.ID); err == nil {
			return fmt.Errorf("User contact method still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyUserContactMethodExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user contact method ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Users.GetContactMethod(rs.Primary.Attributes["user_id"], rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Contact method not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccAddPhoneContactOutsideTerraform(n, p string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		resID := rs.Primary.ID

		if resID == "" {
			return fmt.Errorf("No User Contact Method ID is set")
		}
		userID := rs.Primary.Attributes["user_id"]

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Users.GetContactMethod(userID, rs.Primary.ID)
		if err != nil {
			return err
		}

		found.Address = p
		_, _, err = client.Users.CreateContactMethod(userID, found)
		if err != nil {
			return fmt.Errorf("was not possible to set phone %s contact number outside Terraform state: %v", p, err)
		}

		return nil
	}
}

func testAccCheckPagerDutyUserContactMethodEmailConfig(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id = pagerduty_user.foo.id
  type    = "email_contact_method"
  address = "%[1]v%[2]v"
  label   = "%[1]v"
}
`, username, email)
}

func testAccCheckPagerDutyUserContactMethodEmailConfigUpdated(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id = pagerduty_user.foo.id
  type    = "email_contact_method"
  address = "%[1]v%[2]v"
  label   = "%[1]v"
}
`, username, email)
}

func testAccCheckPagerDutyUserContactMethodPhoneConfig(username, email, phone string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id      = pagerduty_user.foo.id
  type         = "phone_contact_method"
  country_code = "+1"
  address      = "%[3]s"
  label        = "%[1]v"
}
`, username, email, phone)
}

func testAccCheckPagerDutyUserContactMethodSMSConfig(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id      = pagerduty_user.foo.id
  type         = "sms_contact_method"
  country_code = "+1"
  address      = "8458003889"
  label        = "%[1]v"
}
`, username, email)
}

func testAccCheckPagerDutyUserContactMethodSMSConfigUpdated(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id      = pagerduty_user.foo.id
  type         = "sms_contact_method"
  country_code = "+1"
  address      = "6509892965"
  label        = "%[1]v"
}
`, username, email)
}
