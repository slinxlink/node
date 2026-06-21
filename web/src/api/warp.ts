import { request } from '@/util/request'

export const getWarp = () => request('/api/warp')
export const deleteWarp = () => request('/api/warp', {
    method: 'DELETE'
})
export const registerWarp = () => request('/api/warp/register', {
    method: 'POST'
})
export const refreshWarp = () => request('/api/warp/refresh', {
    method: 'POST'
})
export const setWarpAutoUpdate = (day: number) => request('/api/warp/auto-update', {
    method: 'PUT',
    body: JSON.stringify(day)
})
export const setWarpLicense = (license: string) => request('/api/warp/license', {
    method: 'PUT',
    body: JSON.stringify(license)
})