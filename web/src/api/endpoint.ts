import { request } from '@/util/request'

export const getEndpoints = () => request('/api/endpoint')
export const saveEndpoint = (endpoint: any) => request('/api/endpoint/save', {
    method: 'PUT',
    body: JSON.stringify(endpoint)
})
export const deleteEndpoint = (id: number) => request(`/api/endpoint/${id}`, {
    method: 'DELETE'
})
export const toggleEndpoint = (id: number) => request(`/api/endpoint/${id}/toggle`, {
    method: 'PUT'
})
export const createWarpEndpoint = (data: any) => request('/api/endpoint/warp', {
    method: 'POST',
    body: JSON.stringify(data)
})