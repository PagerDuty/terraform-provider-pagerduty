package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyLicenseAllocation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyLicenseAllocationCreate,
		Read:   resourcePagerDutyLicenseAllocationRead,
		Update: resourcePagerDutyLicenseAllocationUpdate,
		Delete: resourcePagerDutyLicenseAllocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"license_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePagerDutyTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)
	role := d.Get("role").(string)

	log.Printf("[DEBUG] creating user", userID, teamID, role)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.Teams.AddUserWithRole(teamID, userID, role); err != nil {
			if isErrCode(err, 500) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, teamID))

	return fetchPagerDutyTeamMembershipWithRetries(d, meta, genError, 0, d.Get("role").(string))
}

func fetchPagerDutyLicenseAllocation(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID, licenseID := resourcePagerDutyParseColonCompoundID(d.Id())
	role := d.Get("role").(string)

	log.Printf("[DEBUG] Fetching license and role for user: %s", userID)
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		user, err := client.Users.GetWithLicense(userID)
		if err != nil {
			errResp := errCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}
			return nil
		}

		if user.License.ID == licenseID && user.Role == role {
			d.Set("user_id", userID)
			d.Set("license_id", licenseID)
			d.Set("role", role)
		} else {
			log.Printf("[WARN] config for user: %s, license: %s and role: %s does not match fetched license: %s or role: %s", userID, licenseID, role, user.License.ID, user.Role)
		}
		return nil
	})
}

func resourcePagerDutyLicenseAllocationRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyLicenseAllocation(d, meta, handleNotFoundError)
}

func resourcePagerDutyLicenseAllocationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)
	licenseID := d.Get("team_id").(string)
	role := d.Get("role").(string)

	log.Printf("[DEBUG] Updating user: %s with license: %s and role: %s ", userID, licenseID, role)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		licenseReference := &pagerduty.LicenseReference{ID: licenseID, Type: "license_reference"}
		newUser := &pagerduty.User{ID: userID, Role: role, License: licenseReference}
		if _, _, err := client.Users.Update(userID, newUser); err != nil {
			if isErrCode(err, 500) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, licenseID))

	return fetchPagerDutyLicenseAllocation(d, meta, genError)
}

func resourcePagerDutyLicenseAllocationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Users must have a license association, to delete the User you must delete the user resource")

	d.SetId("")

	return nil
}
