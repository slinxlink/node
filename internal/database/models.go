package database

import "time"

type Config struct {
	ID          uint `gorm:"primarykey"`
	SecretKey   string
	Username    string
	Password    string
	Port        int
	Path        string
	IPv4        string
	IPv6        string
	Domain      string
	LogEnable   bool
	LogLevel    string
	LogPath     string
	BoardEnable bool
	StartedAt   time.Time
}

func (Config) TableName() string { return "config" }

type Core struct {
	ID         uint   `gorm:"primarykey"`
	Name       string // sing-box / xray
	BinPath    string // bin/sing-box
	ConfigPath string // data/sing-box.json
	LogEnable  bool
	LogLevel   string // trace / debug / info / warn / error / fatal / panic
	LogPath    string // data/sing-box.log
	Repo       string // github.com/SagerNet/sing-box
	Version    string // 当前版本号
	StartedAt  time.Time
	UpdatedAt  time.Time
}

func (Core) TableName() string { return "core" }

type Stats struct {
	ID        uint `gorm:"primarykey"`
	InboundID int  `gorm:"default:0"`
	Upload    int64
	Download  int64
	Online    int
}

func (Stats) TableName() string { return "stats" }

type Cert struct {
	ID        uint `gorm:"primarykey"`
	Domain    string
	CertPath  string
	KeyPath   string
	Mode      string // dns / http / manual
	Acme      uint
	Dns       uint
	AutoRenew bool
	UpdatedAt time.Time
	ExpireAt  time.Time
}

func (Cert) TableName() string { return "cert" }

type Acme struct {
	ID         uint `gorm:"primarykey"`
	Email      string
	Provider   string // letsencrypt / zerossl
	PrivateKey string
	EabKid     string
	EabHmac    string
}

func (Acme) TableName() string { return "acme" }

type DnsAccount struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Provider string // aliyun / cloudflare
	Key      string // 阿里云 AccessKeyID / Cloudflare Email
	Secret   string // 阿里云 AccessKeySecret / Cloudflare API Token
}

func (DnsAccount) TableName() string { return "dns_account" }

type Inbound struct {
	ID        uint   `gorm:"primarykey"`
	Name      string // 备注名
	Protocol  string // vless / vmess / hysteria
	Port      int    // 监听端口
	Enable    bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 传输层
	Transport      string // RAW / WebSocket
	WsPath         string // WebSocket 路径，如 /ray
	WsHost         string // WebSocket Host header
	WsPingInterval int    // WebSocket 心跳周期，秒，0 表示禁用

	// Hysteria
	UDPTimeout        int    // UDP 空闲超时，秒，0 表示默认
	MasqueradeEnabled bool   // 是否开启伪装
	MasqueradeType    string // default / proxy / file / string
	MasqueradeURL     string // 伪装地址，配合 proxy 类型使用
	RewriteHost       bool   // 是否重写 Host
	IgnoreTLSVerify   bool   // 跳过 TLS 验证
	MasqueradePath    string // 伪装路径，配合 file 类型使用
	MasqueradeCode    int    // 状态码，配合 string 类型使用
	MasqueradeBody    string `gorm:"type:text"` // 伪装内容，配合 string 类型使用

	// 通用 TLS
	TLSType       string `gorm:"default:'none'"` // none / TLS / Reality
	ServerName    string // 域名
	CipherSuites  string // 加密套件，逗号分隔，如 TLS_AES_128_GCM_SHA256
	TLSMinVersion string // TLS 最小版本，如 1.0 / 1.1 / 1.2 / 1.3
	TLSMaxVersion string // TLS 最大版本，如 1.0 / 1.1 / 1.2 / 1.3
	UTLS          string // uTLS 指纹，如 chrome / firefox / safari / ios / android / edge / 360 / qq / random
	Insecure      bool   // 跳过 TLS 验证，下发给客户端订阅使用
	ALPN          string // ALPN，逗号分隔，h3, h2, http/1.1
	Certificates  string // JSON 数组，如 [1, 2, 3]，关联 Certificate 表 ID

	// Reality（VLESS 专用）
	RealityServerName  string // SNI，如 www.amd.com
	RealityServer      string // 伪装目标，如 www.amd.com
	RealityServerPort  int    // 伪装目标端口
	RealityMaxTimeDiff int    // 最大时间差，毫秒，0 表示不限制
	RealityShortIDs    string // JSON 数组，如 ["1b5a14ff", "0478"]
	RealityPrivateKey  string // 生成配置用
	RealityPublicKey   string // 下发订阅用
	Flow               string // xtls-rprx-vision，VLESS+Reality 专用，下发给客户端
}

