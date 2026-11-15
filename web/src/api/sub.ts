import { request } from '@/util/request'
import type { Inbound } from '@/api/inbound'
import type { User } from '@/api/user'

// 后端按整条用户和整条入站来生成链接
export const getUri = (data: { user: User; inbound: Inbound }): Promise<{ uri: string }> =>
    request('/api/sub/uri', {
        method: 'POST',
        body: JSON.stringify(data),
    })

export const getUrl = (token: string): Promise<{ urls: string[] }> =>
    request('/api/sub/url', {
        method: 'POST',
        body: JSON.stringify({ token }),
    })

export const getJson = (data: {
    user: User
    inbound: Inbound
    format: string
}): Promise<{ json?: string }> =>
    request('/api/sub/json', {
        method: 'POST',
        body: JSON.stringify(data),
    })
