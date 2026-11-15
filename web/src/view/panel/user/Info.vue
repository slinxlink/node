<!-- view/panel/user/Info.vue -->
<template>
    <div class="info">
        <Field>
            <template #prefix>启用</template>
            <Stamp :color="user.Enable ? 'green' : 'red'">{{ user.Enable ? '开启' : '关闭' }}</Stamp>
        </Field>
        <Field>
            <template #prefix>名称</template>
            {{ user.Name || '-' }}
        </Field>
        <Field>
            <template #prefix>Token</template>
            <span class="mono">{{ user.Token }}</span>
            <Copy :value="user.Token" />
        </Field>
        <Field>
            <template #prefix>UUID</template>
            <span class="mono">{{ user.UUID }}</span>
            <Copy :value="user.UUID" />
        </Field>
        <Field>
            <template #prefix>Password</template>
            <span class="mono">{{ user.Password }}</span>
            <Copy :value="user.Password" />
        </Field>
        <Field>
            <template #prefix>到期时间</template>
            <Expire :text="formatExpire(user.ExpireAt)" :level="expireLevel(user.ExpireAt)" />
        </Field>
        <Field>
            <template #prefix>创建时间</template>
            {{ formatTime(user.CreatedAt) }}
        </Field>
        <Field>
            <template #prefix>更新时间</template>
            {{ formatTime(user.UpdatedAt) }}
        </Field>
        <Field>
            <template #prefix>关联入站</template>
            <div class="stamps">
                <Stamp v-for="item in related" :key="item.ID" :color="display.protocolColor(item.Protocol)">
                    {{ item.Protocol }}:{{ item.Port }}<template v-if="item.Name"> - {{ item.Name }}</template>
                </Stamp>
                <span v-if="!related.length">-</span>
            </div>
        </Field>
        <Field>
            <Traffic :used="user.Upload + user.Download" :limit="user.TrafficLimit" />
        </Field>

        <template v-if="uris.length">
            <Divider>复制链接</Divider>
            <Link v-for="item in uris" :key="item.name + item.value" v-bind="item" />
        </template>

        <template v-if="urls.length">
            <Divider>订阅信息</Divider>
            <Link v-for="item in urls" :key="item.value" v-bind="item" />
        </template>
    </div>
</template>

<script setup lang="ts">
import Divider from '@/component/layout/Divider.vue'
import Field from '@/component/layout/Field.vue'
import Traffic from '@/component/module/Traffic.vue'
import Copy from '@/component/widget/Copy.vue'
import Stamp from '@/component/widget/Stamp.vue'
import { getJson, getUri, getUrl } from '@/api/sub'
import type { Inbound } from '@/api/inbound'
import type { User } from '@/api/user'
import { expireLevel, formatExpire, formatTime } from '@/util/format'
import { display } from '@/view/panel/inbound/display'
import Expire from '@/view/panel/user/Expire.vue'
import Link from '@/view/panel/user/Link.vue'

// 订阅页里后端已经把链接都备好了，直接给；面板里不给，自己去接口取
export interface LinkData {
    uris: string[]
    urls: string[]
    jsons: string[]
}

// 用户的详细信息：基本资料，每个关联入站的分享链接(可下载 sing-box 配置)，订阅地址
const props = defineProps<{
    user: User
    inbounds: Inbound[]
    data?: LinkData
}>()

interface Item {
    label: string
    color: string
    name: string
    value: string
    download?: string
}

// 用户关联的入站；找不到的(入站已被删)略过
const related = computed(() => {
    let ids: number[] = []

    try {
        ids = JSON.parse(props.user.Inbounds || '[]')
    } catch {
        return []
    }

    return ids.flatMap((id) => props.inbounds.find((item) => item.ID === id) ?? [])
})

const uris = ref<Item[]>([])
const urls = ref<Item[]>([])

const labels = ['SUB', 'CLASH', 'SURGE']
const colors = ['green', 'yellow', 'blue']

// 订阅地址：几个依次是通用、Clash、Surge
function urlItems(list: string[], user: User): Item[] {
    return list.map((value, i) => ({
        label: labels[i] ?? 'SUB',
        color: colors[i] ?? 'text-200',
        name: user.Name || user.Token,
        value,
    }))
}

// 后端给的链接是纯字符串：协议从开头取(hysteria2 去掉 2)，名字是 # 后面那段
function uriItems(data: LinkData): Item[] {
    return data.uris.map((uri, i) => {
        const protocol = (uri.split('://')[0] ?? '').replace('2', '')

        return {
            label: protocol.toUpperCase(),
            color: display.protocolColor(protocol),
            name: decodeURIComponent(uri.split('#')[1] ?? ''),
            value: uri,
            download: data.jsons[i] || undefined,
        }
    })
}

onMounted(async () => {
    const { user, data } = props

    if (data) {
        uris.value = uriItems(data)
        urls.value = urlItems(data.urls, user)

        return
    }

    const [uriResults, jsonResults, urlResult] = await Promise.all([
        Promise.all(related.value.map((inbound) => getUri({ user, inbound }))),
        Promise.all(related.value.map((inbound) => getJson({ user, inbound, format: 'singbox' }))),
        getUrl(user.Token),
    ])

    uris.value = related.value.map((inbound, i) => ({
        label: inbound.Protocol.toUpperCase(),
        color: display.protocolColor(inbound.Protocol),
        name: inbound.Name || String(inbound.Port),
        value: uriResults[i]?.uri ?? '',
        download: jsonResults[i]?.json || undefined,
    }))

    urls.value = urlItems(urlResult.urls, user)
})
</script>

<style scoped>
.info {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    text-align: left;

    :deep(.prefix) {
        flex-basis: 60px;
        align-self: flex-start; /* 行变高(关联入站多行)时，标题固定在第一行，不跟着整行居中 */
        min-height: 40px;
    }

    .mono {
        font-family: monospace;
        font-size: 12px;
        color: var(--text-200);
    }

    .stamps {
        display: flex;
        flex-wrap: wrap;
        gap: 5px;
        padding: 9px 0; /* 一行标签 22px，上下各 9px 正好 40px，和别的行的第一行对齐 */
    }
}
</style>
