// 后台任务(申请证书、同步等)的实时输出，每行是一条 JSON：{ time, level, message }，level 为 DONE 表示结束
export const task = (id: string) => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('token') ?? ''

    return new WebSocket(`${protocol}//${location.host}/api/task/${id}?token=${token}`)
}
