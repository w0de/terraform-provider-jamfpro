package sharedschemas

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressScopeIDs suppresses ordering differences in lists of IDs.
// Due to DiffSuppressFunc limitation w/r/t arrays, invoked as a CustomizeDiff func.
// Suppresses diffs by setting new value to old value.
func DiffSuppressScopeIDs(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	tflog.Info(ctx, "suppressing ID diffs bazxoo")
	for _, key := range diff.UpdatedKeys() {
		if !strings.HasSuffix(key, "_ids") {
			continue
		}

		tflog.Info(ctx, fmt.Sprintf("foobar Scanning for erroneous ID array diffs in %s\n", key))

		old, new := diff.GetChange(key)
		oldType := reflect.ValueOf(old)
		newType := reflect.ValueOf(new)

		if oldType.Kind() != reflect.Array || newType.Kind() != reflect.Array {
			tflog.Info(ctx, fmt.Sprintf("xoobar kinds are wrong for %s\n", key))
			continue
		}

		if oldType.Type().Elem().Kind() != reflect.Int || newType.Type().Elem().Kind() != reflect.Int {
			tflog.Info(ctx, fmt.Sprintf("xoobar elem kinds are wrong for %s\n", key))
			continue
		}

		oldIDs := old.([]int)
		newIDs := new.([]int)

		if len(oldIDs) != len(newIDs) {
			tflog.Info(ctx, fmt.Sprintf("xoobar len mismatch for %s\n", key))
			continue
		}

		sort.Ints(oldIDs)
		sort.Ints(newIDs)
		equivalent := true

		for i, v := range oldIDs {
			if v != newIDs[i] {
				tflog.Info(ctx, fmt.Sprintf("xoobar value mismatch for %s\n", key))
				equivalent = false
				break
			}
		}

		if equivalent {
			tflog.Info(ctx, fmt.Sprintf("bazbar Suppressing ID array diffs in %s\n", key))
			err := diff.SetNew(key, old)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
