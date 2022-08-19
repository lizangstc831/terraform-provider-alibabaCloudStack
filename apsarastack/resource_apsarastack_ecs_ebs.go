package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceApsaraStackEcsEbsStorageSets() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEcsEbsStorageSetsCreate,
		Read:   resourceApsaraStackEcsEbsStorageSetsRead,
		Delete: resourceApsaraStackEcsEbsStorageSetsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"storage_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"maxpartition_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackEcsEbsStorageSetsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//var response map[string]interface{}
	action := "CreateStorageSet"
	response := &datahub.EcsStorageSetsCreate{}

	//request := make(map[string]interface{})
	//conn, err := client.NewEcsClient()
	//if err != nil {
	//	return WrapError(err)
	//}
	StorageSetName := d.Get("storage_set_name").(string)
	MaxPartitionNumber := d.Get("maxpartition_number").(string)
	ZoneId := d.Get("zone_id").(string)

	//zoneid := d.Get("zone_id").(string)
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret":    client.SecretKey,
		"AccessKeyId":        client.AccessKey,
		"Product":            "Ecs",
		"RegionId":           client.RegionId,
		"Department":         client.Department,
		"ResourceGroup":      client.ResourceGroup,
		"Action":             action,
		"Version":            "2014-05-26",
		"StorageSetName":     StorageSetName,
		"MaxPartitionNumber": MaxPartitionNumber,
		"ZoneId":             ZoneId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw, request)
		bresponse := raw.(*responses.CommonResponse)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), response)

		//var response *ecs.CreateCommandResponse
		//response, _ := raw.(*ecs.CreateCommandResponse)
		d.SetId(fmt.Sprint(response.StorageSetId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ecs_command", action, ApsaraStackSdkGoERROR)
	}

	return resourceApsaraStackEcsEbsStorageSetsRead(d, meta)
}
func resourceApsaraStackEcsEbsStorageSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsEbsStorageSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_ecs_command ecsService.DescribeEcsCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("storage_set_name", object.StorageSets.StorageSet[0].StorageSetName)
	d.Set("zone_id", object.StorageSets.StorageSet[0].ZoneId)
	d.Set("maxpartition_number", strconv.Itoa(object.StorageSets.StorageSet[0].StorageSetPartitionNumber))
	return nil
}
func resourceApsaraStackEcsEbsStorageSetsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteStorageSet"
	//var response map[string]interface{}
	//conn, err := client.NewEcsClient()
	//if err != nil {
	//	return WrapError(err)
	//}
	//request := map[string]interface{}{
	//	"CommandId": d.Id(),
	//}
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"StorageSetId":    d.Id(),
	}

	//request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCmdId.NotFound", "InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return nil
}
