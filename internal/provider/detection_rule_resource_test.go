package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func generateTestRule() DetectionRule {
	rule := DetectionRule{
		RuleID: "7CE764F6-36A7-4E72-AB8B-166170CD1C93",
	}
	return rule
}

func TestAccDetectionRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDetectionRuleResourceConfig(generateTestRule()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elastic-siem_detection_rule.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("elastic-siem_detection_rule.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "elastic-siem_detection_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"rule_content"},
			},
			// Update and Read testing
			{
				Config: testAccDetectionRuleResourceConfig(generateTestRule()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elastic-siem_detection_rule.test", "rule_content", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDetectionRuleResourceConfig(ruleContent DetectionRule) string {
	str, err := json.Marshal(ruleContent)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	content := strconv.Quote(string(str))
	return fmt.Sprintf(`%s
resource "elastic-siem_detection_rule" "test" {
  rule_content = %s
}
`, providerConfig, content)
}
