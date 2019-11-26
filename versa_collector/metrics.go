package versa_collector

import "github.com/prometheus/client_golang/prometheus"

var (
	metricsDesc = []*prometheus.Desc{
		versaSitesAvailabilityPercent,
		versaApplicationUsageBandwidthRxBps,
		versaApplicationUsageBandwidthTxBps,
		versaApplicationUsageVolumeRxByte,
		versaApplicationUsageVolumeTxByte,
		versaSiteCircuitBandwidthUsageTxBps,
		versaSiteCircuitBandwidthUsageRxBps,
		versaApplianceCPULoadPercent,
		versaApplianceMemoryLoadPercent,
		versaApplianceDiskLoadPercent,
		versaApplianceSessionsLoad,
		versaSLADelay,
		versaSLAJitterFwd,
		versaSLAJitterRev,
		versaSLALossFwd,
		versaSLALossRev,
	}

	versaSitesAvailabilityPercent = prometheus.NewDesc(
		"versa_analytics_sites_availability_percent",
		"The availability percentage for the particular site",
		[]string{"tenant", "site"},
		nil,
	)

	versaApplicationUsageBandwidthRxBps = prometheus.NewDesc(
		"versa_analytics_application_usage_bandwidth_rx_bps",
		"The application RX bandwidth usage rate in bits per second",
		[]string{"tenant", "site", "app_name", "client_ip", "circuit"},
		nil,
	)

	versaApplicationUsageBandwidthTxBps = prometheus.NewDesc(
		"versa_analytics_application_usage_bandwidth_tx_bps",
		"The application TX bandwidth usage rate in bits per second",
		[]string{"tenant", "site", "app_name", "client_ip", "circuit"},
		nil,
	)

	versaApplicationUsageVolumeRxByte = prometheus.NewDesc(
		"versa_analytics_application_usage_volume_rx_bytes",
		"The application RX volume usage in bytes",
		[]string{"tenant", "site", "app_name", "client_ip", "circuit"},
		nil,
	)

	versaApplicationUsageVolumeTxByte = prometheus.NewDesc(
		"versa_analytics_application_usage_volume_tx_bytes",
		"The application TX volume usage in bytes",
		[]string{"tenant", "site", "app_name", "client_ip", "circuit"},
		nil,
	)

	versaSiteCircuitBandwidthUsageTxBps = prometheus.NewDesc(
		"versa_analytics_site_circuit_usage_bandwidth_tx_bps",
		"The site circuit TX bandwidth usage rate in bits per second",
		[]string{"tenant", "site", "circuit"},
		nil,
	)

	versaSiteCircuitBandwidthUsageRxBps = prometheus.NewDesc(
		"versa_analytics_site_circuit_usage_bandwidth_rx_bps",
		"The site circuit RX bandwidth usage rate in bits per second",
		[]string{"tenant", "site", "circuit"},
		nil,
	)

	versaApplianceCPULoadPercent = prometheus.NewDesc(
		"versa_analytics_appliance_cpu_load_pct",
		"The appliance CPU Load in percentage",
		[]string{"tenant", "site"},
		nil,
	)

	versaApplianceMemoryLoadPercent = prometheus.NewDesc(
		"versa_analytics_appliance_memory_load_pct",
		"The appliance Memory Load in percentage",
		[]string{"tenant", "site"},
		nil,
	)

	versaApplianceDiskLoadPercent = prometheus.NewDesc(
		"versa_analytics_appliance_disk_load_pct",
		"The appliance Disk Load in percentage",
		[]string{"tenant", "site"},
		nil,
	)

	versaApplianceSessionsLoad = prometheus.NewDesc(
		"versa_analytics_appliance_sessions_load",
		"The appliance current sessions",
		[]string{"tenant", "site"},
		nil,
	)

	versaSLADelay = prometheus.NewDesc(
		"versa_analytics_site_slam_delay_ms",
		"The SLA probe delay reported in milliseconds",
		[]string{"tenant", "site", "source_site", "destination_site", "source_circuit", "destination_circuit"},
		nil,
	)

	versaSLAJitterFwd = prometheus.NewDesc(
		"versa_analytics_site_slam_jitter_fwd_ms",
		"The SLA probe forward jitter reported in milliseconds",
		[]string{"tenant", "site", "source_site", "destination_site", "source_circuit", "destination_circuit"},
		nil,
	)
	versaSLAJitterRev = prometheus.NewDesc(
		"versa_analytics_site_slam_jitter_rcv_ms",
		"The SLA probe reverse jitter reported in milliseconds",
		[]string{"tenant", "site", "source_site", "destination_site", "source_circuit", "destination_circuit"},
		nil,
	)
	versaSLALossFwd = prometheus.NewDesc(
		"versa_analytics_site_slam_loss_fwd_pct",
		"The SLA probe forward loss reported in percent",
		[]string{"tenant", "site", "source_site", "destination_site", "source_circuit", "destination_circuit"},
		nil,
	)
	versaSLALossRev = prometheus.NewDesc(
		"versa_analytics_site_slam_loss_rcv_pct",
		"The SLA probe reverse loss reported in percent",
		[]string{"tenant", "site", "source_site", "destination_site", "source_circuit", "destination_circuit"},
		nil,
	)
)
