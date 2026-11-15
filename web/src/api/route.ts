import { request } from '@/util/request'

// 和后端 database.Rule 对应：同一个 Sort 的几行合起来是一条规则
export interface Rule {
    ID: number
    Sort: number
    Index: number
    Key: string
    Value: string
}

export const getRule = (): Promise<Rule[]> =>
    request('/api/rule').then((res: { data: Rule[] }) => res.data)

// 把这些入站(端口)的流量绑定到 tag 对应的端点出去
export const saveRule = (tag: string, inbounds: string[]) =>
    request('/api/rule', {
        method: 'PUT',
        body: JSON.stringify({ tag, inbounds }),
    })
