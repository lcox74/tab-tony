package zerotier

type NetworkAPI struct {
	Id                    string `json:"id"`
	Description           string `json:"description"`
	RuleSource            string `json:"rulesSource"`
	OwnerId               string `json:"ownerId"`
	OnlineMemberCount     int    `json:"onlineMemberCount"`
	AuthorisedMemberCount int    `json:"authorizedMemberCount"`
	TotalMemberCount      int    `json:"totalMemberCount"`
	Config                struct {
		NetworkId         string        `json:"id,omitempty"`
		CreationTime      int64         `json:"creationTime,omitempty"`
		DNS               []DnsConfig   `json:"dns"`
		EnableBroadcast   bool          `json:"enableBroadcast"`
		IpAssignmentPools []IpPoolRange `json:"ipAssignmentPools"`
		LastModified      int64         `json:"lastModified,omitempty"`
		MTU               int           `json:"mtu"`
		MulticasetLimit   int           `json:"multicastLimit"`
		Name              string        `json:"name"`
		Private           bool          `json:"private"`
		Routes            []NetRoute    `json:"routes"`
		Rules             []NetRule     `json:"rules"`
		Tags              []struct {
		} `json:"tags"`
		Ipv4AssignMode struct {
			ZeroTier bool `json:"zt"`
		} `json:"v4AssignMode"`
		Ipv6AssignMode struct {
			RFC4193  bool `json:"6plane"`
			ZT6Plane bool `json:"rfc4193"`
			ZeroTier bool `json:"zt"`
		} `json:"v6AssignMode"`
	} `json:"config"`
	CapabilitiesByName struct{} `json:"capabilitiesByName"`
	TagsByName         struct{} `json:"tagsByName"`
	Permissions        struct{} `json:"permissions"`
}

type DnsConfig struct {
	Domain  string   `json:"domain"`
	Servers []string `json:"servers"`
}

type IpPoolRange struct {
	StartRange string `json:"ipRangeStart"`
	EndRange   string `json:"ipRangeEnd"`
}

type NetRoute struct {
	Target string `json:"target"`
	Via    string `json:"via,omitempty"` // Nullable
}
type NetRule struct {
	EtherType int    `json:"etherType,omitempty"`
	Note      bool   `json:"not,omitempty"`
	Or        bool   `json:"or,omitempty"`
	Type      string `json:"type"`
}


type NetworkMemberAPI struct {
	Id              string   `json:"id"`
	Address         string   `json:"address"`
	Nwid            string   `json:"nwid"`
	Authorized      bool     `json:"authorized"`
	ActiveBridge    bool     `json:"activeBridge"`
	Identity        string   `json:"identity"`
	IpAssignments   []string `json:"ipAssignments"`
	Revision        int      `json:"revision"`
	VMajor          int      `json:"vMajor"`
	VMinor          int      `json:"vMinor"`
	VRev            int      `json:"vRev"`
	VProto          int      `json:"vProto"`
}