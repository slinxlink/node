import { request } from '@/util/request'

export const generatePort = () => request('/api/generate/port')
export const generateRealityTarget = () => request('/api/generate/reality-target')
export const generateShortIDs = () => request('/api/generate/shortids')
export const generateRealityKeyPair = () => request('/api/generate/reality-keypair')
export const generateToken = () => request('/api/generate/token')
export const generateUUID = () => request('/api/generate/uuid')
export const generatePassword = () => request('/api/generate/password')