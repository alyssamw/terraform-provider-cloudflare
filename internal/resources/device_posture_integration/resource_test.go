package device_posture_integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stainless-sdks/cloudflare-terraform/internal/acctest"
	"github.com/stainless-sdks/cloudflare-terraform/internal/consts"
	"github.com/stainless-sdks/cloudflare-terraform/internal/utils"
)

func TestAccCloudflareDevicePostureIntegrationCreate(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the Access
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		t.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("cloudflare_device_posture_integration.%s", rnd)
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	clientID := os.Getenv("CLOUDFLARE_WORKSPACE_ONE_CLIENT_ID")
	clientSecret := os.Getenv("CLOUDFLARE_WORKSPACE_ONE_CLIENT_SECRET")
	apiURL := os.Getenv("CLOUDFLARE_WORKSPACE_ONE_API_URL")
	authURL := os.Getenv("CLOUDFLARE_WORKSPACE_ONE_AUTH_URL")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_WorkspaceOne(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareDevicePostureIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareDevicePostureIntegration(rnd, accountID, clientID, clientSecret, apiURL, authURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, consts.AccountIDSchemaKey, accountID),
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "type", "workspace_one"),
					resource.TestCheckResourceAttr(name, "interval", "24h"),
					resource.TestCheckResourceAttr(name, "config.0.auth_url", authURL),
					resource.TestCheckResourceAttr(name, "config.0.api_url", apiURL),
					resource.TestCheckResourceAttr(name, "config.0.client_id", clientID),
				),
			},
		},
	})
}

func testAccCloudflareDevicePostureIntegration(rnd, accountID, clientID, clientSecret, apiURL, authURL string) string {
	return fmt.Sprintf(`
resource "cloudflare_device_posture_integration" "%[1]s" {
	account_id                = "%[2]s"
	name                      = "%[1]s"
	type                      = "workspace_one"
	interval                  = "24h"
	config {
		api_url       =  "%[5]s"
		auth_url      =  "%[6]s"
		client_id     =  "%[3]s"
		client_secret =  "%[4]s"
	}
}
`, rnd, accountID, clientID, clientSecret, apiURL, authURL)
}

func testAccCheckCloudflareDevicePostureIntegrationDestroy(s *terraform.State) error {
	client, clientErr := acctest.SharedV1Client() // TODO(terraform): replace with SharedV2Clent
	if clientErr != nil {
		tflog.Error(context.TODO(), fmt.Sprintf("failed to create Cloudflare client: %s", clientErr))
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_device_posture_integration" {
			continue
		}

		_, err := client.DevicePostureIntegration(context.Background(), rs.Primary.Attributes[consts.AccountIDSchemaKey], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Device Posture Integration still exists")
		}
	}

	return nil
}
