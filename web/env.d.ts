/// <reference types="vite/client" />

declare interface Window {
    __PANEL_PATH__: string

    // 订阅页：后端把这个用户的信息直接写进页面，没有这个用户或已停用时是 null
    __SUB_DATA__?: {
        user: import('@/api/user').User
        inbounds: import('@/api/inbound').Inbound[]
        uris: string[]
        urls: string[]
        jsons: string[]
    } | null
}
