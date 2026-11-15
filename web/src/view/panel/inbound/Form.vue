<!-- view/panel/inbound/Form.vue -->
<template>
    <div class="form">
        <Tab v-model:value="tab" :items="tabs" />

        <template v-if="tab === 'basic'">
            <Field>
                <template #prefix>启用</template>
                <Toggle v-model:value="form.Enable" />
            </Field>
            <Field>
                <template #prefix>备注</template>
                <Input v-model:value="form.Name" autocomplete="off" />
            </Field>
            <Field>
                <template #prefix>协议</template>
                <Select v-model:value="form.Protocol" :items="protocols" />
            </Field>
            <Field>
                <template #prefix>端口</template>
                <Input
                    v-model:value="form.Port"
                    type="number"
                    :min="1"
                    :max="65535"
                    placeholder="1 - 65535"
                />
            </Field>
            <Field>
                <template #prefix>
                    <span>流量上限</span>
                    <Tip msg="GB[br]0 = 不限" />
                </template>
                <Input v-model:value="form.TrafficLimit" type="number" :min="0" />
            </Field>
            <ResetCycle v-model:cycle="form.ResetCycle" v-model:day="form.ResetDay" />
        </template>

        <template v-if="tab === 'transport'">
            <template v-if="isTcp">
                <Field>
                    <template #prefix>传输</template>
                    <Select v-model:value="form.Transport" :items="transports" />
                </Field>
                <template v-if="form.Transport === 'websocket'">
                    <Field>
                        <template #prefix>路径</template>
                        <Input v-model:value="form.WsPath" placeholder="/path" />
                    </Field>
                    <Field>
                        <template #prefix>Host</template>
                        <Input v-model:value="form.WsHost" placeholder="example.com" />
                    </Field>
                    <Field>
                        <template #prefix>心跳周期</template>
                        <Input
                            v-model:value="form.WsPingInterval"
                            type="number"
                            :min="0"
                            placeholder="秒，0 表示禁用"
                        />
                    </Field>
                </template>
            </template>

            <template v-if="form.Protocol === 'hysteria'">
                <Field>
                    <template #prefix>版本</template>
                    <Select :value="'2'" :items="[{ label: '2', value: '2' }]" disabled />
                </Field>
                <Field>
                    <template #prefix>UDP 超时 (s)</template>
                    <Input
                        v-model:value="form.UDPTimeout"
                        type="number"
                        :min="0"
                        placeholder="秒，0 表示默认"
                    />
                </Field>
                <Field>
                    <template #prefix>伪装</template>
                    <Toggle v-model:value="form.MasqueradeEnabled" />
                </Field>
                <template v-if="form.MasqueradeEnabled">
                    <Field>
                        <template #prefix>类型</template>
                        <Select v-model:value="form.MasqueradeType" :items="masqueradeTypes" />
                    </Field>
                    <template v-if="form.MasqueradeType === 'proxy'">
                        <Field>
                            <template #prefix>Upstream URL</template>
                            <Input
                                v-model:value="form.MasqueradeURL"
                                placeholder="https://www.example.com"
                            />
                        </Field>
                        <Field>
                            <template #prefix>重写 Host</template>
                            <Toggle v-model:value="form.RewriteHost" />
                        </Field>
                        <Field>
                            <template #prefix>跳过 TLS</template>
                            <Toggle v-model:value="form.IgnoreTLSVerify" />
                        </Field>
                    </template>
                    <Field v-if="form.MasqueradeType === 'file'">
                        <template #prefix>目录路径</template>
                        <Input v-model:value="form.MasqueradePath" placeholder="/var/www/html" />
                    </Field>
                    <template v-if="form.MasqueradeType === 'string'">
                        <Field>
                            <template #prefix>状态码</template>
                            <Input
                                v-model:value="form.MasqueradeCode"
                                type="number"
                                :min="100"
                                :max="599"
                                placeholder="200"
                            />
                        </Field>
                        <Field>
                            <template #prefix>内容</template>
                            <!-- 回车在这里是换行，不能让弹窗当成确认 -->
                            <Textarea
                                v-model:value="form.MasqueradeBody"
                                placeholder="<html>...</html>"
                                @keydown.enter.stop
                            />
                        </Field>
                    </template>
                </template>
                <Field>
                    <template #prefix>混淆</template>
                    <Select v-model:value="form.ObfsType" :items="obfsTypes" />
                </Field>
                <template v-if="form.ObfsType">
                    <Field>
                        <template #prefix>混淆密码</template>
                        <Input v-model:value="form.ObfsPassword">
                            <template #trailing>
                                <Button variant="text" :tabindex="-1" @click="genObfsPassword">
                                    <Icon name="casino" />
                                </Button>
                            </template>
                        </Input>
                    </Field>
                    <Field v-if="form.ObfsType === 'gecko'">
                        <template #prefix>包大小</template>
                        <Input
                            v-model:value="form.ObfsMinPacketSize"
                            type="number"
                            :min="1"
                            placeholder="最小"
                        />
                        <Input
                            v-model:value="form.ObfsMaxPacketSize"
                            type="number"
                            :min="1"
                            placeholder="最大"
                        />
                    </Field>
                </template>
            </template>

            <template v-if="form.Protocol === 'tuic'">
                <Field>
                    <template #prefix>版本</template>
                    <Select :value="'5'" :items="[{ label: '5', value: '5' }]" disabled />
                </Field>
                <Field>
                    <template #prefix>拥塞控制</template>
                    <Select v-model:value="form.TuicCongestionControl" :items="congestions" />
                </Field>
                <Field>
                    <template #prefix>认证超时 (s)</template>
                    <Input
                        v-model:value="form.TuicAuthTimeout"
                        type="number"
                        :min="1"
                        placeholder="3"
                    />
                </Field>
                <Field>
                    <template #prefix>零 RTT 握手</template>
                    <Toggle v-model:value="form.TuicZeroRTT" />
                </Field>
                <Field>
                    <template #prefix>心跳间隔 (s)</template>
                    <Input
                        v-model:value="form.TuicHeartbeat"
                        type="number"
                        :min="1"
                        placeholder="10"
                    />
                </Field>
                <Field>
                    <template #prefix>UDP 转发</template>
                    <Select v-model:value="form.TuicUDPRelayMode" :items="relayModes" />
                </Field>
            </template>

            <template v-if="isUdp">
                <Field>
                    <template #prefix>端口跳跃</template>
                    <Toggle v-model:value="form.HopEnabled" />
                </Field>
                <template v-if="form.HopEnabled">
                    <Field>
                        <template #prefix>端口范围</template>
                        <Input v-model:value="form.HopPort" placeholder="10000-11000" />
                    </Field>
                    <Field>
                        <template #prefix>跳跃间隔 (s)</template>
                        <Input v-model:value="form.HopInterval" placeholder="5-10" />
                    </Field>
                </template>
            </template>

            <template v-if="form.Protocol === 'anytls'">
                <Field>
                    <template #prefix>检测间隔 (s)</template>
                    <Input
                        v-model:value="form.AnyTLSIdleSessionCheckInterval"
                        type="number"
                        :min="1"
                        placeholder="30"
                    />
                </Field>
                <Field>
                    <template #prefix>空闲超时 (s)</template>
                    <Input
                        v-model:value="form.AnyTLSIdleSessionTimeout"
                        type="number"
                        :min="1"
                        placeholder="30"
                    />
                </Field>
                <Field>
                    <template #prefix>最小空闲会话</template>
                    <Input
                        v-model:value="form.AnyTLSMinIdleSession"
                        type="number"
                        :min="1"
                        placeholder="1"
                    />
                </Field>
            </template>
        </template>

        <template v-if="tab === 'security'">
            <Field>
                <template #prefix>安全</template>
                <Choice
                    v-model:value="form.TLSType"
                    :items="securities"
                    :disabled="lockedSecurity"
                />
            </Field>

            <template v-if="form.TLSType === 'tls'">
                <Field>
                    <template #prefix>SNI</template>
                    <Input v-model:value="form.ServerName" placeholder="example.com" />
                </Field>
                <Field>
                    <template #prefix>加密套件</template>
                    <Select v-model:value="form.CipherSuites" :items="cipherSuites" />
                </Field>
                <Field>
                    <template #prefix>uTLS</template>
                    <Select v-model:value="form.UTLS" :items="utls" />
                </Field>
                <Field>
                    <template #prefix>ALPN</template>
                    <TagSelect
                        v-model:value="form.ALPN"
                        :items="alpnItems"
                        :disabled="form.Protocol === 'hysteria'"
                    />
                </Field>
                <Field>
                    <template #prefix>版本范围</template>
                    <Select v-model:value="form.TLSMinVersion" :items="tlsVersions" />
                    <Select v-model:value="form.TLSMaxVersion" :items="tlsVersions" />
                </Field>
                <Field>
                    <template #prefix>跳过验证</template>
                    <Toggle v-model:value="form.Insecure" />
                </Field>
                <Field>
                    <template #prefix>证书</template>
                    <Select v-model:value="form.Cert" :items="certItems" placeholder="选择证书" />
                </Field>
                <Field>
                    <template #prefix>ECH</template>
                    <Toggle v-model:value="form.ECHEnabled" />
                </Field>
                <template v-if="form.ECHEnabled">
                    <Field>
                        <template #prefix>私钥</template>
                        <Textarea v-model:value="form.ECHKey" readonly>
                    <template #trailing>
                        <Button variant="text" :tabindex="-1" @click="genECH">
                            <Icon name="casino" />
                        </Button>
                    </template>
                </Textarea>
                    </Field>
                    <Field>
                        <template #prefix>配置</template>
                        <Textarea v-model:value="form.ECHConfig" readonly>
                    <template #trailing>
                        <Button variant="text" :tabindex="-1" @click="genECH">
                            <Icon name="casino" />
                        </Button>
                    </template>
                </Textarea>
                    </Field>
                </template>
            </template>

            <template v-if="form.TLSType === 'reality'">
                <Field>
                    <template #prefix>SNI</template>
                    <Input v-model:value="form.RealityServerName">
                        <template #trailing>
                            <Button variant="text" :tabindex="-1" @click="genTarget">
                                <Icon name="casino" />
                            </Button>
                        </template>
                    </Input>
                </Field>
                <Field>
                    <template #prefix>目标</template>
                    <Input v-model:value="target" placeholder="www.example.com:443">
                        <template #trailing>
                            <Button variant="text" :tabindex="-1" @click="genTarget">
                                <Icon name="casino" />
                            </Button>
                        </template>
                    </Input>
                </Field>
                <Field>
                    <template #prefix>uTLS</template>
                    <Select v-model:value="form.UTLS" :items="utls" />
                </Field>
                <Field>
                    <template #prefix>最大时间差</template>
                    <Input
                        v-model:value="form.RealityMaxTimeDiff"
                        type="number"
                        :min="0"
                        placeholder="毫秒，0 不限制"
                    />
                </Field>
                <Field>
                    <template #prefix>短 ID</template>
                    <TagSelect v-model:value="form.RealityShortIDs" :items="shortIDItems">
                        <template #trailing>
                            <Button variant="text" :tabindex="-1" @click="genShortIDs">
                                <Icon name="casino" />
                            </Button>
                        </template>
                    </TagSelect>
                </Field>
                <Field>
                    <template #prefix>公钥</template>
                    <Input v-model:value="form.RealityPublicKey" readonly>
                        <template #trailing>
                            <Button variant="text" :tabindex="-1" @click="genKeyPair">
                                <Icon name="casino" />
                            </Button>
                        </template>
                    </Input>
                </Field>
                <Field>
                    <template #prefix>私钥</template>
                    <Input v-model:value="form.RealityPrivateKey" readonly>
                        <template #trailing>
                            <Button variant="text" :tabindex="-1" @click="genKeyPair">
                                <Icon name="casino" />
                            </Button>
                        </template>
                    </Input>
                </Field>
                <Field>
                    <template #prefix>Flow</template>
                    <Select v-model:value="form.Flow" :items="flows" />
                </Field>
            </template>
        </template>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import ResetCycle from '@/component/module/ResetCycle.vue'
