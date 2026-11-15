import { request } from '@/util/request'

// 和后端 database.Cert 对应；Mode 是 dns / http / manual，Acme、Dns 是账号 ID
export interface Cert {
    ID: number
    Domain: string
    CertPath: string
    KeyPath: string
    Mode: string
    Acme: number
    Dns: number
    AutoRenew: boolean
    UpdatedAt: string
    ExpireAt: string
}

// 和后端 database.Acme 对应
export interface Acme {
    ID: number
    Email: string
    Provider: string
    PrivateKey: string
    EabKid: string
    EabHmac: string
}

// 和后端 database.DnsAccount 对应
export interface DnsAccount {
    ID: number
    Name: string
    Provider: string
    Key: string
    Secret: string
}

// 证书
export const getCert = (): Promise<Cert[]> => request('/api/cert')

// 保存时提交的内容：新增不带 ID；手动上传时还要带公钥、私钥的文本
export type CertData = Omit<Cert, 'ID' | 'CertPath' | 'KeyPath' | 'UpdatedAt' | 'ExpireAt'> &
    Partial<Pick<Cert, 'ID' | 'CertPath' | 'KeyPath' | 'UpdatedAt' | 'ExpireAt'>> & {
        CertContent?: string
        KeyContent?: string
    }

export const saveCert = (data: CertData): Promise<Cert> =>
    request('/api/cert', { method: 'POST', body: JSON.stringify(data) })

export const deleteCert = (id: number) => request(`/api/cert/${id}`, { method: 'DELETE' })

export const toggleCertRenew = (id: number) => request(`/api/cert/${id}/renew`, { method: 'PUT' })

// 申请证书是后台任务，返回任务号，进度用 task() 连接实时读取
export const applyCert = (id: number): Promise<{ task_id?: string }> =>
    request(`/api/cert/${id}/apply`, { method: 'POST' })

export const getCertContent = (id: number): Promise<{ cert: string; key: string }> =>
    request(`/api/cert/${id}/content`)

// ACME 账号
export const getAcme = (): Promise<Acme[]> => request('/api/acme')

export type AcmeData = Omit<Acme, 'ID' | 'PrivateKey'> & Partial<Pick<Acme, 'ID' | 'PrivateKey'>>

// 新增账号要向 ACME 服务商注册，也是后台任务
export const saveAcme = (data: AcmeData): Promise<{ task_id?: string }> =>
    request('/api/acme', { method: 'POST', body: JSON.stringify(data) })

export const deleteAcme = (id: number) => request(`/api/acme/${id}`, { method: 'DELETE' })

// DNS 账号
export const getDns = (): Promise<DnsAccount[]> => request('/api/dns')

export type DnsData = Omit<DnsAccount, 'ID'> & Partial<Pick<DnsAccount, 'ID'>>

export const saveDns = (data: DnsData) =>
    request('/api/dns', { method: 'POST', body: JSON.stringify(data) })

export const deleteDns = (id: number) => request(`/api/dns/${id}`, { method: 'DELETE' })
