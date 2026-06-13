package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/util"
)

const uaBrowser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64"

var unlockClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

var followClient = &http.Client{
	Timeout: 15 * time.Second,
}

// Disney+ 需要的 Media_Cookie，启动时拉取
var mediaCookie string

func init() {
	go func() {
		resp, err := http.Get("https://raw.githubusercontent.com/1-stream/RegionRestrictionCheck/main/cookies")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		mediaCookie = string(body)
	}()
}

type UnlockResult struct {
	Platform string
	Status   string
	Region   string
}

var platforms = []string{
	"apple", "bing", "google_play", "steam", "reddit", "wikipedia",
	"netflix", "disney", "prime_video", "youtube_premium", "dazn", "tvbanywhere", "iqiyi",
	"claude", "chatgpt", "gemini",
}

func FetchUnlockInfo() ([]database.Unlock, error) {
	ipv4, _ := util.GetPublicIPs()

	results := make([]UnlockResult, len(platforms))
	var wg sync.WaitGroup

	for i, platform := range platforms {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			results[idx] = checkPlatform(p)
		}(i, platform)
	}

	wg.Wait()

	var unlocks []database.Unlock
	for _, r := range results {
		record := database.Unlock{
			IP:        ipv4,
			IPVersion: "v4",
			Platform:  r.Platform,
			Status:    r.Status,
			Region:    r.Region,
			UpdatedAt: time.Now(),
		}

		var existing database.Unlock
		database.DB.Where(database.Unlock{
			IP:        ipv4,
			IPVersion: "v4",
			Platform:  r.Platform,
		}).First(&existing)

		if existing.ID == 0 {
			database.DB.Create(&record)
		} else {
			record.ID = existing.ID
			database.DB.Save(&record)
		}

		unlocks = append(unlocks, record)
	}

	return unlocks, nil
}

func checkPlatform(platform string) UnlockResult {
	switch platform {
	case "apple":
		return checkApple()
	case "bing":
		return checkBing()
	case "google_play":
		return checkGooglePlay()
	case "steam":
		return checkSteam()
	case "reddit":
		return checkReddit()
	case "wikipedia":
		return checkWikipedia()
	case "netflix":
		return checkNetflix()
	case "disney":
		return checkDisney()
	case "prime_video":
		return checkPrimeVideo()
	case "youtube_premium":
		return checkYouTubePremium()
	case "dazn":
		return checkDazn()
	case "tvbanywhere":
		return checkTVBAnywhere()
	case "iqiyi":
		return checkIqiyi()
	case "claude":
		return checkClaude()
	case "chatgpt":
		return checkChatGPT()
	case "gemini":
		return checkGemini()
	default:
		return UnlockResult{Platform: platform, Status: "false"}
	}
}

// ────────────────────────────────────────────────
// 常用平台
// ────────────────────────────────────────────────

