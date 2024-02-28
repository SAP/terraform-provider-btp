package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

/*
IN CASE OF ANY CHANGES TO THE SCHEMA, PLEASE UPDATE THE DOCUMENTATION IN THE TEMPLATE FILE templates/data-sources/globalaccount_with_hierarchy.md.tmpl
*/
var subaccountSchemaAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of the subaccount.",
		Computed:            true,
		Validators: []validator.String{
			uuidvalidator.ValidUUID(),
		},
	},
	"created_by": schema.StringAttribute{
		MarkdownDescription: "The details of the user that created the subaccount.",
		Computed:            true,
	},
	"created_date": schema.StringAttribute{
		MarkdownDescription: "The date and time when the subaccount was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
		Computed:            true,
	},
	"last_modified": schema.StringAttribute{
		MarkdownDescription: "The date and time when the subaccount was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "A descriptive name of the subaccount for customer-facing UIs.",
		Computed:            true,
	},
	"parent_id": schema.StringAttribute{
		MarkdownDescription: "The ID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.",
		Computed:            true,
	},
	"parent_name": schema.StringAttribute{
		MarkdownDescription: "The name of the subaccount's parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the name of the global account.",
		Computed:            true,
	},
	"parent_type": schema.StringAttribute{
		MarkdownDescription: "The type of the subaccount's parent entity. If the subaccount is located directly in the global account (not in a directory), then this will have the value 'Global Account' .",
		Computed:            true,
	},
	"region": schema.StringAttribute{
		MarkdownDescription: "The region in which the subaccount was created.",
		Computed:            true,
	},
	"state": schema.StringAttribute{
		MarkdownDescription: "The current state of the subaccount. Possible values are: \n" +
			getFormattedValueAsTableRow("state", "description") +
			getFormattedValueAsTableRow("---", "---") +
			getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
			getFormattedValueAsTableRow("`STARTED`", "CRUD operation on an entity has started.") +
			getFormattedValueAsTableRow("`CANCELLED`", "The operation or processing was canceled by the operator.") +
			getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the entity is in progress.") +
			getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
			getFormattedValueAsTableRow("`CREATING`", "Creating entity operation is in progress.") +
			getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the entity was not created or was created but cannot be used.") +
			getFormattedValueAsTableRow("`UPDATING`", "Updating entity operation is in progress.") +
			getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the entity was not updated.") +
			getFormattedValueAsTableRow("`DELETING`", "Deleting entity operation is in progress.") +
			getFormattedValueAsTableRow("`DELETION_FAILED`", "The delete operation failed, and the entity was not deleted.") +
			getFormattedValueAsTableRow("`MOVING`", "Moving entity operation is in progress.") +
			getFormattedValueAsTableRow("`MOVE_FAILED`", "Entity could not be moved to a different location.") +
			getFormattedValueAsTableRow("`PENDING REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
			getFormattedValueAsTableRow("`MIGRATING`", "Migrating entity from Neo to Cloud Foundry."),
		Computed: true,
	},
	"subdomain": schema.StringAttribute{
		MarkdownDescription: "The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.",
		Computed:            true,
	},
	"type": schema.StringAttribute{
		MarkdownDescription: "The type of resource, in this case this will be 'Subaccount'.",
		Computed:            true,
	},
}

