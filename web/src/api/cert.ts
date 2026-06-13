import { request } from '@/util/request'

// 证书
export const getCert = () => request('/api/cert')
export const saveCert = (data: any) => request('/api/cert', { method: 'POST', body: JSON.stringify(data) })
export const deleteCert = (id: number) => request(`/api/cert/${id}`, { method: 'DELETE' })
export const applyCert = (id: number) => request(`/api/cert/${id}/apply`, { method: 'POST' })
export const getCertContent = (id: number) => request(`/api/cert/${id}/content`)

// ACME 账号
export const getAcme = () => request('/api/acme')
export const saveAcme = (data: any) => request('/api/acme', { method: 'POST', body: JSON.stringify(data) })
export const deleteAcme = (id: number) => request(`/api/acme/${id}`, { method: 'DELETE' })

// DNS 账号
export const getDns = () => request('/api/dns')
export const saveDns = (data: any) => request('/api/dns', { method: 'POST', body: JSON.stringify(data) })
export const deleteDns = (id: number) => request(`/api/dns/${id}`, { method: 'DELETE' })