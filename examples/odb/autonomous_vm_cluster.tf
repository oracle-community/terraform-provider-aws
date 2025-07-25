//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

//Autonomous VM Cluster with default maintenance window
resource "aws_odb_cloud_autonomous_vm_cluster" "test" {
  cloud_exadata_infrastructure_id         = "<exadata_infra_id>" //refer your exadata infra id
  odb_network_id                          = "<odb_net_id>" //refer_your_odb_net_id
  display_name             				= "Ofake-avmc-my_avmc"
  autonomous_data_storage_size_in_tbs     = 5
  memory_per_oracle_compute_unit_in_gbs   = 2
  total_container_databases               = 1
  cpu_core_count_per_node                 = 40
  license_model                                = "LICENSE_INCLUDED"
  db_servers								   = [] //ids of db server. refer your exa infra
  scan_listener_port_tls = 8561
  scan_listener_port_non_tls = 1024
  maintenance_window = {
    preference = "NO_PREFERENCE"
    days_of_week =	[]
    hours_of_day =	[]
    months = []
    weeks_of_month =[]
    lead_time_in_weeks = 0
  }

}
