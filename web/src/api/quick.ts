import { request } from '@/util/request'

export const quick = () => request('/api/quick', {
    method: 'POST'
})