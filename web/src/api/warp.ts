import { request } from '@/util/request'

// 和后端 database.Warp 对应：库里存的账户
export interface Warp {
    ID: number
    AccessToken: string
    DeviceID: string
    LicenseKey: string
    PublicKey: string
    PrivateKey: string
    AutoUpdate: number
    CreatedAt: string
    UpdatedAt: string
}

// 和后端 service.WarpData 对应：注册、刷新时拿到的账户状态，不落库
export interface WarpData {
    Address: string
    Reserved: string
    PublicKey: string
    PrivateKey: string
    PeerAddress: string
    PeerPort: number
    PeerPublicKey: string
    DeviceName: string
    DeviceModel: string
    DeviceEnabled: boolean
    AccountType: string
    Role: string
    WarpPlusData: number
    Quota: number
    Usage: number
}

// 没有账户时 data 是空
export const getWarp = (): Promise<Warp | null> =>
    request('/api/warp').then((res: { data: Warp | null }) => res.data)

export const deleteWarp = () =>
    request('/api/warp', {
        method: 'DELETE',
    })

export const registerWarp = (): Promise<WarpData> =>
    request('/api/warp/register', {
        method: 'POST',
    }).then((res: { data: WarpData }) => res.data)

export const refreshWarp = (): Promise<WarpData> =>
    request('/api/warp/refresh', {
        method: 'POST',
    }).then((res: { data: WarpData }) => res.data)

export const setWarpAutoUpdate = (day: number): Promise<Warp> =>
    request('/api/warp/auto-update', {
        method: 'PUT',
        body: JSON.stringify(day),
    }).then((res: { data: Warp }) => res.data)

export const setWarpLicense = (license: string): Promise<WarpData> =>
    request('/api/warp/license', {
        method: 'PUT',
        body: JSON.stringify(license),
    }).then((res: { data: WarpData }) => res.data)
