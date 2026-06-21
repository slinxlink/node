import { request } from '@/util/request'

export const getIP = () => request('/api/detect/ip')
export const fetchIP = (source: string) => request('/api/detect/ip/fetch', {
    method: 'POST',
    body: JSON.stringify({ source })
})

export const getUnlock = () => request('/api/detect/unlock')
export const fetchUnlock = () => request('/api/detect/unlock/fetch', {
    method: 'POST'
})

export const getBackRoute = () => request('/api/detect/back-route')
export const fetchBackRoute = () => request('/api/detect/back-route/fetch', {
    method: 'POST'
})