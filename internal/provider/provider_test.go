// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"time": providerserver.NewProtocol5WithError(New()),
	}
}

func providerVersion080() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"time": {
			VersionConstraint: "0.8.0",
			Source:            "hashicorp/time",
		},
	}
}

func testCheckAttributeValuesDiffer(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) == testStringValue(j) {
			return fmt.Errorf("attribute values are the same")
		}

		return nil
	}
}

func testCheckAttributeValuesSame(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) != testStringValue(j) {
			return fmt.Errorf("attribute values are different")
		}

		return nil
	}
}

//nolint:unparam
func testExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

// Certain testing requires time differences that are too fast for unit testing.
// Sleeping for a second or two seems pragmatic in our testing.
func testSleep(seconds int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(seconds) * time.Second)

		return nil
	}
}

func testStringValue(sPtr *string) string {
	if sPtr == nil {
		return ""
	}

	return *sPtr
}
