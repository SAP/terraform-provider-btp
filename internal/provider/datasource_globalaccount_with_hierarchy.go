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

var subaccountHierarchySchemaAttributes = map[string]schema.Attribute{
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
		MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
		Computed:            true,
	},
	"last_modified": schema.StringAttribute{
		MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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
		MarkdownDescription: "The name of the subaccount's parent entity. Typically could be the global account or a directory.",
		Computed:            true,
	},
	"parent_type": schema.StringAttribute{
		MarkdownDescription: "The type of the subaccount's parent entity. Typically could be the global account or a directory.",
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
		MarkdownDescription: "The type of resource",
	},
}

func directorySchema (level int) schema.NestedAttributeObject{
	if level > 1 {
		return schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the directory.",
							Required:            true,
							Validators: []validator.String{
								uuidvalidator.ValidUUID(),
							},
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "The details of the user that created the directory.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"directories": schema.ListNestedAttribute{
							NestedObject: directorySchema(level-1),
							MarkdownDescription: "The list of directories contained in this directory",
						},
						"directory_type": schema.StringAttribute{
							MarkdownDescription: "The type of the directory",
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
							MarkdownDescription: "The name of the directory's parent entity.",
						},
						"parent_type": schema.StringAttribute{
							MarkdownDescription: "The type of the directory's parent entity.",
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
								Attributes: subaccountHierarchySchemaAttributes,
							},
							MarkdownDescription: "The subaccounts contained in this directory.",
						},
						"subdomain": schema.StringAttribute{
							MarkdownDescription: "This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of resource.",
							Computed:            true,
						},
			},
		}
	}

	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The details of the user that created the directory.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"directory_type": schema.StringAttribute{
				MarkdownDescription: "The type of the directory",
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
				MarkdownDescription: "The name of the directory's parent entity.",
			},
			"parent_type": schema.StringAttribute{
				MarkdownDescription: "The type of the directory's parent entity.",
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
					Attributes: subaccountHierarchySchemaAttributes,
				},
				MarkdownDescription: "The subaccounts contained in this directory.",
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of resource.",
				Computed:            true,
			},
},
	}
}

func directoriesSchema(level int) schema.NestedAttributeObject{
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the directory.",
						Required:            true,
						Validators: []validator.String{
							uuidvalidator.ValidUUID(),
						},
					},
					"created_by": schema.StringAttribute{
						MarkdownDescription: "The details of the user that created the directory.",
						Computed:            true,
					},
					"created_date": schema.StringAttribute{
						MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
						Computed:            true,
					},
					"directories": schema.ListNestedAttribute{
						NestedObject: directorySchema(level),
						MarkdownDescription: "The list of directories contained in this directory",
					},
					"directory_type": schema.StringAttribute{
						MarkdownDescription: "The type of the directory",
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
						MarkdownDescription: "The name of the directory's parent entity.",
					},
					"parent_type": schema.StringAttribute{
						MarkdownDescription: "The type of the directory's parent entity.",
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
							Attributes: subaccountHierarchySchemaAttributes,
						},
						MarkdownDescription: "The subaccounts contained in this directory.",
					},
					"subdomain": schema.StringAttribute{
						MarkdownDescription: "This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of resource.",
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
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a global account.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"directories": schema.ListNestedAttribute{
				NestedObject: directoriesSchema(5),
				MarkdownDescription: "The directories contained in the global account",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the global account.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region in which the resource is created.",
				Computed:            true,
			},
			"subaccounts": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: subaccountHierarchySchemaAttributes,
				},
				MarkdownDescription: "The subaccounts contained in the globalaccount",
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "The subdomain is part of the path used to access the authorization tenant of the global account.",
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
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the resource.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The user that created the resource.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountWithHierarchyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.GlobalAccount.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Global Account", fmt.Sprintf("%s", err))
		return
	}

	data.ID = types.StringValue(cliRes.Guid)
	data.CommercialModel = types.StringValue(cliRes.CommercialModel)
	data.ConsumptionBased = types.BoolValue(cliRes.ConsumptionBased)
	data.ContractStatus = types.StringValue(cliRes.ContractStatus)
	data.CostObjectId = stringNullIfEmpty(cliRes.CostObjectId)
	data.CostObjectType = stringNullIfEmpty(cliRes.CostObjectType)
	data.CreatedDate = timeToValue(cliRes.CreatedDate.Time())

	data.CrmCustomerId = stringNullIfEmpty(cliRes.CrmCustomerId)
	data.CrmTenantId = stringNullIfEmpty(cliRes.CrmTenantId)

	data.Description = types.StringValue(cliRes.Description)
	data.DisplayName = types.StringValue(cliRes.DisplayName)
	data.ExpiryDate = timeToValue(cliRes.ExpiryDate.Time())
	data.GeoAccess = types.StringValue(cliRes.GeoAccess)
	data.LicenseType = types.StringValue(cliRes.LicenseType)
	data.LastModified = timeToValue(cliRes.ModifiedDate.Time())
	data.State = types.StringValue(cliRes.EntityState)
	data.Origin = types.StringValue(cliRes.Origin)
	data.RenewalDate = timeToValue(cliRes.RenewalDate.Time())

	data.ServiceId = stringNullIfEmpty(cliRes.ServiceId)

	data.Subdomain = types.StringValue(cliRes.Subdomain)
	data.Usage = stringNullIfEmpty(cliRes.UseFor)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

