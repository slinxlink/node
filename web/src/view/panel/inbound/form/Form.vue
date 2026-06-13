<template>
    <div class="inbound-form">
        <!-- 基础配置 -->
        <Section title="基础配置">
            <div class="form-row">
                <span class="form-label">启用</span>
                <Toggle v-model="form.Enable" />
            </div>
            <div class="form-row">
                <span class="form-label">备注</span>
                <Input v-model="form.Name" />
            </div>
            <div class="form-row">
                <span class="form-label">协议</span>
                <Select v-model="form.Protocol" :options="[
                    { label: 'vless', value: 'vless' },
                    { label: 'vmess', value: 'vmess' },
                    { label: 'hysteria', value: 'hysteria' },
                ]" />
            </div>
            <div class="form-row quarter">
                <span class="form-label">端口</span>
                <Input v-model="form.Port" type="number" :min="1" :max="65535" placeholder="1 - 65535" />
            </div>
        </Section>

        <!-- 传输 -->
        <Section title="传输">
            <!-- VLESS / VMess -->
            <template v-if="form.Protocol === 'vless' || form.Protocol === 'vmess'">
                <div class="form-row half">
                    <span class="form-label">传输</span>
                    <Select v-model="form.Transport" :options="[
                        { label: 'RAW', value: 'raw' },
                        { label: 'WebSocket', value: 'websocket' },
                    ]" />
                </div>
                <template v-if="form.Transport === 'websocket'">
                    <div class="form-row">
                        <span class="form-label">路径</span>
                        <Input v-model="form.WsPath" placeholder="/path" />
                    </div>
                    <div class="form-row">
                        <span class="form-label">Host</span>
                        <Input v-model="form.WsHost" placeholder="example.com" />
                    </div>
                    <div class="form-row quarter">
                        <span class="form-label">心跳周期</span>
                        <Input v-model="form.WsPingInterval" type="number" :min="0" placeholder="秒，0 表示禁用" />
                    </div>
                </template>
            </template>

            <!-- Hysteria -->
            <template v-if="form.Protocol === 'hysteria'">
                <div class="form-row quarter">
                    <span class="form-label">版本</span>
                    <Select :model-value="'2'" :options="[{ label: '2', value: '2' }]" disabled />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">UDP 超时 (s)</span>
                    <Input v-model="form.UDPTimeout" type="number" :min="0" placeholder="秒，0 表示默认" />
                </div>
                <div class="form-row">
                    <span class="form-label">伪装</span>
                    <Toggle v-model="form.MasqueradeEnabled" />
                </div>
                <template v-if="form.MasqueradeEnabled">
                    <div class="form-row half">
                        <span class="form-label">类型</span>
                        <Select v-model="form.MasqueradeType" :options="[
                            { label: '默认 (404)', value: 'default' },
                            { label: '反向代理', value: 'proxy' },
                            { label: '静态文件', value: 'file' },
                            { label: '自定义内容', value: 'string' },
                        ]" />
                    </div>
                    <template v-if="form.MasqueradeType === 'proxy'">
                        <div class="form-row">
                            <span class="form-label">Upstream URL</span>
                            <Input v-model="form.MasqueradeURL" placeholder="https://www.example.com" />
                        </div>
                        <div class="form-row">
                            <span class="form-label">重写 Host</span>
                            <Toggle v-model="form.RewriteHost" />
                        </div>
                        <div class="form-row">
                            <span class="form-label">跳过 TLS</span>
                            <Toggle v-model="form.IgnoreTLSVerify" />
                        </div>
                    </template>
                    <template v-if="form.MasqueradeType === 'file'">
                        <div class="form-row">
                            <span class="form-label">目录路径</span>
                            <Input v-model="form.MasqueradePath" placeholder="/var/www/html" />
                        </div>
                    </template>
                    <template v-if="form.MasqueradeType === 'string'">
                        <div class="form-row quarter">
                            <span class="form-label">状态码</span>
                            <Input v-model="form.MasqueradeCode" type="number" :min="100" :max="599" placeholder="200" />
                        </div>
                        <div class="form-row">
                            <span class="form-label">内容</span>
                            <Text v-model="form.MasqueradeBody" placeholder="<html>...</html>" />
                        </div>
                    </template>
                </template>
            </template>
        </Section>

        <!-- 安全 -->
        <Section title="安全">
            <div class="form-row">
                <span class="form-label">安全</span>
                    <RadioGroup
                        v-model="form.TLSType"
                        :options="[
                            { label: '无', value: 'none' },
                            { label: 'TLS', value: 'tls' },
                            { label: 'Reality', value: 'reality' },
                        ]"
                        :disabled="form.Protocol === 'hysteria'"
                    />
            </div>
            <!-- tls 专用 -->
            <template v-if="form.TLSType === 'tls'">
                <div class="form-row">
                    <span class="form-label">SNI</span>
                    <Input v-model="form.ServerName" placeholder="example.com" />
                </div>
                <div class="form-row">
                    <span class="form-label">加密套件</span>
                    <Select v-model="form.CipherSuites" :options="[
                        { label: '自动', value: '' },
                        { label: 'ECDHE_ECDSA_AES_128_CBC', value: 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA' },
                        { label: 'ECDHE_ECDSA_AES_256_CBC', value: 'TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA' },
                        { label: 'ECDHE_RSA_AES_128_CBC', value: 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA' },
                        { label: 'ECDHE_RSA_AES_256_CBC', value: 'TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA' },
                        { label: 'ECDHE_ECDSA_AES_128_GCM', value: 'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256' },
                        { label: 'ECDHE_ECDSA_AES_256_GCM', value: 'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384' },
                        { label: 'ECDHE_RSA_AES_128_GCM', value: 'TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256' },
                        { label: 'ECDHE_RSA_AES_256_GCM', value: 'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384' },
                        { label: 'ECDHE_ECDSA_CHACHA20_POLY1305', value: 'TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256' },
                        { label: 'ECDHE_RSA_CHACHA20_POLY1305', value: 'TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256' },
                    ]" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">uTLS</span>
                    <Select v-model="form.UTLS" :options="utlsOptions" />
                </div>
                <div class="form-row">
                    <span class="form-label">ALPN</span>
                    <MultiSelect
                        v-if="form.Protocol !== 'hysteria'"
                        v-model="form.ALPN"
                        :options="[
                            { label: 'h3', value: 'h3' },
                            { label: 'h2', value: 'h2' },
                            { label: 'http/1.1', value: 'http/1.1' },
                        ]"
                    />
                    <MultiSelect
                        v-else
                        v-model="form.ALPN"
                        :options="[{ label: 'h3', value: 'h3' }]"
                        disabled
                    />
                </div>
                <div class="form-row">
                    <span class="form-label">最小/最大版本</span>
                    <div class="select-pair">
                        <Select v-model="form.TLSMinVersion" :options="tlsVersions" />
                        <Select v-model="form.TLSMaxVersion" :options="tlsVersions" />
                    </div>
                </div>
                <div class="form-row">
                    <span class="form-label">跳过验证</span>
                    <Toggle v-model="form.Insecure" />
                </div>
                <div class="form-row">
                    <span class="form-label">证书</span>
                    <Select :model-value="form.Certs[0]" @update:model-value="val => form.Certs[0] = val" :options="certOptions" placeholder="选择证书" />
                </div>
            </template>


            <template v-if="form.TLSType === 'reality'">
                <div class="form-row">
                    <span class="form-label">SNI</span>
                    <RefreshBtn @click="genTarget">
                        <Input v-model="form.RealityServerName" />
                    </RefreshBtn>
                </div>
                <div class="form-row">
                    <span class="form-label">目标</span>
                    <RefreshBtn @click="genTarget">
                        <Input :model-value="form.RealityServer + ':' + form.RealityServerPort"
                            @update:model-value="val => {
                                const parts = String(val).split(':')
                                form.RealityServer = parts[0]
                                form.RealityServerPort = parseInt(parts[1] ?? '443') || 443
                            }"
                            placeholder="www.example.com:443" />
                    </RefreshBtn>
                </div>
                <div class="form-row quarter">
                    <span class="form-label">uTLS</span>
                    <Select v-model="form.UTLS" :options="utlsOptions" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">最大时间差</span>
                    <Input v-model="form.RealityMaxTimeDiff" type="number" :min="0" placeholder="毫秒，0 不限制" />
                </div>
                <div class="form-row">
                    <span class="form-label">短 ID</span>
                    <RefreshBtn @click="genShortIDs">
                        <MultiSelect
                            v-model="form.RealityShortIDs"
                            :options="shortIDOptions"
                        />
                    </RefreshBtn>
                </div>
                <div class="form-row">
                    <span class="form-label">公钥</span>
                    <Input v-model="form.RealityPublicKey" readonly />
                </div>
                <div class="form-row">
                    <span class="form-label">私钥</span>
                    <Input v-model="form.RealityPrivateKey" readonly />
                </div>
                <div class="form-row">
                    <span class="form-label"></span>
                    <button class="action-btn" @click="genKeyPair">生成密钥对</button>
                </div>
                <div class="form-row half">
                    <span class="form-label">Flow</span>
                    <Select v-model="form.Flow" :options="[
                        { label: '无', value: '' },
                        { label: 'xtls-rprx-vision', value: 'xtls-rprx-vision' },
                    ]" />
                </div>
            </template>
        </Section>
    </div>
