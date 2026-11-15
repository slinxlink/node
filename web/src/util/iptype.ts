const orgTypeMap: Record<string, string> = {
    hosting: '机房',
    isp: '住宅',
    business: '企业',
    education: '教育',
    government: '政府',
    banking: '金融',
    mobile: '移动',
    proxy: '代理',
}

export function translateType(type?: string): string {
    if (!type) return '-'

    return orgTypeMap[type] || type
}

// 家宽、教育、政府是好的，机房、企业、金融、代理是差的，移动居中；返回 color.css 里的令牌名，没有归类的返回空
export function typeColor(type?: string): string {
    if (!type) return ''

    if (['isp', 'education', 'government'].includes(type)) return 'green'
    if (['hosting', 'business', 'banking', 'proxy'].includes(type)) return 'red'
    if (type === 'mobile') return 'yellow'

    return ''
}
