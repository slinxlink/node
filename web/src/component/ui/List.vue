<template>
    <div class="list">
        <table>
            <thead>
                <tr>
                    <th v-for="header in headers" :key="header">{{ header }}</th>
                </tr>
            </thead>
            <tbody>
                <slot />
            </tbody>
        </table>
    </div>
</template>

<script setup lang="ts">
defineProps<{ headers: string[] }>()
</script>

<style>
.list {
    overflow-x: auto;
    overflow-y: visible;
    border-radius: 20px;
    background-color: var(--color-bg-dark);
    padding: 20px;
    width: 100%;

    table {
        min-width: 100%;
        border-collapse: collapse;
        font-size: var(--font-size);
        white-space: nowrap;
        table-layout: auto;
        
        thead {
            th {
                text-align: left;
                padding-bottom: 10px;
                color: var(--color-text-dark);
                font-weight: bold;
                border-bottom: 1px solid var(--color-bg);
                max-width: 300px;
                overflow: hidden;
                text-overflow: ellipsis;
            }
        }

        tbody {
            td {
                padding: 10px 0;
                border-bottom: 1px solid var(--color-bg);
                color: var(--color-text-light);
                vertical-align: middle;
                max-width: 300px;
                overflow: hidden;
                text-overflow: ellipsis;
            }

            tr:last-child td {
                padding-bottom: 0;
                border-bottom: none;
            }
        }

        th, td {
            &:first-child {
                padding-right: 10px;
            }
            &:not(:first-child):not(:last-child) {
                padding-right: 25px;
            }
            &:nth-last-child(2) {
                width: 100px;
            }
            &:last-child {
                width: 150px;
            }
        }

        .muted {
            color: var(--color-text-dark);
        }

        .actions {
            display: flex;
            gap: 5px;
            align-items: center;
        }

        .icon-btn {
            background-color: transparent;
            padding: 5px;
            border-radius: 5px;
            color: var(--color-text-dark);

            &:hover {
                background-color: var(--color-bg);
                color: var(--color-text-light);
            }

            &.danger:hover {
                color: var(--color-red);
            }
        }

        .empty {
            text-align: center;
            color: var(--color-text-dark);
            padding: 40px !important;
        }
    }
}
</style>