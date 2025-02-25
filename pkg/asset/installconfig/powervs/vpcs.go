package powervs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func GetPermittedNetwork(zoneID string, dnsCRN crn.CRN, region string) (string, error) {
	var vpcChoice string
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	options, err := client.GetDNSInstancePermittedNetworks(context.TODO(), dnsCRN.ServiceInstance, zoneID)
	if err != nil {
		return "", err
	}
	allVPCs := []vpcv1.VPC{}
	if vpcs, detailedResponse, vpcErr := client.vpcAPI.ListVpcs(client.vpcAPI.NewListVpcsOptions()); vpcErr != nil {
		if detailedResponse.GetStatusCode() != http.StatusNotFound {
			return "", vpcErr
		}
	} else if vpcs != nil {
		allVPCs = append(allVPCs, vpcs.Vpcs...)
	}
	var vpcNamesList []string
	for _, vpc := range allVPCs {
		for _, vpcCRN := range options {
			if strings.Compare(vpcCRN, *vpc.CRN) == 0 {
				vpcNamesList = append(vpcNamesList, *vpc.Name)
			}
		}
	}
	// if err = survey.AskOne(&survey.Select{
	// 	Message: "VPC",
	// 	Help:    "The VPC of the cluster. If you don't see your intended VPC listed, add the VPC to Permitted Networks of the DNS Zone and rerun the installer.",
	// 	Options: vpcNamesList,
	// },
	// 	&vpcChoice,
	// 	survey.WithValidator(func(ans interface{}) error {
	// 		choice := ans.(core.OptionAnswer).Value
	// 		i := sort.SearchStrings(vpcNamesList, choice)
	// 		if i == len(vpcNamesList) || vpcNamesList[i] != choice {
	// 			return fmt.Errorf("invalid VPC %q", choice)
	// 		}
	// 		return nil
	// 	}),
	// ); err != nil {
	// 	return "", fmt.Errorf("failed UserInput: %w", err)
	// }

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Permitted Network",
				Help:    "The VPC of the cluster. If you don't see your intended VPC listed, add the VPC to Permitted Networks of the DNS Zone and rerun the installer.",
				Default: "",
				Options: vpcNamesList,
			},
		},
	}, &vpcChoice)
	if err != nil {
		return "", fmt.Errorf("survey.ask failed with: %w", err)
	}
	return vpcChoice, nil
}