import Button from '@/component/widget/Button.vue'
import Choice from '@/component/widget/Choice.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Tab from '@/component/widget/Tab.vue'
import type { TabItem } from '@/component/widget/Tab.vue'
import TagSelect from '@/component/widget/TagSelect.vue'
import Textarea from '@/component/widget/Textarea.vue'
import Tip from '@/component/widget/Tip.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { getCert } from '@/api/cert'
import type { Cert } from '@/api/cert'
import {
    generateECHKeyPair,
    generatePassword,
    generatePort,
    generateRealityKeyPair,
    generateRealityTarget,
    generateShortIDs,
} from '@/api/generate'
import { saveInbound } from '@/api/inbound'
import type { Inbound, InboundData } from '@/api/inbound'
import { toast } from '@/composable/toast'
import {
    alpns,
    cipherSuites,
    congestions,
    flows,
    masqueradeTypes,
    obfsTypes,
    protocols,
    relayModes,
    securities,
    tlsVersions,
    transports,
    utls,
} from '@/view/panel/inbound/options'

const props = defineProps<{
    inbound?: Inbound // 不传是添加，传了是编辑
}>()

const tabs: [TabItem, ...TabItem[]] = [
    { label: '基础', value: 'basic' },
    { label: '传输', value: 'transport' },
    { label: '安全', value: 'security' },
]

