import { request } from '@/util/request'

export const login = (username: string, password: string) =>
    request('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
    })

export const changeCredentials = (data: {
    old_password: string
    new_username?: string
    new_password?: string
}) => request('/api/auth', { method: 'PUT', body: JSON.stringify(data) })

export const logout = () => request('/api/auth/logout', { method: 'POST' })