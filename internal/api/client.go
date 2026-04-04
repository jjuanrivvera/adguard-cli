package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

func NewClient(baseURL, username, password string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) do(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + "/control" + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		resp.Body.Close()
		return nil, fmt.Errorf("authentication failed (HTTP %d)", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func (c *Client) get(endpoint string, target any) error {
	resp, err := c.do("GET", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) getRaw(endpoint string) ([]byte, error) {
	resp, err := c.do("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) post(endpoint string, body any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = strings.NewReader(string(data))
	}

	resp, err := c.do("POST", endpoint, reader)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) postAndDecode(endpoint string, body any, target any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = strings.NewReader(string(data))
	}

	resp, err := c.do("POST", endpoint, reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// Ping checks connectivity to the AdGuard Home instance.
func (c *Client) Ping() error {
	var status map[string]any
	return c.get("/status", &status)
}

// --- Status ---

type ServerStatus struct {
	DNSAddresses    []string `json:"dns_addresses"`
	DNSPort         int      `json:"dns_port"`
	HTTPPort        int      `json:"http_port"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	Running         bool     `json:"running"`
	Version         string   `json:"version"`
}

func (c *Client) GetStatus() (*ServerStatus, error) {
	var s ServerStatus
	return &s, c.get("/status", &s)
}

func (c *Client) SetProtection(enabled bool) error {
	body := map[string]any{
		"enabled":  enabled,
		"duration": 0,
	}
	return c.post("/dns_config", body)
}

// --- Stats ---

type Stats struct {
	NumDNSQueries           int      `json:"num_dns_queries"`
	NumBlockedFiltering     int      `json:"num_blocked_filtering"`
	NumReplacedSafebrowsing int      `json:"num_replaced_safebrowsing"`
	NumReplacedParental     int      `json:"num_replaced_parental"`
	AvgProcessingTime       float64  `json:"avg_processing_time"`
	TopQueriedDomains       []map[string]int `json:"top_queried_domains"`
	TopClients              []map[string]int `json:"top_clients"`
	TopBlockedDomains       []map[string]int `json:"top_blocked_domains"`
	TimeUnits               string   `json:"time_units"`
}

func (c *Client) GetStats() (*Stats, error) {
	var s Stats
	return &s, c.get("/stats", &s)
}

func (c *Client) ResetStats() error {
	return c.post("/stats_reset", nil)
}

// --- Clients ---

type ClientEntry struct {
	Name                    string   `json:"name"`
	IDs                     []string `json:"ids"`
	Tags                    []string `json:"tags"`
	UseGlobalSettings       bool     `json:"use_global_settings"`
	UseGlobalBlockedServices bool    `json:"use_global_blocked_services"`
	BlockedServices         []string `json:"blocked_services"`
	FilteringEnabled        bool     `json:"filtering_enabled"`
	ParentalEnabled         bool     `json:"parental_enabled"`
	SafebrowsingEnabled     bool     `json:"safebrowsing_enabled"`
	Upstreams               []string `json:"upstreams"`
}

type ClientsResponse struct {
	Clients     []ClientEntry `json:"clients"`
	AutoClients []struct {
		IP     string `json:"ip"`
		Name   string `json:"name"`
		Source string `json:"source"`
	} `json:"auto_clients"`
}

func (c *Client) GetClients() (*ClientsResponse, error) {
	var cr ClientsResponse
	return &cr, c.get("/clients", &cr)
}

func (c *Client) AddClient(client ClientEntry) error {
	return c.post("/clients/add", client)
}

func (c *Client) UpdateClient(name string, client ClientEntry) error {
	body := map[string]any{
		"name": name,
		"data": client,
	}
	return c.post("/clients/update", body)
}

func (c *Client) DeleteClient(name string) error {
	body := map[string]any{
		"name": name,
	}
	return c.post("/clients/delete", body)
}

func (c *Client) FindClient(ip string) (map[string]any, error) {
	var result map[string]any
	return result, c.get(fmt.Sprintf("/clients/find?ip0=%s", ip), &result)
}

// --- Blocked Services ---

type BlockedServicesResponse struct {
	Schedule map[string]any `json:"schedule"`
	IDs      []string       `json:"ids"`
}

func (c *Client) GetBlockedServices() (*BlockedServicesResponse, error) {
	var bs BlockedServicesResponse
	return &bs, c.get("/blocked_services/get", &bs)
}

type ServiceEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetAllServices() ([]ServiceEntry, error) {
	var services []ServiceEntry
	return services, c.get("/blocked_services/all", &services)
}

func (c *Client) SetBlockedServices(services BlockedServicesResponse) error {
	return c.post("/blocked_services/set", services)
}

// --- DNS Rewrites ---

type RewriteEntry struct {
	Domain string `json:"domain"`
	Answer string `json:"answer"`
}

func (c *Client) GetRewrites() ([]RewriteEntry, error) {
	var rewrites []RewriteEntry
	return rewrites, c.get("/rewrite/list", &rewrites)
}

func (c *Client) AddRewrite(entry RewriteEntry) error {
	return c.post("/rewrite/add", entry)
}

func (c *Client) DeleteRewrite(entry RewriteEntry) error {
	return c.post("/rewrite/delete", entry)
}

// --- Query Log ---

type QueryLogEntry struct {
	Answer     []map[string]any `json:"answer"`
	Client     string           `json:"client"`
	ClientInfo map[string]any   `json:"client_info"`
	Elapsed    string           `json:"elapsedMs"`
	Question   map[string]any   `json:"question"`
	Reason     string           `json:"reason"`
	Status     string           `json:"status"`
	Time       string           `json:"time"`
}

type QueryLogResponse struct {
	Data   []QueryLogEntry `json:"data"`
	Oldest string          `json:"oldest"`
}

func (c *Client) GetQueryLog(limit int) (*QueryLogResponse, error) {
	var ql QueryLogResponse
	return &ql, c.get(fmt.Sprintf("/querylog?limit=%d", limit), &ql)
}

// --- Filtering ---

type FilterStatus struct {
	Enabled  bool     `json:"enabled"`
	Interval int      `json:"interval"`
	Filters  []Filter `json:"filters"`
}

type Filter struct {
	ID         int    `json:"id"`
	Enabled    bool   `json:"enabled"`
	LastUpdated string `json:"last_updated"`
	Name       string `json:"name"`
	RulesCount int    `json:"rules_count"`
	URL        string `json:"url"`
}

func (c *Client) GetFiltering() (*FilterStatus, error) {
	var fs FilterStatus
	return &fs, c.get("/filtering/status", &fs)
}

func (c *Client) AddFilter(name, url string, enabled bool) error {
	body := map[string]any{
		"name":    name,
		"url":     url,
		"enabled": enabled,
	}
	return c.post("/filtering/add_url", body)
}

func (c *Client) RemoveFilter(url string) error {
	body := map[string]any{
		"url": url,
	}
	return c.post("/filtering/remove_url", body)
}

func (c *Client) RefreshFilters() error {
	body := map[string]any{
		"whitelist": false,
	}
	return c.post("/filtering/refresh", body)
}

func (c *Client) SetFilteringRules(rules string) error {
	return c.post("/filtering/set_rules", map[string]any{"rules": rules})
}

func (c *Client) CheckHost(host string) (map[string]any, error) {
	var result map[string]any
	return result, c.get(fmt.Sprintf("/filtering/check_host?name=%s", host), &result)
}

// --- DHCP ---

type DHCPConfigV4 struct {
	GatewayIP     string `json:"gateway_ip"`
	SubnetMask    string `json:"subnet_mask"`
	RangeStart    string `json:"range_start"`
	RangeEnd      string `json:"range_end"`
	LeaseDuration int    `json:"lease_duration"`
}

type DHCPConfigV6 struct {
	RangeStart    string `json:"range_start"`
	LeaseDuration int    `json:"lease_duration"`
}

type DHCPStaticLease struct {
	MAC      string `json:"mac"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

type DHCPStatus struct {
	InterfaceName string            `json:"interface_name"`
	V4            DHCPConfigV4      `json:"v4"`
	V6            DHCPConfigV6      `json:"v6"`
	Leases        []DHCPStaticLease `json:"leases"`
	StaticLeases  []DHCPStaticLease `json:"static_leases"`
	Enabled       bool              `json:"enabled"`
}

func (c *Client) GetDHCPStatus() (*DHCPStatus, error) {
	var s DHCPStatus
	return &s, c.get("/dhcp/status", &s)
}

func (c *Client) SetDHCPConfig(cfg map[string]any) error {
	return c.post("/dhcp/set_config", cfg)
}

func (c *Client) AddStaticLease(lease DHCPStaticLease) error {
	return c.post("/dhcp/add_static_lease", lease)
}

func (c *Client) RemoveStaticLease(lease DHCPStaticLease) error {
	return c.post("/dhcp/remove_static_lease", lease)
}

func (c *Client) ResetDHCP() error {
	return c.post("/dhcp/reset", nil)
}

func (c *Client) ResetDHCPLeases() error {
	return c.post("/dhcp/reset_leases", nil)
}

func (c *Client) GetDHCPInterfaces() (map[string]any, error) {
	var result map[string]any
	return result, c.get("/dhcp/interfaces", &result)
}

// --- TLS ---

type TLSStatus struct {
	Enabled         bool     `json:"enabled"`
	ForceHTTPS      bool     `json:"force_https"`
	PortHTTPS       int      `json:"port_https"`
	PortDNSOverTLS  int      `json:"port_dns_over_tls"`
	PortDNSOverQUIC int      `json:"port_dns_over_quic"`
	CertificatePath string   `json:"certificate_path"`
	PrivateKeyPath  string   `json:"private_key_path"`
	ValidCert       bool     `json:"valid_cert"`
	ValidKey        bool     `json:"valid_key"`
	ValidPair       bool     `json:"valid_pair"`
	ServerName      string   `json:"server_name"`
	DNSNames        []string `json:"dns_names"`
	NotBefore       string   `json:"not_before"`
	NotAfter        string   `json:"not_after"`
}

func (c *Client) GetTLSStatus() (*TLSStatus, error) {
	var s TLSStatus
	return &s, c.get("/tls/status", &s)
}

func (c *Client) SetTLSConfig(cfg map[string]any) error {
	return c.post("/tls/configure", cfg)
}

// --- SafeBrowsing ---

type ToggleStatus struct {
	Enabled bool `json:"enabled"`
}

func (c *Client) GetSafeBrowsingStatus() (*ToggleStatus, error) {
	var s ToggleStatus
	return &s, c.get("/safebrowsing/status", &s)
}

func (c *Client) SetSafeBrowsing(enabled bool) error {
	if enabled {
		return c.post("/safebrowsing/enable", nil)
	}
	return c.post("/safebrowsing/disable", nil)
}

// --- Parental ---

func (c *Client) GetParentalStatus() (*ToggleStatus, error) {
	var s ToggleStatus
	return &s, c.get("/parental/status", &s)
}

func (c *Client) SetParental(enabled bool) error {
	if enabled {
		return c.post("/parental/enable", nil)
	}
	return c.post("/parental/disable", nil)
}

// --- SafeSearch ---

type SafeSearchConfig struct {
	Enabled    bool `json:"enabled"`
	Bing       bool `json:"bing"`
	DuckDuckGo bool `json:"duckduckgo"`
	Ecosia     bool `json:"ecosia"`
	Google     bool `json:"google"`
	Pixabay    bool `json:"pixabay"`
	Yandex     bool `json:"yandex"`
	YouTube    bool `json:"youtube"`
}

func (c *Client) GetSafeSearchStatus() (*SafeSearchConfig, error) {
	var s SafeSearchConfig
	return &s, c.get("/safesearch/status", &s)
}

func (c *Client) SetSafeSearch(cfg SafeSearchConfig) error {
	return c.post("/safesearch/settings", cfg)
}

// --- Access Control ---

type AccessList struct {
	AllowedClients    []string `json:"allowed_clients"`
	DisallowedClients []string `json:"disallowed_clients"`
	BlockedHosts      []string `json:"blocked_hosts"`
}

func (c *Client) GetAccessList() (*AccessList, error) {
	var a AccessList
	return &a, c.get("/access/list", &a)
}

func (c *Client) SetAccessList(a AccessList) error {
	return c.post("/access/set", a)
}

// --- DNS Config ---

type DNSConfig struct {
	UpstreamDNS           []string `json:"upstream_dns"`
	BootstrapDNS          []string `json:"bootstrap_dns"`
	FallbackDNS           []string `json:"fallback_dns"`
	ProtectionEnabled     bool     `json:"protection_enabled"`
	RateLimit             int      `json:"ratelimit"`
	BlockingMode          string   `json:"blocking_mode"`
	EDNSCSEnabled         bool     `json:"edns_cs_enabled"`
	DNSSECEnabled         bool     `json:"dnssec_enabled"`
	DisableIPv6           bool     `json:"disable_ipv6"`
	UpstreamMode          string   `json:"upstream_mode"`
	CacheSize             int      `json:"cache_size"`
	CacheMinTTL           int      `json:"cache_ttl_min"`
	CacheMaxTTL           int      `json:"cache_ttl_max"`
	CacheOptimistic       bool     `json:"cache_optimistic"`
	UsePrivatePTRResolvers bool    `json:"use_private_ptr_resolvers"`
}

func (c *Client) GetDNSConfig() (*DNSConfig, error) {
	var d DNSConfig
	return &d, c.get("/dns_info", &d)
}

func (c *Client) SetDNSConfig(cfg map[string]any) error {
	return c.post("/dns_config", cfg)
}

func (c *Client) ClearCache() error {
	return c.post("/cache_clear", nil)
}

func (c *Client) ClearQueryLog() error {
	return c.post("/querylog_clear", nil)
}

// --- Version ---

type VersionInfo struct {
	Version         string `json:"version"`
	NewVersion      string `json:"new_version"`
	Announcement    string `json:"announcement"`
	AnnouncementURL string `json:"announcement_url"`
	CanAutoUpdate   bool   `json:"can_autoupdate"`
}

func (c *Client) GetVersionInfo() (*VersionInfo, error) {
	var v VersionInfo
	err := c.postAndDecode("/version.json", map[string]any{"recheck_now": false}, &v)
	return &v, err
}

func (c *Client) Update() error {
	return c.post("/update", nil)
}
