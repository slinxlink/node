export const task = (id: string) => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''
    return new WebSocket(`${protocol}//${location.host}/api/task/${id}?token=${token}`)
}