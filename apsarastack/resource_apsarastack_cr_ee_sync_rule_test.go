package apsarastack

//import (
//	"fmt"
//	"os"
//	"testing"
//
//	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
//	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//)
//
//func TestAccApsaraStackCrEESyncRule_Basic(t *testing.T) {
//	region := os.Getenv("APSARASTACK_REGION")
//	sourceInstanceId, targetInstanceId := getCrEESyncRuleTestEnv(t)
//	resourceId := "apsarastack_cr_ee_sync_rule.default"
//	ra := resourceAttrInit(resourceId, nil)
//	serviceFunc := func() interface{} {
//		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
//	}
//	var v *cr_ee.SyncRulesItem
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEESyncRule")
//	rac := resourceAttrCheckInit(rc, ra)
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := acctest.RandIntRange(1000000, 9999999)
//	name := fmt.Sprintf("tf-testacc-cr-sync-%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEESyncRuleConfigDependence)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//			getCrEESyncRuleTestEnv(t)
//		},
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id":           "${apsarastack_cr_ee_namespace.source_ns.instance_id}",
//					"namespace_name":        "${apsarastack_cr_ee_namespace.source_ns.name}",
//					"name":                  "${var.name}",
//					"target_region_id":      region,
//					"target_instance_id":    "${apsarastack_cr_ee_namespace.target_ns.instance_id}",
//					"target_namespace_name": "${apsarastack_cr_ee_namespace.target_ns.name}",
//					"tag_filter":            ".*",
//					"depends_on": []string{
//						"apsarastack_cr_ee_repo.source_repo",
//						"apsarastack_cr_ee_repo.target_repo",
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_id":           sourceInstanceId,
//						"namespace_name":        name,
//						"name":                  name,
//						"target_region_id":      region,
//						"target_instance_id":    targetInstanceId,
//						"target_namespace_name": name,
//						"tag_filter":            ".*",
//						"rule_id":               CHECKSET,
//						"sync_direction":        "FROM",
//						"sync_scope":            "NAMESPACE",
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//			getCrEESyncRuleTestEnv(t)
//		},
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id":           "${apsarastack_cr_ee_namespace.source_ns.instance_id}",
//					"namespace_name":        "${apsarastack_cr_ee_namespace.source_ns.name}",
//					"name":                  "${var.name}",
//					"target_region_id":      region,
//					"target_instance_id":    "${apsarastack_cr_ee_namespace.target_ns.instance_id}",
//					"target_namespace_name": "${apsarastack_cr_ee_namespace.target_ns.name}",
//					"tag_filter":            ".*",
//					"repo_name":             "${var.name}",
//					"target_repo_name":      "${var.name}",
//					"depends_on": []string{
//						"apsarastack_cr_ee_repo.source_repo",
//						"apsarastack_cr_ee_repo.target_repo",
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_id":           sourceInstanceId,
//						"namespace_name":        name,
//						"name":                  name,
//						"target_region_id":      region,
//						"target_instance_id":    targetInstanceId,
//						"target_namespace_name": name,
//						"tag_filter":            ".*",
//						"rule_id":               CHECKSET,
//						"sync_direction":        "FROM",
//						"sync_scope":            "REPO",
//						"repo_name":             name,
//						"target_repo_name":      name,
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}
//
//func TestAccApsaraStackCrEESyncRule_Multi(t *testing.T) {
//	region := os.Getenv("APSARASTACK_REGION")
//	sourceInstanceId, targetInstanceId := getCrEESyncRuleTestEnv(t)
//	resourceId := "apsarastack_cr_ee_sync_rule.rule2"
//	ra := resourceAttrInit(resourceId, nil)
//	serviceFunc := func() interface{} {
//		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
//	}
//
//	var v *cr_ee.SyncRulesItem
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEESyncRule")
//	rac := resourceAttrCheckInit(rc, ra)
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := acctest.RandIntRange(1000000, 9999999)
//	name := fmt.Sprintf("tf-testacc-cr-sync-%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEESyncRuleConfigMultiDependence)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//			getCrEESyncRuleTestEnv(t)
//		},
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id":           "${apsarastack_cr_ee_namespace.source_ns.instance_id}",
//					"namespace_name":        "${apsarastack_cr_ee_namespace.source_ns.name}",
//					"name":                  "${var.name}2",
//					"target_region_id":      region,
//					"target_instance_id":    "${apsarastack_cr_ee_namespace.target_ns.instance_id}",
//					"target_namespace_name": "${apsarastack_cr_ee_namespace.target_ns.name}",
//					"tag_filter":            ".*",
//					"repo_name":             "${var.name}2",
//					"target_repo_name":      "${var.name}2",
//					"depends_on": []string{
//						"apsarastack_cr_ee_sync_rule.rule0",
//						"apsarastack_cr_ee_sync_rule.rule1",
//						"apsarastack_cr_ee_repo.source_repo",
//						"apsarastack_cr_ee_repo.target_repo",
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_id":           sourceInstanceId,
//						"namespace_name":        name,
//						"name":                  name + fmt.Sprint(2),
//						"target_region_id":      region,
//						"target_instance_id":    targetInstanceId,
//						"target_namespace_name": name,
//						"tag_filter":            ".*",
//						"rule_id":               CHECKSET,
//						"sync_direction":        "FROM",
//						"sync_scope":            "REPO",
//						"repo_name":             name + fmt.Sprint(2),
//						"target_repo_name":      name + fmt.Sprint(2),
//					}),
//				),
//			},
//		},
//	})
//}
//
//func resourceCrEESyncRuleConfigDependence(name string) string {
//	sourceInstanceId := os.Getenv("CR_EE_TEST_SOURCE_INSTANCE_ID")
//	targetInstanceId := os.Getenv("CR_EE_TEST_TARGET_INSTANCE_ID")
//	configTemplate := `
//variable "name" {
//	default = "%s"
//}
//
//resource "apsarastack_cr_ee_namespace" "source_ns" {
//	instance_id = "%s"
//	name = "${var.name}"
//	auto_create	= true
//	default_visibility = "PRIVATE"
//}
//
//resource "apsarastack_cr_ee_namespace" "target_ns" {
//	instance_id = "%s"
//	name = "${var.name}"
//	auto_create	= true
//	default_visibility = "PRIVATE"
//}
//
//resource "apsarastack_cr_ee_repo" "source_repo" {
//	instance_id = "${apsarastack_cr_ee_namespace.source_ns.instance_id}"
//	namespace = "${apsarastack_cr_ee_namespace.source_ns.name}"
//	name = "${var.name}"
//	summary = "test"
//	repo_type = "PRIVATE"
//	detail = "test"
//}
//
//resource "apsarastack_cr_ee_repo" "target_repo" {
//	instance_id = "${apsarastack_cr_ee_namespace.target_ns.instance_id}"
//	namespace = "${apsarastack_cr_ee_namespace.target_ns.name}"
//	name = "${var.name}"
//	summary = "test"
//	repo_type = "PRIVATE"
//	detail = "test"
//}
//`
//	return fmt.Sprintf(configTemplate, name, sourceInstanceId, targetInstanceId)
//}
//
//func resourceCrEESyncRuleConfigMultiDependence(name string) string {
//	region := os.Getenv("APSARASTACK_REGION")
//	sourceInstanceId := os.Getenv("CR_EE_TEST_SOURCE_INSTANCE_ID")
//	targetInstanceId := os.Getenv("CR_EE_TEST_TARGET_INSTANCE_ID")
//	configTemplate := `
//variable "region" {
//	default = "%s"
//}
//
//variable "name" {
//	default = "%s"
//}
//
//resource "apsarastack_cr_ee_namespace" "source_ns" {
//	instance_id = "%s"
//	name = "${var.name}"
//	auto_create	= true
//	default_visibility = "PRIVATE"
//}
//
//resource "apsarastack_cr_ee_namespace" "target_ns" {
//	instance_id = "%s"
//	name = "${var.name}"
//	auto_create	= true
//	default_visibility = "PRIVATE"
//}
//
//resource "apsarastack_cr_ee_repo" "source_repo" {
//	count = 3
//	instance_id = "${apsarastack_cr_ee_namespace.source_ns.instance_id}"
//	namespace = "${apsarastack_cr_ee_namespace.source_ns.name}"
//	name = "${var.name}${count.index}"
//	summary = "test"
//	repo_type = "PRIVATE"
//	detail = "test"
//}
//
//resource "apsarastack_cr_ee_repo" "target_repo" {
//	count = 3
//	instance_id = "${apsarastack_cr_ee_namespace.target_ns.instance_id}"
//	namespace = "${apsarastack_cr_ee_namespace.target_ns.name}"
//	name = "${var.name}${count.index}"
//	summary = "test"
//	repo_type = "PRIVATE"
//	detail = "test"
//}
//
//resource "apsarastack_cr_ee_sync_rule" "rule0" {
//	instance_id = "${apsarastack_cr_ee_namespace.source_ns.instance_id}"
//	namespace_name = "${apsarastack_cr_ee_namespace.source_ns.name}"
//	name = "${var.name}0"
//	target_region_id = "${var.region}"
//	target_instance_id = "${apsarastack_cr_ee_namespace.target_ns.instance_id}"
//	target_namespace_name = "${apsarastack_cr_ee_namespace.target_ns.name}"
//	tag_filter = ".*"
//	repo_name = "${var.name}0"
//	target_repo_name = "${var.name}0"
//	depends_on = [
//		apsarastack_cr_ee_repo.source_repo,
//		apsarastack_cr_ee_repo.target_repo
//	]
//}
//
//resource "apsarastack_cr_ee_sync_rule" "rule1" {
//	instance_id = "${apsarastack_cr_ee_namespace.source_ns.instance_id}"
//	namespace_name = "${apsarastack_cr_ee_namespace.source_ns.name}"
//	name = "${var.name}1"
//	target_region_id = "${var.region}"
//	target_instance_id = "${apsarastack_cr_ee_namespace.target_ns.instance_id}"
//	target_namespace_name = "${apsarastack_cr_ee_namespace.target_ns.name}"
//	tag_filter = ".*"
//	repo_name = "${var.name}1"
//	target_repo_name = "${var.name}1"
//	depends_on = [
//		apsarastack_cr_ee_sync_rule.rule0,
//		apsarastack_cr_ee_repo.source_repo,
//		apsarastack_cr_ee_repo.target_repo
//	]
//}
//`
//	return fmt.Sprintf(configTemplate, region, name, sourceInstanceId, targetInstanceId)
//}
//
//func getCrEESyncRuleTestEnv(t *testing.T) (string, string) {
//	sourceInstanceId := os.Getenv("CR_EE_TEST_SOURCE_INSTANCE_ID")
//	targetInstanceId := os.Getenv("CR_EE_TEST_TARGET_INSTANCE_ID")
//	if sourceInstanceId == "" || targetInstanceId == "" {
//		t.Skipf("Skipping cr ee test case without default instances")
//	}
//
//	return sourceInstanceId, targetInstanceId
//}
