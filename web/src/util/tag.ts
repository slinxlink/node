export const protocolMap: Record<string, string> = {
    vless: 'primary',
    vmess: 'green',
    hysteria: 'blue',
    trojan: 'purple',
    tuic: 'orange',
}

export function protocol(value: string): string {
    return protocolMap[value] ?? 'gray'
}