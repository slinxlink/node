<template>
    <div class="card md">
        <div class="header">
            <div class="title">SLINX node</div>
        </div>
        <div class="actions">
            <a class="action-btn" href="https://github.com/slinxlink/node" target="_blank" rel="noopener noreferrer">
                <img class="svg" src="@/asset/image/social/github.svg" alt="GitHub" />
                GitHub 仓库
            </a>
            <a class="action-btn" href="https://t.me/slinxlink" target="_blank" rel="noopener noreferrer">
                <img class="svg" src="@/asset/image/social/telegram.svg" alt="Telegram" />
                加入频道
            </a>
            <button class="action-btn version" :disabled="!hasUpdate" :class="{ 'has-update': hasUpdate }" @click="openUpdate">
                <i class="icon">new_releases</i>
                {{ hasUpdate ? '有新版本' : version }}
            </button>
        </div>
    </div>
    <Drawer v-model="showUpdate" title="发现新版本" saveText="立即更新" @save="handleUpdate">
        <Update :current-version="version" :latest-version="latestVersion" :changelog="changelog" />
    </Drawer>
</template>

<script setup lang="ts">
import Drawer from '@/component/Drawer.vue'
import Update from '@/component/widget/Update.vue'
import { getVersion, checkUpdate, doUpdate } from '@/api/update'

const modal = inject<any>('modal')
const version = ref('...')
const hasUpdate = ref(false)
const latestVersion = ref('')
const changelog = ref('')
const showUpdate = ref(false)

onMounted(async () => {
    const [versionRes, updateRes] = await Promise.all([
        getVersion(),
        checkUpdate(),
    ])
    version.value = versionRes.version
    hasUpdate.value = updateRes.has_update
    latestVersion.value = updateRes.latest_version
    changelog.value = updateRes.changelog
})

function openUpdate() {
    showUpdate.value = true
}

async function handleUpdate() {
    modal.value?.show('confirm', '面板在更新完成后将会自动重启\n确认更新？', async () => {
        modal.value?.update('warn', '更新中...')
        try {
            await doUpdate()
            let countdown = 3
            modal.value?.update('success', `更新成功，等待重启中...\n${countdown}`)
            const timer = setInterval(() => {
                countdown--
                if (countdown <= 0) {
                    clearInterval(timer)
                    window.location.reload()
                } else {
                    modal.value?.update('success', `更新成功，等待重启中...\n${countdown}`)
                }
            }, 1000)
        } catch (err: any) {
            modal.value?.update('error', err?.error ?? '更新失败')
        }
    })
}
</script>

<style scoped>
.header {
    margin-bottom: 10px;
}
.actions {
    display: flex;
    margin-top: auto;
    margin-left: calc(-1 * var(--card-padding, 20px));
    margin-right: calc(-1 * var(--card-padding, 20px));
    margin-bottom: calc(-1 * var(--card-padding, 20px));
    border-top: 1px solid var(--color-bg);
    overflow: hidden;

    .action-btn {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 5px;
        padding: 10px;
        background: transparent;
        border-right: 1px solid var(--color-bg);
        color: var(--color-text);
        font-size: var(--font-size-sm);

        img {
            height: 15px;
            width: 15px;
        }

        &:first-child { border-radius: 0 0 0 20px; }
        &:nth-child(2) { border-radius: 0; }
        &:last-child { 
            border-right: none; 
            border-radius: 0 0 20px 0;
        }

        &:hover {
            background: var(--color-primary);
            color: var(--color-text-light);
        }

        &:disabled {
            opacity: 1;
        }

        &.has-update {
            color: var(--color-yellow);
        }
    }
}
</style>