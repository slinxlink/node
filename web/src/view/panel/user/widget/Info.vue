<template>
    <div class="info">
        <FormRow title="启用">
            <span class="tag" :class="user.Enable ? 'green' : 'red'">{{ user.Enable ? '开启' : '关闭' }}</span>
        </FormRow>
        <FormRow title="名称">
            <span>{{ user.Name || '-' }}</span>
        </FormRow>
        <FormRow title="Token">
            <span class="mono">{{ user.Token }}</span>
            <Copy :value="user.Token" size="sm" />
        </FormRow>
        <FormRow title="UUID">
            <span class="mono">{{ user.UUID }}</span>
            <Copy :value="user.UUID" size="sm" />
        </FormRow>
        <FormRow title="认证">
            <span class="mono">{{ user.Password }}</span>
            <Copy :value="user.Password" size="sm" />
        </FormRow>
        <FormRow title="创建时间">
            <span>{{ formatTime(user.CreatedAt) }}</span>
        </FormRow>
        <FormRow title="更新时间">
            <span>{{ formatTime(user.UpdatedAt) }}</span>
        </FormRow>
        <FormRow title="关联入站">
            <div class="tags">
                <span
                    v-for="ib in inboundTags"
                    :key="ib.id"
                    class="tag"
                    :class="ib.color"
                >
                    {{ ib.port }}
                </span>
            </div>
        </FormRow>
        <div class="divider">复制链接</div>
        <Link v-for="uri in uris" :key="uri.value" :label="uri.label" :color="uri.color" :name="uri.name" :value="uri.value" />
        <div class="divider">订阅信息</div>
        <Link v-for="url in urls" :key="url.value" :label="url.label" :color="url.color" :name="url.name" :value="url.value" />
    </div>
</template>

<script setup lang="ts">
import FormRow from '@/component/ui/FormRow.vue'
import Copy from '@/component/widget/Copy.vue'
import Link from '@/component/widget/Link.vue'
import { formatTime } from '@/util/format.ts'
import { getUri, getUrl } from '@/api/sub'

const props = defineProps<{
    user: any
    inbounds: any[]
}>()

const colorMap: Record<string, string> = {
    vless: 'primary',
    vmess: 'green',
    hysteria: 'blue',
}

const inboundTags = computed(() => {
    try {
        const ids: number[] = JSON.parse(props.user.Inbounds || '[]')
        return ids.map(id => {
            const ib = props.inbounds.find(i => i.ID === id)
            return ib ? { id, port: ib.Port, name: ib.Name, protocol: ib.Protocol, color: colorMap[ib.Protocol] ?? 'gray', inbound: ib } : null
        }).filter((ib): ib is { id: number, port: any, name: string, protocol: string, color: string, inbound: any } => ib !== null)
    } catch {
        return []
    }
})

const uris = ref<{ label: string, color: string, name: string, value: string }[]>([])
const urls = ref<{ label: string, color: string, name: string, value: string }[]>([])

onMounted(async () => {
    const [uriResults, urlResult] = await Promise.all([
        Promise.all(inboundTags.value.map(ib => getUri({ user: props.user, inbound: ib.inbound }))),
        getUrl(props.user.Token)
    ])
    
    uris.value = inboundTags.value.map((ib, i) => ({
        label: ib.protocol.toUpperCase(),
        color: ib.color,
        name: ib.name,
        value: uriResults[i].uri,
    }))

    urls.value = urlResult.urls.map((url: string, i: number) => ({
        label: i === 0 ? 'SUB' : 'CLASH',
        color: i === 0 ? 'green' : 'yellow',
        name: props.user.Name || props.user.Token,
        value: url,
    }))
})
</script>

<style scoped>
.info {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 20px;
}

.mono {
    font-family: monospace;
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
}
</style>