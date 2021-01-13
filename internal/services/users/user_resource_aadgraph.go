package users

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/helpers/aadgraph"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
	"github.com/terraform-providers/terraform-provider-azuread/internal/utils"
)

func userResourceCreateAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.AadClient

	upn := d.Get("user_principal_name").(string)
	mailNickName := d.Get("mail_nickname").(string)

	// default mail nickname to the first part of the UPN (matches the portal)
	if mailNickName == "" {
		mailNickName = strings.Split(upn, "@")[0]
	}

	userCreateParameters := graphrbac.UserCreateParameters{
		AccountEnabled: utils.Bool(d.Get("account_enabled").(bool)),
		DisplayName:    utils.String(d.Get("display_name").(string)),
		MailNickname:   &mailNickName,
		PasswordProfile: &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                     utils.String(d.Get("password").(string)),
		},
		UserPrincipalName:    &upn,
		AdditionalProperties: map[string]interface{}{},
	}

	if v, ok := d.GetOk("given_name"); ok {
		userCreateParameters.GivenName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("surname"); ok {
		userCreateParameters.Surname = utils.String(v.(string))
	}

	if v, ok := d.GetOk("usage_location"); ok {
		userCreateParameters.UsageLocation = utils.String(v.(string))
	}

	if v, ok := d.GetOk("immutable_id"); ok {
		userCreateParameters.ImmutableID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("job_title"); ok {
		userCreateParameters.AdditionalProperties["jobTitle"] = v.(string)
	}

	if v, ok := d.GetOk("department"); ok {
		userCreateParameters.AdditionalProperties["department"] = v.(string)
	}

	if v, ok := d.GetOk("company_name"); ok {
		userCreateParameters.AdditionalProperties["companyName"] = v.(string)
	}

	if v, ok := d.GetOk("physical_delivery_office_name"); ok {
		userCreateParameters.AdditionalProperties["physicalDeliveryOfficeName"] = v.(string)
	}

	if v, ok := d.GetOk("street_address"); ok {
		userCreateParameters.AdditionalProperties["streetAddress"] = v.(string)
	}

	if v, ok := d.GetOk("city"); ok {
		userCreateParameters.AdditionalProperties["city"] = v.(string)
	}

	if v, ok := d.GetOk("state"); ok {
		userCreateParameters.AdditionalProperties["state"] = v.(string)
	}

	if v, ok := d.GetOk("country"); ok {
		userCreateParameters.AdditionalProperties["country"] = v.(string)
	}

	if v, ok := d.GetOk("postal_code"); ok {
		userCreateParameters.AdditionalProperties["postalCode"] = v.(string)
	}

	if v, ok := d.GetOk("mobile"); ok {
		userCreateParameters.AdditionalProperties["mobile"] = v.(string)
	}

	user, err := client.Create(ctx, userCreateParameters)
	if err != nil {
		return tf.ErrorDiagF(err, "Creating user %q", upn)
	}

	if user.ObjectID == nil || *user.ObjectID == "" {
		return tf.ErrorDiagF(errors.New("API returned group with nil object ID"), "Bad API Response")
	}

	d.SetId(*user.ObjectID)

	_, err = aadgraph.WaitForCreationReplication(ctx, d.Timeout(schema.TimeoutCreate), func() (interface{}, error) {
		return client.Get(ctx, *user.ObjectID)
	})

	if err != nil {
		return tf.ErrorDiagF(err, "Waiting for user %q with object ID: %q", upn, *user.ObjectID)
	}

	return userResourceReadAadGraph(ctx, d, meta)
}

func userResourceUpdateAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.AadClient

	var userUpdateParameters graphrbac.UserUpdateParameters

	if d.HasChange("display_name") {
		userUpdateParameters.DisplayName = utils.String(d.Get("display_name").(string))
	}

	if d.HasChange("given_name") {
		userUpdateParameters.GivenName = utils.String(d.Get("given_name").(string))
	}

	if d.HasChange("surname") {
		userUpdateParameters.Surname = utils.String(d.Get("surname").(string))
	}

	if d.HasChange("mail_nickname") {
		userUpdateParameters.MailNickname = utils.String(d.Get("mail_nickname").(string))
	}

	if d.HasChange("account_enabled") {
		userUpdateParameters.AccountEnabled = utils.Bool(d.Get("account_enabled").(bool))
	}

	if d.HasChange("password") {
		userUpdateParameters.PasswordProfile = &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                     utils.String(d.Get("password").(string)),
		}
	}

	if d.HasChange("usage_location") {
		userUpdateParameters.UsageLocation = utils.String(d.Get("usage_location").(string))
	}

	if d.HasChange("immutable_id") {
		userUpdateParameters.ImmutableID = utils.String(d.Get("immutable_id").(string))
	}

	additionalProperties := map[string]interface{}{}

	if d.HasChange("job_title") {
		additionalProperties["jobTitle"] = d.Get("job_title").(string)
	}

	if d.HasChange("department") {
		additionalProperties["department"] = d.Get("department").(string)
	}

	if d.HasChange("company_name") {
		additionalProperties["companyName"] = d.Get("company_name").(string)
	}

	if d.HasChange("physical_delivery_office_name") {
		additionalProperties["physicalDeliveryOfficeName"] = d.Get("physical_delivery_office_name").(string)
	}

	if d.HasChange("street_address") {
		additionalProperties["streetAddress"] = d.Get("street_address").(string)
	}

	if d.HasChange("city") {
		additionalProperties["city"] = d.Get("city").(string)
	}

	if d.HasChange("state") {
		additionalProperties["state"] = d.Get("state").(string)
	}

	if d.HasChange("country") {
		additionalProperties["country"] = d.Get("country").(string)
	}

	if d.HasChange("postal_code") {
		additionalProperties["postalCode"] = d.Get("postal_code").(string)
	}

	if d.HasChange("mobile") {
		additionalProperties["mobile"] = d.Get("mobile").(string)
	}

	if len(additionalProperties) > 0 {
		userUpdateParameters.AdditionalProperties = additionalProperties
	}

	if _, err := client.Update(ctx, d.Id(), userUpdateParameters); err != nil {
		return tf.ErrorDiagF(err, "Updating User with object ID: %q", d.Id())
	}

	return userResourceReadAadGraph(ctx, d, meta)
}

func userResourceReadAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.AadClient

	objectId := d.Id()

	user, err := client.Get(ctx, objectId)
	if err != nil {
		if utils.ResponseWasNotFound(user.Response) {
			log.Printf("[DEBUG] User with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving user with object ID: %q", objectId)
	}

	if dg := tf.Set(d, "object_id", user.ObjectID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "immutable_id", user.ImmutableID); dg != nil {
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

	return nil
}

func userResourceDeleteAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.AadClient

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return tf.ErrorDiagF(err, "Deleting user with object ID: %q", d.Id())
		}
	}

	return nil
}
