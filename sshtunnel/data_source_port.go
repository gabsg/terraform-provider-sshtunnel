package sshtunnel

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maclarensg/terraform-provider-sshtunnel/connect"
)

func dataSourcePort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePortRead,
		Schema: map[string]*schema.Schema{
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*connect.ConnectionInfo)

	time.Sleep(3 * time.Second)

	if c.Tunnel == nil {
		return diag.FromErr(errors.New("tunnel is Nil"))
	}

	if err := d.Set("port", c.Tunnel.Local.Port); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