func (Inbound) TableName() string { return "inbound" }

type User struct {
	ID        uint   `gorm:"primarykey"`
	Name      string // 备注名
	Token     string `gorm:"unique"` // 订阅 token
	Inbounds  string // JSON 数组 [1, 3, 5]，绑定的入站 ID
	UUID      string // VLESS / VMess 用
	Password  string // Hysteria 用
	Enable    bool   `gorm:"default:true"`
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string { return "user" }

type Board struct {
	ID           uint   `gorm:"primarykey"`
	Name         string // 备注名，如 "slinx.link"
	Host         string // 面板地址，如 "https://slinx.link"
	NodeID       int    // 节点ID
	Key          string // 通讯密钥
	Inbound      int    // 对接的入站ID
	Enable       bool   `gorm:"default:true"`
	Type         string `gorm:"default:'SLINX'"` // SSPanel / v2board ...
	SyncInterval int    `gorm:"default:60"`      // 同步间隔，秒

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Board) TableName() string { return "board" }

type BoardUser struct {
	ID        uint `gorm:"primarykey"`
	BoardID   uint // 关联 Panel 表
	UserID    int  // 面板那边的用户 ID
	UUID      string
	Passwd    string
	Upload    int64
	Download  int64
	AliveIP   int // 当前在线 IP 数
	UpdatedAt time.Time
}

func (BoardUser) TableName() string { return "board_user" }

type SystemLog struct {
	ID        uint `gorm:"primarykey"`
	CPU       float64
	RAM       float64
	Load      float64
	Upload    int64
	Download  int64
	CreatedAt time.Time
}

func (SystemLog) TableName() string { return "system_log" }

type IP struct {
	ID              uint      `gorm:"primarykey"`
	IP              string    // 用于缓存验证
	Source          string    // 数据来源平台
	IPVersion       string    // v4 / v6
	ASN             string    // ASN号
	ASNOrg          string    // ASN组织
	RegisterCountry string    // 注册地国家代码
	Country         string    // 使用地国家代码
	City            string    // 城市
	IPType          string    // IP类型 住宅/机房/企业等
	OrgType         string    // 组织类型 hosting/isp/business...
	UpdatedAt       time.Time // 最后更新时间
}

func (IP) TableName() string { return "ip" }

type Unlock struct {
	ID        uint   `gorm:"primarykey"`
	IP        string // 出口IP，用于缓存验证
	IPVersion string // v4 / v6
	Platform  string // 平台名 netflix / disney / youtube 等
	Status    string // true / false / reject
	Region    string // 国家代码 JP / US 等，没有就空
	UpdatedAt time.Time
}

func (Unlock) TableName() string { return "unlock" }

type Route struct {
	ID        uint   `gorm:"primarykey"`
	City      string // shanghai / beijing / guangzhou
	Telecom   string // 线路类型，如 "电信163" / "电信CN2GT" / "电信CN2GIA"
	Unicom    string // "联通4837" / "联通9929"
	Mobile    string // "移动普通" / "移动CMI" / "移动CMIN2"
	UpdatedAt time.Time
}

func (Route) TableName() string { return "route" }
