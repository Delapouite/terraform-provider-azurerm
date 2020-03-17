package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLPath() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.URLPathOperatorAny),
					string(cdn.URLPathOperatorBeginsWith),
					string(cdn.URLPathOperatorContains),
					string(cdn.URLPathOperatorEndsWith),
					string(cdn.URLPathOperatorEqual),
					string(cdn.URLPathOperatorGreaterThan),
					string(cdn.URLPathOperatorGreaterThanOrEqual),
					string(cdn.URLPathOperatorLessThan),
					string(cdn.URLPathOperatorLessThanOrEqual),
				}, false),
			},

			"negate_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},

			"transforms": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(cdn.Lowercase),
						string(cdn.Uppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionURLPath(upc map[string]interface{}) *cdn.DeliveryRuleURLPathCondition {
	requestURICondition := cdn.DeliveryRuleURLPathCondition{
		Name: cdn.NameURLPath,
		Parameters: &cdn.URLPathMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlPathMatchConditionParameters"),
			Operator:        cdn.URLPathOperator(upc["operator"].(string)),
			NegateCondition: utils.Bool(upc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(upc["match_values"].(*schema.Set).List()),
		},
	}

	if rawTransforms := upc["transforms"].([]interface{}); len(rawTransforms) != 0 {
		transforms := make([]cdn.Transform, 0)
		for _, t := range rawTransforms {
			transforms = append(transforms, cdn.Transform(t.(string)))
		}
		requestURICondition.Parameters.Transforms = &transforms
	}

	return &requestURICondition
}

func FlattenArmCdnEndpointConditionURLPath(upc *cdn.DeliveryRuleURLPathCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := upc.Parameters; params != nil {
		res["operator"] = string(params.Operator)

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}

		if params.Transforms != nil {
			transforms := make([]string, 0)
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
			res["transforms"] = &transforms
		}
	}

	return res
}
