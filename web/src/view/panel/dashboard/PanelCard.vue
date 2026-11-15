<template>
    <Box title="SLINX node">
        <template #footer>
            <Anchor href="https://github.com/slinxlink/node" external>
                <Image class="invert" :src="github" alt="GitHub" />
                GitHub 仓库
            </Anchor>
            <Anchor href="https://t.me/slinxlink" external>
                <Image :src="telegram" alt="Telegram" />
                加入频道
            </Anchor>
            <Button :disabled="!hasUpdate" :class="{ update: hasUpdate }" @click="openUpdate">
                <Icon name="new_releases" />
                {{ hasUpdate ? '有新版本' : version }}
            </Button>
        </template>
    </Box>
</template>

<script setup lang="ts">
import github from '@/asset/image/social/github.svg'
import telegram from '@/asset/image/social/telegram.svg'
import Box from '@/component/layout/Box.vue'
import Anchor from '@/component/widget/Anchor.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Image from '@/component/widget/Image.vue'
import Update from '@/view/panel/dashboard/Update.vue'
import { getVersion, checkUpdate, doUpdate } from '@/api/update'
import { modal } from '@/composable/modal'
import { toast } from '@/composable/toast'

const version = ref('...')
const hasUpdate = ref(false)
const latest = ref('')
const changelog = ref('')

onMounted(async () => {
    try {
        const [v, u] = await Promise.all([getVersion(), checkUpdate()])
        version.value = v.version
        hasUpdate.value = u.has_update
        latest.value = u.latest_version
        changelog.value = u.changelog
    } catch {
        // 检查更新依赖 GitHub，失败了就当没有新版本
    }
})

function openUpdate() {
    modal.show({
        title: '发现新版本',
        width: '500px',
        component: Update,
        props: {
            currentVersion: version.value,
            latestVersion: latest.value,
            changelog: changelog.value,
        },
        onConfirm: async () => {
            toast.loading('更新中...')
            try {
                await doUpdate()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '更新失败')
                return false
            }
            toast.success('更新成功，等待重启中...')
            setTimeout(() => window.location.reload(), 3000)
        },
    })
}
</script>

<style scoped>
img {
    width: 15px;
    height: 15px;
}

.update {
    color: var(--yellow);
}
</style>
