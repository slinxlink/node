import { request } from '@/util/request'

interface Usage {
    percent: number
    used: number
    total: number
}

export interface SystemStatus {
    cpu: { percent: number; cores: number; mhz: number }
    ram: Usage
    swap: Usage
    disk: Usage
    connections: { tcp: number; udp: number }
    network: { upload: number; download: number }
}

export interface Stats {
    Upload: number
    Download: number
}

export interface SystemLog {
    CPU: number
    RAM: number
    Load: number
    Upload: number
    Download: number
    CreatedAt: string
}

export const getSystem = (): Promise<SystemStatus> => request('/api/system/status')

export const getStats = (): Promise<Stats> => request('/api/stats')

export const getSystemLog = (): Promise<SystemLog[]> => request('/api/system/log')
