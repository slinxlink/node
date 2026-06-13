import { getConfig } from '@/api/config'

export const configStore = reactive({
    BoardEnable: false,
    Username: '',
    SubPath: '/link',
})

export async function loadConfig() {
    const res = await getConfig()
    configStore.BoardEnable = res.BoardEnable
    configStore.Username = res.Username
    configStore.SubPath = res.SubPath
}