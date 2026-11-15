<template>
    <Transition name="mask">
        <div v-if="state.visible && state.type === 'loading'" class="mask"></div>
    </Transition>
    <div class="toast">
        <Transition name="slide">
            <div v-if="state.visible" class="pill">
                <Icon :name="icon" :class="[state.type, { spin: state.type === 'loading' }]" />
                <span class="message">{{ state.message }}</span>
            </div>
        </Transition>
        <Transition name="fade">
            <ul v-if="state.visible && state.detail.length">
                <li v-for="(item, index) in shown" :key="index">
                    <Icon :name="levelIcon[item.level]" :class="item.level" />
                    <span>{{ item.message }}</span>
                </li>
                <li v-if="rest > 0">
                    <span>… 还有 {{ rest }} 条</span>
                </li>
            </ul>
        </Transition>
    </div>
</template>

<script setup lang="ts">
import { toastState as state } from '@/composable/toast'
import Icon from '@/component/widget/Icon.vue'

const limit = 10

const levelIcon = { warn: 'error', error: 'cancel' }

const icon = computed(() => {
    switch (state.type) {
        case 'success':
            return 'check_circle'
        case 'warn':
            return 'error'
        case 'error':
            return 'cancel'
        case 'loading':
            return 'progress_activity'
    }
})

const shown = computed(() => state.detail.slice(0, limit))

const rest = computed(() => state.detail.length - shown.value.length)
</script>

<style scoped>
.mask {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    -webkit-backdrop-filter: blur(5px);
    backdrop-filter: blur(5px);
    z-index: 98;

    &.mask-enter-active,
    &.mask-leave-active {
        transition: opacity 0.3s;
    }

    &.mask-enter-from,
    &.mask-leave-to {
        opacity: 0;
    }
}

.toast {
    position: fixed;
    top: 10px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 5px;
    width: max-content;
    max-width: calc(100% - 20px);
    color: var(--text-100);
    font-size: 12px;
    z-index: 99;

    .pill {
        display: flex;
        align-items: center;
        gap: 5px;
        height: 30px;
        padding: 5px 10px 5px 5px;
        border-radius: 999px;
        background-color: var(--background-300);
        box-shadow: var(--shadow);

        &.slide-enter-active,
        &.slide-leave-active {
            transition: transform 0.3s;
        }

        &.slide-enter-from,
        &.slide-leave-to {
            transform: translateY(-40px);
        }
    }

    .icon {
        font-size: 20px;

        &.success {
            color: var(--green);
        }

        &.warn {
            color: var(--yellow);
        }

        &.error {
            color: var(--red);
        }

        &.loading {
            color: var(--text-300);
        }
    }

    .spin {
        animation: spin 1s linear infinite;
    }

    ul {
        display: flex;
        flex-direction: column;
        max-width: 500px;
        margin: 0;
        padding: 5px 10px;
        border-radius: 15px;
        background-color: var(--background-300);
        box-shadow: var(--shadow);
        list-style: none;
        line-height: 1.6;

        &.fade-enter-active {
            transition: opacity 0.3s 0.3s;
        }

        &.fade-leave-active {
            transition: opacity 0.3s;
        }

        &.fade-enter-from,
        &.fade-leave-to {
            opacity: 0;
        }

        li {
            display: flex;
            align-items: flex-start;
            gap: 5px;
            overflow-wrap: anywhere;

            .icon {
                margin-top: 0.3em;
                font-size: 12px;
            }
        }
    }
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}
</style>
