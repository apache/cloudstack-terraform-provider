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
			"allocation_state": {
				Description: "Allocation state of this cluster for allocation of new resources",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cluster_name": {
				Description: "the cluster name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Description: "Type of the cluster: CloudManaged, ExternalManaged",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"guest_vswitch_name": {
				Description: "Name of virtual switch used for guest traffic in the cluster. This would override zone wide traffic label setting.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"guest_vswitch_type": {
				Description: "Type of virtual switch used for guest traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"hypervisor": {
				Description: "hypervisor type of the cluster: XenServer,KVM,VMware,Hyperv,BareMetal,Simulator,Ovm3",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ovm3_cluster": {
				Description: "Ovm3 native OCFS2 clustering enabled for cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ovm3_pool": {
				Description: "Ovm3 native pooling enabled for cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ovm3_vip": {
				Description: "Ovm3 vip to use for pool (and cluster)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ovm3vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Description: "the password for the host",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"public_vswitch_name": {
				Description: "Name of virtual switch used for public traffic in the cluster. This would override zone wide traffic label setting.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_vswitch_type": {
				Description: "Type of virtual switch used for public traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pod_id": {
				Description: "The Pod ID for the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Description: "the URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"username": {
				Description: "the username for the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vsm_ip_address": {
				Description: "the ipaddress of the VSM associated with this cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vsm_password": {
				Description: "the password for the VSM associated with this cluster",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"vsm_username": {
				Description: "the username for the VSM associated with this cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"zone_id": {
				Description: "the Zone ID for the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"zone_name": {
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
		"allocation_state":        cluster.Allocationstate,
		"cluster_name":            cluster.Name,
		"name":                    cluster.Name,
		"cluster_type":            cluster.Clustertype,
		"hypervisor":              cluster.Hypervisortype,
		"ovm3_vip":                cluster.Ovm3vip,
		"ovm3vip":                 cluster.Ovm3vip,
		"pod_id":                  cluster.Podid,
		"pod_name":                cluster.Podname,
		"zone_id":                 cluster.Zoneid,
		"zone_name":               cluster.Zonename,
		"managed_state":           cluster.Managedstate,
		"cpu_overcommit_ratio":    cluster.Cpuovercommitratio,
		"memory_overcommit_ratio": cluster.Memoryovercommitratio,
		"arch":                    cluster.Arch,
		"capacity":                dsFlattenClusterCapacity(cluster.Capacity),
	}

	// Set fields that may not be available in all cluster responses to empty strings
	// These are typically only available during cluster creation/configuration
	emptyStringFields := []string{
		"guest_vswitch_name",
		"guest_vswitch_type",
		"ovm3_cluster",
		"ovm3_pool",
		"password",
		"public_vswitch_name",
		"public_vswitch_type",
		"url",
		"username",
		"vsm_ip_address",
		"vsm_password",
		"vsm_username",
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			log.Printf("[WARN] Error setting %s: %s", k, err)
		}
	}

	for _, field := range emptyStringFields {
		if err := d.Set(field, ""); err != nil {
			log.Printf("[WARN] Error setting %s: %s", field, err)
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
