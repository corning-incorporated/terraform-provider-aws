// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package route53resolver

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_route53_resolver_firewall_rules", name="Firewall Rules")
func dataSourceResolverFirewallRules() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceResolverFirewallFirewallRulesRead,

		Schema: map[string]*schema.Schema{
			names.AttrAction: {
				Type:     schema.TypeString,
				Optional: true,
			},
			"firewall_rule_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firewall_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrAction: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_override_dns_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_override_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_override_ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"block_response": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrCreationTime: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_request_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_domain_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_rule_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrPriority: {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			names.AttrPriority: {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceResolverFirewallFirewallRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).Route53ResolverClient(ctx)

	firewallRuleGroupID := d.Get("firewall_rule_group_id").(string)
	rules, err := findFirewallRules(ctx, conn, firewallRuleGroupID, func(rule awstypes.FirewallRule) bool {
		if v, ok := d.GetOk(names.AttrAction); ok && string(rule.Action) != v.(string) {
			return false
		}

		if v, ok := d.GetOk(names.AttrPriority); ok && aws.ToInt32(rule.Priority) != int32(v.(int)) {
			return false
		}

		return true
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Route53 Resolver Firewall Rules (%s): %s", firewallRuleGroupID, err)
	}

	if n := len(rules); n == 0 {
		return sdkdiag.AppendErrorf(diags, "no Route53 Resolver Firewall Rules matched")
	} else if n > 1 {
		return sdkdiag.AppendErrorf(diags, "%d Route53 Resolver Firewall Rules matched; use additional constraints to reduce matches to a single Firewall Rule", n)
	}

	if err := d.Set("firewall_rules", flattenFirewallRules(rules)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting firewall_rules: %s", err)
	}

	d.SetId(firewallRuleGroupID)

	return diags
}

func flattenFirewallRules(apiObjects []awstypes.FirewallRule) []any {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []any

	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenFirewallRule(apiObject))
	}

	return tfList
}

func flattenFirewallRule(apiObject awstypes.FirewallRule) map[string]any {
	tfMap := map[string]any{
		names.AttrAction:          apiObject.Action,
		"block_override_dns_type": apiObject.BlockOverrideDnsType,
		"block_response":          apiObject.BlockResponse,
	}

	if apiObject.BlockOverrideDomain != nil {
		tfMap["block_override_domain"] = aws.ToString(apiObject.BlockOverrideDomain)
	}
	if apiObject.BlockOverrideTtl != nil {
		tfMap["block_override_ttl"] = aws.ToInt32(apiObject.BlockOverrideTtl)
	}
	if apiObject.CreationTime != nil {
		tfMap[names.AttrCreationTime] = aws.ToString(apiObject.CreationTime)
	}
	if apiObject.CreatorRequestId != nil {
		tfMap["creator_request_id"] = aws.ToString(apiObject.CreatorRequestId)
	}
	if apiObject.FirewallDomainListId != nil {
		tfMap["firewall_domain_list_id"] = aws.ToString(apiObject.FirewallDomainListId)
	}
	if apiObject.FirewallRuleGroupId != nil {
		tfMap["firewall_rule_group_id"] = aws.ToString(apiObject.FirewallRuleGroupId)
	}
	if apiObject.ModificationTime != nil {
		tfMap["modification_time"] = aws.ToString(apiObject.ModificationTime)
	}
	if apiObject.Name != nil {
		tfMap[names.AttrName] = aws.ToString(apiObject.Name)
	}
	if apiObject.Priority != nil {
		tfMap[names.AttrPriority] = aws.ToInt32(apiObject.Priority)
	}
	return tfMap
}
