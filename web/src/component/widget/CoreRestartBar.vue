<template>
    <div class="save-bar">
        <span class="tip">
            <i class="icon">report</i>
            此处的所有更改都需要保存并重启核心才能生效
        </span>
        <div class="actions">
            <button :disabled="!dirty" @click="handleSave">保存</button>
            <button :disabled="!saved" @click="handleRestart">重启核心</button>
        </div>
    </div>
</template>

<script setup lang="ts">
const modal = inject<any>('modal')

const props = defineProps<{
    dirty: boolean
    onSave: () => Promise<void>
    onRestart: () => Promise<void>
}>()

const emit = defineEmits(['update:dirty'])

const saved = ref(true)

async function handleSave() {
    try {
        await props.onSave()
        saved.value = true
        emit('update:dirty', false)
    } catch {
        // 失败不重置
    }
}

async function handleRestart() {
    try {
        await props.onRestart()
        modal.value?.show('success', '核心已重启')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

watch(() => props.dirty, (val) => {
    if (val) saved.value = false
})
</script>

<style scoped>
.save-bar {
    display: flex;
    flex-direction: row;
    gap: 10px;
    padding: 20px;
    background-color: var(--color-bg-dark);
    border-radius: 20px;
    width: 100%;

    .actions {
        display: flex;
        flex-direction: row;
        gap: 10px;
        margin-left: auto;

        button {
            font-size: var(--font-size-sm);
            padding: 9px 20px;
            line-height: 1;
            border: 1px solid transparent;

            &:disabled {
                cursor: not-allowed;
                background-color: var(--color-bg-dark);
                border: 1px solid var(--color-bg-light);
            }
        }
    }
}

.tip {
    display: flex;
    align-items: center;
    justify-self: start;
    gap: 5px;
    font-size: var(--font-size-sm);
    color: var(--color-yellow);

    .icon {
        font-size: var(--font-size-sm);
    }
}

@media (max-width: 768px) {
    .save-bar {
        flex-direction: column;

        .actions {
            margin-left: 0;
        }
    }
}
</style>