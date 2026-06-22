import { request } from '@/util/request'

export const getUri = (data: { user: any, inbound: any }) => request('/api/sub/uri', {
    method: 'POST',
    body: JSON.stringify(data)
})

export const getUrl = (token: string) => request('/api/sub/url', {
    method: 'POST',
    body: JSON.stringify({ token })
})

export const getJson = (data: { user: any, inbound: any, format: string }) => request('/api/sub/json', {
    method: 'POST',
    body: JSON.stringify(data)
})