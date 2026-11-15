<!-- view/panel/endpoint/Route.vue -->
<template>
    <div class="route">
        <Field>
            <template #prefix>绑定入站</template>
            <TagSelect v-model:value="selected" :items="items" placeholder="未绑定" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import TagSelect from '@/component/widget/TagSelect.vue'
import { getInbounds } from '@/api/inbound'
import type { Inbound } from '@/api/inbound'
import type { Endpoint } from '@/api/endpoint'
import { getRule, saveRule } from '@/api/route'
import { toast } from '@/composable/toast'

const props = defineProps<{
    endpoint: Endpoint
}>()

const inbounds = ref<Inbound[]>([])
const selected = ref<string[]>([])

// 只列出启用的入站，规则里按端口存
const items = computed(() =>
    inbounds.value
        .filter((item) => item.Enable)
        .map((item) => ({
            label: `${item.Protocol}:${item.Port}${item.Name ? ` - ${item.Name}` : ''}`,
            value: String(item.Port),
        })),
)

// 找出出站是这个端点的那条规则，再取同一组里的入站
onMounted(async () => {
    inbounds.value = await getInbounds()

    const rules = await getRule()
    const outbound = rules.find(
        (item) => item.Key === 'outbound' && item.Value === JSON.stringify(props.endpoint.Tag),
    )

    if (!outbound) return

    const inbound = rules.find((item) => item.Sort === outbound.Sort && item.Key === 'inbound')

    if (inbound) {
        selected.value = JSON.parse(inbound.Value)
    }
})

// 成功返回 true，失败返回 false（弹窗保持打开）
async function submit(): Promise<boolean> {
    try {
        await saveRule(props.endpoint.Tag, selected.value)
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
.route {
    width: 100%;
    text-align: left;

    :deep(.prefix) {
        flex-basis: 100px;
        align-self: flex-start;
        min-height: 40px;
    }

    :deep(.tag-select) {
        margin: 3px 0;
    }
}
</style>
