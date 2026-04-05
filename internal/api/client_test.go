package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServer creates a mock AdGuard Home server with the given handler map.
// Keys are "METHOD /path", values are response bodies (JSON-serializable).
func newTestServer(t *testing.T, handlers map[string]any) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()

	for pattern, response := range handlers {
		resp := response // capture
		method := ""
		path := pattern
		for i, c := range pattern {
			if c == ' ' {
				method = pattern[:i]
				path = pattern[i+1:]
				break
			}
		}

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if method != "" && r.Method != method {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Check basic auth
			user, pass, ok := r.BasicAuth()
			if !ok || user != "admin" || pass != "secret" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if resp == nil {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{}"))
				return
			}
			_ = json.NewEncoder(w).Encode(resp)
		})
	}

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	// Client expects /control/ prefix — strip it from base URL since server has the full paths
	client := NewClient(server.URL, "admin", "secret")
	return server, client
}

func TestGetStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/status": ServerStatus{
			Version:           "v0.107.57",
			Running:           true,
			ProtectionEnabled: true,
			DNSPort:           53,
			HTTPPort:          80,
			DNSAddresses:      []string{"127.0.0.1"},
		},
	})

	status, err := client.GetStatus()
	require.NoError(t, err)
	assert.Equal(t, "v0.107.57", status.Version)
	assert.True(t, status.Running)
	assert.True(t, status.ProtectionEnabled)
	assert.Equal(t, 53, status.DNSPort)
	assert.Equal(t, 80, status.HTTPPort)
}

func TestGetStatus_AuthFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/control/status", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient(server.URL, "wrong", "creds")
	_, err := client.GetStatus()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestPing(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/status": map[string]any{"running": true},
	})

	err := client.Ping()
	assert.NoError(t, err)
}

func TestGetStats(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/stats": Stats{
			NumDNSQueries:       100000,
			NumBlockedFiltering: 15000,
			AvgProcessingTime:   0.125,
			TimeUnits:           "hours",
		},
	})

	stats, err := client.GetStats()
	require.NoError(t, err)
	assert.Equal(t, 100000, stats.NumDNSQueries)
	assert.Equal(t, 15000, stats.NumBlockedFiltering)
	assert.InDelta(t, 0.125, stats.AvgProcessingTime, 0.001)
}

func TestResetStats(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/stats_reset": nil,
	})

	err := client.ResetStats()
	assert.NoError(t, err)
}

func TestGetClients(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/clients": ClientsResponse{
			Clients: []ClientEntry{
				{
					Name:                     "Desktop",
					IDs:                      []string{"192.168.1.50", "10.0.0.2"},
					UseGlobalBlockedServices: false,
					BlockedServices:          []string{},
					FilteringEnabled:         true,
				},
				{
					Name:                     "Laptop",
					IDs:                      []string{"192.168.1.51"},
					UseGlobalBlockedServices: true,
					FilteringEnabled:         true,
				},
			},
		},
	})

	resp, err := client.GetClients()
	require.NoError(t, err)
	assert.Len(t, resp.Clients, 2)
	assert.Equal(t, "Desktop", resp.Clients[0].Name)
	assert.Equal(t, []string{"192.168.1.50", "10.0.0.2"}, resp.Clients[0].IDs)
	assert.False(t, resp.Clients[0].UseGlobalBlockedServices)
	assert.Empty(t, resp.Clients[0].BlockedServices)
}

func TestAddClient(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/clients/add": nil,
	})

	err := client.AddClient(ClientEntry{
		Name: "Test Client",
		IDs:  []string{"192.168.0.99"},
	})
	assert.NoError(t, err)
}

func TestDeleteClient(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/clients/delete": nil,
	})

	err := client.DeleteClient("Test Client")
	assert.NoError(t, err)
}

func TestUpdateClient(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/clients/update": nil,
	})

	err := client.UpdateClient("Old Name", ClientEntry{
		Name: "New Name",
		IDs:  []string{"192.168.0.99"},
	})
	assert.NoError(t, err)
}

