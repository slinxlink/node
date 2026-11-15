import { request } from '@/util/request'

export interface Config {
    ID: number
    Username: string

    Port: number
    Path: string
    Domain: string
    IPv4: string
    IPv6: string

    SubEnable: boolean
    SubPath: string
    SubPort: number
    RulesetAutoUpdate: boolean

    LogEnable: boolean
    LogLevel: string
    LogPath: string

    BoardEnable: boolean

    BBR: boolean

    Repo: string
    StartedAt: string
    Dir: string
}

export const getConfig = (): Promise<Config> => request('/api/config')

export const updateConfig = (data: Partial<Config>) =>
    request('/api/config', {
        method: 'PUT',
        body: JSON.stringify(data),
    })

export const resetConfig = () =>
    request('/api/config/reset', {
        method: 'POST',
    })

export const Log = () => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''
    return new WebSocket(`${protocol}//${location.host}/api/log/slinx?token=${token}`)
}
