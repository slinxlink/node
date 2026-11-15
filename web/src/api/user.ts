import { request } from '@/util/request'

// 和后端 database.User 对应；Inbounds 在库里是 JSON 数组字符串，里面是入站 ID
export interface User {
    ID: number
    Enable: boolean
    Name: string
    Token: string
    Inbounds: string
    UUID: string
    Password: string
    Upload: number
    Download: number
    TrafficLimit: number
    ResetCycle: string
    ResetDay: number
    ResetAt: string
    AliveIP: string
    ExpireAt: string
    CreatedAt: string
    UpdatedAt: string
}

export const getUsers = (): Promise<User[]> => request('/api/user')

export const deleteUser = (id: number) =>
    request(`/api/user/${id}`, {
        method: 'DELETE',
    })

export const toggleUser = (id: number) =>
    request(`/api/user/${id}/toggle`, {
        method: 'PUT',
    })

export const resetUserTraffic = (id: number) =>
    request(`/api/user/${id}/reset`, {
        method: 'PUT',
    })

// 保存时提交的内容：新增不带 ID；编辑要把原来的整条带上，后端是整条覆盖，少了的字段会被清掉
export type UserData = Omit<User, 'ID' | 'CreatedAt' | 'UpdatedAt' | 'ExpireAt' | 'Upload' | 'Download' | 'AliveIP' | 'ResetAt'> &
    Partial<Pick<User, 'ID' | 'CreatedAt' | 'UpdatedAt' | 'ExpireAt'>>

export const saveUser = (data: UserData) =>
    request('/api/user/save', {
        method: 'PUT',
        body: JSON.stringify(data),
    })
