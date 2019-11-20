package versa_collector

import (
	"strings"
	"sync"

	"github.com/lucabrasi83/peppamon_versa/logging"
	"github.com/lucabrasi83/peppamon_versa/versa_client"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	appUsageVolumeBytesLimit = 1000

	appUsageRateBpsLimit = 100
)

type VersaAnalyticsExporter struct {
	mu                   sync.Mutex
	VersaAnalyticsClient *versa_client.VersaAnalyticsClient
	Metrics              []prometheus.Metric
}

func NewVersaAnalyticsExporter() *VersaAnalyticsExporter {
	return &VersaAnalyticsExporter{
		VersaAnalyticsClient: versa_client.NewVersaAnalyticsClient(),
		Metrics:              nil,
	}
}

func (v *VersaAnalyticsExporter) Describe(ch chan<- *prometheus.Desc) {

	for _, desc := range v.Metrics {
		ch <- desc.Desc()
	}
}

func (v *VersaAnalyticsExporter) Collect(ch chan<- prometheus.Metric) {

	logging.PeppaMonLog("info", "Started Versa Analytics metrics scraping")

	// Bootstrap Versa Login and Tenant List building
	err := v.VersaAnalyticsClient.Login()

	if err != nil {
		return
	}

	err = v.VersaAnalyticsClient.GetTenantList()

	if err != nil {
		return
	}

	v.launchMetricsCollection()

	for _, metric := range v.Metrics {
		ch <- metric
	}

	// Empty out Metrics slice once scraping is done
	v.mu.Lock()
	v.Metrics = nil
	v.mu.Unlock()

	logging.PeppaMonLog("info", "Completed Versa Analytics metrics scraping")
}

func (v *VersaAnalyticsExporter) launchMetricsCollection() {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		v.versaSitesAvailabilityMetric()
	}()

	go func() {
		defer wg.Done()
		v.versaApplicationUsageRateMetric()
	}()

	go func() {
		defer wg.Done()
		v.versaSiteCircuitsUsageMetric()
	}()

	go func() {
		defer wg.Done()
		v.versaApplianceComputeUsageMetric()
	}()

	wg.Wait()
}