const tab = ref('basic')

const GB = 1024 ** 3

// 输入框给的都是字符串，数字项在表单里按字符串存，提交时再转成数字
const form = ref({
    Enable: true,
    Name: '',
    Protocol: 'vless',
    Port: '0',
    TrafficLimit: '0',
    ResetCycle: '',
    ResetDay: '1',

    Transport: 'raw',
    WsPath: '/',
    WsHost: '',
    WsPingInterval: '0',

    UDPTimeout: '60',
    MasqueradeEnabled: false,
    MasqueradeType: 'default',
    MasqueradeURL: '',
    RewriteHost: false,
    IgnoreTLSVerify: false,
    MasqueradePath: '',
    MasqueradeCode: '200',
    MasqueradeBody: '',
    ObfsType: '',
    ObfsPassword: '',
    ObfsMinPacketSize: '512',
    ObfsMaxPacketSize: '1200',

    TuicCongestionControl: 'bbr',
    TuicAuthTimeout: '3',
    TuicZeroRTT: false,
    TuicHeartbeat: '10',
    TuicUDPRelayMode: 'native',

    HopEnabled: false,
    HopPort: '10000-11000',
    HopInterval: '5-10',

    AnyTLSIdleSessionCheckInterval: '30',
    AnyTLSIdleSessionTimeout: '30',
    AnyTLSMinIdleSession: '1',

    TLSType: 'none',
    ServerName: '',
    CipherSuites: '',
    ALPN: ['h2', 'http/1.1'],
    TLSMinVersion: '1.2',
    TLSMaxVersion: '1.3',
    UTLS: 'chrome',
    Insecure: false,
    Cert: null as number | null,
    ECHEnabled: false,
    ECHKey: '',
    ECHConfig: '',

    RealityServerName: '',
    RealityServer: '',
    RealityServerPort: '443',
    RealityMaxTimeDiff: '0',
    RealityShortIDs: [] as string[],
    RealityPublicKey: '',
    RealityPrivateKey: '',
    Flow: 'xtls-rprx-vision',
})

