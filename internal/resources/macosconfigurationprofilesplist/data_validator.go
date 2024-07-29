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
	if diff.HasChange("scope") {
		new, old := diff.GetChange("scope")
		new, err := suppressIDOrderDiffs(ctx, "scope", new, old)
		if err != nil {
			return err
		}

		if new != nil {
			diff.SetNew("scope", new)
		}
	}

	if diff.HasChange("exclusions") {
		new, old := diff.GetChange("exclusions")
		new, err := suppressIDOrderDiffs(ctx, "exclusions", new, old)
		if err != nil {
			return err
		}

		if new != nil {
			diff.SetNew("exclusions", new)
		}
	}

	return nil
}

// DiffSuppressScopeIDs suppresses ordering differences in lists of IDs.
// Due to DiffSuppressFunc limitation w/r/t arrays, invoked as a CustomizeDiff func.
// Suppresses diffs by setting new value to old value.
func suppressIDOrderDiffs(ctx context.Context, schema string, newResource interface{}, oldResource interface{}) (any, error) {
	tflog.Info(ctx, fmt.Sprintf("foobar suppressing ID diffs in %s\n", schema))
	updated := false
	new := newResource.([]interface{})[0].(map[string]any)
	old := oldResource.([]interface{})[0].(map[string]any)

	for key, newIDs := range new {
		newVal := newIDs.([]int)
		if !strings.HasSuffix(key, "_ids") {
			continue
		}

		valType := reflect.ValueOf(newVal)
		tflog.Info(ctx, fmt.Sprintf("foobar %s: %s, %s\n", schema, key, valType.Kind()))

		oldVal := old[key].([]int)
		tflog.Info(ctx, fmt.Sprintf("foobar Scanning for erroneous ID array diffs in %s\n", key))

		oldType := reflect.ValueOf(oldVal)
		newType := reflect.ValueOf(newVal)

		if oldType.Kind() != reflect.Slice || newType.Kind() != reflect.Slice {
			tflog.Info(ctx, fmt.Sprintf("foobar kinds are wrong for %s\n", key))
			continue
		}

		// oldIDs := oldVal.([]interface)
		// newIDs := newVal.([]interface)

		if len(oldVal) != len(newVal) {
			tflog.Info(ctx, fmt.Sprintf("foobar len mismatch for %s\n", key))
			continue
		}

		sort.Ints(oldVal)
		sort.Ints(newVal)
		equivalent := true

		for i, v := range oldVal {
			if v != newVal[i] {
				tflog.Info(ctx, fmt.Sprintf("foobar value mismatch for %s\n", key))
				equivalent = false
				break
			}
		}

		if equivalent {
			tflog.Info(ctx, fmt.Sprintf("foobar Suppressing ID array diffs in %s\n", key))
		}
	}

	if updated {
		return new, nil
	}

	return nil, nil
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
