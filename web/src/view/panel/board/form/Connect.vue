<template>
    <div class="board-connect">
        <div class="section-body">
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
        </div>
    </div>
</template>

<script setup lang="ts">
import Input from '@/component/ui/Input.vue'

const form = defineModel<any>({
    default: () => ({ Host: '', NodeID: 0, Key: '' }),
})
</script>

<style scoped>
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

    &.quarter :deep(.input-wrap) {
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
