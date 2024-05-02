package sharedschemas

import (
	"fmt"

	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// type PayloadContent struct {
// 	PayloadDescription  string `plist:"PayloadDescription,attr"`
// 	PayloadDisplayName  string `plist:"PayloadDisplayName,attr"`
// 	PayloadIdentifier   string `plist:"PayloadIdentifier,attr"`
// 	PayloadOrganization string `plist:"PayloadOrganization,attr"`
// 	PayloadType         string `plist:"PayloadType,attr"`
// 	PayloadUUID         string `plist:"PayloadUUID,attr"`
// 	PayloadVersion      int    `plist:"PayloadVersion,attr"`
// }

type Payload struct {
	PayloadContent      interface{} `plist:"PayloadContent,attr"`
	PayloadDisplayName  string      `plist:"PayloadDisplayName,attr"`
	PayloadOrganization string      `plist:"PayloadOrganization,attr"`
	PayloadType         string      `plist:"PayloadType,attr"`
	PayloadIdentifier   string      `plist:"PayloadIdentifier,attr"`
	PayloadVersion      int         `plist:"PayloadVersion,attr"`
}

type Payloads struct {
	PayloadContent      []interface{} `plist:"PayloadContent,attr"`
	PayloadDescription  string        `plist:"PayloadDescription,attr"`
	PayloadDisplayName  string        `plist:"PayloadDisplayName,attr"`
	PayloadOrganization string        `plist:"PayloadOrganization,attr"`
	PayloadUUID         string        `plist:"PayloadUUID,attr"`
	PayloadType         string        `plist:"PayloadType,attr"`
	PayloadIdentifier   string        `plist:"PayloadIdentifier,attr"`
	PayloadVersion      int           `plist:"PayloadVersion,attr"`
	PayloadEnabled      bool          `plist:"PayloadEnabled,attr"`
	PayloadScope        string        `plist:"PayloadScope,attr"`
}

func GetSharedSchemaPayload() *schema.Schema {
	out := &schema.Schema{
		Type:             schema.TypeString,
		ValidateFunc:     ValidatePayloadPlist,
		DiffSuppressFunc: SuppressFormatDiff,
	}

	return out
}

func UnmarshalPayloads(val []byte) (Payloads, error) {
	var payloads Payloads
	if _, err := plist.Unmarshal(val, &payloads); err != nil {
		return payloads, err
	}

	return payloads, nil
}

func SuppressFormatDiff(key string, old string, new string, d *schema.ResourceData) bool {
	var oldPayload map[string]interface{}
	ov := util.GetString(old)
	plist.Unmarshal([]byte(ov), &oldPayload)
	var newPayload map[string]interface{}
	nv := util.GetString(new)
	plist.Unmarshal([]byte(nv), &newPayload)
	oldPayloadFormatted, _ := plist.MarshalIndent(oldPayload, plist.XMLFormat, "  ")
	newPayloadFormatted, _ := plist.MarshalIndent(newPayload, plist.XMLFormat, "  ")

	return string(oldPayloadFormatted) == string(newPayloadFormatted)
}

func ValidatePayloadPlist(val interface{}, key string) (warns []string, errs []error) {
	var payload map[string]interface{}
	v := util.GetString(val)
	if _, err := plist.Unmarshal([]byte(v), &payload); err != nil {
		errs = append(errs, fmt.Errorf("%q has an invalid profile payload: %v", key, err))
	}

	return warns, errs
}
