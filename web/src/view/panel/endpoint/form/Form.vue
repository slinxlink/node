<template>
    <div class="endpoint">
        <div class="form-row">
            <span class="form-label">启用</span>
            <Toggle v-model="form.Enable" />
        </div>
        <div class="form-row half">
            <span class="form-label">标签</span>
            <Input v-model="form.Tag" placeholder="tag" :disabled="form.Tag === 'warp'" />
        </div>
        <div class="form-row half">
            <span class="form-label">协议</span>
            <Select v-model="form.Type" :options="[
                { label: 'WireGuard', value: 'wireguard' },
            ]" />
        </div>

        <template v-if="form.Type === 'wireguard'">
            <div class="form-row">
                <span class="form-label">本地 IPv4 地址</span>
                <Input v-model="form.Address[0]" placeholder="0.0.0.0/32" />
            </div>
            <div class="form-row">
                <span class="form-label">本地 IPv6 地址</span>
                <Input v-model="form.Address[1]" placeholder="::/128" />
            </div>
            <div class="form-row">
                <span class="form-label">公钥</span>
                <Input v-model="form.PublicKey" readonly />
            </div>
            <div class="form-row">
                <span class="form-label">私钥</span>
                <Input v-model="form.PrivateKey" readonly />
            </div>
            <div class="form-row">
                <span class="form-label"></span>
                <button class="action-btn" @click="genKeyPair">生成密钥对</button>
            </div>
            <div class="form-row quarter">
                <span class="form-label">MTU</span>
                <Input v-model="form.MTU" type="number" :min="0" placeholder="1408" />
            </div>
            <div class="form-row">
                <span class="form-label">对端地址</span>
                <Input :model-value="form.PeerAddress + ':' + form.PeerPort"
                    @update:model-value="val => {
                        const parts = String(val).split(':')
                        form.PeerAddress = parts[0]
                        form.PeerPort = parseInt(parts[1] ?? '0') || 0
                    }"
                    placeholder="engage.cloudflareclient.com:2408" />
            </div>
            <div class="form-row">
                <span class="form-label">Reserved</span>
                <Input v-model="reserved" placeholder="0,0,0" />
            </div>
            <div class="form-row">
                <span class="form-label">对端公钥</span>
                <Input v-model="form.PeerPublicKey" />
            </div>
            <div class="form-row">
                <span class="form-label">允许的IP</span>
                <Input v-model="form.AllowedIPs" placeholder="0.0.0.0/0,::/0" disabled />
            </div>
            <div class="form-row quarter">
                <span class="form-label">UDP 超时 (s)</span>
                <Input v-model="form.UDPTimeout" type="number" :min="0" placeholder="秒，0 表示默认" />
            </div>
            <div class="form-row quarter">
                <span class="form-label">Workers</span>
                <Input v-model="form.Workers" type="number" :min="0" placeholder="0 表示默认" />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'
import { generateWireguardKeyPair } from '@/api/generate'

const form = defineModel<any>({ default: () => ({
    Enable: true,
    Tag: '',
    Type: 'wireguard',
    Address: ['', ''],
    PublicKey: '',
    PrivateKey: '',
    MTU: 1408,
    Reserved: '',
    PeerAddress: '',
    PeerPort: 0,
    PeerPublicKey: '',
    AllowedIPs: '0.0.0.0/0,::/0',
    UDPTimeout: 0,
    Workers: 0,
}) })

const reserved = computed({
    get() {
        if (!form.value.Reserved) return ''
        try {
            return JSON.parse(form.value.Reserved).join(',')
        } catch {
            return form.value.Reserved
        }
    },
    set(val: string) {
        const nums = val.split(',').map(s => parseInt(s.trim())).filter(n => !isNaN(n))
        form.value.Reserved = JSON.stringify(nums)
    }
})

async function genKeyPair() {
    const res = await generateWireguardKeyPair()
    form.value.PrivateKey = res.private_key
    form.value.PublicKey = res.public_key
}
</script>

<style scoped>
.endpoint {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 20px;
}

:deep(.section) {
    border-radius: 0;

    .header {
        background-color: var(--color-bg);
    }
    .body {
        .content {
            display: flex;
            flex-direction: column;
            gap: 10px;
            padding: 20px;
        }
    }
}

.form-row {
    display: flex;
    align-items: center;
    gap: 16px;

    button {
        font-size: var(--font-size-sm);
    }

    &.half :deep(.input-wrap),
    &.half :deep(.select) {
        max-width: 50%;
    }

    &.quarter :deep(.input-wrap),
    &.quarter :deep(.select) {
        max-width: 25%;
    }
}

.form-label {
    width: 100px;
    flex-shrink: 0;
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
    text-align: right;
}

.action-btn {
    gap: 5px;
    padding: 8px 16px;
}

.pair {
    display: flex;
    gap: 2px;
    flex: 1;

    &:deep(.input-wrap) {
        &:first-child .input {
            border-radius: 999px 0 0 999px;
        }
        &:last-child .input {
            border-radius: 0 999px 999px 0;
        }
    }
}
</style>