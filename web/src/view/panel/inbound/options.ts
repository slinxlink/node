// view/panel/inbound/options.ts
// 入站表单里各个下拉和单选的选项

export const protocols = [
    { label: 'VLESS', value: 'vless' },
    { label: 'VMess', value: 'vmess' },
    { label: 'Hysteria', value: 'hysteria' },
    { label: 'Trojan', value: 'trojan' },
    { label: 'TUIC', value: 'tuic' },
    { label: 'AnyTLS', value: 'anytls' },
]

export const transports = [
    { label: 'RAW', value: 'raw' },
    { label: 'WebSocket', value: 'websocket' },
]

export const masqueradeTypes = [
    { label: '默认 (404)', value: 'default' },
    { label: '反向代理', value: 'proxy' },
    { label: '静态文件', value: 'file' },
    { label: '自定义内容', value: 'string' },
]

export const obfsTypes = [
    { label: '不启用', value: '' },
    { label: 'Salamander', value: 'salamander' },
    { label: 'Gecko (实验性)', value: 'gecko' },
]

export const congestions = [
    { label: 'Cubic', value: 'cubic' },
    { label: 'BBR', value: 'bbr' },
    { label: 'New Reno', value: 'new_reno' },
]

export const relayModes = [
    { label: '自动', value: '' },
    { label: 'Native', value: 'native' },
    { label: 'QUIC', value: 'quic' },
]

export const securities = [
    { label: '无', value: 'none' },
    { label: 'TLS', value: 'tls' },
    { label: 'Reality', value: 'reality' },
]

export const tlsVersions = ['1.0', '1.1', '1.2', '1.3'].map((v) => ({ label: v, value: v }))

export const utls = [
    { label: 'Chrome', value: 'chrome' },
    { label: 'Firefox', value: 'firefox' },
    { label: 'Safari', value: 'safari' },
    { label: 'iOS', value: 'ios' },
    { label: 'Android', value: 'android' },
    { label: 'Edge', value: 'edge' },
    { label: '360', value: '360' },
    { label: 'QQ', value: 'qq' },
    { label: 'Random', value: 'random' },
    { label: 'Randomized', value: 'randomized' },
]

export const alpns = [
    { label: 'h3', value: 'h3' },
    { label: 'h2', value: 'h2' },
    { label: 'http/1.1', value: 'http/1.1' },
]

// 加密套件：标签是缩写，值是 Go 标准库里的名字
const suites: [string, string][] = [
    ['ECDHE_ECDSA_AES_128_CBC', 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA'],
    ['ECDHE_ECDSA_AES_256_CBC', 'TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA'],
    ['ECDHE_RSA_AES_128_CBC', 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA'],
    ['ECDHE_RSA_AES_256_CBC', 'TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA'],
    ['ECDHE_ECDSA_AES_128_GCM', 'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256'],
    ['ECDHE_ECDSA_AES_256_GCM', 'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384'],
    ['ECDHE_RSA_AES_128_GCM', 'TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256'],
    ['ECDHE_RSA_AES_256_GCM', 'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384'],
    ['ECDHE_ECDSA_CHACHA20_POLY1305', 'TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256'],
    ['ECDHE_RSA_CHACHA20_POLY1305', 'TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256'],
]

export const cipherSuites = [
    { label: '自动', value: '' },
    ...suites.map(([label, value]) => ({ label, value })),
]

export const flows = [
    { label: '无', value: '' },
    { label: 'xtls-rprx-vision', value: 'xtls-rprx-vision' },
]
