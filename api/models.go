package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timestamp struct {
	time.Time
}

type Date struct {
	Timestamp
}

const (
	timestampLayout = "2006-01-02 15:04:05"
	dateLayout      = "2006-01-02"
)

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	switch s {
	case "", "null":
		t.Time = time.Time{}
		return nil
	}

	if parsed, err := time.Parse(timestampLayout, s); err == nil {
		t.Time = parsed
		return nil
	}

	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		t.Time = parsed
		return nil
	}

	return fmt.Errorf("cannot parse Timestamp %q: must be %q or RFC3339", s, timestampLayout)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(t.Format(timestampLayout))
}

func (t Timestamp) String() string {
	return t.Format(timestampLayout)
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	switch s {
	case "", "null":
		d.Time = time.Time{}
		return nil
	}

	// 1) Try date-only
	if parsed, err := time.Parse(dateLayout, s); err == nil {
		d.Time = parsed
		return nil
	}

	// 2) Try the timestamp layout
	if parsed, err := time.Parse(timestampLayout, s); err == nil {
		d.Time = parsed
		return nil
	}

	// 3) Try full RFC3339
	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		d.Time = parsed
		return nil
	}

	return fmt.Errorf(
		"cannot parse Date %q: must be %q, %q, or RFC3339",
		s, dateLayout, timestampLayout,
	)
}

// MarshalJSON always emits your date layout
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(d.Format(dateLayout))
}

// String prints using your date layout
func (d Date) String() string {
	return d.Format(dateLayout)
}

type CloudServerFlavour struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Cores       int    `json:"cores"`
	RamSize     int    `json:"ramSize"`
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

type Customer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type CloudServerImage struct {
	Id          int    `json:"id"`
	OSDistro    string `json:"osDistro"`
	OSVersion   string `json:"osVersion"`
	ReleaseDate Date   `json:"releaseDate"`
	Name        string `json:"name"`
	Size        Size64 `json:"size"`
	VirtualSize Size64 `json:"virtualSize"`
	MinRAM      int    `json:"minRAM"`
	MinDisk     int    `json:"minDisk"`
}

type CloudServer struct {
	Id               int              `json:"id"`
	ContractId       int              `json:"contractId"`
	Name             string           `json:"name"`
	IPv4             string           `json:"ipv4"`
	IPv6             string           `json:"ipv6"`
	Status           string           `json:"status"`
	AvailabilityZone string           `json:"availabilityZone"`
	Image            CloudServerImage `json:"image"`
}

type CloudServerActionResponse struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ChangeFlavourRequest struct {
	FlavourId int `json:"flavourId"`
}

type ChangeFlavourResponse struct {
	ServerId  int    `json:"serverId"`
	FlavourId int    `json:"flavourId"`
	Message   string `json:"message"`
}

type VMNetwork struct {
	Network   string `json:"network"`
	FixedIPv4 string `json:"fixed_ipv4"`
	FixedIPv6 string `json:"fixed_ipv6"`
}

type ProductPlan struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p Product) PlansString() string {
	b, err := json.Marshal(p.Plans)
	if err != nil {
		return fmt.Sprintf("Unable to marshal plans for product %d: %v", p.Id, err)
	}
	return string(b)
}

type Product struct {
	Id    int           `json:"id"`
	Name  string        `json:"name"`
	Plans []ProductPlan `json:"plans"`
}

type CloudServerOrder struct {
	ProductId        int         `json:"productId"`
	ProductPlanId    int         `json:"productPlanId"`
	ImageId          int         `json:"imageId"`
	Password         string      `json:"password"`
	AvailabilityZone string      `json:"availabilityZone"`
	Name             string      `json:"name"`
	SshKey           string      `json:"sshKey"`
	StorageSize      string      `json:"storageSize"`
	Networks         []VMNetwork `json:"networks"`
}

type CloudServerOrderResponse struct {
	Id         int    `json:"id"`
	ContractId int    `json:"contractId"`
	OrderId    string `json:"orderId"`
}

type AvailabilityZone struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

type VirtualServerActionResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type VirtualServerAction string

const (
	VirtualServerSoftReboot VirtualServerAction = "soft-reboot"
	VirtualServerHardReboot VirtualServerAction = "hard-reboot"
	VirtualServerPowerOff   VirtualServerAction = "power-off"
	VirtualServerPowerOn    VirtualServerAction = "power-on"
	VirtualServerReset      VirtualServerAction = "reset"
)

func (a VirtualServerAction) String() string {
	return string(a)
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Subnet struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	IpVersion       int              `json:"ipVersion"`
	Cidr            string           `json:"cidr"`
	AllocationPools []AllocationPool `json:"allocationPools"`
}

func (v VirtualNetwork) SubnetsString() string {
	b, err := json.Marshal(v.Subnets)
	if err != nil {
		return fmt.Sprintf("Unable to marshal subnets for Virtual Network %s: %v", v.Id, err)
	}
	return string(b)
}

type VirtualNetwork struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Subnets []Subnet `json:"subnets"`
}

type AttachedNetwork struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

type DetachVirtualNetworkRequest struct {
	NetworkId string `json:"networkId"`
}

type AttachVirtualNetworkRequest struct {
	NetworkId string `json:"networkId"`
	IPv4      string `json:"ipv4"`
	IPv6      string `json:"ipv6"`
}

type DetachVirtualNetworkResponse struct {
	ServerId int    `json:"serverId"`
	Message  string `json:"message"`
}

type AttachVirtualNetworkResponse struct {
	ServerId  int    `json:"serverId"`
	NetworkId string `json:"networkId"`
	Message   string `json:"message"`
}

type ResetCloudServerRequest struct {
	ImageId  int    `json:"imageId"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Backwards compatible workaround for API returning strings instead of numbers in some areas
// TODO: Remove once the API publishes numbers
type Size64 uint64

func (s *Size64) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("parsing Size64 from string %q: %w", str, err)
		}
		*s = Size64(v)
		return nil
	}
	var v uint64
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("parsing Size64 from number %s: %w", string(b), err)
	}
	*s = Size64(v)
	return nil
}

func (s Size64) MarshalJSON() ([]byte, error) {
	// always emit as a number
	return []byte(strconv.FormatUint(uint64(s), 10)), nil
}
