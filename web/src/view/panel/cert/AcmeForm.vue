<!-- view/panel/cert/AcmeForm.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>服务商</template>
            <Select v-model:value="form.Provider" :items="providers" />
        </Field>
        <Field>
            <template #prefix>邮箱</template>
            <Input v-model:value="form.Email" placeholder="example@email.com" autocomplete="off" />
        </Field>
        <template v-if="form.Provider === 'zerossl'">
            <Field>
                <template #prefix>EAB Kid</template>
                <Input v-model:value="form.EabKid" placeholder="EAB Kid" autocomplete="off" />
            </Field>
            <Field>
                <template #prefix>EAB Hmac</template>
                <Input v-model:value="form.EabHmac" placeholder="EAB Hmac Key" autocomplete="off" />
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
import { saveAcme } from '@/api/cert'
import type { Acme } from '@/api/cert'
import { toast } from '@/composable/toast'

const props = defineProps<{
    acme?: Acme // 不传是添加，传了是编辑
}>()

const emit = defineEmits<{
    changed: []
}>()

const providers = [
    { label: "Let's Encrypt", value: 'letsencrypt' },
    { label: 'ZeroSSL', value: 'zerossl' },
]

const form = ref({
    ID: undefined as number | undefined,
    PrivateKey: undefined as string | undefined,
    Provider: 'letsencrypt',
    Email: '',
    EabKid: '',
    EabHmac: '',
})

// 新账号要向服务商注册，是后台任务：任务号有值就显示实时输出；running 是还在跑，done 是已经成功
const taskId = ref('')
const running = ref(false)
const done = ref(false)

let finish: ((ok: boolean) => void) | null = null

function onDone() {
    running.value = false
    done.value = true
    emit('changed')
    toast.success('账号已注册')
    finish?.(true)
}

function onError() {
    running.value = false
    emit('changed')
    finish?.(false)
}

onMounted(() => {
    const acme = props.acme

    if (!acme) return

    form.value = {
        ID: acme.ID,
        PrivateKey: acme.PrivateKey,
        Provider: acme.Provider,
        Email: acme.Email,
        EabKid: acme.EabKid,
        EabHmac: acme.EabHmac,
    }
})

// 成功返回 true 关掉弹窗；返回 false 保持打开(失败，或者注册跑完了让人看输出)
async function submit(): Promise<boolean> {
    if (done.value) return true
    if (running.value) return false

    try {
        const res = await saveAcme({ ...form.value } as Parameters<typeof saveAcme>[0])

        // 没有任务号：直接保存完了
        if (!res.task_id) {
            toast.success('保存成功')
            emit('changed')

            return true
        }

        running.value = true
        done.value = false
        taskId.value = ''
        await nextTick()

        await new Promise<boolean>((resolve) => {
            finish = resolve
            taskId.value = res.task_id as string
        })

        return false
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')
        running.value = false

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
    }
}
</style>
