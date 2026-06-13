import { request } from '@/util/request'

export const getVersion = () => request('/api/update/version')
export const checkUpdate = () => request('/api/update/check')
export const doUpdate = () => request('/api/update', { method: 'POST' })