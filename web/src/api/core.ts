import { request } from '@/util/request'

export interface Core {
    ID: number
    Name: string
    BinPath: string
    ConfigPath: string
    LogEnable: boolean
    LogLevel: string
    LogPath: string
    Version: string
    StartedAt: string
    UpdatedAt: string
}

export interface CoreProcess {
    memory: number
    threads: number
}

export const getCore = (): Promise<Core> => request('/api/core')

export const getCoreStatus = (): Promise<{ status: string }> => request('/api/core/status')
export const startCore = () => request('/api/core/start', { method: 'POST' })
export const stopCore = () => request('/api/core/stop', { method: 'POST' })
export const restartCore = () => request('/api/core/restart', { method: 'POST' })

// 核心没运行时接口会报错
export const getCoreProcess = (): Promise<CoreProcess> => request('/api/core/process')

export const updateCore = (data: Core) =>
    request('/api/core', {
        method: 'PUT',
        body: JSON.stringify(data),
    })

export const resetCore = () => request('/api/core/reset', { method: 'POST' })

// 核心实时日志的连接，断了由调用方重连
export const coreLog = () => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''

    return new WebSocket(`${protocol}//${location.host}/api/log/core?token=${token}`)
}

// 核心当前的配置文件，是纯文本的 JSON，不是接口的统一格式
export const getCoreConfig = (): Promise<string> => {
    const token = localStorage.getItem('token')

    return fetch('/api/core/config', {
        headers: { Authorization: `Bearer ${token}` },
    }).then((res) => res.text())
}
