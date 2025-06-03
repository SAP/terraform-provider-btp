package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newAccountsDirectoryFacade(cliClient *v2Client) accountsDirectoryFacade {
	return accountsDirectoryFacade{cliClient: cliClient}
}

type accountsDirectoryFacade struct {
	cliClient *v2Client
}

func (f *accountsDirectoryFacade) getCommand() string {
	return "accounts/directory"
}

func (f *accountsDirectoryFacade) Get(ctx context.Context, directoryId string, adminDirectoryId string) (cis.DirectoryResponseObject, CommandResponse, error) {
	if adminDirectoryId == "" {
		return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
			"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
			"directoryID":   directoryId,
		}))

	} else {
		return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
			"globalAccount":  f.cliClient.GetGlobalAccountSubdomain(),
			"directoryID":    directoryId,
			"adminDirectory": adminDirectoryId,
		}))
	}
}

type DirectoryCreateInput struct {
	DisplayName      string              `btpcli:"displayName"`
	Description      *string             `btpcli:"description"`
	ParentID         *string             `btpcli:"parentID"`
	Subdomain        *string             `btpcli:"subdomain"`
	Labels           map[string][]string `btpcli:"labels"`
	Globalaccount    string              `btpcli:"globalAccount"`
	Features         []string            `btpcli:"directoryFeatures"`
	AdminDirectoryId string              `btpcli:"adminDirectory"`
}

type DirectoryUpdateInput struct {
	DirectoryId      string              `btpcli:"directoryID"`
	Globalaccount    string              `btpcli:"globalAccount"`
	DisplayName      *string             `btpcli:"displayName"`
	Description      *string             `btpcli:"description"`
	Labels           map[string][]string `btpcli:"labels"`
	AdminDirectoryId string              `btpcli:"adminDirectory"`
}

type DirectoryEnableInput struct {
	DirectoryId   string   `btpcli:"directoryID"`
	Globalaccount string   `btpcli:"globalAccount"`
	Subdomain     *string  `btpcli:"subdomain"`
	Features      []string `btpcli:"directoryFeatures"`
	Confirm       bool     `btpcli:"confirm"`
}

func (f *accountsDirectoryFacade) Create(ctx context.Context, args *DirectoryCreateInput) (cis.DirectoryResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.DirectoryResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsDirectoryFacade) Update(ctx context.Context, args *DirectoryUpdateInput) (cis.DirectoryResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.DirectoryResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *accountsDirectoryFacade) Enable(ctx context.Context, args *DirectoryEnableInput) (cis.DirectoryResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()
	args.Confirm = true

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.DirectoryResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewEnableRequest(f.getCommand(), params))
}

func (f *accountsDirectoryFacade) Delete(ctx context.Context, directoryId string, adminDirectoryId string) (cis.DirectoryResponseObject, CommandResponse, error) {
	if adminDirectoryId == "" {
		return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
			"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
			"directoryID":   directoryId,
			"forceDelete":   "true",
			"confirm":       "true",
		}))
	} else {
		return doExecute[cis.DirectoryResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
			"globalAccount":  f.cliClient.GetGlobalAccountSubdomain(),
			"directoryID":    directoryId,
			"forceDelete":    "true",
			"confirm":        "true",
			"adminDirectory": adminDirectoryId,
		}))
	}

}
