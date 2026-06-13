import { request } from '@/util/request'

export const getCore = () => request('/api/core')
export const updateCore = (data: any) => request('/api/core', {
    method: 'PUT',
    body: JSON.stringify(data)
})
export const resetCore = () => request('/api/core/reset', { method: 'POST' })

export const getCoreStatus = () => request('/api/core/status')
export const startCore = () => request('/api/core/start', { method: 'POST' })
export const stopCore = () => request('/api/core/stop', { method: 'POST' })
export const restartCore = () => request('/api/core/restart', { method: 'POST' })

export const Log = () => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''
    return new WebSocket(`${protocol}//${location.host}/api/log/core?token=${token}`)
}

export const getCoreConfig = () => {
    const token = localStorage.getItem('token')
    return fetch('/api/core/config', {
        headers: { Authorization: `Bearer ${token}` }
    }).then(r => r.text())
}

export const getCoreProcess = () => request('/api/core/process')