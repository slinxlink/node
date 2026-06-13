import { request } from '@/util/request'

export const restartPanel = () => request('/api/restart', { method: 'POST' })