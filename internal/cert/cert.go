package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	legolog "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/task"
	"github.com/seekky/slinx-node/internal/util"
)

var certDir string

func init() {
	wd, _ := os.Getwd()
	certDir = wd + "/cert"
}

// ── 申请证书 ──────────────────────────────────────────────────────────────────

func ApplyCert(cert *database.Cert, t *task.Task) {
	defer t.Done()

	certPath := certDir + "/" + cert.Domain + "/cert.pem"
	if data, err := os.ReadFile(certPath); err == nil {
		if block, _ := pem.Decode(data); block != nil {
			if x509cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				if x509cert.NotAfter.After(time.Now().Add(10 * 24 * time.Hour)) {
					t.Log("ERROR", fmt.Sprintf("证书 %s 有效期内，如需强制续签请等待到期前 10 天", cert.Domain))
					return
				}
			}
		}
	}

	switch cert.Mode {
	case "dns":
		applyDNS(cert, t)
	case "http":
		applyHTTP(cert, t)
	default:
		t.Log("ERROR", fmt.Sprintf("不支持的申请方式: %s", cert.Mode))
	}
}

func applyHTTP(cert *database.Cert, t *task.Task) {
	client, err := buildClient(cert, t)
	if err != nil {
		t.Log("ERROR", err.Error())
		return
	}

	t.Log("INFO", "启动 HTTP 验证服务器...")
	provider := http01.NewProviderServer("", "80")
	if err := client.Challenge.SetHTTP01Provider(provider); err != nil {
		t.Log("ERROR", "配置 HTTP 验证失败: "+err.Error())
		return
	}

	t.Log("INFO", "申请证书: "+cert.Domain)
	req := certificate.ObtainRequest{Domains: []string{cert.Domain}, Bundle: true}
	res, err := client.Certificate.Obtain(req)
	if err != nil {
		t.Log("ERROR", "申请失败: "+err.Error())
		return
	}

	saveCertFiles(cert, res, t)
}

func applyDNS(cert *database.Cert, t *task.Task) {
	client, err := buildClient(cert, t)
	if err != nil {
		t.Log("ERROR", err.Error())
		return
	}

	t.Log("INFO", "读取 DNS 账号...")
	var dns database.DnsAccount
	if err := database.DB.First(&dns, cert.Dns).Error; err != nil {
		t.Log("ERROR", "DNS 账号不存在")
		return
	}

	t.Log("INFO", "配置 DNS 验证...")
	if err := setDNSProvider(client, &dns); err != nil {
		t.Log("ERROR", "配置 DNS 失败: "+err.Error())
		return
	}

	t.Log("INFO", "申请证书: "+cert.Domain)
	req := certificate.ObtainRequest{Domains: []string{cert.Domain}, Bundle: true}
	res, err := client.Certificate.Obtain(req)
	if err != nil {
		t.Log("ERROR", "申请失败: "+err.Error())
		return
	}

	saveCertFiles(cert, res, t)
}

func buildClient(cert *database.Cert, t *task.Task) (*lego.Client, error) {
	t.Log("INFO", "读取 ACME 账号...")
	var acme database.Acme
	if err := database.DB.First(&acme, cert.Acme).Error; err != nil {
		return nil, fmt.Errorf("ACME 账号不存在")
	}
	if acme.PrivateKey == "" {
		return nil, fmt.Errorf("ACME 账号未注册，请先注册")
	}

	t.Log("INFO", "解码私钥...")
	privateKey, err := decodePrivateKey(acme.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("解码私钥失败: %w", err)
	}

	t.Log("INFO", "初始化 ACME 客户端...")
	user := &acmeUser{Email: acme.Email, key: privateKey}
	config := lego.NewConfig(user)
	caURL, err := caURL(acme.Provider)
	if err != nil {
		return nil, err
	}
	config.CADirURL = caURL
	config.Certificate.KeyType = certcrypto.RSA2048
	config.Certificate.Timeout = 10 * time.Minute

	legolog.Logger = task.NewLogger(t)

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %w", err)
	}

	t.Log("INFO", "验证账号注册信息...")
	reg, err := client.Registration.ResolveAccountByKey()
	if err != nil {
		return nil, fmt.Errorf("获取账号信息失败: %w", err)
	}
	user.Registration = reg

	return client, nil
}