func TestGetBlockedServices(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/blocked_services/get": BlockedServicesResponse{
			IDs: []string{"youtube", "tiktok", "facebook"},
		},
	})

	resp, err := client.GetBlockedServices()
	require.NoError(t, err)
	assert.Equal(t, []string{"youtube", "tiktok", "facebook"}, resp.IDs)
}

func TestGetAllServices(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/blocked_services/all": []ServiceEntry{
			{ID: "youtube", Name: "YouTube"},
			{ID: "tiktok", Name: "TikTok"},
			{ID: "facebook", Name: "Facebook"},
		},
	})

	services, err := client.GetAllServices()
	require.NoError(t, err)
	assert.Len(t, services, 3)
	assert.Equal(t, "youtube", services[0].ID)
}

func TestSetBlockedServices(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/blocked_services/set": nil,
	})

	err := client.SetBlockedServices(BlockedServicesResponse{
		IDs: []string{"youtube", "tiktok"},
	})
	assert.NoError(t, err)
}

func TestGetRewrites(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/rewrite/list": []RewriteEntry{
			{Domain: "home.example.com", Answer: "192.168.0.1"},
			{Domain: "nas.example.com", Answer: "192.168.0.105"},
		},
	})

	rewrites, err := client.GetRewrites()
	require.NoError(t, err)
	assert.Len(t, rewrites, 2)
	assert.Equal(t, "home.example.com", rewrites[0].Domain)
	assert.Equal(t, "192.168.0.1", rewrites[0].Answer)
}

func TestAddRewrite(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/rewrite/add": nil,
	})

	err := client.AddRewrite(RewriteEntry{Domain: "test.local", Answer: "10.0.0.1"})
	assert.NoError(t, err)
}

func TestDeleteRewrite(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/rewrite/delete": nil,
	})

	err := client.DeleteRewrite(RewriteEntry{Domain: "test.local", Answer: "10.0.0.1"})
	assert.NoError(t, err)
}

func TestGetQueryLog(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/querylog": QueryLogResponse{
			Data: []QueryLogEntry{
				{
					Client: "192.168.1.50",
					Status: "NotFilteredNotFound",
					Reason: "",
					Time:   "2026-04-03T15:00:00Z",
					Question: map[string]any{
						"name": "example.com",
						"type": "A",
					},
				},
			},
			Oldest: "2026-04-01T00:00:00Z",
		},
	})

	ql, err := client.GetQueryLog(10)
	require.NoError(t, err)
	assert.Len(t, ql.Data, 1)
	assert.Equal(t, "192.168.1.50", ql.Data[0].Client)
}

func TestGetFiltering(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/filtering/status": FilterStatus{
			Enabled:  true,
			Interval: 24,
			Filters: []Filter{
				{ID: 1, Name: "AdGuard DNS", Enabled: true, RulesCount: 55000, URL: "https://adguard.com/filter.txt"},
			},
		},
	})

	fs, err := client.GetFiltering()
	require.NoError(t, err)
	assert.True(t, fs.Enabled)
	assert.Len(t, fs.Filters, 1)
	assert.Equal(t, 55000, fs.Filters[0].RulesCount)
}

func TestAddFilter(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/filtering/add_url": nil,
	})

	err := client.AddFilter("OISD", "https://oisd.nl", true)
	assert.NoError(t, err)
}

func TestRemoveFilter(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/filtering/remove_url": nil,
	})

	err := client.RemoveFilter("https://oisd.nl")
	assert.NoError(t, err)
}

func TestRefreshFilters(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/filtering/refresh": nil,
	})

	err := client.RefreshFilters()
	assert.NoError(t, err)
}

func TestCheckHost(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/filtering/check_host": map[string]any{
			"reason":    "FilteredBlockedService",
			"filter_id": -2,
		},
	})

	result, err := client.CheckHost("youtube.com")
	require.NoError(t, err)
	assert.Equal(t, "FilteredBlockedService", result["reason"])
}

func TestGetDHCPStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/dhcp/status": DHCPStatus{
			Enabled:       false,
			InterfaceName: "eth0",
			V4:            DHCPConfigV4{GatewayIP: "192.168.0.1", SubnetMask: "255.255.255.0"},
			Leases:        []DHCPStaticLease{},
			StaticLeases:  []DHCPStaticLease{{MAC: "AA:BB:CC:DD:EE:FF", IP: "192.168.0.50", Hostname: "server"}},
		},
	})

	s, err := client.GetDHCPStatus()
	require.NoError(t, err)
	assert.False(t, s.Enabled)
	assert.Equal(t, "eth0", s.InterfaceName)
	assert.Len(t, s.StaticLeases, 1)
	assert.Equal(t, "server", s.StaticLeases[0].Hostname)
}

func TestAddStaticLease(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dhcp/add_static_lease": nil,
	})

	err := client.AddStaticLease(DHCPStaticLease{MAC: "AA:BB:CC:DD:EE:FF", IP: "192.168.0.50", Hostname: "test"})
	assert.NoError(t, err)
}

func TestRemoveStaticLease(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dhcp/remove_static_lease": nil,
	})

	err := client.RemoveStaticLease(DHCPStaticLease{MAC: "AA:BB:CC:DD:EE:FF", IP: "192.168.0.50", Hostname: "test"})
	assert.NoError(t, err)
}

func TestGetTLSStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/tls/status": TLSStatus{
			Enabled:    false,
			PortHTTPS:  443,
			ValidCert:  false,
			ServerName: "dns.example.com",
		},
	})

	s, err := client.GetTLSStatus()
	require.NoError(t, err)
	assert.False(t, s.Enabled)
	assert.Equal(t, 443, s.PortHTTPS)
	assert.Equal(t, "dns.example.com", s.ServerName)
}

func TestGetSafeBrowsingStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/safebrowsing/status": ToggleStatus{Enabled: true},
	})

	s, err := client.GetSafeBrowsingStatus()
	require.NoError(t, err)
	assert.True(t, s.Enabled)
}

func TestSetSafeBrowsing_Enable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/safebrowsing/enable": nil,
	})

	err := client.SetSafeBrowsing(true)
	assert.NoError(t, err)
}

func TestSetSafeBrowsing_Disable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/safebrowsing/disable": nil,
	})

	err := client.SetSafeBrowsing(false)
	assert.NoError(t, err)
}

func TestGetParentalStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/parental/status": ToggleStatus{Enabled: false},
	})

	s, err := client.GetParentalStatus()
	require.NoError(t, err)
	assert.False(t, s.Enabled)
}

func TestGetSafeSearchStatus(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/safesearch/status": SafeSearchConfig{
			Enabled: true,
			Google:  true,
			YouTube: true,
			Bing:    false,
		},
	})

	s, err := client.GetSafeSearchStatus()
	require.NoError(t, err)
	assert.True(t, s.Enabled)
	assert.True(t, s.Google)
	assert.True(t, s.YouTube)
	assert.False(t, s.Bing)
}

func TestGetAccessList(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/access/list": AccessList{
			AllowedClients:    []string{"0.0.0.0/0"},
			DisallowedClients: []string{},
			BlockedHosts:      []string{"version.bind", "id.server"},
		},
	})

	a, err := client.GetAccessList()
	require.NoError(t, err)
	assert.Equal(t, []string{"0.0.0.0/0"}, a.AllowedClients)
	assert.Empty(t, a.DisallowedClients)
	assert.Len(t, a.BlockedHosts, 2)
}

func TestGetDNSConfig(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/dns_info": DNSConfig{
			UpstreamDNS:  []string{"94.140.14.140"},
			BootstrapDNS: []string{"9.9.9.10"},
			CacheSize:    4194304,
			RateLimit:    20,
			BlockingMode: "default",
		},
	})

	d, err := client.GetDNSConfig()
	require.NoError(t, err)
	assert.Equal(t, []string{"94.140.14.140"}, d.UpstreamDNS)
	assert.Equal(t, 4194304, d.CacheSize)
	assert.Equal(t, 20, d.RateLimit)
}

func TestClearCache(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/cache_clear": nil,
	})

	err := client.ClearCache()
	assert.NoError(t, err)
}

