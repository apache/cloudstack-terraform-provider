package cloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func dataSourceCloudstackTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackTemplateRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"template_filter": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"format": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hypervisor": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceCloudstackTemplateRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cloudstack.ListTemplatesParams{}
	p.SetListall(true)
	p.SetTemplatefilter(d.Get("template_filter").(string))

	csTemplates, err := cs.Template.ListTemplates(&p)
	if err != nil {
		return fmt.Errorf("Failed to list templates: %s", err)
	}

	filters := d.Get("filter")
	var templates []*cloudstack.Template

	for _, t := range csTemplates.Templates {
		match, err := applyFilters(t, filters.(*schema.Set))
		if err != nil {
			return err
		}

		if match {
			templates = append(templates, t)
		}
	}

	if len(templates) == 0 {
		return fmt.Errorf("No template is matching with the specified regex")
	}

	template, err := latestTemplate(templates)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected template: %s\n", template.Displaytext)

	return templateDescriptionAttributes(d, template)
}

func templateDescriptionAttributes(d *schema.ResourceData, template *cloudstack.Template) error {
	d.SetId(template.Id)
	d.Set("template_id", template.Id)
	d.Set("account", template.Account)
	d.Set("created", template.Created)
	d.Set("display_text", template.Displaytext)
	d.Set("format", template.Format)
	d.Set("hypervisor", template.Hypervisor)
	d.Set("name", template.Name)
	d.Set("size", template.Size)

	tags := make(map[string]interface{})
	for _, tag := range template.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func latestTemplate(templates []*cloudstack.Template) (*cloudstack.Template, error) {
	var latest time.Time
	var template *cloudstack.Template

	for _, t := range templates {
		created, err := time.Parse("2006-01-02T15:04:05-0700", t.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a template: %s", err)
		}

		if created.After(latest) {
			latest = created
			template = t
		}
	}

	return template, nil
}

func applyFilters(template *cloudstack.Template, filters *schema.Set) (bool, error) {
	var templateJSON map[string]interface{}
	t, _ := json.Marshal(template)
	json.Unmarshal(t, &templateJSON)

	for _, f := range filters.List() {
		m := f.(map[string]interface{})

		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}

		templateField := templateJSON[m["name"].(string)].(string)
		if !r.MatchString(templateField) {
			return false, nil
		}

	}

	return true, nil
}
