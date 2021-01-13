package applications

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/helpers/aadgraph"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
	"github.com/terraform-providers/terraform-provider-azuread/internal/utils"
)

func applicationDataSourceReadAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.AadClient

	var app graphrbac.Application

	if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		resp, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return tf.ErrorDiagPathF(nil, "object_id", "Application with object ID %q was not found", objectId)
			}

			return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving Application with object ID %q", objectId)
		}

		app = resp
	} else {
		var fieldName, fieldValue string
		if applicationId, ok := d.Get("application_id").(string); ok && applicationId != "" {
			fieldName = "appId"
			fieldValue = applicationId
		} else if displayName, ok := d.Get("display_name").(string); ok && displayName != "" {
			fieldName = "displayName"
			fieldValue = displayName
		} else if name, ok := d.Get("name").(string); ok && name != "" {
			fieldName = "displayName"
			fieldValue = name
		} else {
			return tf.ErrorDiagF(nil, "One of `object_id`, `application_id` or `displayName` must be specified")
		}

		filter := fmt.Sprintf("%s eq '%s'", fieldName, fieldValue)

		resp, err := client.ListComplete(ctx, filter)
		if err != nil {
			return tf.ErrorDiagF(err, "Listing applications for filter %q", filter)
		}

		values := resp.Response().Value
		if values == nil {
			return tf.ErrorDiagF(fmt.Errorf("nil values for applications matching filter: %q", filter), "Bad API response")
		}
		if len(*values) == 0 {
			return tf.ErrorDiagF(fmt.Errorf("No applications found matching filter: %q", filter), "Application not found")
		}
		if len(*values) > 1 {
			return tf.ErrorDiagF(fmt.Errorf("Found multiple applications matching filter: %q", filter), "Multiple applications found")
		}

		app = (*values)[0]
		switch fieldName {
		case "appId":
			if app.AppID == nil {
				return tf.ErrorDiagF(fmt.Errorf("nil AppID for applications matching filter: %q", filter), "Bad API Response")
			}
			if *app.AppID != fieldValue {
				return tf.ErrorDiagF(fmt.Errorf("AppID does not match (%q != %q) for applications matching filter: %q", *app.AppID, fieldValue, filter), "Bad API Response")
			}
		case "displayName":
			if app.DisplayName == nil {
				return tf.ErrorDiagF(fmt.Errorf("nil displayName for applications matching filter: %q", filter), "Bad API Response")
			}
			if *app.DisplayName != fieldValue {
				return tf.ErrorDiagF(fmt.Errorf("DisplayName does not match (%q != %q) for applications matching filter: %q", *app.DisplayName, fieldValue, filter), "Bad API Response")
			}
		}
	}

	if app.ObjectID == nil {
		return tf.ErrorDiagF(fmt.Errorf("ObjectID returned for application is nil"), "Bad API Response")
	}

	d.SetId(*app.ObjectID)

	if dg := tf.Set(d, "object_id", app.ObjectID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "application_id", app.AppID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "name", app.DisplayName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "homepage", app.Homepage); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "logout_url", app.LogoutURL); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "available_to_other_tenants", app.AvailableToOtherTenants); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "oauth2_allow_implicit_flow", app.Oauth2AllowImplicitFlow); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "identifier_uris", tf.FlattenStringSlicePtr(app.IdentifierUris)); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "reply_urls", tf.FlattenStringSlicePtr(app.ReplyUrls)); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "required_resource_access", flattenApplicationRequiredResourceAccessAad(app.RequiredResourceAccess)); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "optional_claims", flattenApplicationOptionalClaimsAad(app.OptionalClaims)); dg != nil {
		return dg
	}

	var appType string
	if v := app.PublicClient; v != nil && *v {
		appType = "native"
	} else {
		appType = "webapp/api"
	}

	if dg := tf.Set(d, "type", appType); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "app_roles", aadgraph.FlattenAppRoles(app.AppRoles)); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "group_membership_claims", app.GroupMembershipClaims); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "oauth2_permissions", aadgraph.FlattenOauth2Permissions(app.Oauth2Permissions)); dg != nil {
		return dg
	}

	owners, err := aadgraph.ApplicationAllOwners(ctx, client, d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "owners", "Could not retrieve owners for application with object ID %q", *app.ObjectID)
	}
	if dg := tf.Set(d, "owners", owners); dg != nil {
		return dg
	}

	return nil
}
