package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	legolog "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
	"github.com/seekky/slinx-node/internal/database"
	"github.com/seekky/slinx-node/internal/task"
	"github.com/seekky/slinx-node/internal/util"
)

// ── ACME 用户 ─────────────────────────────────────────────────────────────────

type acmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *acmeUser) GetEmail() string                        { return u.Email }
func (u *acmeUser) GetRegistration() *registration.Resource { return u.Registration }
func (u *acmeUser) GetPrivateKey() crypto.PrivateKey        { return u.key }

// ── 入口 ──────────────────────────────────────────────────────────────────────

func RegisterAcme(a *database.Acme, t *task.Task) {
	defer t.Done()
	var err error
	switch a.Provider {
	case "letsencrypt":
		err = registerLetsEncrypt(a, t)
	case "zerossl":
		err = registerZeroSSL(a, t)
	default:
		t.Log("ERROR", fmt.Sprintf("不支持的服务商: %s", a.Provider))
		return
	}
	if err != nil {
		t.Log("ERROR", err.Error())
	}
}

// ── Let's Encrypt ─────────────────────────────────────────────────────────────

func registerLetsEncrypt(a *database.Acme, t *task.Task) error {
	t.Log("INFO", "生成私钥...")
	privateKey, err := generateKey()
	if err != nil {
		return err
	}

	t.Log("INFO", "连接 Let's Encrypt...")
	user := &acmeUser{Email: a.Email, key: privateKey}
	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryProduction
	config.Certificate.KeyType = certcrypto.RSA2048

	legolog.Logger = task.NewLogger(t)

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	t.Log("INFO", "注册账号...")
	reg, err := client.Registration.Register(registration.RegisterOptions{
		TermsOfServiceAgreed: true,
	})
	if err != nil {
		return err
	}

	user.Registration = reg
	return savePrivateKey(a, privateKey, t)
}

// ── ZeroSSL ───────────────────────────────────────────────────────────────────

func registerZeroSSL(a *database.Acme, t *task.Task) error {
	t.Log("INFO", "生成私钥...")
	privateKey, err := generateKey()
	if err != nil {
		return err
	}

	t.Log("INFO", "连接 ZeroSSL...")
	user := &acmeUser{Email: a.Email, key: privateKey}
	config := lego.NewConfig(user)
	config.CADirURL = "https://acme.zerossl.com/v2/DV90"
	config.Certificate.KeyType = certcrypto.RSA2048

	legolog.Logger = task.NewLogger(t)

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	t.Log("INFO", "注册账号 (EAB)...")
	reg, err := client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
		TermsOfServiceAgreed: true,
		Kid:                  a.EabKid,
		HmacEncoded:          a.EabHmac,
	})
	if err != nil {
		return err
	}

	user.Registration = reg
	return savePrivateKey(a, privateKey, t)
}

// ── 公共工具 ──────────────────────────────────────────────────────────────────

func generateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func savePrivateKey(a *database.Acme, key *ecdsa.PrivateKey, t *task.Task) error {
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	keyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}))
	a.PrivateKey = keyPEM
	database.DB.Save(a)
	util.Info("[acme] 注册成功: %s", a.Email)
	t.Log("INFO", "注册成功")
	return nil
}

func caURL(provider string) (string, error) {
	switch provider {
	case "letsencrypt":
		return lego.LEDirectoryProduction, nil
	case "zerossl":
		return "https://acme.zerossl.com/v2/DV90", nil
	default:
		return "", fmt.Errorf("不支持的服务商: %s", provider)
	}
}
