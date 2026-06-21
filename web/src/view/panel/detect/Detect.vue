<template>
    <main>
        <div class="page">
            <div class="body">
                <IPDetect />
                <UnlockDetect />
                <BackRouteDetect />
            </div>
        </div>
    </main>
</template>

<script setup lang="ts">
import IPDetect from '@/view/panel/detect/widget/IP.vue'
import UnlockDetect from '@/view/panel/detect/widget/Unlock.vue'
import BackRouteDetect from '@/view/panel/detect/widget/BackRoute.vue'
</script>

<style scoped>
:deep(.detect) {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    padding: 20px;
    background-color: var(--color-bg-dark);
    border-radius: 20px;
    box-shadow: var(--box-shadow);

    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: 10px;

        .title-wrap {
            display: flex;
            align-items: flex-end;
            gap: 5px;

            .title {
                font-weight: bold;
                color: var(--color-text-light);
            }

            .updated {
                font-size: var(--font-size-sm);
                color: var(--color-text-dark);
            }
        }

        .action {
            display: flex;
            align-items: center;
            gap: 5px;
            margin-left: auto;

            > button {
                gap: 5px;
                padding: 5px 15px 5px 10px;
                font-size: var(--font-size-sm);
                background-color: var(--color-bg-light);
                color: var(--color-text-light);

                &:hover:not(:disabled) {
                    background-color: var(--color-bg);
                }

                .icon {
                    font-size: 12px;
                }

                @media (max-width: 768px) {
                    padding: 6px;

                    span {
                        display: none;
                    }
                }
            }
        }
    }
}
</style>