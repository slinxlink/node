import { request } from '@/util/request'

export const getRule = () => request('/api/rule')
export const saveRule = (tag: string, inbounds: string[]) => request('/api/rule', {
    method: 'PUT',
    body: JSON.stringify({ tag, inbounds })
})