func (v *VersaAnalyticsExporter) versaSitesAvailabilityMetric() {
	sitesAvail, err := v.VersaAnalyticsClient.GetSitesAvailability()

	if err != nil {
		return
	}

	for _, tenant := range sitesAvail {
		if len(tenant.SitesList) > 0 {
			for _, site := range tenant.SitesList {
				metric := prometheus.MustNewConstMetric(
					versaSitesAvailabilityPercent,
					prometheus.GaugeValue,
					site.AvailabilityPct,
					tenant.TenantName, site.SiteName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()
			}
		}
	}

}

func (v *VersaAnalyticsExporter) versaApplicationUsageRateMetric() {
	appUsage, err := v.VersaAnalyticsClient.GetSitesApplicationUsageRate()

	if err != nil {
		return
	}

	for _, tenant := range appUsage {
		for _, siteUsage := range tenant.Data {

			if len(siteUsage.Data) == 0 {
				continue
			}

			metricTokens := strings.Split(siteUsage.Name, ",")

			siteName := metricTokens[0]
			appName := metricTokens[1]
			ipAddress := metricTokens[2]
			circuitName := metricTokens[3]

			appUsageRate := siteUsage.Data[0][1].(float64)

			// Filter app usage rate to avoid metric high cardinality
			if appUsageRate < appUsageRateBpsLimit {
				continue
			}

			switch siteUsage.Metric {
			case "bw-rx":
				metric :=
					prometheus.MustNewConstMetric(
						versaApplicationUsageBandwidthRxBps,
						prometheus.GaugeValue,
						appUsageRate,
						tenant.TenantName, siteName, appName, ipAddress, circuitName,
					)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "bw-tx":
				metric := prometheus.MustNewConstMetric(
					versaApplicationUsageBandwidthTxBps,
					prometheus.GaugeValue,
					appUsageRate,
					tenant.TenantName, siteName, appName, ipAddress, circuitName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()
			}

		}
	}
}

func (v *VersaAnalyticsExporter) versaApplicationUsageVolumeMetric() {
	appUsage, err := v.VersaAnalyticsClient.GetSitesApplicationUsageVolume()

	if err != nil {
		return
	}

	for _, tenant := range appUsage {
		for _, siteUsage := range tenant.Data {

			if len(siteUsage.Data) == 0 {
				continue
			}

			metricTokens := strings.Split(siteUsage.Name, ",")

			siteName := metricTokens[0]
			appName := metricTokens[1]
			ipAddress := metricTokens[2]
			circuitName := metricTokens[3]

			appUsageRate := siteUsage.Data[0][1].(float64)

			// Filter app usage volume to avoid metric high cardinality
			if appUsageRate < appUsageVolumeBytesLimit {
				continue
			}

			switch siteUsage.Metric {
			case "volume-rx":
				metric :=
					prometheus.MustNewConstMetric(
						versaApplicationUsageVolumeRxByte,
						prometheus.CounterValue,
						appUsageRate,
						tenant.TenantName, siteName, appName, ipAddress, circuitName,
					)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "volume-tx":
				metric := prometheus.MustNewConstMetric(
					versaApplicationUsageVolumeTxByte,
					prometheus.CounterValue,
					appUsageRate,
					tenant.TenantName, siteName, appName, ipAddress, circuitName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()
			}

		}
	}
}

func (v *VersaAnalyticsExporter) versaSiteCircuitsUsageMetric() {
	tenantCircuitUsage, err := v.VersaAnalyticsClient.GetSitesCircuitBandwidthUsage()

	if err != nil {
		return
	}

	for _, tenant := range tenantCircuitUsage {
		for _, siteUsage := range tenant.Data {

			if len(siteUsage.Data) == 0 {
				continue
			}

			metricTokens := strings.Split(siteUsage.Name, ",")

			siteName := metricTokens[0]
			circuitName := metricTokens[1]

			circuitUsageRate := siteUsage.Data[0][1].(float64)

			if circuitUsageRate == 0 {
				continue
			}

			switch siteUsage.Metric {
			case "bw-rx":
				metric :=
					prometheus.MustNewConstMetric(
						versaSiteCircuitBandwidthUsageRxBps,
						prometheus.GaugeValue,
						circuitUsageRate,
						tenant.TenantName, siteName, circuitName,
					)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "bw-tx":
				metric := prometheus.MustNewConstMetric(
					versaSiteCircuitBandwidthUsageTxBps,
					prometheus.GaugeValue,
					circuitUsageRate,
					tenant.TenantName, siteName, circuitName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()
			}

		}
	}
}

func (v *VersaAnalyticsExporter) versaApplianceComputeUsageMetric() {
	applianceComputePerfUsage, err := v.VersaAnalyticsClient.GetApplianceComputePerf()

	if err != nil {
		return
	}

	for _, tenant := range applianceComputePerfUsage {
		for _, applianceUsage := range tenant.Data {

			if len(applianceUsage.Data) == 0 {
				continue
			}

			siteName := applianceUsage.Name

			performanceUsageMetric := applianceUsage.Data[0][1].(float64)

			if performanceUsageMetric == 0 {
				continue
			}

			switch applianceUsage.Metric {
			case "cpuload":
				metric :=
					prometheus.MustNewConstMetric(
						versaApplianceCPULoadPercent,
						prometheus.GaugeValue,
						performanceUsageMetric,
						tenant.TenantName, siteName,
					)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "memload":
				metric := prometheus.MustNewConstMetric(
					versaApplianceMemoryLoadPercent,
					prometheus.GaugeValue,
					performanceUsageMetric,
					tenant.TenantName, siteName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "diskload":
				metric := prometheus.MustNewConstMetric(
					versaApplianceDiskLoadPercent,
					prometheus.GaugeValue,
					performanceUsageMetric,
					tenant.TenantName, siteName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()

			case "sessload":
				metric := prometheus.MustNewConstMetric(
					versaApplianceSessionsLoad,
					prometheus.GaugeValue,
					performanceUsageMetric,
					tenant.TenantName, siteName,
				)
				v.mu.Lock()
				v.Metrics = append(v.Metrics, metric)
				v.mu.Unlock()
			}

		}
	}
}