</template>

<script setup lang="ts">
import Section from '@/component/ui/Section.vue'
import Toggle from '@/component/ui/Toggle.vue'
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'
import Text from '@/component/ui/Text.vue'
import RadioGroup from '@/component/ui/Radio.vue'
import MultiSelect from '@/component/ui/MultiSelect.vue'
import RefreshBtn from '@/component/ui/RefreshBtn.vue'
import { generateRealityTarget, generateRealityKeyPair, generateShortIDs } from '@/api/generate'
import { getCert } from '@/api/cert'

const form = defineModel<any>({ default: () => ({
    Enable: true,
    Name: '',
    Protocol: 'vless',
    Port: 0,
}) })

const tlsVersions = [
    { label: '1.0', value: '1.0' },
    { label: '1.1', value: '1.1' },
    { label: '1.2', value: '1.2' },
    { label: '1.3', value: '1.3' },
]

const utlsOptions = [
    { label: 'chrome', value: 'chrome' },
    { label: 'firefox', value: 'firefox' },
    { label: 'safari', value: 'safari' },
    { label: 'ios', value: 'ios' },
    { label: 'android', value: 'android' },
    { label: 'edge', value: 'edge' },
    { label: '360', value: '360' },
    { label: 'qq', value: 'qq' },
    { label: 'random', value: 'random' },
    { label: 'randomized', value: 'randomized' },
]

