import { request } from '@/util/request'

export const refreshRuleset = () => request('/api/ruleset/refresh', {
    method: 'POST'
})