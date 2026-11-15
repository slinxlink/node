// 证书页用到的显示：各种代码对应的中文，证书的状态

const modes: Record<string, string> = { dns: 'DNS 验证', http: 'HTTP 验证', manual: '手动上传' }
const acmeProviders: Record<string, string> = { letsencrypt: "Let's Encrypt", zerossl: 'ZeroSSL' }
const dnsProviders: Record<string, string> = { aliyun: '阿里云', cloudflare: 'Cloudflare' }

export const display = {
    mode: (mode: string) => modes[mode] ?? mode,
    acmeProvider: (provider: string) => acmeProviders[provider] ?? provider,
    dnsProvider: (provider: string) => dnsProviders[provider] ?? provider,

    // 按过期时间分：没有时间是未知，过期了是已过期，不到 10 天是即将过期，其余正常；颜色是 color.css 里的令牌名
    status: (expireAt: string) => {
        if (!expireAt) return { text: '未知', color: 'text-200' }

        const days = (new Date(expireAt).getTime() - Date.now()) / 86400000

        if (days < 0) return { text: '已过期', color: 'red' }
        if (days < 10) return { text: '即将过期', color: 'yellow' }

        return { text: '正常', color: 'green' }
    },
}