/*
This function is designed to retrieve the schema for all directories, accommodating nested directories, up to 5 levels
(according to the btpcli).The schema is constructed recursively, allowing directories to have directories underneath.
All directories share the same parameters, except those at the last level, which cannot contain sub-directories. This
recursive strategy defines the directories with a consistent set of attributes based on the specified level.

IN CASE OF ANY CHANGES TO THE SCHEMA, PLEASE UPDATE THE DOCUMENTATION IN THE TEMPLATE FILE templates/data-sources/globalaccount_with_hierarchy.md.tmpl
*/
func directorySchemaObject(level int) schema.NestedAttributeObject {
	if level > 1 {
		return schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The ID of the directory.",
					Computed:            true,
					Validators: []validator.String{
						uuidvalidator.ValidUUID(),
					},
				},
				"created_by": schema.StringAttribute{
					MarkdownDescription: "The details of the user that created the directory.",
					Computed:            true,
				},
				"created_date": schema.StringAttribute{
					MarkdownDescription: "The date and time when the directory was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
					Computed:            true,
				},
				"directories": schema.ListNestedAttribute{
					NestedObject:        directorySchemaObject(level - 1),
					MarkdownDescription: "The list of directories contained in this directory.",
					Computed:            true,
					Optional:            true,
				},
				"features": schema.SetAttribute{
					ElementType: types.StringType,
					MarkdownDescription: "The features that are enabled for the directory. Possible values are: \n" +
						getFormattedValueAsTableRow("value", "description") +
						getFormattedValueAsTableRow("---", "---") +
						getFormattedValueAsTableRow("`DEFAULT` ", "All directories have the following basic feature enabled:"+
							"<br> 1. Group and filter subaccounts for reports and filters "+
							"<br> 2. Monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model)"+
							"<br> 3. Set custom properties and tags to the directory for identification and reporting purposes.") +
						getFormattedValueAsTableRow("`ENTITLEMENTS`", "Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory.") +
						getFormattedValueAsTableRow("`AUTHORIZATIONS`", "Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the `ENTITLEMENTS` feature."),
					Computed: true,
				},
				"last_modified": schema.StringAttribute{
					MarkdownDescription: "The date and time when the directory was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
					Computed:            true,
				},
				"name": schema.StringAttribute{
					MarkdownDescription: "The display name of the directory.",
					Computed:            true,
				},
				"parent_id": schema.StringAttribute{
					MarkdownDescription: "The ID of the directory's parent entity. Typically this is the global account.",
					Computed:            true,
				},
				"parent_name": schema.StringAttribute{
					MarkdownDescription: "The name of the directory's parent entity. Typically this is the global account.",
					Computed:            true,
				},
				"parent_type": schema.StringAttribute{
					MarkdownDescription: "The type of the directory's parent entity. Typically this will have the value 'Global Account'.",
					Computed:            true,
				},
				"state": schema.StringAttribute{
					MarkdownDescription: "The current state of the directory. Possible values are: \n" +
						getFormattedValueAsTableRow("state", "description") +
						getFormattedValueAsTableRow("---", "---") +
						getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
						getFormattedValueAsTableRow("`STARTED`", "CRUD operation on an entity has started.") +
						getFormattedValueAsTableRow("`CANCELLED`", "The operation or processing was canceled by the operator.") +
						getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the entity is in progress.") +
						getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
						getFormattedValueAsTableRow("`CREATING`", "Creating entity operation is in progress.") +
						getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the entity was not created or was created but cannot be used.") +
						getFormattedValueAsTableRow("`UPDATING`", "Updating entity operation is in progress.") +
						getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the entity was not updated.") +
						getFormattedValueAsTableRow("`DELETING`", "Deleting entity operation is in progress.") +
						getFormattedValueAsTableRow("`DELETION_FAILED`", "The delete operation failed, and the entity was not deleted.") +
						getFormattedValueAsTableRow("`MOVING`", "Moving entity operation is in progress.") +
						getFormattedValueAsTableRow("`MOVE_FAILED`", "Entity could not be moved to a different location.") +
						getFormattedValueAsTableRow("`PENDING REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
						getFormattedValueAsTableRow("`MIGRATING`", "Migrating entity from Neo to Cloud Foundry."),
					Computed: true,
				},
				"subaccounts": schema.ListNestedAttribute{
					NestedObject: schema.NestedAttributeObject{
						Attributes: subaccountSchemaAttributes,
					},
					MarkdownDescription: "The subaccounts contained in this directory.",
					Computed:            true,
					Optional:            true,
				},
				"subdomain": schema.StringAttribute{
					MarkdownDescription: "This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.",
					Computed:            true,
					Optional:            true,
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "The type of resource, in this case will have the value 'Directory'.",
					Computed:            true,
				},
			},
		}
	}

	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The details of the user that created the directory.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the directory was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"features": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "The features that are enabled for the directory. Possible values are: \n" +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`DEFAULT` ", "All directories have the following basic feature enabled:"+
						"<br> 1. Group and filter subaccounts for reports and filters "+
						"<br> 2. Monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model)"+
						"<br> 3. Set custom properties and tags to the directory for identification and reporting purposes.") +
					getFormattedValueAsTableRow("`ENTITLEMENTS`", "Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory.") +
					getFormattedValueAsTableRow("`AUTHORIZATIONS`", "Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the `ENTITLEMENTS` feature."),
				Computed: true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the directory.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory's parent entity. Typically this is the global account.",
				Computed:            true,
			},
			"parent_name": schema.StringAttribute{
				MarkdownDescription: "The name of the directory's parent entity. Typically this is the global account.",
				Computed:            true,
			},
			"parent_type": schema.StringAttribute{
				MarkdownDescription: "The type of the directory's parent entity. Typically this will have the value 'Global Account'.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the directory. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`STARTED`", "CRUD operation on an entity has started.") +
					getFormattedValueAsTableRow("`CANCELLED`", "The operation or processing was canceled by the operator.") +
					getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the entity is in progress.") +
					getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating entity operation is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the entity was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", "Updating entity operation is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the entity was not updated.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting entity operation is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The delete operation failed, and the entity was not deleted.") +
					getFormattedValueAsTableRow("`MOVING`", "Moving entity operation is in progress.") +
					getFormattedValueAsTableRow("`MOVE_FAILED`", "Entity could not be moved to a different location.") +
					getFormattedValueAsTableRow("`PENDING REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
					getFormattedValueAsTableRow("`MIGRATING`", "Migrating entity from Neo to Cloud Foundry."),
				Computed: true,
			},
			"subaccounts": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: subaccountSchemaAttributes,
				},
				MarkdownDescription: "The subaccounts contained in this directory.",
				Computed:            true,
				Optional:            true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.",
				Computed:            true,
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of resource, in this case will have the value 'Directory'.",
				Computed:            true,
			},
		},
	}
}

