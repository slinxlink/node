import { request } from '@/util/request'
import type { WarpData } from '@/api/warp'

// 和后端 database.Endpoint 对应；Address 在库里是 JSON 数组字符串，Reserved 同理
export interface Endpoint {
    ID: number
    Enable: boolean
    Tag: string
    Type: string
    MTU: number
    Address: string
    PrivateKey: string
    PublicKey: string
    UDPTimeout: number
    Workers: number
    PeerAddress: string
    PeerPort: number
    PeerPublicKey: string
    AllowedIPs: string
    Reserved: string
    CreatedAt: string
    UpdatedAt: string
}

export const getEndpoints = (): Promise<Endpoint[]> =>
    request('/api/endpoint').then((res: { data: Endpoint[] }) => res.data)

export const deleteEndpoint = (id: number) =>
    request(`/api/endpoint/${id}`, {
        method: 'DELETE',
    })

export const toggleEndpoint = (id: number) =>
    request(`/api/endpoint/${id}/toggle`, {
        method: 'PUT',
    })

// 保存时提交的内容：新增不带 ID；编辑要把原来的整条带上，后端是整条覆盖
export type EndpointData = Omit<Endpoint, 'ID' | 'CreatedAt' | 'UpdatedAt'> &
    Partial<Pick<Endpoint, 'ID' | 'CreatedAt' | 'UpdatedAt'>>

export const saveEndpoint = (data: EndpointData) =>
    request('/api/endpoint/save', {
        method: 'PUT',
        body: JSON.stringify(data),
    })

// 把 WARP 账户应用到 tag 为 warp 的端点；没有就新建
export const createWarpEndpoint = (data: WarpData) =>
    request('/api/endpoint/warp', {
        method: 'POST',
        body: JSON.stringify(data),
    })
