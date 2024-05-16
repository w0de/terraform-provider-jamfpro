package sharedschemas

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetSharedSchemaPayload() *schema.Schema {
	out := &schema.Schema{
		Type:         schema.TypeString,
		StateFunc:    ExtractPayload,
		ValidateFunc: ValidatePayload,
		// DiffSuppressFunc: SuppressFormatDiff,
	}

	return out
}

func ExtractPayload(payload any) string {
	_, payloads, _ := state.UnmarshalProfileAndPayloads(payload.(string))

	return payloads[0]
}

func ValidatePayload(payload interface{}, key string) (warns []string, errs []error) {
	if _, payloads, err := state.UnmarshalProfileAndPayloads(payload.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%q has an invalid profile payload: %v", key, err))
	} else if len(payloads) != 1 {
		errs = append(errs, fmt.Errorf("%q does not contain a PayloadContent array containing one payload (found %d payloads)", key, len(payloads)))
	}

	return warns, errs
}

// func SuppressFormatDiff(key string, old string, new string, d *schema.ResourceData) bool {
// 	var oldPayload map[string]interface{}
// 	ov := util.GetString(old)
// 	plist.Unmarshal([]byte(ov), &oldPayload)
// 	var newPayload map[string]interface{}
// 	nv := util.GetString(new)
// 	plist.Unmarshal([]byte(nv), &newPayload)
// 	oldPayloadFormatted, _ := plist.MarshalIndent(oldPayload, plist.XMLFormat, "  ")
// 	newPayloadFormatted, _ := plist.MarshalIndent(newPayload, plist.XMLFormat, "  ")

// 	return string(oldPayloadFormatted) == string(newPayloadFormatted)
// }

// func ValidatePayloadPlist(val interface{}, key string) (warns []string, errs []error) {
// 	var payload map[string]interface{}
// 	v := util.GetString(val)
// 	if _, err := plist.Unmarshal([]byte(v), &payload); err != nil {
// 		errs = append(errs, fmt.Errorf("%q has an invalid profile payload: %v", key, err))
// 	}

// 	return warns, errs
// }
