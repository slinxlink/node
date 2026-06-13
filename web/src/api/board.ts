import { request } from '@/util/request'

export const getBoards = () => request('/api/board')
export const saveBoard = (data: any) => request('/api/board/save', {
    method: 'PUT',
    body: JSON.stringify(data)
})
export const deleteBoard = (id: number) => request(`/api/board/${id}`, {
    method: 'DELETE'
})
export const toggleBoard = (id: number) => request(`/api/board/${id}/toggle`, {
    method: 'PUT'
})

export const getBoardUser = (id: number) => request(`/api/board/${id}/user`)