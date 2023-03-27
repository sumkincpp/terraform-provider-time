// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timemodifier

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ReplaceIfOutdated(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}

	rotationTimestamp, err := time.Parse(time.RFC3339, req.StateValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"replaceIfOutdated plan modifier error",
			"The rotation rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	now := time.Now().UTC()

	resp.RequiresReplace = now.After(rotationTimestamp)
}