func checkApple() UnlockResult {
	resp, err := followClient.Get("https://gspe1-ssl.ls.apple.com/pep/gcc")
	if err != nil {
		return UnlockResult{Platform: "apple", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	region := strings.TrimSpace(string(body))
	if region == "" {
		return UnlockResult{Platform: "apple", Status: "false"}
	}
	return UnlockResult{Platform: "apple", Status: "true", Region: strings.ToUpper(region)}
}

func checkBing() UnlockResult {
	req, _ := http.NewRequest("GET", "https://www.bing.com/search?q=curl", nil)
	req.Header.Set("User-Agent", uaBrowser)
	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "bing", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	// 检查是否被限制（Bing 在中国显示特定内容）
	if strings.Contains(content, "cn.bing.com") {
		return UnlockResult{Platform: "bing", Status: "reject", Region: "CN"}
	}
	region := extractField(content, `Region:"`, `"`)
	if strings.Contains(content, `sj_cook.set("SRCHHPGUSR","HV"`) {
		return UnlockResult{Platform: "bing", Status: "reject", Region: strings.ToUpper(region)}
	}
	return UnlockResult{Platform: "bing", Status: "true", Region: strings.ToUpper(region)}
}

func checkGooglePlay() UnlockResult {
	req, _ := http.NewRequest("GET", "https://play.google.com/", nil)
	req.Header.Set("User-Agent", uaBrowser)
	req.Header.Set("Accept-Language", "en-US;q=0.9")

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "google_play", Status: "false"}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	region := extractField(content, `<div class="yVZQTb">`, `<`)
	if idx := strings.Index(region, " ("); idx != -1 {
		region = region[:idx]
	}
	if region == "" {
		return UnlockResult{Platform: "google_play", Status: "false"}
	}

	return UnlockResult{Platform: "google_play", Status: "true", Region: regionToCountry(region)}
}

func checkSteam() UnlockResult {
	resp, err := followClient.Get("https://store.steampowered.com/api/appdetails?appids=761830&filters=price_overview")
	if err != nil {
		return UnlockResult{Platform: "steam", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	currency := extractField(string(body), `"currency":"`, `"`)
	if currency == "" {
		return UnlockResult{Platform: "steam", Status: "false"}
	}
	region := currencyToCountry(currency)
	return UnlockResult{Platform: "steam", Status: "true", Region: region}
}

func checkReddit() UnlockResult {
	req, _ := http.NewRequest("GET", "https://www.reddit.com/", nil)
	req.Header.Set("User-Agent", uaBrowser)
	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "reddit", Status: "false"}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return UnlockResult{Platform: "reddit", Status: "true"}
	}
	return UnlockResult{Platform: "reddit", Status: "false"}
}

func checkWikipedia() UnlockResult {
	req, _ := http.NewRequest("GET", "https://zh.wikipedia.org/w/index.php?title=Wikipedia%3A%E6%B2%99%E7%9B%92&action=edit", nil)
	req.Header.Set("User-Agent", uaBrowser)
	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "wikipedia", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	if strings.Contains(strings.ToLower(content), "banned") {
		return UnlockResult{Platform: "wikipedia", Status: "false"}
	}
	return UnlockResult{Platform: "wikipedia", Status: "true"}
}

// ────────────────────────────────────────────────
// 流媒体
// ────────────────────────────────────────────────

func checkNetflix() UnlockResult {
	cookie := `flwssn=d2c72c47-49e9-48da-b7a2-2dc6d7ca9fcf; nfvdid=BQFmAAEBEMZa4XMYVzVGf9-kQ1HXumtAKsCyuBZU4QStC6CGEGIVznjNuuTerLAG8v2-9V_kYhg5uxTB5_yyrmqc02U5l1Ts74Qquezc9AE-LZKTo3kY3g%3D%3D; SecureNetflixId=v%3D3%26mac%3DAQEAEQABABSQHKcR1d0sLV0WTu0lL-BO63TKCCHAkeY.%26dt%3D1745376277212; NetflixId=v%3D3%26ct%3DBgjHlOvcAxLAAZuNS4_CJHy9NKJPzUV-9gElzTlTsmDS1B59TycR-fue7f6q7X9JQAOLttD7OnlldUtnYWXL7VUfu9q4pA0gruZKVIhScTYI1GKbyiEqKaULAXOt0PHQzgRLVTNVoXkxcbu7MYG4wm1870fZkd5qrDOEseZv2WIVk4xIeNL87EZh1vS3RZU3e-qWy2tSmfSNUC-FVDGwxbI6-hk3Zg2MbcWYd70-ghohcCSZp5WHAGXg_xWVC7FHM3aOUVTGwRCU1RgGIg4KDKGr_wsTRRw6HWKqeA..`

	makeReq := func(url string) (string, error) {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", uaBrowser)
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("cookie", cookie)
		resp, err := followClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return string(body), nil
	}

	content1, err1 := makeReq("https://www.netflix.com/title/81280792")
	content2, err2 := makeReq("https://www.netflix.com/title/70143836")

	if err1 != nil || err2 != nil {
		return UnlockResult{Platform: "netflix", Status: "false"}
	}

	result1 := strings.Contains(content1, "Oh no!")
	result2 := strings.Contains(content2, "Oh no!")

	if result1 && result2 {
		return UnlockResult{Platform: "netflix", Status: "reject"}
	}

	if !result1 || !result2 {
		region := extractField(content1, `"country":"`, `"`)
		return UnlockResult{Platform: "netflix", Status: "true", Region: strings.ToUpper(region)}
	}

	return UnlockResult{Platform: "netflix", Status: "false"}
}

func checkDisney() UnlockResult {
	// Step 1: 获取 assertion
	req1, _ := http.NewRequest("POST", "https://disney.api.edge.bamgrid.com/devices",
		strings.NewReader(`{"deviceFamily":"browser","applicationRuntime":"chrome","deviceProfile":"windows","attributes":{}}`))
	req1.Header.Set("authorization", "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84")
	req1.Header.Set("content-type", "application/json; charset=UTF-8")
	req1.Header.Set("User-Agent", uaBrowser)

	resp1, err := followClient.Do(req1)
	if err != nil {
		return UnlockResult{Platform: "disney", Status: "false"}
	}
	defer resp1.Body.Close()
	body1, _ := io.ReadAll(resp1.Body)
	content1 := string(body1)

	if strings.Contains(strings.ToLower(content1), "403 error") {
		return UnlockResult{Platform: "disney", Status: "false"}
	}

	assertion := extractField(content1, `"assertion":"`, `"`)
	if assertion == "" {
		return UnlockResult{Platform: "disney", Status: "false"}
	}

	// Step 2: 获取 token
	preDisneyCookie := getLine(mediaCookie, 1)
	disneyCookie := strings.ReplaceAll(preDisneyCookie, "DISNEYASSERTION", assertion)

	req2, _ := http.NewRequest("POST", "https://disney.api.edge.bamgrid.com/token",
		strings.NewReader(disneyCookie))
	req2.Header.Set("authorization", "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84")
	req2.Header.Set("User-Agent", uaBrowser)

	resp2, err := followClient.Do(req2)
	if err != nil {
		return UnlockResult{Platform: "disney", Status: "false"}
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	content2 := string(body2)

	if strings.Contains(content2, "forbidden-location") || strings.Contains(strings.ToLower(content2), "403 error") {
		return UnlockResult{Platform: "disney", Status: "false"}
	}

	// Step 3: GraphQL 获取地区
	refreshToken := extractField(content2, `"refresh_token":"`, `"`)
	fakeContent := getLine(mediaCookie, 8)
	disneyContent := strings.ReplaceAll(fakeContent, "ILOVEDISNEY", refreshToken)

	req3, _ := http.NewRequest("POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
		strings.NewReader(disneyContent))
	req3.Header.Set("authorization", "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84")
	req3.Header.Set("User-Agent", uaBrowser)

	resp3, err := followClient.Do(req3)
	if err != nil {
		return UnlockResult{Platform: "disney", Status: "false"}
	}
	defer resp3.Body.Close()
	body3, _ := io.ReadAll(resp3.Body)
	content3 := string(body3)

	// Step 4: 检查 preview/unavailable
	req4, _ := http.NewRequest("GET", "https://disneyplus.com", nil)
	req4.Header.Set("User-Agent", uaBrowser)
	resp4, err := followClient.Do(req4)
	isUnavailable := false
	if err == nil {
		defer resp4.Body.Close()
		finalURL := resp4.Request.URL.String()
		isUnavailable = strings.Contains(finalURL, "preview") || strings.Contains(finalURL, "unavailable")
	}

	region := extractField(content3, `"countryCode":"`, `"`)
	inSupported := extractField(content3, `"inSupportedLocation":`, `,`)

	if region == "" {
		return UnlockResult{Platform: "disney", Status: "false"}
	}
	if isUnavailable {
		return UnlockResult{Platform: "disney", Status: "false"}
	}
	if inSupported == "false" {
		return UnlockResult{Platform: "disney", Status: "restricted", Region: strings.ToUpper(region)}
	}
	if inSupported == "true" {
		return UnlockResult{Platform: "disney", Status: "true", Region: strings.ToUpper(region)}
	}

	return UnlockResult{Platform: "disney", Status: "false"}
}

func checkPrimeVideo() UnlockResult {
	req, _ := http.NewRequest("GET", "https://www.primevideo.com", nil)
	req.Header.Set("User-Agent", uaBrowser)

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "prime_video", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	isBlocked := strings.Contains(content, "isServiceRestricted")
	region := extractField(content, `"currentTerritory":"`, `"`)

	if isBlocked {
		return UnlockResult{Platform: "prime_video", Status: "false"}
	}
	if region != "" {
		return UnlockResult{Platform: "prime_video", Status: "true", Region: strings.ToUpper(region)}
	}

	return UnlockResult{Platform: "prime_video", Status: "false"}
}

func checkYouTubePremium() UnlockResult {
	req, _ := http.NewRequest("GET", "https://www.youtube.com/premium", nil)
	req.Header.Set("User-Agent", uaBrowser)
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "YSC=FSCWhKo2Zgw; VISITOR_PRIVACY_METADATA=CgJERRIEEgAgYQ%3D%3D; PREF=f7=4000; __Secure-YEC=CgtRWTBGTFExeV9Iayjele2yBjIKCgJERRIEEgAgYQ%3D%3D; SOCS=CAISOAgDEitib3FfaWRlbnRpdHlmcm9udGVuZHVpc2VydmVyXzIwMjQwNTI2LjAxX3AwGgV6aC1DTiACGgYIgMnpsgY; VISITOR_INFO1_LIVE=Di84mAIbgKY; __Secure-BUCKET=CGQ")

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "youtube_premium", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	if strings.Contains(content, "www.google.cn") {
		return UnlockResult{Platform: "youtube_premium", Status: "false", Region: "CN"}
	}

	if strings.Contains(content, "Premium is not available in your country") {
		return UnlockResult{Platform: "youtube_premium", Status: "false"}
	}

	region := extractField(content, `"INNERTUBE_CONTEXT_GL":"`, `"`)
	if region == "" {
		region = "UNKNOWN"
	}

	if strings.Contains(content, "ad-free") {
		return UnlockResult{Platform: "youtube_premium", Status: "true", Region: region}
	}

	return UnlockResult{Platform: "youtube_premium", Status: "false"}
}

func checkDazn() UnlockResult {
	body := `{"Version":"2","LandingPageKey":"generic","Languages":"zh-CN","Platform":"web","Manufacturer":"","PromoCode":"","PlatformAttributes":{}}`
	req, _ := http.NewRequest("POST", "https://startup.core.indazn.com/misl/v5/Startup", strings.NewReader(body))
	req.Header.Set("User-Agent", uaBrowser)
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://www.dazn.com")
	req.Header.Set("referer", "https://www.dazn.com/")
	req.Header.Set("x-session-id", "fd264e77-79d5-480c-a514-a275b649da14")

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "dazn", Status: "false"}
	}
	defer resp.Body.Close()
	body2, _ := io.ReadAll(resp.Body)
	content := string(body2)

	if strings.Contains(strings.ToLower(content), "security policy has been breached") {
		return UnlockResult{Platform: "dazn", Status: "false"}
	}

	isAllowed := extractField(content, `"isAllowed":`, `,`)
	region := extractField(content, `"GeolocatedCountry":"`, `"`)

	if isAllowed == "true" {
		return UnlockResult{Platform: "dazn", Status: "true", Region: strings.ToUpper(region)}
	}
	return UnlockResult{Platform: "dazn", Status: "false"}
}

