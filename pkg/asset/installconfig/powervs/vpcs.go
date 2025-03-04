package powervs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func GetUserSelectedPermittedNetwork(zoneID string, dnsCRN crn.CRN, region string) (string, error) {
	var vpcChoice string
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	permittedNetworkCRNs, err := client.GetDNSInstancePermittedNetworks(context.TODO(), dnsCRN.ServiceInstance, zoneID)
	if err != nil {
		return "", err
	}
	allVPCs := []vpcv1.VPC{}
	var vpcs *vpcv1.VPCCollection
	var detailedResponse *core.DetailedResponse
	if vpcs, detailedResponse, err = client.vpcAPI.ListVpcs(client.vpcAPI.NewListVpcsOptions()); err != nil {
		if detailedResponse.GetStatusCode() != http.StatusNotFound {
			return "", err
		}
	} else if vpcs != nil {
		allVPCs = append(allVPCs, vpcs.Vpcs...)
	}
	var vpcNames []string
	for _, vpc := range allVPCs {
		for _, vpcCRN := range permittedNetworkCRNs {
			if strings.Compare(vpcCRN, *vpc.CRN) == 0 {
				vpcNames = append(vpcNames, *vpc.Name)
			}
		}
	}
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Permitted Network",
				Help:    "The VPC of the cluster. If you don't see your intended VPC listed, add the VPC to Permitted Networks of the DNS Zone and rerun the installer.",
				Default: "",
				Options: vpcNames,
			},
		},
	}, &vpcChoice)
	if err != nil {
		return "", fmt.Errorf("survey.ask failed with: %w", err)
	}
	return vpcChoice, nil
}
