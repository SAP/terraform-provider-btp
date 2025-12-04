package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/destinations"
)

func newConnectivityDestinationCertificatesFacade(cliClient *v2Client) connectivityDestinationCertificatesFacade {
	return connectivityDestinationCertificatesFacade{cliClient: cliClient}
}

type connectivityDestinationCertificatesFacade struct {
	cliClient *v2Client
}

func (f *connectivityDestinationCertificatesFacade) getCommand() string {
	return "connectivity/destination-certificate"
}

type FileInput struct {
	Filename           string `btpcli:"filename"`
	CertificateContent string `btpcli:"value"`
}

type DestinationCertificateGetInput struct {
	SubaccountId      string `btpcli:"subaccount"`
	ServiceInstanceId string `btpcli:"serviceInstance,omitempty"`
	CertificateName   string `btpcli:"certName"`
}

func (f *connectivityDestinationCertificatesFacade) Get(ctx context.Context, args *DestinationCertificateGetInput) (destinations.DestinationCertificateResponseObject, CommandResponse, error) {

	params := map[string]string{
		"subaccount": args.SubaccountId,
		"certName":   args.CertificateName,
	}

	if len(args.ServiceInstanceId) > 0 {
		params["serviceInstance"] = args.ServiceInstanceId
	}

	return doExecute[destinations.DestinationCertificateResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))
}

func (f *connectivityDestinationCertificatesFacade) Delete(ctx context.Context, args *DestinationCertificateGetInput) (CommandResponse, error) {
	params := map[string]string{
		"subaccount": args.SubaccountId,
		"certName":   args.CertificateName,
	}

	if len(args.ServiceInstanceId) > 0 {
		params["serviceInstance"] = args.ServiceInstanceId
	}

	return f.cliClient.Execute(ctx, NewDeleteRequest(f.getCommand(), params))
}

func (f *connectivityDestinationCertificatesFacade) List(ctx context.Context, subaccountId string, serviceInstanceId string) (map[string][]destinations.DestinationCertificateResponseObject, CommandResponse, error) {

	res := map[string][]destinations.DestinationCertificateResponseObject{}

	params := map[string]string{
		"subaccount": subaccountId,
		"namesOnly":  "false",
	}

	subaccountCerts, _, err := doExecute[[]destinations.DestinationCertificateResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
	if err != nil {
		return nil, CommandResponse{}, err
	}

	res["subaccount"] = subaccountCerts

	if len(serviceInstanceId) > 0 {
		params["serviceInstance"] = serviceInstanceId

		serviceInstanceCerts, _, err := doExecute[[]destinations.DestinationCertificateResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
		if err != nil {
			return nil, CommandResponse{}, err
		}

		res["serviceInstance"] = serviceInstanceCerts
	}

	return res, CommandResponse{}, nil
}
