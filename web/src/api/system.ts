import { request } from '@/util/request'

export const getSystem = () => request('/api/system/status')

export const getStats = () => request('/api/stats')

export const getSystemLog = () => request('/api/system/log')