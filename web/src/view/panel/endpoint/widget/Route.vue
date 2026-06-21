<template>
    <div class="route">
        <div class="form-row">
            <span class="form-label">绑定入站</span>
            <MultiSelect
                v-model="selected"
                :options="inboundOptions"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
import MultiSelect from '@/component/ui/MultiSelect.vue'
import { getInbounds } from '@/api/inbound'
import { getRule, saveRule } from '@/api/route'


const modal = inject<any>('modal')

const props = defineProps<{ endpoint: any }>()

const inbounds = ref<any[]>([])
const selected = ref<string[]>([])

const inboundOptions = computed(() =>
    inbounds.value
        .filter(ib => ib.Enable)
        .map(ib => ({
            label: String(ib.Port),
            value: String(ib.Port),
        }))
)

onMounted(async () => {
    inbounds.value = await getInbounds()

    const res = await getRule()
    const rules = res.data as { Sort: number, Key: string, Value: string }[]

    const tagJSON = JSON.stringify(props.endpoint.Tag)
    const outboundRow = rules.find(r => r.Key === 'outbound' && r.Value === tagJSON)
    if (outboundRow) {
        const inboundRow = rules.find(r => r.Sort === outboundRow.Sort && r.Key === 'inbound')
        if (inboundRow) {
            selected.value = JSON.parse(inboundRow.Value)
        }
    }
})

async function save() {
    try {
        await saveRule(props.endpoint.Tag, selected.value)
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

defineExpose({ save })
</script>

<style scoped>
.route {
    padding: 20px;
    width: 100%;
    
    .form-row {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .form-label {
        width: 100px;
        flex-shrink: 0;
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
        text-align: right;
    }
}

</style>