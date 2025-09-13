package collector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tailscale.com/client/tailscale/v2"
)

func TestTailscaleCollector_WithMockServer(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/tailnet/example.com/devices":
			devices := map[string][]tailscale.Device{
				"devices": {
				{
					ID:        "device1",
					NodeID:    "nDevice1CNTRL",
					Name:      "test-device-1.example.ts.net",
					Hostname:  "test-device-1",
					User:      "user1@example.com",
					OS:        "linux",
					Addresses: []string{"100.64.1.1"},
					Authorized: true,
					IsExternal: false,
					LastSeen: tailscale.Time{Time: time.Now().Add(-1 * time.Minute)},
					Created:  tailscale.Time{Time: time.Now().Add(-24 * time.Hour)},
					Expires:  tailscale.Time{Time: time.Now().Add(30 * 24 * time.Hour)},
					UpdateAvailable:           false,
					KeyExpiryDisabled:         false,
					BlocksIncomingConnections: false,
					EnabledRoutes:             []string{"10.0.0.0/16"},
					AdvertisedRoutes:          []string{"10.0.0.0/16", "192.168.1.0/24"},
					Tags:                      []string{"tag:server"},
					ClientConnectivity: &tailscale.ClientConnectivity{
						Endpoints:             []string{"192.168.1.100:41641"},
						DERP:                  "1",
						MappingVariesByDestIP: false,
						DERPLatency: map[string]tailscale.DERPRegion{
							"nyc": {LatencyMilliseconds: 25.5},
							"sfo": {LatencyMilliseconds: 75.2},
						},
					},
				},
				{
					ID:        "device2",
					NodeID:    "nDevice2CNTRL",
					Name:      "test-device-2.example.ts.net",
					Hostname:  "test-device-2",
					User:      "user2@example.com",
					OS:        "windows",
					Addresses: []string{"100.64.1.2"},
					Authorized: false,
					IsExternal: true,
					LastSeen:   tailscale.Time{Time: time.Now().Add(-10 * time.Minute)},
					Created:    tailscale.Time{Time: time.Now().Add(-48 * time.Hour)},
					Expires:    tailscale.Time{Time: time.Now().Add(60 * 24 * time.Hour)},
					UpdateAvailable:           true,
					KeyExpiryDisabled:         true,
					BlocksIncomingConnections: true,
					EnabledRoutes:             []string{},
					AdvertisedRoutes:          []string{},
					Tags:                      []string{"tag:client"},
				},
			}
			json.NewEncoder(w).Encode(devices)

		case "/api/v2/tailnet/example.com/users":
			users := []tailscale.User{
				{
					ID:          "user1",
					LoginName:   "user1@example.com",
					DisplayName: "User One",
					Role:        tailscale.UserRoleAdmin,
					Status:      tailscale.UserStatusActive,
					Type:        tailscale.UserTypeMember,
					Created:     time.Now().Add(-30 * 24 * time.Hour),
					LastSeen:    time.Now().Add(-1 * time.Hour),
				},
				{
					ID:          "user2",
					LoginName:   "user2@example.com",
					DisplayName: "User Two",
					Role:        tailscale.UserRoleMember,
					Status:      tailscale.UserStatusActive,
					Type:        tailscale.UserTypeShared,
					Created:     time.Now().Add(-15 * 24 * time.Hour),
					LastSeen:    time.Now().Add(-2 * time.Hour),
				},
			}
			json.NewEncoder(w).Encode(users)

		case "/api/v2/tailnet/example.com/keys":
			keys := []tailscale.Key{
				{
					ID:          "key1",
					Key:         "tskey-auth-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Description: "Test Auth Key",
					Created:     time.Now().Add(-7 * 24 * time.Hour),
					Expires:     time.Now().Add(83 * 24 * time.Hour),
				},
				{
					ID:          "key2",
					Key:         "tskey-api-xxxxxxxxxxxx-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
					Description: "Test API Key",
					Created:     time.Now().Add(-14 * 24 * time.Hour),
					Expires:     time.Now().Add(76 * 24 * time.Hour),
					Revoked:     time.Now().Add(-1 * 24 * time.Hour),
				},
			}
			json.NewEncoder(w).Encode(keys)

		case "/api/v2/tailnet/example.com/dns/preferences":
			prefs := tailscale.DNSPreferences{
				MagicDNS: true,
			}
			json.NewEncoder(w).Encode(prefs)

		case "/api/v2/tailnet/example.com/dns/nameservers":
			nameservers := []string{"8.8.8.8", "1.1.1.1"}
			json.NewEncoder(w).Encode(nameservers)

		case "/api/v2/tailnet/example.com/dns/searchpaths":
			searchPaths := []string{"example.com", "internal.com"}
			json.NewEncoder(w).Encode(searchPaths)

		case "/api/v2/tailnet/example.com/settings":
			settings := tailscale.TailnetSettings{
				DevicesApprovalOn:                        true,
				DevicesAutoUpdatesOn:                     false,
				DevicesKeyDurationDays:                   180,
				UsersApprovalOn:                          false,
				UsersRoleAllowedToJoinExternalTailnets:   tailscale.RoleAllowedToJoinExternalTailnetsAdmin,
				NetworkFlowLoggingOn:                     false,
				RegionalRoutingOn:                        true,
				ACLsExternallyManagedOn:                  false,
				PostureIdentityCollectionOn:              false,
			}
			json.NewEncoder(w).Encode(settings)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create collector with custom HTTP client that points to our mock server
	httpClient := &http.Client{
		Transport: &testTransport{baseURL: server.URL},
	}

	// Create collector using the proper constructor
	collector, err := NewTailscaleCollector(nil, httpClient, "example.com")
	require.NoError(t, err)

	// Test metric collection
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Gather metrics
	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Test 1: API should be up
	upMetric := findMetricFamily(metrics, "tailscale_up")
	require.NotNil(t, upMetric, "tailscale_up metric should exist")
	assert.Equal(t, 1.0, upMetric.GetMetric()[0].GetGauge().GetValue(), "API should be up")

	// Test 2: Device info metrics
	deviceInfoMetrics := findMetricFamily(metrics, "tailscale_device_info")
	require.NotNil(t, deviceInfoMetrics, "device info metrics should exist")
	assert.Len(t, deviceInfoMetrics.GetMetric(), 2, "should have 2 device info metrics")

	// Test 3: Online/offline status
	deviceOnlineMetrics := findMetricFamily(metrics, "tailscale_device_online")
	require.NotNil(t, deviceOnlineMetrics, "device online metrics should exist")

	onlineCount := 0
	for _, metric := range deviceOnlineMetrics.GetMetric() {
		if metric.GetGauge().GetValue() == 1.0 {
			onlineCount++
		}
	}
	assert.Equal(t, 1, onlineCount, "only one device should be online")

	// Test 4: Authorization status
	deviceAuthorizedMetrics := findMetricFamily(metrics, "tailscale_device_authorized")
	require.NotNil(t, deviceAuthorizedMetrics, "device authorized metrics should exist")

	authorizedCount := 0
	for _, metric := range deviceAuthorizedMetrics.GetMetric() {
		if metric.GetGauge().GetValue() == 1.0 {
			authorizedCount++
		}
	}
	assert.Equal(t, 1, authorizedCount, "only one device should be authorized")

	// Test 5: User metrics
	usersTotal := findMetricFamily(metrics, "tailscale_users_total")
	require.NotNil(t, usersTotal, "users total metric should exist")
	assert.Equal(t, 2.0, usersTotal.GetMetric()[0].GetGauge().GetValue(), "should have 2 users")

	// Test 6: Key metrics
	keysTotal := findMetricFamily(metrics, "tailscale_keys_total")
	require.NotNil(t, keysTotal, "keys total metric should exist")
	assert.Equal(t, 2.0, keysTotal.GetMetric()[0].GetGauge().GetValue(), "should have 2 keys")

	// Test 7: DNS metrics
	dnsNameservers := findMetricFamily(metrics, "tailscale_dns_nameservers")
	require.NotNil(t, dnsNameservers, "DNS nameservers metric should exist")
	assert.Equal(t, 2.0, dnsNameservers.GetMetric()[0].GetGauge().GetValue(), "should have 2 nameservers")

	dnsSearchPaths := findMetricFamily(metrics, "tailscale_dns_search_paths")
	require.NotNil(t, dnsSearchPaths, "DNS search paths metric should exist")
	assert.Equal(t, 2.0, dnsSearchPaths.GetMetric()[0].GetGauge().GetValue(), "should have 2 search paths")
}

func TestTailscaleCollector_APIFailure(t *testing.T) {
	// Test with server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create collector with client that points to error server
	collector, err := NewTailscaleCollector(nil, &http.Client{
		Transport: &testTransport{baseURL: server.URL},
	}, "example.com")
	require.NoError(t, err)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Should have 'up' metric set to 0
	upMetric := findMetricFamily(metrics, "tailscale_up")
	require.NotNil(t, upMetric, "up metric should exist even on API failure")
	assert.Equal(t, 0.0, upMetric.GetMetric()[0].GetGauge().GetValue(), "API should be down")

	// Should not have device metrics when API fails
	deviceMetrics := findMetricFamily(metrics, "tailscale_device_info")
	assert.Nil(t, deviceMetrics, "should not have device metrics when API fails")
}

func TestTailscaleCollector_EmptyResponses(t *testing.T) {
	// Test with empty responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/tailnet/example.com/devices":
			devices := map[string][]tailscale.Device{
				"devices": {},
			}
			json.NewEncoder(w).Encode(devices)

		case "/api/v2/tailnet/example.com/users":
			users := map[string][]tailscale.User{
				"users": {},
			}
			json.NewEncoder(w).Encode(users)

		case "/api/v2/tailnet/example.com/keys":
			keys := map[string][]tailscale.Key{
				"keys": {},
			}
			json.NewEncoder(w).Encode(keys)

		case "/api/v2/tailnet/example.com/dns/preferences":
			prefs := tailscale.DNSPreferences{
				MagicDNS: false,
			}
			json.NewEncoder(w).Encode(prefs)

		case "/api/v2/tailnet/example.com/dns/nameservers":
			nameservers := map[string][]string{
				"dns": {},
			}
			json.NewEncoder(w).Encode(nameservers)

		case "/api/v2/tailnet/example.com/dns/searchpaths":
			searchPaths := map[string][]string{
				"searchPaths": {},
			}
			json.NewEncoder(w).Encode(searchPaths)

		case "/api/v2/tailnet/example.com/settings":
			settings := tailscale.TailnetSettings{
				DevicesApprovalOn:                        false,
				DevicesAutoUpdatesOn:                     false,
				DevicesKeyDurationDays:                   180,
				UsersApprovalOn:                          false,
				UsersRoleAllowedToJoinExternalTailnets:   tailscale.RoleAllowedToJoinExternalTailnetsNone,
				NetworkFlowLoggingOn:                     false,
				RegionalRoutingOn:                        false,
				ACLsExternallyManagedOn:                  false,
				PostureIdentityCollectionOn:              false,
			}
			json.NewEncoder(w).Encode(settings)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create collector with client that points to mock server
	collector, err := NewTailscaleCollector(nil, &http.Client{
		Transport: &testTransport{baseURL: server.URL},
	}, "example.com")
	require.NoError(t, err)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Should have 'up' metric set to 1
	upMetric := findMetricFamily(metrics, "tailscale_up")
	require.NotNil(t, upMetric)
	assert.Equal(t, 1.0, upMetric.GetMetric()[0].GetGauge().GetValue(), "API should be up")

	// Should have zero counts for all resources
	usersTotal := findMetricFamily(metrics, "tailscale_users_total")
	require.NotNil(t, usersTotal)
	assert.Equal(t, 0.0, usersTotal.GetMetric()[0].GetGauge().GetValue(), "should have 0 users")

	keysTotal := findMetricFamily(metrics, "tailscale_keys_total")
	require.NotNil(t, keysTotal)
	assert.Equal(t, 0.0, keysTotal.GetMetric()[0].GetGauge().GetValue(), "should have 0 keys")
}

func TestTailscaleCollector_MetricLabels(t *testing.T) {
	// Test that metric labels are correctly set
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/tailnet/example.com/devices":
			devices := map[string][]tailscale.Device{
				"devices": {
					{
						ID:            "test-device",
						NodeID:        "nTestCNTRL",
						Name:          "test.example.ts.net",
						Hostname:      "test-host",
						User:          "test@example.com",
						OS:            "darwin",
						Addresses:     []string{"100.64.1.100"},
						ClientVersion: "1.50.0",
						Authorized:    true,
						IsExternal:    false,
						LastSeen:      tailscale.Time{Time: time.Now().Add(-30 * time.Second)},
					},
				},
			}
			json.NewEncoder(w).Encode(devices)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create collector with client that points to mock server
	client := &tailscale.Client{
		APIKey:  "test-key",
		Tailnet: "example.com",
		HTTP: &http.Client{
			Transport: &testTransport{baseURL: server.URL},
		},
	}

	collector, err := NewTailscaleCollector(nil, &http.Client{
		Transport: &testTransport{baseURL: server.URL},
	}, "example.com")
	require.NoError(t, err)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Check device info metric labels
	deviceInfoMetrics := findMetricFamily(metrics, "tailscale_device_info")
	require.NotNil(t, deviceInfoMetrics)
	require.Len(t, deviceInfoMetrics.GetMetric(), 1)

	metric := deviceInfoMetrics.GetMetric()[0]
	labels := make(map[string]string)
	for _, labelPair := range metric.GetLabel() {
		labels[labelPair.GetName()] = labelPair.GetValue()
	}

	// Verify all expected labels are present and correct
	assert.Equal(t, "test-device", labels["id"])
	assert.Equal(t, "test.example.ts.net", labels["name"])
	assert.Equal(t, "test-host", labels["hostname"])
	assert.Equal(t, "test@example.com", labels["user"])
	assert.Equal(t, "darwin", labels["os"])
	assert.Equal(t, "1.50.0", labels["client_version"])
	assert.Equal(t, "100.64.1.100", labels["tailscale_ip"])
}

// Helper functions
// testTransport is a custom HTTP transport that redirects all requests to our test server
type testTransport struct {
	baseURL string
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create a new request URL by combining test server base URL with the original path
	testURL := t.baseURL + req.URL.Path
	if req.URL.RawQuery != "" {
		testURL += "?" + req.URL.RawQuery
	}

	// Parse the new URL
	newURL, err := req.URL.Parse(testURL)
	if err != nil {
		return nil, err
	}

	// Create a new request with the test server URL
	newReq := req.Clone(req.Context())
	newReq.URL = newURL

	// Use default transport to make the actual request
	return http.DefaultTransport.RoundTrip(newReq)
}

func findMetricFamily(metrics []*dto.MetricFamily, name string) *dto.MetricFamily {
	for _, mf := range metrics {
		if mf.GetName() == name {
			return mf
		}
	}
	return nil
}


