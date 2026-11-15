// 国家代码转成旗帜表情：每个字母加上区域指示符的偏移
export function countryFlag(code?: string): string {
    if (!code) return ''

    return code.toUpperCase().replace(/./g, (c) => String.fromCodePoint(c.charCodeAt(0) + 127397))
}
