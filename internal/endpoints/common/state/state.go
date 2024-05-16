// common/state/state.go
// This package contains shared / common resource functions
package state

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
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

// Struct to mirror MacOS .plist conifguration profile data with bucket for unexpected values
type ConfigurationProfile struct {
	PayloadContent     []PayloadContentListItem
	PayloadDisplayName string
	PayloadIdentifier  string
	PayloadType        string
	PayloadUUID        string
	PayloadVersion     int
	UnexpectedValues   map[string]interface{} `mapstructure:",remain"`
}

// Struct to mirror xml payload item with key for all dynamic values
type PayloadContentListItem struct {
	// PayloadVersion        int
	// PayloadType           string
	PayloadDisplayName  string
	PayloadOrganization string
	PayloadIdentifier   string
	PayloadUUID         string
	PayloadEnabled      bool
	NonComputedValues   map[string]interface{} `mapstructure:",remain"`
}

// plistDataToStruct takes xml .plist bytes data and returns ConfigurationProfile
func StructToPayloadData(payload PayloadContentListItem) (string, error) {
	plistData, err := plist.MarshalIndent(payload.NonComputedValues, plist.XMLFormat, "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	return string(plistData), nil
}

func ConfigurationFilePlistToStructFromFile(filepath string) (*ConfigurationProfile, error) {
	plistFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer plistFile.Close()

	xmlData, err := io.ReadAll(plistFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read plist/xml file: %v", err)
	}

	return plistDataToStruct(xmlData)
}

func ConfigurationProfilePlistToStructFromString(plistData string) (*ConfigurationProfile, error) {
	return plistDataToStruct([]byte(plistData))
}

func plistDataToStruct(plistBytes []byte) (*ConfigurationProfile, error) {
	var unmarshalledPlist map[string]interface{}
	_, err := plist.Unmarshal(plistBytes, &unmarshalledPlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plist/xml: %v", err)
	}

	var out ConfigurationProfile
	err = mapstructure.Decode(unmarshalledPlist, &out)
	if err != nil {
		return nil, fmt.Errorf("(mapstructure) failed to map unmarshaled configuration profile to struct: %v", err)
	}

	return &out, nil
}

func UnmarshalProfileAndPayloads(plistData string) (*ConfigurationProfile, []string, error) {
	profile, err := ConfigurationProfilePlistToStructFromString(plistData)
	if err != nil {
		return nil, nil, err
	}

	payloads := make([]string, len(profile.PayloadContent))
	for i, v := range profile.PayloadContent {
		if payload, err := StructToPayloadData(v); err != nil {
			return nil, nil, err
		} else {
			payloads[i] = payload
		}
	}

	sort.Strings(payloads)

	return profile, payloads, nil
}

func UpdateConfigurationProfilePayloads(d *schema.ResourceData, plistData string) error {
	profile, payloads, err := UnmarshalProfileAndPayloads(plistData)
	if err != nil {
		return err
	}

	if err := d.Set("payloads", payloads); err != nil {
		return err
	}

	if err := d.Set("plist", plistData); err != nil {
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
// 	payloads := make([]PayloadContentListItem, 0, 100)
// 	for _, val := range rawPayloads {
// 		profile, err := ConfigurationProfilePlistToStructFromString(val)
// 		if err != nil {
// 			return err
// 		}

// 		for _, v := range profile.PayloadContent {
// 			payloads := append(payloads, v)
// 		}
// 	}

// }
