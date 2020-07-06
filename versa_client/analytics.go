package versa_client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"time"

	"github.com/lucabrasi83/peppamon_versa/logging"
)

const (
	longReportPrecision   = "15minutesAgo"
	mediumReportPrecision = "5minutesAgo"
	shortReportPrecision  = "1minutesAgo"
)

type VersaAnalyticsClient struct {
	Hostname   string
	Protocol   string
	Username   string
	Password   string
	HttpClient *http.Client
	Tenants    VersaTenantList
}

type VersaTenantList []struct {
	TenantName string `json:"name"`
}

type VersaSitesAvailability struct {
	TenantName string
	SitesList  []struct {
		SiteName        string
		AvailabilityPct float64
	}
}

type VersaApplicationUsageRate struct {
	TenantName string
	QTime      int `json:"qTime"`
	Data       []struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Metric     string          `json:"metric"`
		MetricName string          `json:"metricName"`
		Label      string          `json:"label"`
		Data       [][]interface{} `json:"data"`
	} `json:"data"`
}

type VersaApplicationUsageVolume struct {
	TenantName string
	QTime      int `json:"qTime"`
	Data       []struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Metric     string          `json:"metric"`
		MetricName string          `json:"metricName"`
		Label      string          `json:"label"`
		Data       [][]interface{} `json:"data"`
	} `json:"data"`
}

type VersaSiteSLAMetrics struct {
	TenantName string
	QTime      int `json:"qTime"`
	Data       []struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Metric     string          `json:"metric"`
		MetricName string          `json:"metricName"`
		Label      string          `json:"label"`
		Data       [][]interface{} `json:"data"`
	} `json:"data"`
}

type VersaSiteBandwidthUsage struct {
	TenantName string
	QTime      int `json:"qTime"`
	Data       []struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Metric     string          `json:"metric"`
		MetricName string          `json:"metricName"`
		Label      string          `json:"label"`
		Data       [][]interface{} `json:"data"`
	} `json:"data"`
}

type VersaAppliancePerformance struct {
	TenantName string
	QTime      int `json:"qTime"`
	Data       []struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Metric     string          `json:"metric"`
		MetricName string          `json:"metricName"`
		Label      string          `json:"label"`
		Data       [][]interface{} `json:"data"`
	} `json:"data"`
}

func NewVersaAnalyticsClient() *VersaAnalyticsClient {
	cookieJar, _ := cookiejar.New(nil)

	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	versaHTTPClient := &http.Client{
		Timeout:   10 * time.Minute,
		Jar:       cookieJar,
		Transport: httpTransport,
	}

	return &VersaAnalyticsClient{
		Hostname:   os.Getenv("PEPPAMON_VERSA_ANALYTICS_HOSTNAME"),
		Username:   os.Getenv("PEPPAMON_VERSA_ANALYTICS_USERNAME"),
		Password:   os.Getenv("PEPPAMON_VERSA_ANALYTICS_PASSWORD"),
		Protocol:   "https",
		HttpClient: versaHTTPClient,
	}
}

func (v *VersaAnalyticsClient) Login() error {

	url := fmt.Sprintf("%s://%s/versa/login?username=%s&password=%s", v.Protocol, v.Hostname, v.Username, v.Password)

	queryTitle := "Versa Analytics Login"

	httpNewReq, err := http.NewRequest("POST", url, nil)

	if err != nil {
		logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
		return err
	}

	cookieRes, err := v.HttpClient.Do(httpNewReq)

	if err != nil {
		logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
		return err
	}

	if cookieRes.StatusCode != http.StatusOK || cookieRes.StatusCode > http.StatusAccepted {
		logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
			cookieRes.StatusCode, queryTitle)
		return fmt.Errorf("versa analytics responded with HTTP error code %v for %v", cookieRes.StatusCode, queryTitle)
	}
	return nil
}

func (v *VersaAnalyticsClient) GetTenantList() error {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Versa Tenants")

	url := fmt.Sprintf("%s://%s/versa/analytics/v1.0.0/data/provider/features/SDWAN/tenants?count=-1", v.Protocol, v.Hostname)

	queryTitle := "Get Tenants List"

	httpNewReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
		return err
	}

	httpNewReq.Header.Add("Content-Type", "application/json")

	tenantsRes, err := v.HttpClient.Do(httpNewReq)

	if err != nil {
		logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
		return err
	}

	if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
		logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
			tenantsRes.StatusCode, queryTitle)
		return fmt.Errorf("versa analytics responded with HTTP error code %v for %v", tenantsRes.StatusCode, queryTitle)
	}

	defer func() {

		errBodyClose := tenantsRes.Body.Close()

		if errBodyClose != nil {
			logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
		}
	}()

	var tenantList VersaTenantList

	err = json.NewDecoder(tenantsRes.Body).Decode(&tenantList)

	if err != nil {
		logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
		return err
	}

	v.Tenants = tenantList

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Versa Tenants")

	return nil

}

