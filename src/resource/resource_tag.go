package resource

import (
	"context"
	"demo-xml/src/class"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scottdware/go-panos"
)

func RessourceTagSchemas() map[string]*schema.Schema {
	return map[string]*schema.Schema {
		"name": {
			Type: schema.TypeString,
			Required: true,
			Description: "The administrative tag's name",
		},
		"color": {
			Type: schema.TypeString,
			Required: true,
			Description: "The administrative tag's color",
		},
		"comments": {
			Type: schema.TypeString,
			Optional: true,
			Description: "The administrative tag's comment",
		},
		"vsys": {
			Type: schema.TypeString,
			Optional: true,
			Description: "The administrative tag's vsys",
			Default: "vsys1",
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	// Create a new tag
	tag := class.CreateTagClass(d.Get("name").(string), d.Get("color").(string), d.Get("comments").(string))

	var diags diag.Diagnostics
	err := tag.Add(pan)
	if(err != nil) {
		diags = append(diags, diag.Diagnostic {
			Severity: diag.Error,
			Summary: "Cannot Create tag",
			Detail: fmt.Sprintf("Error append during creation of the tag: %s, with: %s", err, tag.ToString()),
		})
	}
	d.SetId(tag.Name)
	return diags
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	var diags diag.Diagnostics

	tagName := d.Id()

	exist := class.CheckIfTagExist(tagName,pan)

	if(!exist) {
		// Check if ressource exist
		d.SetId("")
		return diags
	}

	tag, err := class.SearchTag(tagName,pan)
	if(err != nil) {
		return diag.FromErr(err)
	}
	
	d.Set("name", tag.Name)
	d.Set("color", tag.Color)
	d.Set("comments", tag.Comments)
	d.Set("vsys", "vsys1")

	return diags
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)

	tagName := d.Id()

	tag, err := class.SearchTag(tagName,pan)
	if(err != nil) {
		return diag.FromErr(err)
	}
	tag.Name = d.Get("name").(string)
	tag.Color = d.Get("color").(string)
	tag.Comments = d.Get("comments").(string)
	if(tag.Name != tagName) {
		d.SetId(tag.Name)
	}
		

	err = tag.Edit(pan)
	if(err != nil) {
		return diag.FromErr(err)
	}

	return resourceTagRead(ctx, d, m)
}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pan := m.(*panos.PaloAlto)
	var diags diag.Diagnostics

	tagName := d.Id()

	tag, err := class.SearchTag(tagName,pan)
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

func ResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext: resourceTagRead,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: RessourceTagSchemas(),
	}
}