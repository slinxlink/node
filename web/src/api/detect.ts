import { request } from '@/util/request'

// 和后端 database.IP 对应
export interface IPInfo {
    ID: number
    IP: string
    Source: string
    IPVersion: string
    ASN: string
    ASNOrg: string
    Country: string
    City: string
    Latitude: number // 纬度和经度都是 0 表示没查到
    Longitude: number
    IPType: string
    Native: string // native 原生、broadcast 广播，空是查不出来
    UpdatedAt: string
}

// 和后端 database.Unlock 对应；Status 是 true / false / reject
export interface Unlock {
    ID: number
    IP: string
    IPVersion: string
    Platform: string
    Status: string
    Region: string
    UpdatedAt: string
}

// 和后端 database.BackRoute 对应
export interface BackRoute {
    ID: number
    City: string
    Telecom: string
    Unicom: string
    Mobile: string
    UpdatedAt: string
}

export const getIP = (): Promise<IPInfo[]> =>
    request('/api/detect/ip').then((res: { data?: IPInfo[] }) => res.data ?? [])

export const fetchIP = (source: string): Promise<IPInfo> =>
    request('/api/detect/ip/fetch', {
        method: 'POST',
        body: JSON.stringify({ source }),
    }).then((res: { data: IPInfo }) => res.data)

export const getUnlock = (): Promise<Unlock[]> =>
    request('/api/detect/unlock').then((res: { data?: Unlock[] }) => res.data ?? [])

export const fetchUnlock = (): Promise<Unlock[]> =>
    request('/api/detect/unlock/fetch', {
        method: 'POST',
    }).then((res: { data?: Unlock[] }) => res.data ?? [])

export const getBackRoute = (): Promise<BackRoute[]> =>
    request('/api/detect/back-route').then((res: { data?: BackRoute[] }) => res.data ?? [])

export const fetchBackRoute = (): Promise<BackRoute[]> =>
    request('/api/detect/back-route/fetch', {
        method: 'POST',
    }).then((res: { data?: BackRoute[] }) => res.data ?? [])
