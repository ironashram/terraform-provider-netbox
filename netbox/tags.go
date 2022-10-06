package netbox

import (
	"fmt"
	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/extras"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const tagsKey = "tags"
const tagsIds = "tag_ids"
const tagsNames = "tag_names"

var tagsSchema = &schema.Schema{
	Type: schema.TypeSet,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Optional: true,
	Set:      schema.HashString,
}

var tagsSchemaRead = &schema.Schema{
	Type: schema.TypeSet,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Computed: true,
	Set:      schema.HashString,
}
var tagIdsSchemaRead = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{
		Type: schema.TypeInt,
	},
}

var tagNamesSchemaRead = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

func getNestedTagListFromResourceDataSet(client *client.NetBoxAPI, d interface{}) ([]*models.NestedTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	tagList := d.(*schema.Set).List()
	var tags []*models.NestedTag
	for _, tag := range tagList {

		tagString := tag.(string)
		params := extras.NewExtrasTagsListParams()
		params.Name = &tagString
		limit := int64(2) // We search for a unique tag. Having two hits suffices to know its not unique.
		params.Limit = &limit
		res, err := client.Extras.ExtrasTagsList(params, nil)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error retrieving tag %s from netbox", tag.(string)),
				Detail:   fmt.Sprintf("API Error trying to retrieve tag %s from netbox", tag.(string)),
			})
		} else {
			payload := res.GetPayload()
			if *payload.Count == int64(1) {
				tags = append(tags, &models.NestedTag{
					Name: payload.Results[0].Name,
					Slug: payload.Results[0].Slug,
				})
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Error retrieving tag %s from netbox", tag.(string)),
					Detail:   fmt.Sprintf("Could not map tag %s to unique tag in netbox", tag.(string)),
				})
			}
		}
	}
	return tags, diags
}

func getTagListFromNestedTagList(nestedTags []*models.NestedTag) []string {
	var tagNames []string
	for _, nestedTag := range nestedTags {
		tagNames = append(tagNames, *nestedTag.Name)
	}
	return tagNames
}

func getTagIdsListFromNestedTagList(nestedTags []*models.NestedTag) []int64 {
	var tagIds []int64
	for _, nestedTag := range nestedTags {
		tagIds = append(tagIds, nestedTag.ID)
	}
	return tagIds
}
