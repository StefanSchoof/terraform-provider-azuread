package applications

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/internal/services/applications/parse"
	"github.com/terraform-providers/terraform-provider-azuread/internal/tf"
	"github.com/terraform-providers/terraform-provider-azuread/internal/validate"
)

func applicationPasswordResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: applicationPasswordResourceCreate,
		ReadContext:   applicationPasswordResourceRead,
		DeleteContext: applicationPasswordResourceDelete,

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PasswordID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"application_object_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"key_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"description"},
				ForceNew:      true, // TODO: ForceNew not required with MS Graph, remove in v2.0
			},

			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"display_name"},
				Deprecated:    "this property has been renamed to `display_name` and will be removed in version 2.0 of the provider",
				ForceNew:      true,
			},

			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 863), // Encrypted secret cannot be empty and can be at most 1024 bytes.
			},

			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"end_date_relative"},
				ValidateFunc: validation.IsRFC3339Time,
			},

			"end_date_relative": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"end_date"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceApplicationPasswordInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceApplicationPasswordInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

// TODO: MS Graph does not seem to support user-supplied passwords for applications. Needs investigation.
// Attempting to add a password to an application returns an error:
// "New password credentials must be generated using service actions"

func applicationPasswordResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
	//	return applicationPasswordResourceCreateMsGraph(ctx, d, meta)
	//}
	return applicationPasswordResourceCreateAadGraph(ctx, d, meta)
}

func applicationPasswordResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
	//	return applicationPasswordResourceReadMsGraph(ctx, d, meta)
	//}
	return applicationPasswordResourceReadAadGraph(ctx, d, meta)
}

func applicationPasswordResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
	//	return applicationPasswordResourceDeleteMsGraph(ctx, d, meta)
	//}
	return applicationPasswordResourceDeleteAadGraph(ctx, d, meta)
}

func resourceApplicationPasswordInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application_object_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"key_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 863),
			},

			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"end_date_relative"},
				ValidateFunc: validation.IsRFC3339Time,
			},

			"end_date_relative": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"end_date"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceApplicationPasswordInstanceStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Println("[DEBUG] Migrating ID from v0 to v1 format")
	newId, err := parse.OldPasswordID(rawState["id"].(string))
	if err != nil {
		return rawState, fmt.Errorf("generating new ID: %s", err)
	}

	rawState["id"] = newId.String()
	return rawState, nil
}
