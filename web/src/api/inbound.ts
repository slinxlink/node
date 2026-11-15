import { request } from '@/util/request'

// 和后端 database.Inbound 一一对应；ALPN、Certs、RealityShortIDs 在库里是字符串
export interface Inbound {
    ID: number
    Enable: boolean
    Name: string
    Protocol: string
    Port: number
    CreatedAt: string
    Upload: number
    Download: number
    TrafficLimit: number
    ResetCycle: string
    ResetDay: number
    ResetAt: string

    Transport: string
    WsPath: string
    WsHost: string
    WsPingInterval: number

    UDPTimeout: number
    MasqueradeEnabled: boolean
    MasqueradeType: string
    MasqueradeURL: string
    RewriteHost: boolean
    IgnoreTLSVerify: boolean
    MasqueradePath: string
    MasqueradeCode: number
    MasqueradeBody: string
    ObfsType: string
    ObfsPassword: string
    ObfsMinPacketSize: number
    ObfsMaxPacketSize: number

    TuicCongestionControl: string
    TuicAuthTimeout: number
    TuicZeroRTT: boolean
    TuicHeartbeat: number
    TuicUDPRelayMode: string

    HopEnabled: boolean
    HopPort: string
    HopInterval: string

    AnyTLSIdleSessionCheckInterval: number
    AnyTLSIdleSessionTimeout: number
    AnyTLSMinIdleSession: number

    TLSType: string
    ServerName: string
    CipherSuites: string
    TLSMinVersion: string
    TLSMaxVersion: string
    UTLS: string
    Insecure: boolean
    ALPN: string // 逗号分隔
    Certs: string // JSON 数组，证书 ID
    ECHEnabled: boolean
    ECHKey: string
    ECHConfig: string

    RealityServerName: string
    RealityServer: string
    RealityServerPort: number
    RealityMaxTimeDiff: number
    RealityShortIDs: string // JSON 数组
    RealityPrivateKey: string
    RealityPublicKey: string
    Flow: string
}

export const resetInboundTraffic = (id: number) =>
    request(`/api/inbound/${id}/reset`, {
        method: 'PUT',
    })

// 保存时提交的内容：新增不带 ID，编辑带 ID
export type InboundData = Omit<Inbound, 'ID' | 'CreatedAt' | 'Upload' | 'Download' | 'ResetAt'> & { ID?: number }

export const getInbounds = (): Promise<Inbound[]> => request('/api/inbound')

export const saveInbound = (data: InboundData) =>
    request('/api/inbound/save', {
        method: 'PUT',
        body: JSON.stringify(data),
    })

export const deleteInbound = (id: number) =>
    request(`/api/inbound/${id}`, {
        method: 'DELETE',
    })

export const toggleInbound = (id: number) =>
    request(`/api/inbound/${id}/toggle`, {
        method: 'PUT',
    })
