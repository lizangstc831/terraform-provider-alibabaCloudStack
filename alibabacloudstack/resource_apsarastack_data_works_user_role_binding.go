package alibabacloudstack

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksUserRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDataWorksUserRoleBindingCreate,
		Read:   resourceAlibabacloudStackDataWorksUserRoleBindingRead,
		Update: resourceAlibabacloudStackDataWorksUserRoleBindingUpdate,
		Delete: resourceAlibabacloudStackDataWorksUserRoleBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"role_project_owner", "role_project_admin", "role_project_dev", "role_project_pe", "role_project_deploy", "role_project_guest", "role_project_security"}, false),
			},
		},
	}
}

func resourceAlibabacloudStackDataWorksUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AddProjectMemberToRole"
	request := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v.(string)
	}

	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v.(string)
	}

	if v, ok := d.GetOk("role_code"); ok {
		request["RoleCode"] = v.(string)
	}

	request["RegionId"] = "default"
	request["ClientToken"] = fmt.Sprint(uuid.NewRandom())
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_data_works_folder", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RoleCode"], ":", request["ProjectId"], ":", request["UserId"]))

	return resourceAlibabacloudStackDataWorksUserRoleBindingRead(d, meta)
}
func resourceAlibabacloudStackDataWorksUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksUserRoleBinding(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksUserRoleBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("user_id", parts[2])
	d.Set("project_id", parts[1])
	d.Set("role_code", parts[0])

	return nil
}
func resourceAlibabacloudStackDataWorksUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	// 没有对应 API
	return resourceAlibabacloudStackDataWorksUserRoleBindingRead(d, meta)
}
func resourceAlibabacloudStackDataWorksUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "RemoveProjectMemberFromRole"
	var response map[string]interface{}
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ProjectId": parts[1],
		"UserId":    parts[2],
		"RoleCode":  parts[0],
	}

	request["RegionId"] = "default"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
