package sharedschemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetSharedSchemaPayload() *schema.Schema {
	out := &schema.Schema{
		Type: schema.TypeString,
		// ValidateFunc: ValidatePayloadPlist,
		// DiffSuppressFunc: SuppressFormatDiff,
	}

	return out
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
