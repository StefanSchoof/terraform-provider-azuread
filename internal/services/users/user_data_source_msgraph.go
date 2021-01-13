package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/models"

	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
)

func userDataSourceReadMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.MsClient

	var user models.User

	if upn, ok := d.Get("user_principal_name").(string); ok && upn != "" {
		filter := fmt.Sprintf("userPrincipalName eq '%s'", upn)
		users, _, err := client.List(ctx, filter)
		if err != nil {
			return tf.ErrorDiagF(err, "Finding user with UPN: %q", upn)
		}
		if users == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}

		count := len(*users)
		if count > 1 {
			return tf.ErrorDiagPathF(nil, "user_principal_name", "More than one user found with UPN: %q", upn)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "user_principal_name", "User with UPN %q was not found", upn)
		}

		user = (*users)[0]
	} else if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		u, status, err := client.Get(ctx, objectId)
		if err != nil {
			if status == http.StatusNotFound {
				return tf.ErrorDiagPathF(nil, "object_id", "User not found with object ID: %q", objectId)
			}
			return tf.ErrorDiagF(err, "Retrieving user with object ID: %q", objectId)
		}
		if u == nil {
			return tf.ErrorDiagPathF(nil, "object_id", "User not found with object ID: %q", objectId)
		}
		user = *u
	} else if mailNickname, ok := d.Get("mail_nickname").(string); ok && mailNickname != "" {
		filter := fmt.Sprintf("mailNickname eq '%s'", mailNickname)
		users, _, err := client.List(ctx, filter)
		if err != nil {
			return tf.ErrorDiagF(err, "Finding user with email alias: %q", mailNickname)
		}
		if users == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}

		count := len(*users)
		if count > 1 {
			return tf.ErrorDiagPathF(nil, "mail_nickname", "More than one user found with email alias: %q", upn)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "mail_nickname", "User not found with email alias: %q", upn)
		}

		user = (*users)[0]
	} else {
		return tf.ErrorDiagF(nil, "One of `object_id`, `user_principal_name` or `mail_nickname` must be supplied")
	}

	if user.ID == nil {
		return tf.ErrorDiagF(errors.New("API returned user with nil object ID"), "Bad API Response")
	}

	d.SetId(*user.ID)

	if dg := tf.Set(d, "object_id", user.ID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "immutable_id", user.OnPremisesImmutableId); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_immutable_id", user.OnPremisesImmutableId); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_sam_account_name", user.OnPremisesSamAccountName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_user_principal_name", user.OnPremisesUserPrincipalName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "user_principal_name", user.UserPrincipalName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "account_enabled", user.AccountEnabled); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "display_name", user.DisplayName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "given_name", user.GivenName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "surname", user.Surname); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "mail", user.Mail); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "mail_nickname", user.MailNickname); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "usage_location", user.UsageLocation); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "job_title", user.JobTitle); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "department", user.Department); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "company_name", user.CompanyName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "office_location", user.OfficeLocation); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "physical_delivery_office_name", user.OfficeLocation); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "street_address", user.StreetAddress); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "city", user.City); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "state", user.State); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "country", user.Country); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "postal_code", user.PostalCode); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "mobile", user.MobilePhone); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "mobile_phone", user.MobilePhone); dg != nil {
		return dg
	}

	return nil
}
