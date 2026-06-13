<template>
    <div class="board-form">
        <div class="section">
            <div class="section-body">
                <div class="form-row quarter">
                    <span class="form-label">备注</span>
                    <Input v-model="form.Name" />
                </div>
                <div class="form-row">
                    <span class="form-label">面板地址</span>
                    <Input v-model="form.Host" placeholder="https://example.com" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">节点ID</span>
                    <Input v-model="form.NodeID" type="number" />
                </div>
                <div class="form-row">
                    <span class="form-label">通讯密钥</span>
                    <Input v-model="form.Key" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">面板类型</span>
                    <Select v-model="form.Type" :options="[
                        { label: 'SLINX', value: 'SLINX' },
                    ]" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">同步间隔(s)</span>
                    <Input v-model="form.SyncInterval" type="number" />
                </div>
                <div class="form-row quarter">
                    <span class="form-label">关联入站</span>
                    <Select v-model="form.Inbound" :options="inboundOptions" />
                </div>
                <div class="form-row">
                    <span class="form-label">启用</span>
                    <Toggle v-model="form.Enable" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'

const props = defineProps<{
    inbounds: any[]
}>()

const form = defineModel<any>({ default: () => ({
    Enable: true,
    Name: '',
    Host: '',
    NodeID: 0,
    Key: '',
    Inbound: 0,
    Type: 'SLINX',
    SyncInterval: 60,
}) })

const inboundOptions = computed(() =>
    props.inbounds.map(ib => ({
        label: String(ib.Port),
        value: ib.ID,
    }))
)
</script>

<style scoped>
.board-form {
    interpolate-size: allow-keywords;
}

.section-body {
    padding: 16px 20px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.form-row {
    display: flex;
    align-items: center;
    gap: 16px;

    &.half :deep(.input-wrap),
    &.half :deep(.select-wrap) {
        max-width: 50%;
    }

    &.quarter :deep(.input-wrap),
    &.quarter :deep(.select-wrap) {
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
</style>