const certs = ref<Cert[]>([])
const shortIDs = ref<string[]>([])

const certItems = computed(() => certs.value.map((c) => ({ label: c.Domain, value: c.ID })))
const shortIDItems = computed(() => shortIDs.value.map((id) => ({ label: id, value: id })))

// Hysteria 只能用 h3
const alpnItems = computed(() => (form.value.Protocol === 'hysteria' ? alpns.slice(0, 1) : alpns))

const isTcp = computed(() => ['vless', 'vmess', 'trojan'].includes(form.value.Protocol))
const isUdp = computed(() => ['hysteria', 'tuic'].includes(form.value.Protocol))

// 这几种协议只能用 TLS，不让改
const lockedSecurity = computed(() => ['hysteria', 'trojan', 'tuic'].includes(form.value.Protocol))

// 伪装目标在库里是域名和端口两项，表单里合成一个输入框
const target = computed({
    get: () => `${form.value.RealityServer}:${form.value.RealityServerPort}`,
    set: (value: string) => {
        const [host, port] = value.split(':')
        form.value.RealityServer = host ?? ''
        form.value.RealityServerPort = String(parseInt(port ?? '443') || 443)
    },
})

watch(
    () => form.value.Protocol,
    (value) => {
        if (!['hysteria', 'trojan', 'tuic', 'anytls'].includes(value)) return

        form.value.TLSType = 'tls'

        if (value === 'hysteria' || value === 'tuic') form.value.ALPN = ['h3']
    },
)

function fill(ib: Inbound) {
    const { ID, CreatedAt, Upload, Download, ResetAt, TrafficLimit, ALPN, Certs, RealityShortIDs, ...rest } = ib
    void ID
    void CreatedAt
    void Upload
    void Download
    void ResetAt

    const ids: string[] = RealityShortIDs ? JSON.parse(RealityShortIDs) : []
    const list: number[] = Certs ? JSON.parse(Certs) : []

    // 数字项转成字符串，其余原样
    const text = Object.fromEntries(
        Object.entries(rest).map(([key, value]) => [
            key,
            typeof value === 'number' ? String(value) : value,
        ]),
    )

    Object.assign(form.value, text, {
        TrafficLimit: String(Math.round((TrafficLimit / GB) * 100) / 100),
        ALPN: ALPN ? ALPN.split(',') : [],
        Cert: list[0] ?? null,
        RealityShortIDs: ids,
    })

    shortIDs.value = ids
}

// 添加时只需要一个随机端口和默认证书；Reality 的东西等切到 Reality 再生成
async function fresh() {
    const [list, port] = await Promise.all([getCert(), generatePort()])

    certs.value = list
    form.value.Port = String(port.port)
    form.value.Cert = list[0]?.ID ?? 0
}

