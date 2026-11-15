// view/panel/inbound/display.ts
// 入站列表里各项的显示：颜色是 color.css 里的令牌名，给 Stamp 用

const protocolColor: Record<string, string> = {
    vless: 'primary-300',
    vmess: 'green',
    hysteria: 'blue',
    trojan: 'purple',
    tuic: 'orange',
    anytls: 'yellow',
}

const tlsLabel: Record<string, string> = {
    tls: 'TLS',
    reality: 'Reality',
}

export const display = {
    protocolColor: (protocol: string) => protocolColor[protocol] ?? 'text-200',

    // hysteria 和 tuic 跑在 UDP 上，其余看传输层
    transport: (protocol: string, transport: string) => {
        if (protocol === 'hysteria' || protocol === 'tuic') return 'UDP'

        return transport === 'websocket' ? 'WS' : 'TCP'
    },

    // none 表示不加密，不显示
    tls: (value: string) => (value === 'none' ? '' : (tlsLabel[value] ?? value)),
}
