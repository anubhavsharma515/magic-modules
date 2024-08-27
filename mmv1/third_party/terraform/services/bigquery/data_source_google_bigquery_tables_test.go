// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bigquery_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBigqueryTables_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryTables_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.#", "1"), // Check if at least one table is found
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.0", "test_table_1"), // Check for specific table ID
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryTables_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  
  resource "google_bigquery_dataset" "foo" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "US"
    default_table_expiration_ms = 3600000
  }

  resource "google_bigquery_table" "test_table" {
    dataset_id        = google_bigquery_dataset.foo.dataset_id
    table_id          = "test_table_1"
    deletion_protection = false
    schema     = <<EOF
    [
      {
        "name": "name",
        "type": "STRING",
        "mode": "NULLABLE"
      }
    ]
    EOF
  }

  data "google_bigquery_tables" "example" {
    dataset_id = google_bigquery_dataset.foo.dataset_id
  }
`, context)
}

