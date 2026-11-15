<!-- view/panel/user/Alive.vue -->
<template>
    <div class="alive">
        <Stamp v-for="ip in items" :key="ip" color="green">
            <span class="line">{{ ip }}</span>
        </Stamp>
        <span v-if="!items.length" class="none">-</span>
    </div>
</template>

<script setup lang="ts">
import Stamp from '@/component/widget/Stamp.vue'

defineProps<{
    items: string[]
}>()
</script>

<style scoped>
.alive {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 9px 0;
    min-width: 111px;

    .stamp {
        align-self: flex-start;
        min-width: 0;
        width: max-content;
        max-width: var(--cell-max, 100%);

        .line {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    .none {
        color: var(--text-200);
    }
}
</style>
