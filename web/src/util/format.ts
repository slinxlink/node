export function formatBytes(bytes: number = 0): string {
    if (!bytes || bytes <= 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.max(0, Math.min(Math.floor(Math.log(bytes) / Math.log(k)), sizes.length - 1))
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

export function formatUptime(startedAt: string): string {
    if (!startedAt) return '-'
    const diff = Date.now() - new Date(startedAt).getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    if (days > 0) return `${days}d ${hours}h`
    if (hours > 0) return `${hours}h ${mins}m`
    return `${mins}m`
}

export function formatTime(t: string): string {
    if (!t) return '-'
    return new Date(t).toLocaleString('sv-SE', { hour12: false }).replace('T', ' ')
}