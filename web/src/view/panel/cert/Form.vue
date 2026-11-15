<!-- view/panel/cert/Form.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>申请方式</template>
            <Select v-model:value="form.Mode" :items="modes" :disabled="mode === 'reapply'" />
        </Field>

        <template v-if="form.Mode !== 'manual'">
            <Field>
                <template #prefix>域名</template>
                <Input v-model:value="form.Domain" placeholder="example.com" :disabled="mode !== 'add'" autocomplete="off" />
            </Field>
            <Field>
                <template #prefix>ACME 账号</template>
                <Select v-model:value="form.Acme" :items="acmeItems" placeholder="选择 ACME 账号" :disabled="mode === 'reapply'" />
            </Field>
            <Field v-if="form.Mode === 'dns'">
                <template #prefix>DNS 账号</template>
                <Select v-model:value="form.Dns" :items="dnsItems" placeholder="选择 DNS 账号" :disabled="mode === 'reapply'" />
            </Field>
            <Field>
                <template #prefix>自动续签</template>
                <Toggle v-model:value="form.AutoRenew" :disabled="mode === 'reapply'" />
            </Field>
        </template>

        <template v-else>
            <Field>
                <template #prefix>公钥</template>
                <Textarea v-model:value="form.CertContent" placeholder="-----BEGIN CERTIFICATE-----" />
            </Field>
            <Field>
                <template #prefix>私钥</template>
                <Textarea v-model:value="form.KeyContent" placeholder="-----BEGIN PRIVATE KEY-----" />
            </Field>
        </template>

        <Task v-if="taskId" :task-id="taskId" @done="onDone" @error="onError" />
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Task from '@/component/module/Task.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Textarea from '@/component/widget/Textarea.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { applyCert, getCertContent, saveCert } from '@/api/cert'
import type { Acme, Cert, DnsAccount } from '@/api/cert'
import { toast } from '@/composable/toast'
import { display } from '@/view/panel/cert/display'

const props = defineProps<{
    mode: 'add' | 'edit' | 'reapply' // 申请、编辑、重新申请
    acmes: Acme[]
    dnsList: DnsAccount[]
    cert?: Cert
}>()

const emit = defineEmits<{
    changed: [] // 证书有变化(申请成功、上传、保存)，外面要刷新列表
}>()

const modes = [
    { label: 'DNS 验证', value: 'dns' },
    { label: 'HTTP 验证', value: 'http' },
    { label: '手动上传', value: 'manual' },
]

const form = ref({
    ID: undefined as number | undefined,
    Domain: '',
    Mode: 'dns',
    Acme: 0,
    Dns: 0,
    AutoRenew: true,
    CertContent: '',
    KeyContent: '',
})

const acmeItems = computed(() =>
    props.acmes.map((item) => ({
        label: `${item.Email} [${display.acmeProvider(item.Provider)}]`,
        value: item.ID,
    })),
)

const dnsItems = computed(() =>
    props.dnsList.map((item) => ({
        label: `${item.Name} [${display.dnsProvider(item.Provider)}]`,
        value: item.ID,
    })),
)

// 申请是后台任务：任务号一有值就显示实时输出；running 是还在跑，done 是已经成功
const taskId = ref('')
const running = ref(false)
const done = ref(false)

let finish: ((ok: boolean) => void) | null = null

function onDone() {
    running.value = false
    done.value = true
    emit('changed')
    toast.success('申请成功')
    finish?.(true)
}

function onError() {
    running.value = false
    emit('changed')
    finish?.(false)
}

// 发起申请，等任务结束(成功或失败)才返回
async function apply(id: number): Promise<boolean> {
    taskId.value = ''
    await nextTick()

    try {
        const res = await applyCert(id)

        if (!res.task_id) return true

        running.value = true
        done.value = false

        return await new Promise<boolean>((resolve) => {
            finish = resolve
            taskId.value = res.task_id as string
        })
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '申请失败')
        running.value = false

        return false
    }
}

onMounted(async () => {
    const cert = props.cert

    if (!cert) return

    form.value = {
        ID: cert.ID,
        Domain: cert.Domain,
        Mode: cert.Mode,
        Acme: cert.Acme,
        Dns: cert.Dns,
        AutoRenew: cert.AutoRenew,
        CertContent: '',
        KeyContent: '',
    }

    const content = await getCertContent(cert.ID)

    form.value.CertContent = content.cert
    form.value.KeyContent = content.key

    // 重新申请：打开就开始，不用再点确认
    if (props.mode === 'reapply') {
        apply(cert.ID)
    }
})

// 成功返回 true 关掉弹窗；返回 false 保持打开(失败，或者申请在跑、跑完了让人看输出)
async function submit(): Promise<boolean> {
    if (props.mode === 'reapply') return !running.value

    if (done.value) return true
    if (running.value) return false

    try {
        const saved = await saveCert({ ...form.value })

        // 编辑只保存
        if (props.mode === 'edit') {
            toast.success('保存成功')
            emit('changed')

            return true
        }

        form.value.ID = saved.ID

        // 手动上传不用申请，保存就是上传
        if (form.value.Mode === 'manual') {
            toast.success('证书上传成功')
            emit('changed')

            return true
        }

        await apply(saved.ID)

        return false
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }
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

    :deep(.prefix) {
        flex-basis: 100px;
        align-self: flex-start;
        min-height: 40px;
    }

    :deep(textarea) {
        height: 140px;
    }
}
</style>
