import { request } from '@/util/request'

export const getConfig = () => request('/api/config')

export const updateConfig = (data: any) => request('/api/config', {
    method: 'PUT',
    body: JSON.stringify(data)
})

export const resetConfig = () => request('/api/config/reset', {
    method: 'POST'
})

export const Log = () => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''
    return new WebSocket(`${protocol}//${location.host}/api/log/slinx?token=${token}`)
}