func newGlobalaccountWithHierarchyDataSource() datasource.DataSource {
	return &globalaccountWithHierarchyDataSource{}
}

type globalaccountWithHierarchyDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountWithHierarchyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_with_hierarchy", req.ProviderTypeName)
}

func (ds *globalaccountWithHierarchyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountWithHierarchyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	/*
		IN CASE OF ANY CHANGES TO THE SCHEMA, PLEASE UPDATE THE DOCUMENTATION IN THE TEMPLATE FILE templates/data-sources/globalaccount_with_hierarchy.md.tmpl
	*/
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a global account's hierarchy structure

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"directories": schema.ListNestedAttribute{
				/*
					The schema accommodates six levels of directories, despite the btpcli limiting it to five. This adjustment is made to address
					a discrepancy between the schema of the last directory level and the structure named 'directoryHierarchyType' (refer to type_directoryHierarchy.go).
					This structure is utilized for mapping values across all directory levels.
					The inclusion of an extra level helps avoid mismatches, in cases where the structure specifies a 'directories' parameter not present
					in the last level of directories. Consequently, the second-to-last level of directories in the provider aligns with the last level of
					directories in the btpcli, preventing mismatch errors.
				*/
				NestedObject:        directorySchemaObject(6),
				MarkdownDescription: "The directories contained in the global account",
				Computed:            true,
				Optional:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the global account.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the global account. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`STARTED`", "CRUD operation on an entity has started.") +
					getFormattedValueAsTableRow("`CANCELLED`", "The operation or processing was canceled by the operator.") +
					getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the entity is in progress.") +
					getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating entity operation is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the entity was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", " Updating entity operation is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the entity was not updated.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting entity operation is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The delete operation failed, and the entity was not deleted.") +
					getFormattedValueAsTableRow("`MOVING`", "Moving entity operation is in progress.") +
					getFormattedValueAsTableRow("`MOVE_FAILED`", "Entity could not be moved to a different location.") +
					getFormattedValueAsTableRow("`PENDING REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
					getFormattedValueAsTableRow("`MIGRATING`", "Migrating entity from Neo to Cloud Foundry."),
				Computed: true,
			},
			"subaccounts": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: subaccountSchemaAttributes,
				},
				MarkdownDescription: "The subaccounts contained in the globalaccount",
				Computed:            true,
				Optional:            true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "The subdomain is part of the path used to access the authorization tenant of the global account.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the resource, in this case this will have value 'Global Account'.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountWithHierarchyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var globalAccountData globalAccountHierarchyType

	diags := req.Config.Get(ctx, &globalAccountData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.GlobalAccount.GetWithHierarchy(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Global Account", fmt.Sprintf("%s", err))
		return
	}

	globalAccountData, diags = globalAccountHierarchyValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &globalAccountData)
	resp.Diagnostics.Append(diags...)
}
