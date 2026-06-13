import { request } from '@/util/request'

export const getInbounds = () => request('/api/inbound')
export const saveInbound = (data: any) => request('/api/inbound/save', {
    method: 'PUT',
    body: JSON.stringify(data)
})
export const deleteInbound = (id: number) => request(`/api/inbound/${id}`, {
    method: 'DELETE'
})
export const toggleInbound = (id: number) => request(`/api/inbound/${id}/toggle`, {
    method: 'PUT'
})