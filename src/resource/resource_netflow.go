package resource

import (
	"context"
	"demo-xml/src/class"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scottdware/go-panos"
)

/*
 READ DOES NOT WORK
 FUNC EXIT DOES NOT WORK
*/

func RessourceNetflowSchemas() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"name": {
			Type: schema.TypeString,
			Required: true,
			Description: "The administrative Netflow's name",
		},
		"template_refresh_rate": {
			Type: schema.TypeSet,
			Optional: true,
			Description: "Netflow template refresh rate",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"minutes": {
						Type: schema.TypeInt,
						Optional: true,
						Default: 5,
						Description: "Minutes beetwen each refresh",
					},
					"paquets": {
						Type: schema.TypeInt,
						Optional: true,
						Default: 5,
						Description: "Seconds beetwen each refresh",
					},
				},
			},
		},
		"active_timeout": {
			Type: schema.TypeInt,
			Optional: true,
			Default: 5,
			Description: "Minutes before a flow is considered inactive",
		},
		"server": {
			Type: schema.TypeList,
			Required: true,
			Description: "The server to send the netflow",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type: schema.TypeString,
						Required: true,
						Description: "The name of the server",
					},
					"port": {
						Type: schema.TypeInt,
						Required: true,
						Description: "The port of the server",
					},
					"hostname": {
						Type: schema.TypeString,
						Required: true,
						Description: "The hostname of the server",
					},
				},
			},
		},
	}
}

func resourceNetflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	// Create a new netflow
	name := d.Get("name").(string)
	templateRefreshRateSet := d.Get("template_refresh_rate").(*schema.Set)
	
	templateRefreshRateMap := templateRefreshRateSet.List()[0].(map[string]interface{})
	//var minutes int = int(templateRefreshRateMap["minutes"].(int))
	//var packets int = int(templateRefreshRateMap["paquets"].(int))
	//return  diag.FromErr(errors.New(fmt.Sprintf("%i",toto)))
	
	templateRefreshRate := class.TemplateRefresh{
		Minutes: int(templateRefreshRateMap["minutes"].(int)),
		Packets:  int(templateRefreshRateMap["paquets"].(int)),
	}


	activeTimeout := d.Get("active_timeout").(int)

	serverList := d.Get("server").([]interface{})

	var serverListClass []class.NetflowServer
	for _, server := range serverList {
		serverMap := server.(map[string]interface{})
		serverListClass = append(serverListClass, class.NetflowServer{
			Name: serverMap["name"].(string),
			Port: serverMap["port"].(int),
			Host: serverMap["hostname"].(string),
		})
	}

	netflow := class.NetflowProfile{
		Name: name,
		TemplateRefresh: templateRefreshRate,
		ActiveTimeout: activeTimeout,
		Server: class.ListOfServer{ NetflowServer: serverListClass },
	}
	

	var diags diag.Diagnostics	
	err := netflow.Add(pan)
	if(err != nil) {
		diags = append(diags, diag.Diagnostic {
			Severity: diag.Error,
			Summary: "Cannot Create Netflow Profile",
			Detail: fmt.Sprintf("Error append during creation of the tag: %s, with: %s", err, netflow.ToString()),
		})
	}
	d.SetId(name)
	return diags
}

func resourceNetflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	var diags diag.Diagnostics

	netflowName := d.Id()
	exist := class.CheckIfNetflowExist(netflowName, pan)

	if(!exist) {
		// Check if ressource exist
		d.SetId("")
		return diags
	}

	netflow, err := class.SearchNetflow(netflowName,pan)
	if(err != nil) {
		return diag.FromErr(err)
	}
	
	d.Set("name", netflow.Name)

	//templateRefreshRate := schema.Set{}

	template := map[string]interface{}{
		"minutes": netflow.TemplateRefresh.Minutes,
		"packets": netflow.TemplateRefresh.Packets,
	}

	templateFlatten := schema.NewSet(schema.HashResource(RessourceNetflowSchemas()["template_refresh_rate"].Elem.(*schema.Resource)), []interface{}{template})



	d.Set("template_refresh_rate", templateFlatten)
	d.Set("active_timeout", netflow.ActiveTimeout)

	var serverList []map[string]interface{}
	for _, server := range netflow.Server.NetflowServer {
		serverList = append(serverList, map[string]interface{}{
			"name": server.Name,
			"port": server.Port,
			"hostname": server.Host,
		})
	}
	d.Set("server", serverList)

	if(netflow.Name != netflowName) {
		d.SetId(netflow.Name)
	}

	return diags
}

func resourceNetflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)

	netflowName := d.Id()

	netflow, err := class.SearchNetflow(netflowName,pan)
	if(err != nil) {
		return diag.FromErr(err)
	}

	if(d.HasChange("name")) {
		netflowToDelete, err := class.SearchNetflow(netflowName, pan)
		if err != nil {
			return diag.FromErr(err)
		}
		errorDelete := netflowToDelete.Delete(pan)
		if errorDelete != nil {
			return diag.FromErr(errorDelete)
		}
		//return diag.Errorf("Cannot change the name of a netflow profile")
		netflow.Name = d.Get("name").(string)
		d.SetId(netflow.Name)
	}

	if(d.HasChange("template_refresh_rate")) {
		templateRefreshRateSet := d.Get("template_refresh_rate").(*schema.Set)
		templateRefreshRateMap := templateRefreshRateSet.List()[0].(map[string]interface{})		
		templateRefreshRate := class.TemplateRefresh{
			Minutes: int(templateRefreshRateMap["minutes"].(int)),
			Packets:  int(templateRefreshRateMap["paquets"].(int)),
		}
		netflow.TemplateRefresh = templateRefreshRate
	}

	if(d.HasChange("active_timeout")) {
		netflow.ActiveTimeout = d.Get("active_timeout").(int)
	}

	if(d.HasChange("server")) {
		serverList := d.Get("server").([]interface{})
		var serverListClass []class.NetflowServer
		for _, server := range serverList {
			serverMap := server.(map[string]interface{})
			serverListClass = append(serverListClass, class.NetflowServer{
				Name: serverMap["name"].(string),
				Port: serverMap["port"].(int),
				Host: serverMap["hostname"].(string),
			})
		}
		netflow.Server.NetflowServer = serverListClass
	}

	err = netflow.Edit(pan)
	if(err != nil) {
		return diag.FromErr(err)
	}

	return resourceNetflowRead(ctx, d, m)
}

func resourceNetflowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	var diags diag.Diagnostics

	netflowName := d.Id()

	tag, err := class.SearchNetflow(netflowName,pan)
	if(err != nil) {
		return diag.FromErr(err)
	}

	err = tag.Delete(pan)
	if(err != nil) {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func ResourceNetflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetflowCreate,
		ReadContext: resourceNetflowRead,
		UpdateContext: resourceNetflowUpdate,
		DeleteContext: resourceNetflowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: RessourceNetflowSchemas(),
	}
}