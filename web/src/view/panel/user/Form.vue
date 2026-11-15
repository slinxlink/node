<!-- view/panel/user/Form.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>启用</template>
            <Toggle v-model:value="form.Enable" />
        </Field>
        <Field>
            <template #prefix>用户名</template>
            <Input v-model:value="form.Name" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>Token</template>
            <Input v-model:value="form.Token" autocomplete="off">
                <template #trailing>
                    <Button variant="text" :tabindex="-1" @click="genToken">
                        <Icon name="casino" />
                    </Button>
                </template>
            </Input>
        </Field>
        <Field>
            <template #prefix>UUID</template>
            <Input v-model:value="form.UUID" autocomplete="off">
                <template #trailing>
                    <Button variant="text" :tabindex="-1" @click="genUUID">
                        <Icon name="casino" />
                    </Button>
                </template>
            </Input>
        </Field>
        <Field>
            <template #prefix>认证</template>
            <Input v-model:value="form.Password" autocomplete="off">
                <template #trailing>
                    <Button variant="text" :tabindex="-1" @click="genPassword">
                        <Icon name="casino" />
                    </Button>
                </template>
            </Input>
        </Field>
        <Field>
            <template #prefix>关联入站</template>
            <TagSelect v-model:value="form.Inbounds" :items="items" placeholder="未关联" />
        </Field>
        <Field>
            <template #prefix>
                <span>流量上限</span>
                <Tip msg="GB[br]0 = 不限" />
            </template>
            <Input v-model:value="form.TrafficLimit" type="number" :min="0" />
        </Field>
        <ResetCycle v-model:cycle="form.ResetCycle" v-model:day="form.ResetDay" />
        <Field>
            <template #prefix>
                <span>有效期</span>
                <Tip msg="留空为长期有效" />
            </template>
            <DateTime v-model:value="form.ExpireAt" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import ResetCycle from '@/component/module/ResetCycle.vue'
import Button from '@/component/widget/Button.vue'
import DateTime from '@/component/widget/DateTime.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import TagSelect from '@/component/widget/TagSelect.vue'
import Tip from '@/component/widget/Tip.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { generatePassword, generateToken, generateUUID } from '@/api/generate'
import type { Inbound } from '@/api/inbound'
import { saveUser } from '@/api/user'
import type { User } from '@/api/user'
import { toast } from '@/composable/toast'
import { formatTime } from '@/util/format'

const props = defineProps<{
    inbounds: Inbound[]
    user?: User // 不传是添加，传了是编辑
}>()

const GB = 1024 ** 3

const NEVER = '0001-01-01T00:00:00Z'

function toText(iso: string): string {
    return !iso || iso.startsWith('0001') ? '' : formatTime(iso)
}

function toISO(text: string): string {
    return text ? new Date(text.replace(' ', 'T')).toISOString() : NEVER
}

const form = ref({
    Enable: true,
    Name: '',
    Token: '',
    UUID: '',
    Password: '',
    Inbounds: [] as number[],
    TrafficLimit: '0',
    ResetCycle: '',
    ResetDay: '1',
    ExpireAt: '',
})

// 关联入站的选项：和用户列表里一样写成「协议:端口 - 名称」
const items = computed(() =>
    props.inbounds.map((item) => ({
        label: `${item.Protocol}:${item.Port}${item.Name ? ` - ${item.Name}` : ''}`,
        value: item.ID,
    })),
)

async function genToken() {
    form.value.Token = (await generateToken()).token
}

async function genUUID() {
    form.value.UUID = (await generateUUID()).uuid
}

async function genPassword() {
    form.value.Password = (await generatePassword()).password
}

onMounted(() => {
    if (props.user) {
        const { Enable, Name, Token, UUID, Password, Inbounds, TrafficLimit, ResetCycle, ResetDay, ExpireAt } = props.user

        let ids: number[] = []

        try {
            ids = JSON.parse(Inbounds || '[]')
        } catch {
            ids = []
        }

        form.value = {
            Enable,
            Name,
            Token,
            UUID,
            Password,
            Inbounds: ids,
            TrafficLimit: String(Math.round((TrafficLimit / GB) * 100) / 100),
            ResetCycle,
            ResetDay: String(ResetDay || 1),
            ExpireAt: toText(ExpireAt),
        }

        return
    }

    // 添加：三个值各自生成，哪个先回来哪个先填
    genToken()
    genUUID()
    genPassword()
})

// 成功返回 true，失败返回 false（弹窗保持打开）；Token 重复等检查在后端
async function submit(): Promise<boolean> {
    // 编辑时把原来的整条带上：后端整条覆盖，不带的字段(创建时间、到期时间)会被清掉
    const data = {
        ...props.user,
        ...form.value,
        Inbounds: JSON.stringify(form.value.Inbounds),
        TrafficLimit: Math.round(Number(form.value.TrafficLimit) * GB),
        ResetDay: Number(form.value.ResetDay),
        ExpireAt: toISO(form.value.ExpireAt),
    }

    try {
        await saveUser(data)
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

    /* 标题列拉宽，长一点的名字也放得下；行变高(标签换行)时标题固定在第一行 */
    :deep(.prefix) {
        flex-basis: 100px;
        justify-content: space-between;
        align-self: flex-start;
        min-height: 40px;
    }

    :deep(.tag-select) {
        margin: 3px 0;
    }
}
</style>
