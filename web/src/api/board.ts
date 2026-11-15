import { request } from '@/util/request'

// 和后端 database.Board 对应
export interface Board {
    ID: number
    Enable: boolean
    Name: string
    Host: string
    NodeID: number
    Key: string
    Inbound: number
    Type: string
    SyncInterval: number
    CreatedAt: string
    UpdatedAt: string
}

// 和后端 database.BoardUser 对应：面板那边的用户；AliveIP 是 JSON 数组字符串
export interface BoardUser {
    ID: number
    BoardID: number
    UserID: number
    UUID: string
    Passwd: string
    Upload: number
    Download: number
    AliveIP: string
    UpdatedAt: string
}

export const getBoards = (): Promise<Board[]> => request('/api/board')

// 保存时提交的内容：新增不带 ID；编辑要把原来的整条带上
export type BoardData = Omit<Board, 'ID' | 'CreatedAt' | 'UpdatedAt'> &
    Partial<Pick<Board, 'ID' | 'CreatedAt' | 'UpdatedAt'>>

export const saveBoard = (data: BoardData) =>
    request('/api/board/save', {
        method: 'PUT',
        body: JSON.stringify(data),
    })

export const connectBoard = (data: { Host: string; NodeID: number; Key: string }) =>
    request('/api/board/connect', {
        method: 'POST',
        body: JSON.stringify(data),
    })

export const deleteBoard = (id: number) =>
    request(`/api/board/${id}`, {
        method: 'DELETE',
    })

export const toggleBoard = (id: number) =>
    request(`/api/board/${id}/toggle`, {
        method: 'PUT',
    })

export const getBoardUser = (id: number): Promise<BoardUser[]> => request(`/api/board/${id}/user`)
