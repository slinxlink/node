<template>
    <div class="section">
        <div class="header" @click="toggle">
            <div class="title">{{ title }}</div>
            <i class="icon">{{ open ? 'expand_less' : 'expand_more' }}</i>
        </div>
        <div class="body" ref="bodyRef">
            <div class="content">
                <slot />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{ 
    title: string
    defaultOpen?: boolean
}>(), {
    defaultOpen: true
})

const open = ref(props.defaultOpen)
const bodyRef = ref<HTMLElement>()

function toggle() {
    const el = bodyRef.value
    if (!el) return
    if (open.value) {
        el.style.height = el.scrollHeight + 'px'
        requestAnimationFrame(() => {
            el.style.height = '0'
        })
    } else {
        el.style.height = el.scrollHeight + 'px'
        el.addEventListener('transitionend', () => {
            el.style.height = 'auto'
        }, { once: true })
    }
    open.value = !open.value
}

onMounted(() => {
    const el = bodyRef.value
    if (!el) return
    if (!open.value) {
        el.style.height = '0'
    }
})
</script>

<style scoped>
.section {
    background: var(--color-bg-dark);
    border-radius: 20px;
    overflow: hidden;
    width: 100%;

    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px 20px;
        cursor: pointer;
        user-select: none;
        color: var(--color-text-light);
        font-weight: 500;
        transition: background-color 0.3s ease;

        .icon {
            color: var(--color-text-dark);
            font-size: 20px;
            transition: transform 0.3s ease;
        }

        &:hover {
            background-color: var(--color-primary);
            .icon { color: var(--color-text-light); }
        }
        &:active {
            background-color: var(--color-primary-dark);
        }
    }

    .body {
        overflow: hidden;
        height: auto;
        transition: height 0.3s ease;

        .content {
            display: flex;
            flex-direction: column;
            width: 100%;
            padding: 0 20px 10px;
            gap: 10px;
        }
    }

    &:has([style*="height: 0"]) .icon {
        transform: rotate(180deg);
    }
}
</style>