const certOptions = ref<any[]>([])
const shortIDOptions = ref<any[]>([])

onMounted(async () => {
    const certs = await getCert()
    certOptions.value = certs.map((c: any) => ({ label: c.Domain, value: c.ID }))
})

watch(() => form.value.Protocol, (val) => {
    if (val === 'hysteria') {
        form.value.TLSType = 'tls'
        form.value.ALPN = ['h3']
    }
}, { immediate: true })

watch(() => form.value.RealityShortIDs, (val) => {
    if (val?.length) {
        shortIDOptions.value = val.map((id: string) => ({ label: id, value: id }))
    }
}, { immediate: true })

async function genTarget() {
    const res = await generateRealityTarget()
    form.value.RealityServerName = res.domain
    form.value.RealityServer = res.domain
    form.value.RealityServerPort = 443
}

async function genKeyPair() {
    const res = await generateRealityKeyPair()
    form.value.RealityPrivateKey = res.private_key
    form.value.RealityPublicKey = res.public_key
}

async function genShortIDs() {
    const res = await generateShortIDs()
    form.value.RealityShortIDs = res.short_ids
    shortIDOptions.value = res.short_ids.map((id: string) => ({ label: id, value: id }))
}
</script>

<style scoped>
.inbound-form {
    display: flex;
    flex-direction: column;
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

.select-pair {
    display: flex;
    gap: 2px;
    flex: 1;

    &:deep(.select) {
        &:first-child .btn {
            border-radius: 999px 0 0 999px;
        }
        &:last-child .btn {
            border-radius: 0 999px 999px 0;
        }
    }
}
</style>