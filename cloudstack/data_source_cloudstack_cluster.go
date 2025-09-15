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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackCluster() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackClusterRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor": {
				Type:     schema.TypeString,
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
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_overcommit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_overcommit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ovm3vip": {
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
		},
	}
}

func dsFlattenClusterCapacity(capacity []cloudstack.ClusterCapacity) []map[string]interface{} {
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

func datasourceCloudStackClusterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Cluster.NewListClustersParams()

	csClusters, err := cs.Cluster.ListClusters(p)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %s", err)
	}

	filters := d.Get("filter")

	for _, cluster := range csClusters.Clusters {
		match, err := applyClusterFilters(cluster, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			return clusterDescriptionAttributes(d, cluster)
		}
	}

	return fmt.Errorf("no clusters found")
}

func clusterDescriptionAttributes(d *schema.ResourceData, cluster *cloudstack.Cluster) error {
	d.SetId(cluster.Id)

	fields := map[string]interface{}{
		"id":                      cluster.Id,
		"name":                    cluster.Name,
		"cluster_type":            cluster.Clustertype,
		"hypervisor":              cluster.Hypervisortype,
		"pod_id":                  cluster.Podid,
		"pod_name":                cluster.Podname,
		"zone_id":                 cluster.Zoneid,
		"zone_name":               cluster.Zonename,
		"allocation_state":        cluster.Allocationstate,
		"managed_state":           cluster.Managedstate,
		"cpu_overcommit_ratio":    cluster.Cpuovercommitratio,
		"memory_overcommit_ratio": cluster.Memoryovercommitratio,
		"arch":                    cluster.Arch,
		"ovm3vip":                 cluster.Ovm3vip,
		"capacity":                dsFlattenClusterCapacity(cluster.Capacity),
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			log.Printf("[WARN] Error setting %s: %s", k, err)
		}
	}

	return nil
}

func applyClusterFilters(cluster *cloudstack.Cluster, filters *schema.Set) (bool, error) {
	val := reflect.ValueOf(cluster).Elem()

	for _, f := range filters.List() {
		filter := f.(map[string]interface{})
		r, err := regexp.Compile(filter["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(filter["name"].(string), "_", "")
		clusterField := val.FieldByNameFunc(func(fieldName string) bool {
			if strings.EqualFold(fieldName, updatedName) {
				updatedName = fieldName
				return true
			}
			return false
		}).String()

		if r.MatchString(clusterField) {
			return true, nil
		}
	}

	return false, nil
}