func (v *VersaAnalyticsClient) GetSitesAvailability() ([]VersaSitesAvailability, error) {
	logging.PeppaMonLog("info", "Started Batch Job to fetch Sites Availability Metrics")
	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	availabilitySitesSlice := make([]VersaSitesAvailability, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SDWAN/?qt=stats&start-date=5minutesAgo&end-date=today&q=site"+
				"&metrics"+
				"=availability"+
				"&count=-1", v.Protocol, v.Hostname, t.TenantName)

			queryTitle := "Get Sites Availability"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)
				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var sitesAvailabilityStats interface{}

			err = json.NewDecoder(tenantsRes.Body).Decode(&sitesAvailabilityStats)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return
			}

			availabilitySiteObj := VersaSitesAvailability{TenantName: t.TenantName}

			// JSON object parses into a map with string keys
			itemsMap := sitesAvailabilityStats.(map[string]interface{})

			for key, val := range itemsMap {
				if key == "stats" {
					for site, stats := range val.(map[string]interface{}) {
						availabilityPct := stats.(map[string]interface{})["mean"].(float64)

						siteObj := struct {
							SiteName        string
							AvailabilityPct float64
						}{SiteName: site, AvailabilityPct: availabilityPct}
						availabilitySiteObj.SitesList = append(availabilitySiteObj.SitesList, siteObj)

					}
				}
			}
			mu.Lock()
			availabilitySitesSlice = append(availabilitySitesSlice, availabilitySiteObj)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Sites Availability Metrics")
	return availabilitySitesSlice, nil
}

func (v *VersaAnalyticsClient) GetSitesApplicationUsageRate() ([]VersaApplicationUsageRate, error) {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Application Usage Rate Metrics")

	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	applicationUsageSlice := make([]VersaApplicationUsageRate, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SDWAN/?start-date=%s&end-date=today&q=appUser(site,appId,user,"+
				"accCkt)&qt=timeseries&ds=aggregate&gap=1MINUTE&metrics=bw-rx&metrics=bw-tx&count=15000",
				v.Protocol, v.Hostname, t.TenantName, longReportPrecision)

			queryTitle := "Get Application Usage Rate"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)

				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var tenantApplicationUsage VersaApplicationUsageRate

			err = json.NewDecoder(tenantsRes.Body).Decode(&tenantApplicationUsage)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return

			}

			tenantApplicationUsage.TenantName = t.TenantName

			mu.Lock()
			applicationUsageSlice = append(applicationUsageSlice, tenantApplicationUsage)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Application Usage Rate Metrics")
	return applicationUsageSlice, nil
}

func (v *VersaAnalyticsClient) GetSitesApplicationUsageVolume() ([]VersaApplicationUsageVolume, error) {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Application Usage Volume Metrics")

	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	applicationUsageSlice := make([]VersaApplicationUsageVolume, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SDWAN/?start-date=%s&end-date=today&q=appUser(site,appId,user,"+
				"accCkt)&qt=timeseries&ds=aggregate&gap=1MINUTE&metrics=volume-rx&metrics=volume-tx&count=15000",
				v.Protocol, v.Hostname, t.TenantName, longReportPrecision)

			queryTitle := "Get Application Usage Volume"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)

				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var tenantApplicationUsage VersaApplicationUsageVolume

			err = json.NewDecoder(tenantsRes.Body).Decode(&tenantApplicationUsage)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return

			}

			tenantApplicationUsage.TenantName = t.TenantName

			mu.Lock()
			applicationUsageSlice = append(applicationUsageSlice, tenantApplicationUsage)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Application Usage Volume Metrics")
	return applicationUsageSlice, nil
}

