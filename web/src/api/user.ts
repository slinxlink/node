import { request } from '@/util/request'

export const getUsers = () => request('/api/user')
export const saveUser = (data: any) => request('/api/user/save', {
    method: 'PUT',
    body: JSON.stringify(data)
})
export const deleteUser = (id: number) => request(`/api/user/${id}`, {
    method: 'DELETE'
})
export const toggleUser = (id: number) => request(`/api/user/${id}/toggle`, {
    method: 'PUT'
})