func ApplyManual(cert *database.Cert, certContent, keyContent string) error {
	block, _ := pem.Decode([]byte(certContent))
	if block == nil {
		return fmt.Errorf("无法解析证书")
	}
	x509cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析证书失败: %w", err)
	}
	if len(x509cert.DNSNames) > 0 {
		cert.Domain = x509cert.DNSNames[0]
	} else {
		cert.Domain = x509cert.Subject.CommonName
	}

	if err := writeCertFiles(cert, []byte(certContent), []byte(keyContent)); err != nil {
		return err
	}
	util.Info("[cert] 证书上传成功: %s", cert.Domain)
	return nil
}

func setDNSProvider(client *lego.Client, dns *database.DnsAccount) error {
	opts := []dns01.ChallengeOption{
		dns01.AddRecursiveNameservers([]string{"1.1.1.1:53", "8.8.8.8:53"}),
		dns01.DisableAuthoritativeNssPropagationRequirement(),
	}

	switch dns.Provider {
	case "aliyun":
		cfg := alidns.NewDefaultConfig()
		cfg.APIKey = dns.Key
		cfg.SecretKey = dns.Secret
		provider, err := alidns.NewDNSProviderConfig(cfg)
		if err != nil {
			return err
		}
		return client.Challenge.SetDNS01Provider(provider, opts...)
	case "cloudflare":
		cfg := cloudflare.NewDefaultConfig()
		cfg.AuthEmail = dns.Key
		cfg.AuthToken = dns.Secret
		provider, err := cloudflare.NewDNSProviderConfig(cfg)
		if err != nil {
			return err
		}
		return client.Challenge.SetDNS01Provider(provider, opts...)
	default:
		return fmt.Errorf("不支持的 DNS 服务商: %s", dns.Provider)
	}
}

func saveCertFiles(cert *database.Cert, res *certificate.Resource, t *task.Task) {
	t.Log("INFO", "保存证书文件...")
	if err := writeCertFiles(cert, res.Certificate, res.PrivateKey); err != nil {
		t.Log("ERROR", err.Error())
		return
	}
	util.Info("[cert] 证书申请成功: %s", cert.Domain)
	t.Log("INFO", "证书申请成功，过期时间: "+cert.ExpireAt.Format("2006-01-02"))
}

func writeCertFiles(cert *database.Cert, certPEM, keyPEM []byte) error {
	dir := certDir + "/" + cert.Domain
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	certPath := dir + "/cert.pem"
	keyPath := dir + "/key.pem"

	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("保存证书失败: %w", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("保存私钥失败: %w", err)
	}

	expireAt := time.Now().AddDate(0, 3, 0)
	if t2, err := parseCertExpiry(certPEM); err == nil {
		expireAt = t2
	}

	cert.CertPath = certPath
	cert.KeyPath = keyPath
	cert.ExpireAt = expireAt
	database.DB.Save(cert)

	return nil
}

func decodePrivateKey(keyPEM string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, fmt.Errorf("无法解析私钥")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

func parseCertExpiry(certPEM []byte) (time.Time, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return time.Time{}, fmt.Errorf("无法解析证书")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, err
	}
	return cert.NotAfter, nil
}

func ReadCertContent(cert *database.Cert) (map[string]string, error) {
	certPEM, err := os.ReadFile(cert.CertPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败: %w", err)
	}
	keyPEM, err := os.ReadFile(cert.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥失败: %w", err)
	}
	return map[string]string{
		"cert": string(certPEM),
		"key":  string(keyPEM),
	}, nil
}

func RemoveCert(cert *database.Cert) error {
	dir := certDir + "/" + cert.Domain
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("删除证书文件失败: %w", err)
	}
	database.DB.Delete(cert)
	return nil
}
