import router from '@/router'

export const request = async (url: string, options: RequestInit = {}) => {
    const token = localStorage.getItem('token')

    const res = await fetch(url, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
            ...options.headers,
        }
    })

    if (res.status === 401) {
        if (!url.includes('/auth/login')) {
            localStorage.removeItem('token')
            router.push('/login')
            return {}
        }
    }

    if (!res.ok) {
        const data = await res.json()
        throw data
    }

    return res.json()
}