package users_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/terraform-providers/terraform-provider-azuread/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azuread/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/utils"
)

type UserResource struct{}

func TestAccUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_user", "test")
	r := UserResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("force_password_change", "password"),
	})
}

func TestAccUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_user", "test")
	r := UserResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("force_password_change", "password"),
	})
}

func TestAccUser_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_user", "test")
	r := UserResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("force_password_change", "password"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("force_password_change", "password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("force_password_change", "password"),
	})
}

func TestAccUser_threeUsersABC(t *testing.T) {
	dataA := acceptance.BuildTestData(t, "azuread_user", "testA")
	dataB := acceptance.BuildTestData(t, "azuread_user", "testB")
	dataC := acceptance.BuildTestData(t, "azuread_user", "testC")
	r := UserResource{}

	dataA.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.threeUsersABC(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataB.ResourceName).ExistsInAzure(r),
				check.That(dataC.ResourceName).ExistsInAzure(r),
			),
		},
		dataA.ImportStep("force_password_change", "password"),
		dataB.ImportStep("force_password_change", "password"),
		dataC.ImportStep("force_password_change", "password"),
	})
}

func (r UserResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id *string

	switch clients.EnableMsGraphBeta {
	case true:
		user, status, err := clients.Users.MsClient.Get(ctx, state.ID)
		if err != nil {
			if status == http.StatusNotFound {
				return nil, fmt.Errorf("User with object ID %q does not exist", state.ID)
			}
			return nil, fmt.Errorf("failed to retrieve User with object ID %q: %+v", state.ID, err)
		}
		id = user.ID

	case false:
		resp, err := clients.Users.AadClient.Get(ctx, state.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil, fmt.Errorf("User with object ID %q does not exist", state.ID)
			}
			return nil, fmt.Errorf("failed to retrieve User with object ID %q: %+v", state.ID, err)
		}
		id = resp.ObjectID
	}

	return utils.Bool(id != nil && *id == state.ID), nil
}

func (UserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestUser.%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d"
  password            = "%[2]s"
}
`, data.RandomInteger, data.RandomPassword)
}

func (UserResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name   = "acctestUser.%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  force_password_change = true

  display_name    = "acctestUser-%[1]d-DisplayName"
  given_name      = "acctestUser-%[1]d-GivenName"
  surname         = "acctestUser-%[1]d-Surname"
  mail_nickname   = "acctestUser-%[1]d-MailNickname"
  account_enabled = false
  password        = "%[2]s"
  usage_location  = "NO"
  immutable_id    = "%[1]d"

  job_title      = "acctestUser-%[1]d-Job"
  department     = "acctestUser-%[1]d-Dept"
  company_name   = "acctestUser-%[1]d-Company"
  street_address = "acctestUser-%[1]d-Street"
  state          = "acctestUser-%[1]d-State"
  city           = "acctestUser-%[1]d-City"
  country        = "acctestUser-%[1]d-Country"
  postal_code    = "111111"
  mobile         = "(555) 555-5555"

  physical_delivery_office_name = "acctestUser-%[1]d-PDON"
}
`, data.RandomInteger, data.RandomPassword)
}

func (UserResource) threeUsersABC(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "testA" {
  user_principal_name = "acctestUser.%[1]d.A@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-A"
  password            = "%[2]s"
}

resource "azuread_user" "testB" {
  user_principal_name = "acctestUser.%[1]d.B@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-B"
  mail_nickname       = "acctestUser-%[1]d-B"
  password            = "%[2]s"
}

resource "azuread_user" "testC" {
  user_principal_name = "acctestUser.%[1]d.C@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-C"
  password            = "%[2]s"
}
`, data.RandomInteger, data.RandomPassword)
}
