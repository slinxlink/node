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

export const getRoute = () => request('/api/detect/route')
export const fetchRoute = () => request('/api/detect/route/fetch', {
    method: 'POST'
})