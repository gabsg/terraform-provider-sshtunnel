package sshtunnel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maclarensg/terraform-provider-sshtunnel/connect"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"listen_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"jumphost": {
				Type:     schema.TypeString,
				Required: true,
			},
			"jumphost_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"aws_profile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"sshtunnel_port": dataSourcePort(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (client interface{}, diag diag.Diagnostics) {
	listenPort := d.Get("listen_port").(int)
	jumphost := d.Get("jumphost").(string)
	jumphostPort := d.Get("jumphost_port").(int)
	awsprofile := d.Get("aws_profile").(string)
	awsregion := d.Get("aws_region").(string)
	targetHost := d.Get("target_host").(string)
	targetPort := d.Get("target_port").(int)
	user := d.Get("user").(string)

	tunnel := connect.New(listenPort, jumphost, jumphostPort, awsprofile, awsregion, targetHost, targetPort, user)

	return tunnel, diag
}
