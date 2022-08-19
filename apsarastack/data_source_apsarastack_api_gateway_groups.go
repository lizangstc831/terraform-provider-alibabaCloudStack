package apsarastack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackApiGatewayGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackApigatewayGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"billing_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"illegal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceApsaraStackApigatewayGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := cloudapi.CreateDescribeApiGroupsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DescribeApiGroups",
		"Version":         "2016-07-14",
	}
	var allGroups []cloudapi.ApiGroupAttribute

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApiGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_api_gateway_groups", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeApiGroupsResponse)

		allGroups = append(allGroups, response.ApiGroupAttributes.ApiGroupAttribute...)

		if len(allGroups) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredGroups []cloudapi.ApiGroupAttribute
	var gatewayGroupNameRegex *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return WrapError(err)
		}
		gatewayGroupNameRegex = r
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for _, group := range allGroups {
		if gatewayGroupNameRegex != nil && !gatewayGroupNameRegex.MatchString(group.GroupName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[group.GroupId]; !ok {
				continue
			}
		}
		filteredGroups = append(filteredGroups, group)
	}
	return apigatewayGroupsDecriptionAttributes(d, filteredGroups, meta)
}

func apigatewayGroupsDecriptionAttributes(d *schema.ResourceData, groupsSetTypes []cloudapi.ApiGroupAttribute, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	var names []string
	for _, group := range groupsSetTypes {
		mapping := map[string]interface{}{
			"id":             group.GroupId,
			"name":           group.GroupName,
			"region_id":      group.RegionId,
			"sub_domain":     group.SubDomain,
			"description":    group.Description,
			"created_time":   group.CreatedTime,
			"modified_time":  group.ModifiedTime,
			"traffic_limit":  group.TrafficLimit,
			"billing_status": group.BillingStatus,
			"illegal_status": group.IllegalStatus,
		}
		ids = append(ids, group.GroupId)
		s = append(s, mapping)
		names = append(names, group.GroupName)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
