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

// <dict>
//                 <key>PayloadContent</key>
//                 <dict>
//                     <key>com.zscaler.installparams</key>
//                     <dict>
//                         <key>Forced</key>
//                         <array>
//                             <dict>
//                                 <key>mcx_preference_settings</key>
//                                 <dict>
//                                     <key>installation-parameters</key>
//                                     <dict>
//                                         <key>cloudName</key>
//                                         <string>zscalertwo</string>
//                                         <key>enableFips</key>
//                                         <string>0</string>
//                                         <key>hideAppUIOnLaunch</key>
//                                         <string>1</string>
//                                         <key>launchTray</key>
//                                         <string>1</string>
//                                         <key>strictEnforcement</key>
//                                         <string>0</string>
//                                         <key>userDomain</key>
//                                         <string>faire.com</string>
//                                     </dict>
//                                 </dict>
//                             </dict>
//                         </array>
//                     </dict>
//                 </dict>
//                 <key>PayloadDisplayName</key>
//                 <string>Custom Settings</string>
//                 <key>PayloadIdentifier</key>
//                 <string>515E2CFA-B10C-4573-BE91-CE779217C62F</string>
//                 <key>PayloadOrganization</key>
//                 <string>JAMF Software</string>
//                 <key>PayloadType</key>
//                 <string>com.apple.ManagedClient.preferences</string>
//                 <key>PayloadUUID</key>
//                 <string>515E2CFA-B10C-4573-BE91-CE779217C62F</string>
//                 <key>PayloadVersion</key>
//                 <integer>1</integer>
//             </dict>

// <?xml version="1.0" encoding="UTF-8"?>
// <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
// <plist version="1.0">
//     <dict>
//         <key>PayloadContent</key>
//         <array>
//             <dict>
//                 <key>AutoJoin</key>
//                 <true/>
//                 <key>CaptiveBypass</key>
//                 <false/>
//                 <key>EncryptionType</key>
//                 <string>WPA</string>
//                 <key>HIDDEN_NETWORK</key>
//                 <false/>
//                 <key>Interface</key>
//                 <string>BuiltInWireless</string>
//                 <key>Password</key>
//                 <string>welcometofaire</string>
//                 <key>PayloadDescription</key>
//                 <string/>
//                 <key>PayloadDisplayName</key>
//                 <string>WiFi (Faire)</string>
//                 <key>PayloadEnabled</key>
//                 <true/>
//                 <key>PayloadIdentifier</key>
//                 <string>E5C790DB-1EB4-4303-ADB6-1B1C4B607010</string>
//                 <key>PayloadOrganization</key>
//                 <string>Faire Inc</string>
//                 <key>PayloadType</key>
//                 <string>com.apple.wifi.managed</string>
//                 <key>PayloadUUID</key>
//                 <string>E5C790DB-1EB4-4303-ADB6-1B1C4B607010</string>
//                 <key>PayloadVersion</key>
//                 <integer>1</integer>
//                 <key>ProxyType</key>
//                 <string>None</string>
//                 <key>SSID_STR</key>
//                 <string>Faire</string>
//                 <key>SetupModes</key>
//                 <array>
//                 </array>
//             </dict>
//         </array>
//         <key>PayloadDescription</key>
//         <string/>
//         <key>PayloadDisplayName</key>
//         <string>WiFi (Faire)</string>
//         <key>PayloadEnabled</key>
//         <true/>
//         <key>PayloadIdentifier</key>
//         <string>C87E1ED5-FAF1-4F87-BA1E-C65555F38AE1</string>
//         <key>PayloadOrganization</key>
//         <string>Faire Inc</string>
//         <key>PayloadRemovalDisallowed</key>
//         <true/>
//         <key>PayloadScope</key>
//         <string>System</string>
//         <key>PayloadType</key>
//         <string>Configuration</string>
//         <key>PayloadUUID</key>
//         <string>C87E1ED5-FAF1-4F87-BA1E-C65555F38AE1</string>
//         <key>PayloadVersion</key>
//         <integer>1</integer>
//     </dict>
// </plist>
