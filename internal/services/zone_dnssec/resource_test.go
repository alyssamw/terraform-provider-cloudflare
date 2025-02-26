package zone_dnssec_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudflareZoneDNSSECFull(t *testing.T) {
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("cloudflare_zone_dnssec.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareZoneDNSSECResourceConfig(zoneID, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudflareZoneDNSSECDataSourceID(name),
					resource.TestCheckResourceAttrSet(name, consts.ZoneIDSchemaKey),
					resource.TestMatchResourceAttr(name, "status", regexp.MustCompile("active|pending")),
					resource.TestCheckResourceAttrSet(name, "flags"),
					resource.TestCheckResourceAttrSet(name, "algorithm"),
					resource.TestCheckResourceAttrSet(name, "key_type"),
					resource.TestCheckResourceAttrSet(name, "digest_type"),
					resource.TestCheckResourceAttrSet(name, "digest_algorithm"),
					resource.TestCheckResourceAttrSet(name, "digest"),
					resource.TestCheckResourceAttrSet(name, "ds"),
					resource.TestCheckResourceAttrSet(name, "key_tag"),
					resource.TestCheckResourceAttrSet(name, "public_key"),
					resource.TestCheckResourceAttrSet(name, "modified_on"),
				),
			},
		},
	})
}
