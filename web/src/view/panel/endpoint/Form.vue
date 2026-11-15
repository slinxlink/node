<!-- view/panel/endpoint/Form.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>启用</template>
            <Toggle v-model:value="form.Enable" />
        </Field>
        <Field>
            <template #prefix>标签</template>
            <Input v-model:value="form.Tag" placeholder="tag" :disabled="form.Tag === 'warp'" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>协议</template>
            <Select v-model:value="form.Type" :items="types" />
        </Field>
        <Field>
            <template #prefix>本地 IPv4 地址</template>
            <Input v-model:value="form.Address4" placeholder="0.0.0.0/32" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>本地 IPv6 地址</template>
            <Input v-model:value="form.Address6" placeholder="::/128" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>公钥</template>
            <Input v-model:value="form.PublicKey" readonly />
        </Field>
        <Field>
            <template #prefix>私钥</template>
            <Input v-model:value="form.PrivateKey" readonly>
                <template #trailing>
                    <Button variant="text" :tabindex="-1" @click="genKeyPair">
                        <Icon name="casino" />
                    </Button>
                </template>
            </Input>
        </Field>
        <Field>
            <template #prefix>MTU</template>
            <Input v-model:value="form.MTU" type="number" :min="0" placeholder="1408" />
        </Field>
        <Field>
            <template #prefix>对端地址</template>
            <Input v-model:value="peer" placeholder="engage.cloudflareclient.com:2408" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>Reserved</template>
            <Input v-model:value="form.Reserved" placeholder="0,0,0" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>对端公钥</template>
            <Input v-model:value="form.PeerPublicKey" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>允许的 IP</template>
            <Input v-model:value="form.AllowedIPs" placeholder="0.0.0.0/0,::/0" disabled />
        </Field>
        <Field>
            <template #prefix>UDP 超时 (s)</template>
            <Input v-model:value="form.UDPTimeout" type="number" :min="0" placeholder="秒，0 表示默认" />
        </Field>
        <Field>
            <template #prefix>Workers</template>
            <Input v-model:value="form.Workers" type="number" :min="0" placeholder="0 表示默认" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { generateWireguardKeyPair } from '@/api/generate'
import { saveEndpoint } from '@/api/endpoint'
import type { Endpoint } from '@/api/endpoint'
import { toast } from '@/composable/toast'

const props = defineProps<{
    endpoint?: Endpoint // 不传是添加，传了是编辑
}>()

const types = [{ label: 'WireGuard', value: 'wireguard' }]

// 数字都用字符串存着，提交时再转
const form = ref({
    Enable: true,
    Tag: '',
    Type: 'wireguard',
    Address4: '172.16.0.2/32',
    Address6: 'fd00::2/128',
    PublicKey: '',
    PrivateKey: '',
    MTU: '1408',
    PeerAddress: '',
    PeerPort: '2408',
    Reserved: '0,0,0',
    PeerPublicKey: '',
    AllowedIPs: '0.0.0.0/0,::/0',
    UDPTimeout: '0',
    Workers: '0',
})

// 对端地址和端口合成一个输入框：最后一个冒号后面是端口
const peer = computed({
    get: () => `${form.value.PeerAddress}:${form.value.PeerPort}`,
    set: (val: string) => {
        const i = val.lastIndexOf(':')

        form.value.PeerAddress = i < 0 ? val : val.slice(0, i)
        form.value.PeerPort = String(parseInt(i < 0 ? '' : val.slice(i + 1)) || 0)
    },
})

async function genKeyPair() {
    const res = await generateWireguardKeyPair()

    form.value.PrivateKey = res.private_key
    form.value.PublicKey = res.public_key
}

// 库里的 JSON 数组转成输入框里的文字，解析不了就当空
function parseList(text: string): unknown[] {
    try {
        const list = JSON.parse(text || '[]')

        return Array.isArray(list) ? list : []
    } catch {
        return []
    }
}

onMounted(() => {
    const ep = props.endpoint

    if (!ep) {
        genKeyPair()

        return
    }

    const address = parseList(ep.Address)

    form.value = {
        Enable: ep.Enable,
        Tag: ep.Tag,
        Type: ep.Type,
        Address4: String(address[0] ?? ''),
        Address6: String(address[1] ?? ''),
        PublicKey: ep.PublicKey,
        PrivateKey: ep.PrivateKey,
        MTU: String(ep.MTU),
        PeerAddress: ep.PeerAddress,
        PeerPort: String(ep.PeerPort),
        Reserved: parseList(ep.Reserved).join(','),
        PeerPublicKey: ep.PeerPublicKey,
        AllowedIPs: ep.AllowedIPs,
        UDPTimeout: String(ep.UDPTimeout),
        Workers: String(ep.Workers),
    }
})

// 成功返回 true，失败返回 false（弹窗保持打开）
async function submit(): Promise<boolean> {
    const { Address4, Address6, Reserved, ...rest } = form.value

    const reserved = Reserved.split(',')
        .map((item) => parseInt(item.trim()))
        .filter((item) => !isNaN(item))

    // 编辑时把原来的整条带上：后端整条覆盖，不带的字段会被清掉
    const data = {
        ...props.endpoint,
        ...rest,
        MTU: Number(rest.MTU),
        PeerPort: Number(rest.PeerPort),
        UDPTimeout: Number(rest.UDPTimeout),
        Workers: Number(rest.Workers),
        Address: JSON.stringify([Address4, Address6]),
        Reserved: JSON.stringify(reserved),
    }

    try {
        await saveEndpoint(data)
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }

    toast.success('保存成功')

    return true
}

defineExpose({ submit })
</script>

<style scoped>
.form {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    text-align: left;

    /* 标题列拉宽，长一点的名字也放得下 */
    :deep(.prefix) {
        flex-basis: 100px;
    }
}
</style>
