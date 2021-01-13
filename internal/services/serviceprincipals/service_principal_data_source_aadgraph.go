package serviceprincipals

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/internal/clients"
	"github.com/terraform-providers/terraform-provider-azuread/internal/helpers/aadgraph"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
	"github.com/terraform-providers/terraform-provider-azuread/internal/utils"
)

func servicePrincipalDataSourceReadAadGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.AadClient

	var servicePrincipal *graphrbac.ServicePrincipal

	if v, ok := d.GetOk("object_id"); ok {
		// use the object_id to find the Azure AD service principal
		objectId := v.(string)
		sp, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(sp.Response) {
				return tf.ErrorDiagPathF(nil, "object_id", "Service Principal with object ID %q was not found", objectId)
			}

			return tf.ErrorDiagPathF(err, "object_id", "Retrieving service principal with object ID %q", objectId)
		}

		servicePrincipal = &sp
	} else if _, ok := d.GetOk("display_name"); ok {
		// use the display_name to find the Azure AD service principal
		displayName := d.Get("display_name").(string)
		filter := fmt.Sprintf("displayName eq '%s'", displayName)

		sps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return tf.ErrorDiagF(err, "Listing service principals for filter %q", filter)
		}

		for _, sp := range *sps.Response().Value {
			if sp.DisplayName == nil {
				continue
			}

			if *sp.DisplayName == displayName {
				servicePrincipal = &sp
				break
			}
		}

		if servicePrincipal == nil {
			return tf.ErrorDiagF(nil, "No service principal found matching display name: %q", displayName)
		}
	} else {
		// use the application_id to find the Azure AD service principal
		applicationId := d.Get("application_id").(string)
		filter := fmt.Sprintf("appId eq '%s'", applicationId)

		sps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return tf.ErrorDiagF(err, "Listing service principals for filter %q", filter)
		}

		for _, sp := range *sps.Response().Value {
			if sp.AppID == nil {
				continue
			}

			if *sp.AppID == applicationId {
				servicePrincipal = &sp
				break
			}
		}

		if servicePrincipal == nil {
			return tf.ErrorDiagF(nil, "No service principal found for application ID: %q", applicationId)
		}
	}

	if servicePrincipal.ObjectID == nil {
		return tf.ErrorDiagF(errors.New("ObjectID returned for service principal is nil"), "Bad API response")
	}

	d.SetId(*servicePrincipal.ObjectID)

	if dg := tf.Set(d, "object_id", servicePrincipal.ObjectID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "application_id", servicePrincipal.AppID); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "display_name", servicePrincipal.DisplayName); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "app_roles", aadgraph.FlattenAppRoles(servicePrincipal.AppRoles)); dg != nil {
		return dg
	}

	if dg := tf.Set(d, "oauth2_permissions", aadgraph.FlattenOauth2Permissions(servicePrincipal.Oauth2Permissions)); dg != nil {
		return dg
	}

	return nil
}
