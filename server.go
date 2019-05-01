package govultr

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
)

// ServerService is the interface to interact with the server endpoints on the Vultr API
// Link: https://www.vultr.com/api/#server
type ServerService interface {
	ChangeApp(ctx context.Context, vpsID, appID string) error
	ListApps(ctx context.Context, vpsID string) ([]Application, error)
	AppInfo(ctx context.Context, vpsID string) (*ServerAppInfo, error)
	EnableBackup(ctx context.Context, vpsID string) error
	DisableBackup(ctx context.Context, vpsID string) error
	GetBackupSchedule(ctx context.Context, vpsID string) (*BackupSchedule, error)
	SetBackupSchedule(ctx context.Context, vpsID string, backup *BackupSchedule) error
	RestoreBackup(ctx context.Context, vpsID, backupID string) error
	RestoreSnapshot(ctx context.Context, vpsID, snapshotID string) error
	SetLabel(ctx context.Context, vpsID, label string) error
	SetTag(ctx context.Context, vpsID, tag string) error
	Neighbors(ctx context.Context, vpsID string) ([]int, error)
	EnablePrivateNetwork(ctx context.Context, vpsID, networkID string) error
	DisablePrivateNetwork(ctx context.Context, vpsID, networkID string) error
	ListPrivateNetworks(ctx context.Context, vpsID string) ([]PrivateNetwork, error)
	ListUpgradePlan(ctx context.Context, vpsID string) ([]int, error)
	UpgradePlan(ctx context.Context, vpsID, vpsPlanID string) error
	ListOS(ctx context.Context, vpsID string) ([]OS, error)
	ChangeOS(ctx context.Context, vpsID, osID string) error
	IsoAttach(ctx context.Context, vpsID, isoID string) error
	IsoDetach(ctx context.Context, vpsID string) error
	IsoStatus(ctx context.Context, vpsID string) (*ServerIso, error)
	SetFirewallGroup(ctx context.Context, vpsID, firewallGroupID string) error
	GetUserData(ctx context.Context, vpsID string) (*UserData, error)
	SetUserData(ctx context.Context, vpsID, userData string) error
	IPV4Info(ctx context.Context, vpsID string, public bool) ([]IPV4, error)
	IPV6Info(ctx context.Context, vpsID string) ([]IPV6, error)
	AddIPV4(ctx context.Context, vpsID string) error
	DestroyIPV4(ctx context.Context, vpsID, ip string) error
	EnableIPV6(ctx context.Context, vpsID string) error
	Bandwidth(ctx context.Context, vpsID string) ([]map[string]string, error)
}

// ServerServiceHandler handles interaction with the server methods for the Vultr API
type ServerServiceHandler struct {
	client *Client
}

// ServerAppInfo represents information about the application on your VPS
type ServerAppInfo struct {
	AppInfo string `json:"app_info"`
}

// BackupSchedule represents a schedule of a backup that runs on a VPS
type BackupSchedule struct {
	Enabled  bool   `json:"enabled"`
	CronType string `json:"cron_type"`
	NextRun  string `json:"next_scheduled_time_utc"`
	Hour     int    `json:"hour"`
	Dow      int    `json:"dow"`
	Dom      int    `json:"dom"`
}

// PrivateNetwork represents a private network attached to a VPS
type PrivateNetwork struct {
	NetworkID  string `json:"NETWORKID"`
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
}

// ServerIso represents a iso attached to a VPS
type ServerIso struct {
	State string `json:"state"`
	IsoID string `json:"ISOID"`
}

// UserData represents the user data you can give a VPS
type UserData struct {
	UserData string `json:"userdata"`
}

// IPV4 represents IPV4 information for a VPS
type IPV4 struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Type    string `json:"type"`
	Reverse string `json:"reverse"`
}

// IPV6 represents IPV6 information for a VPS
type IPV6 struct {
	IP          string `json:"ip"`
	Network     string `json:"network"`
	NetworkSize string `json:"network_size"`
	Type        string `json:"type"`
}