func checkTVBAnywhere() UnlockResult {
	resp, err := followClient.Get("https://uapisfm.tvbanywhere.com.sg/geoip/check/platform/android")
	if err != nil {
		return UnlockResult{Platform: "tvbanywhere", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return UnlockResult{Platform: "tvbanywhere", Status: "false"}
	}

	region, _ := result["country"].(string)
	allowed, _ := result["allow_in_this_country"].(bool)

	if allowed {
		return UnlockResult{Platform: "tvbanywhere", Status: "true", Region: strings.ToUpper(region)}
	}
	return UnlockResult{Platform: "tvbanywhere", Status: "false"}
}

func checkIqiyi() UnlockResult {
	req, _ := http.NewRequest("GET", "https://www.iq.com/", nil)
	req.Header.Set("User-Agent", uaBrowser)

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "iqiyi", Status: "false"}
	}
	defer resp.Body.Close()

	region := ""
	for _, cookie := range resp.Header["Set-Cookie"] {
		idx := strings.Index(cookie, "mod=")
		if idx != -1 {
			start := idx + 4
			end := strings.IndexAny(cookie[start:], ";& ")
			if end == -1 {
				region = cookie[start:]
			} else {
				region = cookie[start : start+end]
			}
			break
		}
	}

	if region == "" {
		return UnlockResult{Platform: "iqiyi", Status: "false"}
	}

	region = strings.ToUpper(region)
	if region == "NTW" {
		region = "TW"
	}
	if region == "INTL" {
		return UnlockResult{Platform: "iqiyi", Status: "reject"}
	}

	return UnlockResult{Platform: "iqiyi", Status: "true", Region: region}
}

