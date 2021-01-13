package users

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/helpers/aadgraph"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
	"github.com/terraform-providers/terraform-provider-azuread/internal/utils"
)

func userDataSourceReadAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.AadClient

	var user graphrbac.User

	if upn, ok := d.Get("user_principal_name").(string); ok && upn != "" {
		resp, err := client.Get(ctx, upn)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return tf.ErrorDiagPathF(nil, "user_principal_name", "User with UPN %q was not found", upn)
			}

			return tf.ErrorDiagF(err, "Retrieving user with UPN: %q", upn)
		}
		user = resp
	} else if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		u, err := aadgraph.UserGetByObjectId(ctx, client, objectId)
		if err != nil {
			return tf.ErrorDiagF(err, "Finding user with object ID: %q", objectId)
		}
		if u == nil {
			return tf.ErrorDiagPathF(nil, "object_id", "User not found with object ID: %q", objectId)
		}
		user = *u
	} else if mailNickname, ok := d.Get("mail_nickname").(string); ok && mailNickname != "" {
		u, err := aadgraph.UserGetByMailNickname(ctx, client, mailNickname)
		if err != nil {
			return tf.ErrorDiagPathF(err, "mail_nickname", "Finding user with email alias: %q", mailNickname)
		}
		if u == nil {
			return tf.ErrorDiagPathF(nil, "mail_nickname", "User not found with email alias: %q", mailNickname)
		}
		user = *u
	} else {
		return tf.ErrorDiagF(nil, "One of `object_id`, `user_principal_name` and `mail_nickname` must be supplied")
	}

	if user.ObjectID == nil {
		return tf.ErrorDiagF(errors.New("Object ID returned for user is nil"), "Bad API response")
	}

	d.SetId(*user.ObjectID)

	if dg := tf.Set(d, "object_id", user.ObjectID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "immutable_id", user.ImmutableID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_immutable_id", user.ImmutableID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_sam_account_name", user.AdditionalProperties["onPremisesSamAccountName"]); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "onpremises_user_principal_name", user.AdditionalProperties["onPremisesUserPrincipalName"]); dg != nil {
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

	jobTitle := ""
	if v, ok := user.AdditionalProperties["jobTitle"]; ok {
		jobTitle = v.(string)
	}
	if dg := tf.Set(d, "job_title", jobTitle); dg != nil {
		return dg
	}

	dept := ""
	if v, ok := user.AdditionalProperties["department"]; ok {
		dept = v.(string)
	}
	if dg := tf.Set(d, "department", dept); dg != nil {
		return dg
	}

	companyName := ""
	if v, ok := user.AdditionalProperties["companyName"]; ok {
		companyName = v.(string)
	}
	if dg := tf.Set(d, "company_name", companyName); dg != nil {
		return dg
	}

	physDelivOfficeName := ""
	if v, ok := user.AdditionalProperties["physicalDeliveryOfficeName"]; ok {
		physDelivOfficeName = v.(string)
	}
	if dg := tf.Set(d, "physical_delivery_office_name", physDelivOfficeName); dg != nil {
		return dg
	}
	if dg := tf.Set(d, "office_location", physDelivOfficeName); dg != nil {
		return dg
	}

	streetAddress := ""
	if v, ok := user.AdditionalProperties["streetAddress"]; ok {
		streetAddress = v.(string)
	}
	if dg := tf.Set(d, "street_address", streetAddress); dg != nil {
		return dg
	}

	city := ""
	if v, ok := user.AdditionalProperties["city"]; ok {
		city = v.(string)
	}
	if dg := tf.Set(d, "city", city); dg != nil {
		return dg
	}

	state := ""
	if v, ok := user.AdditionalProperties["state"]; ok {
		state = v.(string)
	}
	if dg := tf.Set(d, "state", state); dg != nil {
		return dg
	}

	country := ""
	if v, ok := user.AdditionalProperties["country"]; ok {
		country = v.(string)
	}
	if dg := tf.Set(d, "country", country); dg != nil {
		return dg
	}

	postalCode := ""
	if v, ok := user.AdditionalProperties["postalCode"]; ok {
		postalCode = v.(string)
	}
	if dg := tf.Set(d, "postal_code", postalCode); dg != nil {
		return dg
	}

	mobile := ""
	if v, ok := user.AdditionalProperties["mobile"]; ok {
		mobile = v.(string)
	}
	if dg := tf.Set(d, "mobile", mobile); dg != nil {
		return dg
	}
	if dg := tf.Set(d, "mobile_phone", mobile); dg != nil {
		return dg
	}

	return nil
}
