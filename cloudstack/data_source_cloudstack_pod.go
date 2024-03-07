//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackPod() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackPodRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"pod_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"capacity_total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"capacity_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"percent_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pod_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"for_system_vms": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dsFlattenPodCapacity(capacity []cloudstack.PodCapacity) []map[string]interface{} {
	cap := make([]map[string]interface{}, len(capacity))
	for i, c := range capacity {
		cap[i] = map[string]interface{}{
			"capacity_allocated": c.Capacityallocated,
			"capacity_total":     c.Capacitytotal,
			"capacity_used":      c.Capacityused,
			"cluster_id":         c.Clusterid,
			"cluster_name":       c.Clustername,
			"name":               c.Name,
			"percent_used":       c.Percentused,
			"pod_id":             c.Podid,
			"pod_name":           c.Podname,
			"type":               c.Type,
			"zone_id":            c.Zoneid,
			"zone_name":          c.Zonename,
		}
	}
	return cap
}

func dsFlattenPodIpRanges(ip_ranges []cloudstack.PodIpranges) []map[string]interface{} {
	ranges := make([]map[string]interface{}, len(ip_ranges))
	for i, ip_range := range ip_ranges {
		ranges[i] = map[string]interface{}{
			"end_ip":         ip_range.Endip,
			"for_system_vms": ip_range.Forsystemvms,
			"start_ip":       ip_range.Startip,
			"vlan_id":        ip_range.Vlanid,
		}
	}
	return ranges
}

func datasourceCloudStackPodRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Pod.NewListPodsParams()

	csPods, err := cs.Pod.ListPods(p)
	if err != nil {
		return fmt.Errorf("failed to list pods: %s", err)
	}

	filters := d.Get("filter")

	for _, pod := range csPods.Pods {
		match, err := applyPodFilters(pod, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			return podDescriptionAttributes(d, pod)
		}
	}

	return fmt.Errorf("no pods found")
}

func podDescriptionAttributes(d *schema.ResourceData, pod *cloudstack.Pod) error {
	d.SetId(pod.Id)
	var end_ip string
	if len(pod.Endip) > 0 {
		end_ip = pod.Endip[0]
	}

	fields := map[string]interface{}{
		"pod_id":           pod.Id,
		"name":             pod.Name,
		"allocation_state": pod.Allocationstate,
		"gateway":          pod.Gateway,
		"netmask":          pod.Netmask,
		"start_ip":         pod.Startip[0],
		"vlan_id":          pod.Vlanid[0],
		"zone_id":          pod.Zoneid,
		"zone_name":        pod.Zonename,
		"end_ip":           end_ip,
		"ip_ranges":        dsFlattenPodIpRanges(pod.Ipranges),
		"capacity":         dsFlattenPodCapacity(pod.Capacity),
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			log.Printf("[WARN] Error setting %s: %s", k, err)
		}
	}

	return nil
}

func applyPodFilters(pod *cloudstack.Pod, filters *schema.Set) (bool, error) {
	val := reflect.ValueOf(pod).Elem()

	for _, f := range filters.List() {
		filter := f.(map[string]interface{})
		r, err := regexp.Compile(filter["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(filter["name"].(string), "_", "")
		podField := val.FieldByNameFunc(func(fieldName string) bool {
			if strings.EqualFold(fieldName, updatedName) {
				updatedName = fieldName
				return true
			}
			return false
		}).String()

		if r.MatchString(podField) {
			return true, nil
		}
	}

	return false, nil
}
