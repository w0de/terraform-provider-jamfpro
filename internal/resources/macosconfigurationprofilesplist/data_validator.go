// macosconfigurationprofilesplist_data_validator.go
package macosconfigurationprofilesplist

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/datavalidators"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateDistributionMethod(ctx, diff, i); err != nil {
		return err
	}

	if err := validateMacOSConfigurationProfileLevel(ctx, diff, i); err != nil {
		return err
	}

	if err := validateConfigurationProfileFormatting(ctx, diff, i); err != nil {
		return err
	}

	if err := suppressScopeOrderingDiffs(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

func suppressScopeOrderingDiffs(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	schemata := [2]string{"scope", "exclusions"}
	logctx := tflog.NewSubsystem(ctx, "foobar: suppress sorting diffs in ID lists")
	for _, schema := range schemata {
		tflog.Info(logctx, fmt.Sprintf("scanning: %s\n", schema))
		if !diff.HasChange(schema) {
			tflog.Debug(logctx, fmt.Sprintf("unchanged: %s\n", schema))
			continue
		}

		setNew := false
		oldVal, newVal := diff.GetChange(schema)
		old := oldVal.([]interface{})[0].(map[string]any)
		new := newVal.([]interface{})[0].(map[string]any)

		for key, val := range old {
			if !strings.HasSuffix(key, "_ids") {
				continue
			}

			tflog.Info(logctx, fmt.Sprintf("scanning: %s.%s\n", schema, key))
			oldList := val.([]interface{})
			newList := new[key].([]interface{})

			// Length differs or both empty - order not important
			if len(oldList) != len(newList) || len(oldList) == 0 {
				continue
			}

			tflog.Debug(logctx, fmt.Sprintf("xoobar: (%v, %v) newList:%v, (%v, %v) oldList:%v\n", len(newList), reflect.TypeOf(newList), newList, len(oldList), reflect.TypeOf(oldList), oldList))

			oldIDs, err := convertToIntArray(oldList)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("skipping (old value is not []int): %s.%s", schema, key))
				continue
			}

			newIDs, err := convertToIntArray(newList)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("skipping (new value is not []int): %s.%s", schema, key))
				continue
			}

			sort.Ints(oldIDs)
			sort.Ints(newIDs)

			// Compare sorted IDs
			equivalent := true
			for i, id := range oldIDs {
				if id != newIDs[i] {
					equivalent = false
					break
				}
			}

			// If there's no change between sorted versions, rewrite new to old
			if equivalent {
				tflog.Debug(logctx, fmt.Sprintf("diff.SetNew: %s.%s - \n%v\n->\n%v\n", schema, key, new[key], val))
				tflog.Debug(logctx, fmt.Sprintf("xoobar: newIDs: %v, oldIDs:%v\n", newIDs, oldIDs))
				new[key] = val
				setNew = true
			}
		}

		if setNew {
			if err := diff.SetNew(schema, new); err != nil {
				return err
			}
			tflog.Info(logctx, fmt.Sprintf("diff.SetNew: %s\n", schema))
		}
	}

	return nil
}

func convertToIntArray(input []interface{}) ([]int, error) {
	result := make([]int, len(input))
	for i, v := range input {
		num, ok := v.(int)
		if !ok {
			return nil, fmt.Errorf("element %d is not an int", i)
		}
		result[i] = num
	}
	return result, nil
}

// validateDistributionMethod checks that the 'self_service' block is only used when 'distribution_method' is "Make Available in Self Service".
func validateDistributionMethod(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	distributionMethod, ok := diff.GetOk("distribution_method")

	if !ok {
		return nil
	}

	selfServiceBlockExists := len(diff.Get("self_service").([]interface{})) > 0

	if distributionMethod == "Make Available in Self Service" && !selfServiceBlockExists {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'self_service' block is required when 'distribution_method' is set to 'Make Available in Self Service'", resourceName)
	}

	if distributionMethod != "Make Available in Self Service" && selfServiceBlockExists {
		log.Printf("[WARN] 'jamfpro_macos_configuration_profile.%s': 'self_service' block is not meaningful when 'distribution_method' is set to '%s'", resourceName, distributionMethod)
	}

	return nil
}

// validateMacOSConfigurationProfileLevel validates that the 'PayloadScope' key in the payload matches the 'level' attribute.
func validateMacOSConfigurationProfileLevel(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	plistData, err := plist.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	payloadScope, err := datavalidators.GetPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	if payloadScope != level {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'level' attribute (%s) does not match the 'PayloadScope' in the plist (%s)", resourceName, level, payloadScope)
	}

	return nil
}

// validateConfigurationProfileFormatting validates the indentation of the plist XML.
func validateConfigurationProfileFormatting(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	payloads := diff.Get("payloads").(string)

	if err := datavalidators.CheckPlistIndentationAndWhiteSpace(payloads); err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': %v", resourceName, err)
	}

	return nil
}