// ChangeApp changes the VPS to a different application.
func (s *ServerServiceHandler) ChangeApp(ctx context.Context, vpsID, appID string) error {

	uri := "/v1/server/app_change"

	values := url.Values{
		"SUBID": {vpsID},
		"APPID": {appID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// ListApps retrieves a list of applications to which a virtual machine can be changed.
func (s *ServerServiceHandler) ListApps(ctx context.Context, vpsID string) ([]Application, error) {

	uri := "/v1/server/app_change_list"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var appMap map[string]Application
	err = s.client.DoWithContext(ctx, req, &appMap)

	if err != nil {
		return nil, err
	}

	var appList []Application
	for _, a := range appMap {
		appList = append(appList, a)
	}

	return appList, nil
}

// AppInfo retrieves the application information for a given VPS ID
func (s *ServerServiceHandler) AppInfo(ctx context.Context, vpsID string) (*ServerAppInfo, error) {

	uri := "/v1/server/get_app_info"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	appInfo := new(ServerAppInfo)

	err = s.client.DoWithContext(ctx, req, appInfo)

	if err != nil {
		return nil, err
	}

	return appInfo, nil
}

// EnableBackup enables automatic backups on a given VPS
func (s *ServerServiceHandler) EnableBackup(ctx context.Context, vpsID string) error {

	uri := "/v1/server/backup_enable"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// DisableBackup disable automatic backups on a given VPS
func (s *ServerServiceHandler) DisableBackup(ctx context.Context, vpsID string) error {

	uri := "/v1/server/backup_disable"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// GetBackupSchedule retrieves the backup schedule for a given vps - all time values are in UTC
func (s *ServerServiceHandler) GetBackupSchedule(ctx context.Context, vpsID string) (*BackupSchedule, error) {

	uri := "/v1/server/backup_get_schedule"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return nil, err
	}

	backup := new(BackupSchedule)
	err = s.client.DoWithContext(ctx, req, backup)

	if err != nil {
		return nil, err
	}

	return backup, nil
}

// SetBackupSchedule sets the backup schedule for a given vps - all time values are in UTC
func (s *ServerServiceHandler) SetBackupSchedule(ctx context.Context, vpsID string, backup *BackupSchedule) error {

	uri := "/v1/server/backup_set_schedule"

	values := url.Values{
		"SUBID":     {vpsID},
		"cron_type": {backup.CronType},
		"hour":      {strconv.Itoa(backup.Hour)},
		"dow":       {strconv.Itoa(backup.Dow)},
		"dom":       {strconv.Itoa(backup.Dom)},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// RestoreBackup will restore the specified backup to the given VPS
func (s *ServerServiceHandler) RestoreBackup(ctx context.Context, vpsID, backupID string) error {

	uri := "/v1/server/restore_backup"

	values := url.Values{
		"SUBID":    {vpsID},
		"BACKUPID": {backupID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// RestoreSnapshot will restore the specified snapshot to the given VPS
func (s *ServerServiceHandler) RestoreSnapshot(ctx context.Context, vpsID, snapshotID string) error {

	uri := "/v1/server/restore_snapshot"

	values := url.Values{
		"SUBID":      {vpsID},
		"SNAPSHOTID": {snapshotID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// SetLabel will set a label for a given VPS
func (s *ServerServiceHandler) SetLabel(ctx context.Context, vpsID, label string) error {

	uri := "/v1/server/label_set"

	values := url.Values{
		"SUBID": {vpsID},
		"label": {label},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// SetTag will set a tag for a given VPS
func (s *ServerServiceHandler) SetTag(ctx context.Context, vpsID, tag string) error {

	uri := "/v1/server/tag_set"

	values := url.Values{
		"SUBID": {vpsID},
		"tag":   {tag},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// Neighbors will determine what other vps are hosted on the same physical host as a given vps.
func (s *ServerServiceHandler) Neighbors(ctx context.Context, vpsID string) ([]int, error) {

	uri := "/v1/server/neighbors"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var neighbors []int
	err = s.client.DoWithContext(ctx, req, &neighbors)

	if err != nil {
		return nil, err
	}

	return neighbors, nil
}

// EnablePrivateNetwork enables private networking on a server.
// The server will be automatically rebooted to complete the request.
// No action occurs if private networking was already enabled
func (s *ServerServiceHandler) EnablePrivateNetwork(ctx context.Context, vpsID, networkID string) error {

	uri := "/v1/server/private_network_enable"

	values := url.Values{
		"SUBID":     {vpsID},
		"NETWORKID": {networkID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// DisablePrivateNetwork removes a private network from a server.
// The server will be automatically rebooted to complete the request.
func (s *ServerServiceHandler) DisablePrivateNetwork(ctx context.Context, vpsID, networkID string) error {

	uri := "/v1/server/private_network_disable"

	values := url.Values{
		"SUBID":     {vpsID},
		"NETWORKID": {networkID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// ListPrivateNetworks will list private networks attached to a vps
func (s *ServerServiceHandler) ListPrivateNetworks(ctx context.Context, vpsID string) ([]PrivateNetwork, error) {

	uri := "/v1/server/private_networks"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var networkMap map[string]PrivateNetwork
	err = s.client.DoWithContext(ctx, req, &networkMap)

	if err != nil {
		return nil, err
	}

	var privateNetworks []PrivateNetwork
	for _, p := range networkMap {
		privateNetworks = append(privateNetworks, p)
	}

	return privateNetworks, nil
}

// ListUpgradePlan Retrieve a list of the planIDs for which the vps can be upgraded.
// An empty response array means that there are currently no upgrades available
func (s *ServerServiceHandler) ListUpgradePlan(ctx context.Context, vpsID string) ([]int, error) {

	uri := "/v1/server/upgrade_plan_list"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var plans []int
	err = s.client.DoWithContext(ctx, req, &plans)

	if err != nil {
		return nil, err
	}

	return plans, nil
}

// UpgradePlan will upgrade the plan of a virtual machine.
// The vps will be rebooted upon a successful upgrade.
func (s *ServerServiceHandler) UpgradePlan(ctx context.Context, vpsID, vpsPlanID string) error {

	uri := "/v1/server/upgrade_plan"

	values := url.Values{
		"SUBID":     {vpsID},
		"VPSPLANID": {vpsPlanID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// ListOS retrieves a list of operating systems to which the VPS can be changed to.
func (s *ServerServiceHandler) ListOS(ctx context.Context, vpsID string) ([]OS, error) {

	uri := "/v1/server/os_change_list"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var osMap map[string]OS
	err = s.client.DoWithContext(ctx, req, &osMap)

	if err != nil {
		return nil, err
	}

	var os []OS
	for _, o := range osMap {
		os = append(os, o)
	}

	return os, nil
}

// ChangeOS changes the VPS to a different operating system.
// All data will be permanently lost.
func (s *ServerServiceHandler) ChangeOS(ctx context.Context, vpsID, osID string) error {

	uri := "/v1/server/os_change"

	values := url.Values{
		"SUBID": {vpsID},
		"OSID":  {osID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// IsoAttach will attach an ISO to the given VPS and reboot it
func (s *ServerServiceHandler) IsoAttach(ctx context.Context, vpsID, isoID string) error {

	uri := "/v1/server/iso_attach"

	values := url.Values{
		"SUBID": {vpsID},
		"ISOID": {isoID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// IsoDetach will detach the currently mounted ISO and reboot the server.
func (s *ServerServiceHandler) IsoDetach(ctx context.Context, vpsID string) error {

	uri := "/v1/server/iso_detach"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// IsoStatus retrieves the current ISO state for a given VPS.
// The returned state may be one of: ready | isomounting | isomounted.
func (s *ServerServiceHandler) IsoStatus(ctx context.Context, vpsID string) (*ServerIso, error) {

	uri := "/v1/server/iso_status"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	serverIso := new(ServerIso)
	err = s.client.DoWithContext(ctx, req, serverIso)

	if err != nil {
		return nil, err
	}

	return serverIso, nil
}

// SetFirewallGroup will set, change, or remove the firewall group currently applied to a vps.
//  A value of "0" means "no firewall group"
func (s *ServerServiceHandler) SetFirewallGroup(ctx context.Context, vpsID, firewallGroupID string) error {

	uri := "/v1/server/firewall_group_set"

	values := url.Values{
		"SUBID":           {vpsID},
		"FIREWALLGROUPID": {firewallGroupID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// SetUserData sets the user-data for this subscription.
// User-data is a generic data store, which some provisioning tools and cloud operating systems use as a configuration file.
// It is generally consumed only once after an instance has been launched, but individual needs may vary.
func (s *ServerServiceHandler) SetUserData(ctx context.Context, vpsID, userData string) error {

	uri := "/v1/server/set_user_data"

	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))

	values := url.Values{
		"SUBID":    {vpsID},
		"userdata": {encodedUserData},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// GetUserData retrieves the (base64 encoded) user-data for this VPS
func (s *ServerServiceHandler) GetUserData(ctx context.Context, vpsID string) (*UserData, error) {

	uri := "/v1/server/get_user_data"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	userData := new(UserData)
	err = s.client.DoWithContext(ctx, req, userData)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

// IPV4Info will list the IPv4 information of a virtual machine.
// Public if set to 'true', includes information about the public network adapter (such as MAC address) with the "main_ip" entry.
func (s *ServerServiceHandler) IPV4Info(ctx context.Context, vpsID string, public bool) ([]IPV4, error) {

	uri := "/v1/server/list_ipv4"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)

	if public == true {
		q.Add("public_network", vpsID)
	}

	req.URL.RawQuery = q.Encode()

	var ipMap map[string][]IPV4
	err = s.client.DoWithContext(ctx, req, &ipMap)

	if err != nil {
		return nil, err
	}

	var ipv4 []IPV4
	for _, i := range ipMap {
		ipv4 = i
	}

	return ipv4, nil
}

// IPV6Info will list the IPv6 information of a virtual machine.
// If the virtual machine does not have IPv6 enabled, then an empty array is returned.
func (s *ServerServiceHandler) IPV6Info(ctx context.Context, vpsID string) ([]IPV6, error) {
	uri := "/v1/server/list_ipv6"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var ipMap map[string][]IPV6
	err = s.client.DoWithContext(ctx, req, &ipMap)

	if err != nil {
		return nil, err
	}

	var ipv6 []IPV6
	for _, i := range ipMap {
		ipv6 = i
	}

	return ipv6, nil
}

// AddIPV4 will add a new IPv4 address to a server.
func (s *ServerServiceHandler) AddIPV4(ctx context.Context, vpsID string) error {

	uri := "/v1/server/create_ipv4"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// DestroyIPV4 removes a secondary IPv4 address from a server.
// Your server will be hard-restarted. We suggest halting the machine gracefully before removing IPs.
func (s *ServerServiceHandler) DestroyIPV4(ctx context.Context, vpsID, ip string) error {

	uri := "/v1/server/destroy_ipv4"

	values := url.Values{
		"SUBID": {vpsID},
		"ip":    {ip},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// EnableIPV6 enables IPv6 networking on a server by assigning an IPv6 subnet to it.
func (s *ServerServiceHandler) EnableIPV6(ctx context.Context, vpsID string) error {

	uri := "/v1/server/ipv6_enable"

	values := url.Values{
		"SUBID": {vpsID},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, uri, values)

	if err != nil {
		return err
	}

	err = s.client.DoWithContext(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}

// Bandwidth will get the bandwidth used by a VPS
func (s *ServerServiceHandler) Bandwidth(ctx context.Context, vpsID string) ([]map[string]string, error) {

	uri := "/v1/server/bandwidth"

	req, err := s.client.NewRequest(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", vpsID)
	req.URL.RawQuery = q.Encode()

	var bandwidthMap map[string][][]string
	err = s.client.DoWithContext(ctx, req, &bandwidthMap)

	if err != nil {
		return nil, err
	}

	var bandwidth []map[string]string

	for _, b := range bandwidthMap["incoming_bytes"] {
		inMap := make(map[string]string)
		inMap["date"] = b[0]
		inMap["incoming"] = b[1]
		bandwidth = append(bandwidth, inMap)
	}

	for _, b := range bandwidthMap["outgoing_bytes"] {
		for i := range bandwidth {
			if bandwidth[i]["date"] == b[0] {
				bandwidth[i]["outgoing"] = b[1]
				break
			}
		}
	}

	return bandwidth, nil
}
