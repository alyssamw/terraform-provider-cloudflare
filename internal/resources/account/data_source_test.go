package account_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stainless-sdks/cloudflare-terraform/internal/acctest"
	"github.com/stainless-sdks/cloudflare-terraform/internal/utils"
)

func TestAccCloudflareAccounts(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("data.cloudflare_accounts.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareAccountsConfig(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudflareAccountsSize(name),
				),
			},
		},
	})
}

func testAccCloudflareAccountsConfig(name string) string {
	return fmt.Sprintf(`data "cloudflare_accounts" "%[1]s" { }`, name)
}

func testAccCloudflareAccountsSize(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		var (
			accountsSize int
			err          error
		)

		if accountsSize, err = strconv.Atoi(a["accounts.#"]); err != nil {
			return err
		}

		if accountsSize < 1 {
			return fmt.Errorf("accounts count seems suspicious: %d", accountsSize)
		}

		return nil
	}
}
