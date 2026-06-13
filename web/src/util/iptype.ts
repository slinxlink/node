const orgTypeMap: Record<string, string> = {
    'hosting':    '机房',
    'isp':        '家宽',
    'business':   '企业',
    'education':  '教育',
    'government': '政府',
    'banking':    '金融',
    'mobile':     '移动',
    'proxy':      '代理',
}

export function translateType(type?: string): string {
    if (!type) return '-'
    return orgTypeMap[type] || type
}

export function typeClass(type?: string): string {
    if (!type) return ''
    const green  = ['isp', 'education', 'government']
    const red    = ['hosting', 'business', 'banking', 'proxy']
    const yellow = ['mobile']
    return green.includes(type)  ? 'c-yes'
         : red.includes(type)    ? 'c-no'
         : yellow.includes(type) ? 'c-reject'
         : ''
}