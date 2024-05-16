// common/state/state.go
// This package contains shared / common resource functions
package state

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/tools/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// Helper function to handle "resource not found" errors
func HandleResourceNotFoundError(err error, d *schema.ResourceData) diag.Diagnostics {
	if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "410") {
		d.SetId("") // Remove the resource from Terraform state
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   "The resource was not found on the remote server. It has been removed from the Terraform state.",
			},
		}
	} else {
		// For other errors, return a diagnostic error
		return diag.FromErr(err)
	}
}

// plistDataToStruct takes xml .plist bytes data and returns ConfigurationProfile
func StructToPayloadData(payload utils.PayloadContentListItem) (string, error) {
	plistData, err := plist.MarshalIndent(payload, plist.XMLFormat, "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	return string(plistData), nil
}

func UpdateConfigurationProfilePayloads(d *schema.ResourceData, plistData string) error {
	profile, err := utils.ConfigurationProfilePlistToStructFromString(plistData)
	if err != nil {
		return err
	}

	payloads := make([]string, len(profile.PayloadContent))
	for i, v := range profile.PayloadContent {
		payload, err := StructToPayloadData(v)
		if err != nil {
			return err
		}

		payloads[i] = payload
	}

	sort.Strings(payloads)

	if err := d.Set("payloads", payloads); err != nil {
		return err
	}

	if err := d.Set("identifier", profile.PayloadIdentifier); err != nil {
		return err
	}

	// if err := d.Set("organization", profile.); err != nil {
	// 	return err
	// }

	if err := d.Set("display_name", profile.PayloadDisplayName); err != nil {
		return err
	}

	return nil
}

// func ConstructConfigurationProfile(d *schema.ResourceData, rawPayloads []string) error {
// 	payloads := make([]utils.PayloadContentListItem, 0, 100)
// 	for _, val := range rawPayloads {
// 		profile, err := utils.ConfigurationProfilePlistToStructFromString(val)
// 		if err != nil {
// 			return err
// 		}

// 		for _, v := range profile.PayloadContent {
// 			payloads := append(payloads, v)
// 		}
// 	}

// }