func (v *VersaAnalyticsClient) GetSitesCircuitBandwidthUsage() ([]VersaSiteBandwidthUsage, error) {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Site Circuits Usage Metrics")

	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	applicationUsageSlice := make([]VersaSiteBandwidthUsage, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SDWAN/?start-date=%s&end-date=today&q=linkUsage(site,accCkt)&qt=timeseries&ds=aggregate&gap=1MINUTE&metrics=bw-rx&metrics=bw-tx&count=-1",
				v.Protocol, v.Hostname, t.TenantName, mediumReportPrecision)

			queryTitle := "Get Site Circuits Usage"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)
				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var tenantSiteCircuitUsage VersaSiteBandwidthUsage

			err = json.NewDecoder(tenantsRes.Body).Decode(&tenantSiteCircuitUsage)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return

			}

			tenantSiteCircuitUsage.TenantName = t.TenantName

			mu.Lock()
			applicationUsageSlice = append(applicationUsageSlice, tenantSiteCircuitUsage)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Site Circuits Usage Metrics")
	return applicationUsageSlice, nil
}

func (v *VersaAnalyticsClient) GetSitesSLAMetrics() ([]VersaSiteSLAMetrics, error) {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Site SLA Metrics")

	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	metricsIPSLASlice := make([]VersaSiteSLAMetrics, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SDWAN/?start-date=%s&end-date=today&q=slam(localSite,remoteSite,"+
				"localAccCkt,remoteAccCkt)&qt=timeseries&gap=1MINUTE&ds=aggregate&metrics=fwdDelayVar&metrics"+
				"=revDelayVar&metrics=delay&count=-1&metrics=fwdLossRatio&metrics=revLossRatio",
				v.Protocol, v.Hostname, t.TenantName, longReportPrecision)

			queryTitle := "Get Site SLA Metrics"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)
				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var tenantIPSLAMetrics VersaSiteSLAMetrics

			err = json.NewDecoder(tenantsRes.Body).Decode(&tenantIPSLAMetrics)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return

			}

			tenantIPSLAMetrics.TenantName = t.TenantName

			mu.Lock()
			metricsIPSLASlice = append(metricsIPSLASlice, tenantIPSLAMetrics)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Site SLA Metrics")
	return metricsIPSLASlice, nil
}

func (v *VersaAnalyticsClient) GetApplianceComputePerf() ([]VersaAppliancePerformance, error) {

	logging.PeppaMonLog("info", "Started Batch Job to fetch Appliance Compute Metrics")

	var wg sync.WaitGroup
	wg.Add(len(v.Tenants))

	var mu sync.Mutex

	applicationUsageSlice := make([]VersaAppliancePerformance, 0, len(v.Tenants))

	for _, tenant := range v.Tenants {

		go func(t struct{ TenantName string }) {

			defer wg.Done()

			url := fmt.Sprintf("%s://%s/versa/analytics/v1.0."+
				"0/data/provider/tenants/%s/features/SYSTEM?start-date=%s&q=applMonitor&qt=timeseries&gap=1MINUTE&ds"+
				"=aggregate&count=-1&metrics=CPULOAD&metrics=MEMLOAD&metrics=DISKLOAD&metrics=SESSLOAD",
				v.Protocol, v.Hostname, t.TenantName, longReportPrecision)

			queryTitle := "Get Appliance Compute Performance"

			httpNewReq, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logging.PeppaMonLog("error", "unable to build HTTP request for %v with error %v", queryTitle, err)
				return

			}

			httpNewReq.Header.Add("Content-Type", "application/json")

			tenantsRes, err := v.HttpClient.Do(httpNewReq)

			if err != nil {
				logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
				return

			}

			if tenantsRes.StatusCode != http.StatusOK || tenantsRes.StatusCode > http.StatusAccepted {
				logging.PeppaMonLog("error", "Versa Analytics responded with HTTP error code %v for %v",
					tenantsRes.StatusCode, queryTitle)
				return
			}

			defer func() {

				errBodyClose := tenantsRes.Body.Close()

				if errBodyClose != nil {
					logging.PeppaMonLog("error", "HTTP request for %v failed with error %v", queryTitle, err)
					return
				}
			}()

			var tenantSiteCircuitUsage VersaAppliancePerformance

			err = json.NewDecoder(tenantsRes.Body).Decode(&tenantSiteCircuitUsage)

			if err != nil {
				logging.PeppaMonLog("error", "Unable to decode JSON response from %v with error %v", queryTitle, err)
				return

			}

			tenantSiteCircuitUsage.TenantName = t.TenantName

			mu.Lock()
			applicationUsageSlice = append(applicationUsageSlice, tenantSiteCircuitUsage)
			mu.Unlock()

		}(struct{ TenantName string }(tenant))

	}
	wg.Wait()

	logging.PeppaMonLog("info", "Completed Batch Job to fetch Appliance Compute Metrics")
	return applicationUsageSlice, nil
}
