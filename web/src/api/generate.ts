import { request } from '@/util/request'

export const generatePort = (): Promise<{ port: number }> => request('/api/generate/port')

export const generateRealityTarget = (): Promise<{ domain: string }> =>
    request('/api/generate/reality-target')

export const generateShortIDs = (): Promise<{ short_ids: string[] }> =>
    request('/api/generate/shortids')

export const generateRealityKeyPair = (): Promise<{ private_key: string; public_key: string }> =>
    request('/api/generate/reality-keypair')

export const generatePassword = (): Promise<{ password: string }> =>
    request('/api/generate/password')

export const generateECHKeyPair = (domain: string): Promise<{ ech_key: string; ech_config: string }> =>
    request(`/api/generate/ech-keypair?domain=${domain}`)

export const generateToken = (): Promise<{ token: string }> => request('/api/generate/token')

export const generateUUID = (): Promise<{ uuid: string }> => request('/api/generate/uuid')

export const generateWireguardKeyPair = (): Promise<{ private_key: string; public_key: string }> =>
    request('/api/generate/wireguard-keypair')