// ────────────────────────────────────────────────
// AI 工具
// ────────────────────────────────────────────────
func checkClaude() UnlockResult {
	req, _ := http.NewRequest("GET", "https://claude.ai/", nil)
	req.Header.Set("User-Agent", uaBrowser)
	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "claude", Status: "false"}
	}
	defer resp.Body.Close()
	finalURL := resp.Request.URL.String()
	if strings.HasPrefix(finalURL, "https://claude.ai/") {
		return UnlockResult{Platform: "claude", Status: "true"}
	}
	if finalURL == "https://www.anthropic.com/app-unavailable-in-region" {
		return UnlockResult{Platform: "claude", Status: "false"}
	}
	return UnlockResult{Platform: "claude", Status: "false"}
}

func checkChatGPT() UnlockResult {
	req1, _ := http.NewRequest("GET", "https://api.openai.com/compliance/cookie_requirements", nil)
	req1.Header.Set("Authorization", "Bearer null")
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := followClient.Do(req1)
	if err != nil {
		return UnlockResult{Platform: "chatgpt", Status: "false"}
	}
	defer resp1.Body.Close()
	body1, _ := io.ReadAll(resp1.Body)
	content1 := string(body1)
	req2, _ := http.NewRequest("GET", "https://ios.chat.openai.com/", nil)
	req2.Header.Set("User-Agent", uaBrowser)
	resp2, err := followClient.Do(req2)
	if err != nil {
		return UnlockResult{Platform: "chatgpt", Status: "false"}
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	content2 := string(body2)
	result1 := strings.Contains(strings.ToLower(content1), "unsupported_country")
	result2 := strings.Contains(strings.ToLower(content2), "vpn")
	if !result1 && !result2 {
		return UnlockResult{Platform: "chatgpt", Status: "true"}
	}
	if result1 && result2 {
		return UnlockResult{Platform: "chatgpt", Status: "false"}
	}
	if !result1 && result2 {
		return UnlockResult{Platform: "chatgpt", Status: "reject"}
	}
	if result1 && !result2 {
		return UnlockResult{Platform: "chatgpt", Status: "reject"}
	}
	return UnlockResult{Platform: "chatgpt", Status: "false"}
}

func checkGemini() UnlockResult {
	req, _ := http.NewRequest("GET", "https://gemini.google.com", nil)
	req.Header.Set("User-Agent", uaBrowser)

	resp, err := followClient.Do(req)
	if err != nil {
		return UnlockResult{Platform: "gemini", Status: "false"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	if !strings.Contains(content, "45631641,null,true") {
		return UnlockResult{Platform: "gemini", Status: "false"}
	}

	return UnlockResult{Platform: "gemini", Status: "true"}
}

// ────────────────────────────────────────────────
// 工具函数
// ────────────────────────────────────────────────

func extractField(content, prefix, suffix string) string {
	idx := strings.Index(content, prefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(prefix)
	end := strings.Index(content[start:], suffix)
	if end == -1 {
		return ""
	}
	return content[start : start+end]
}

func getLine(content string, n int) string {
	lines := strings.Split(content, "\n")
	if n <= 0 || n > len(lines) {
		return ""
	}
	return strings.TrimSpace(lines[n-1])
}

func currencyToCountry(currency string) string {
	m := map[string]string{
		"JPY": "JP", "USD": "US", "EUR": "EU",
		"GBP": "GB", "CNY": "CN", "KRW": "KR",
		"HKD": "HK", "TWD": "TW", "SGD": "SG",
		"AUD": "AU", "CAD": "CA", "INR": "IN",
		"BRL": "BR", "MXN": "MX", "RUB": "RU",
		"THB": "TH", "MYR": "MY", "IDR": "ID",
		"PHP": "PH", "VND": "VN", "SAR": "SA",
		"AED": "AE", "TRY": "TR", "NOK": "NO",
		"SEK": "SE", "DKK": "DK", "PLN": "PL",
		"CZK": "CZ", "HUF": "HU", "CHF": "CH",
		"NZD": "NZ", "ZAR": "ZA", "CLP": "CL",
		"ARS": "AR", "COP": "CO", "PEN": "PE",
	}
	if code, ok := m[strings.ToUpper(strings.TrimSpace(currency))]; ok {
		return code
	}
	return currency
}

func regionToCountry(region string) string {
	m := map[string]string{
		"Japan":                "JP",
		"United States":        "US",
		"Hong Kong":            "HK",
		"Taiwan":               "TW",
		"South Korea":          "KR",
		"Singapore":            "SG",
		"Australia":            "AU",
		"Canada":               "CA",
		"United Kingdom":       "GB",
		"Germany":              "DE",
		"France":               "FR",
		"Italy":                "IT",
		"Spain":                "ES",
		"Netherlands":          "NL",
		"Sweden":               "SE",
		"Norway":               "NO",
		"Denmark":              "DK",
		"Finland":              "FI",
		"Poland":               "PL",
		"Czech Republic":       "CZ",
		"Hungary":              "HU",
		"Switzerland":          "CH",
		"Austria":              "AT",
		"Belgium":              "BE",
		"Portugal":             "PT",
		"Ireland":              "IE",
		"New Zealand":          "NZ",
		"Brazil":               "BR",
		"Mexico":               "MX",
		"Argentina":            "AR",
		"Chile":                "CL",
		"Colombia":             "CO",
		"Peru":                 "PE",
		"India":                "IN",
		"Thailand":             "TH",
		"Malaysia":             "MY",
		"Indonesia":            "ID",
		"Philippines":          "PH",
		"Vietnam":              "VN",
		"Saudi Arabia":         "SA",
		"United Arab Emirates": "AE",
		"Turkey":               "TR",
		"China":                "CN",
		"Russia":               "RU",
		"South Africa":         "ZA",
	}
	if code, ok := m[strings.TrimSpace(region)]; ok {
		return code
	}
	return region
}