func TestClearQueryLog(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/querylog_clear": nil,
	})

	err := client.ClearQueryLog()
	assert.NoError(t, err)
}

func TestSetProtection_Enable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dns_config": nil,
	})

	err := client.SetProtection(true)
	assert.NoError(t, err)
}

func TestSetProtection_Disable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dns_config": nil,
	})

	err := client.SetProtection(false)
	assert.NoError(t, err)
}

func TestAPIError_4xx(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/control/clients", func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != "admin" || pass != "secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"message":"bad request"}`, http.StatusBadRequest)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient(server.URL, "admin", "secret")
	_, err := client.GetClients()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error (HTTP 400)")
}

func TestConnectionError(t *testing.T) {
	client := NewClient("http://127.0.0.1:1", "admin", "secret")
	err := client.Ping()
	assert.Error(t, err)
}

func TestSetDHCPConfig(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dhcp/set_config": nil,
	})
	err := client.SetDHCPConfig(map[string]any{"enabled": true})
	assert.NoError(t, err)
}

func TestResetDHCP(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dhcp/reset": nil,
	})
	err := client.ResetDHCP()
	assert.NoError(t, err)
}

func TestResetDHCPLeases(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dhcp/reset_leases": nil,
	})
	err := client.ResetDHCPLeases()
	assert.NoError(t, err)
}

func TestGetDHCPInterfaces(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/dhcp/interfaces": map[string]any{
			"eth0": map[string]any{"name": "eth0", "hardware_address": "AA:BB:CC:DD:EE:FF"},
		},
	})
	ifaces, err := client.GetDHCPInterfaces()
	require.NoError(t, err)
	assert.Contains(t, ifaces, "eth0")
}

func TestSetTLSConfig(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/tls/configure": nil,
	})
	err := client.SetTLSConfig(map[string]any{"enabled": false})
	assert.NoError(t, err)
}

func TestSetParental_Enable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/parental/enable": nil,
	})
	err := client.SetParental(true)
	assert.NoError(t, err)
}

func TestSetParental_Disable(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/parental/disable": nil,
	})
	err := client.SetParental(false)
	assert.NoError(t, err)
}

func TestSetSafeSearch(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/safesearch/settings": nil,
	})
	err := client.SetSafeSearch(SafeSearchConfig{Enabled: true, Google: true})
	assert.NoError(t, err)
}

func TestSetAccessList(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/access/set": nil,
	})
	err := client.SetAccessList(AccessList{
		AllowedClients: []string{"0.0.0.0/0"},
	})
	assert.NoError(t, err)
}

func TestSetDNSConfig(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/dns_config": nil,
	})
	err := client.SetDNSConfig(map[string]any{"cache_size": 8388608})
	assert.NoError(t, err)
}

func TestSetFilteringRules(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/filtering/set_rules": nil,
	})
	err := client.SetFilteringRules("||ads.example.com^")
	assert.NoError(t, err)
}

func TestFindClient(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"GET /control/clients/find": map[string]any{
			"192.168.1.50": map[string]any{"name": "Desktop"},
		},
	})
	result, err := client.FindClient("192.168.1.50")
	require.NoError(t, err)
	assert.Contains(t, result, "192.168.1.50")
}

func TestGetVersionInfo(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/version.json": VersionInfo{
			Version:       "v0.107.57",
			NewVersion:    "",
			CanAutoUpdate: true,
		},
	})
	v, err := client.GetVersionInfo()
	require.NoError(t, err)
	assert.Equal(t, "v0.107.57", v.Version)
	assert.True(t, v.CanAutoUpdate)
}

func TestUpdate(t *testing.T) {
	_, client := newTestServer(t, map[string]any{
		"POST /control/update": nil,
	})
	err := client.Update()
	assert.NoError(t, err)
}

func TestNewClient(t *testing.T) {
	client := NewClient("http://example.com/", "user", "pass")
	assert.Equal(t, "http://example.com", client.baseURL) // trailing slash trimmed
	assert.Equal(t, "user", client.username)
	assert.Equal(t, "pass", client.password)
}