// 切到 Reality 时，目标、密钥对、短 ID 哪个是空的就生成哪个，各自回来各自填；已有内容的(比如编辑)不动
function ensureReality() {
    const f = form.value

    if (!f.RealityServer) genTarget()
    if (!f.RealityPrivateKey) genKeyPair()
    if (!f.RealityShortIDs.length) genShortIDs()
}

watch(
    () => form.value.TLSType,
    (value) => {
        if (value === 'reality') ensureReality()
    },
)

// 目标是从一个不大的名单里随机挑的，挑到和现在一样的就再摇，免得点了没变化；最多摇 10 次
async function genTarget() {
    for (let i = 0; i < 10; i++) {
        const res = await generateRealityTarget()

        if (res.domain === form.value.RealityServerName || res.domain === form.value.RealityServer) continue

        form.value.RealityServerName = res.domain
        form.value.RealityServer = res.domain
        form.value.RealityServerPort = '443'

        return
    }
}

async function genKeyPair() {
    const res = await generateRealityKeyPair()
    form.value.RealityPrivateKey = res.private_key
    form.value.RealityPublicKey = res.public_key
}

async function genShortIDs() {
    const res = await generateShortIDs()
    form.value.RealityShortIDs = res.short_ids
    shortIDs.value = res.short_ids
}

async function genObfsPassword() {
    form.value.ObfsPassword = (await generatePassword()).password
}

async function genECH() {
    const cert = certs.value.find((c) => c.ID === form.value.Cert)

    if (!cert) {
        toast.warning('请先选择证书')

        return
    }

    const res = await generateECHKeyPair(cert.Domain)
    form.value.ECHKey = res.ech_key
    form.value.ECHConfig = res.ech_config
}

onMounted(async () => {
    if (props.inbound) {
        certs.value = await getCert()
        fill(props.inbound)
    } else {
        await fresh()
    }
})

// 成功返回 true，失败返回 false（弹窗保持打开）；各项的范围检查在后端
async function submit(): Promise<boolean> {
    const { Cert, ALPN, RealityShortIDs, ...rest } = form.value

    const data = {
        ...rest,
        Port: Number(rest.Port),
        TrafficLimit: Math.round(Number(rest.TrafficLimit) * GB),
        ResetDay: Number(rest.ResetDay),
        WsPingInterval: Number(rest.WsPingInterval),
        UDPTimeout: Number(rest.UDPTimeout),
        MasqueradeCode: Number(rest.MasqueradeCode),
        ObfsMinPacketSize: Number(rest.ObfsMinPacketSize),
        ObfsMaxPacketSize: Number(rest.ObfsMaxPacketSize),
        TuicAuthTimeout: Number(rest.TuicAuthTimeout),
        TuicHeartbeat: Number(rest.TuicHeartbeat),
        AnyTLSIdleSessionCheckInterval: Number(rest.AnyTLSIdleSessionCheckInterval),
        AnyTLSIdleSessionTimeout: Number(rest.AnyTLSIdleSessionTimeout),
        AnyTLSMinIdleSession: Number(rest.AnyTLSMinIdleSession),
        RealityServerPort: Number(rest.RealityServerPort),
        RealityMaxTimeDiff: Number(rest.RealityMaxTimeDiff),
        ALPN: ALPN.join(','),
        Certs: JSON.stringify([Cert ?? 0]),
        RealityShortIDs: JSON.stringify(RealityShortIDs),
        ID: props.inbound?.ID,
    } satisfies InboundData

    try {
        await saveInbound(data)
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
    margin-top: -15px; /* 抵掉弹窗里标题和内容之间的 15px 间距：最上面是标签，标签自己已经留了 10px */
    text-align: left;

    /* Tab 默认上下留 20px 是给大标题行用的，这里紧凑一点 */
    :deep(.tab) {
        padding: 10px 0;
    }

    /* 标题列拉宽，长一点的名字也放得下 */
    :deep(.prefix) {
        flex-basis: 100px;
        justify-content: space-between;
    }

    /* 行变高(多行输入框、换行的标签选择)时，标题和右边的按钮固定在第一个 40px 里，不跟着整行居中 */
    :deep(.prefix),
    :deep(.suffix) {
        align-self: flex-start;
        min-height: 40px;
    }

    /* 一行 40px、输入框 34px，上下各剩 3px；变高后补上这 3px，行间距才和别的行一样 */
    :deep(.tag-select),
    :deep(.textarea-wrap) {
        margin: 3px 0;
    }

    &:deep(textarea) {
        height: 122px; /* 80px 再加两行(每行 21px) */
    }
}
</style>
