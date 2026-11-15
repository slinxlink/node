import { request } from '@/util/request'

export interface UpdateCheck {
    has_update: boolean
    latest_version: string
    changelog: string
}

export const getVersion = (): Promise<{ version: string }> => request('/api/update/version')
export const checkUpdate = (): Promise<UpdateCheck> => request('/api/update/check')
export const doUpdate = () => request('/api/update', { method: 'POST' })
