package provider

import (
	"context"
	"demo-xml/src/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scottdware/go-panos"
)




func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"excpanos_tag": resource.ResourceTag(),
			"excpanos_netflow": resource.ResourceNetflow(),
		},
		ConfigureContextFunc: ProviderConfigure,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type: schema.TypeString,
				Required: true,
				Description: "The hostname of the Palo Alto",
			},
			"username": {
				Type: schema.TypeString,
				Required: true,
				Description: "The username of the Palo Alto",
			},
			"password": {
				Type: schema.TypeString,
				Required: true,
				Description: "The password of the Palo Alto",
			},
		},
	}
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	var diags diag.Diagnostics

	creds := &panos.AuthMethod{
		Credentials: []string{username, password},	
	}
	
	pan, err := panos.NewSession(hostname, creds)
	if(err != nil) {
		diags = append(diags, diag.Diagnostic {
			Severity: diag.Error,
			Summary: "Unable to connect to the FW",
			Detail: "Error append during connection with the Firewall.",
		})
		return nil, diags
	}
	return pan,diags
}
