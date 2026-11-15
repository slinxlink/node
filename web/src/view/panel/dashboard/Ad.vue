<template>
    <Anchor class="ad" href="https://slinx.link" target="_blank" rel="noopener">
        <div class="head">
            <hgroup>
                <h1>slinx.link</h1>
                <h2>独享节点 · 独立IP · 专属域名 · 无限设备</h2>
            </hgroup>
            <div class="logo">
                <Image :src="logo" alt="SLINX" />
                <span>SLINX</span>
            </div>
        </div>
        <Divider />
        <div class="body">
            <p>厌倦了共享机场？自建节点太麻烦？交给我们，省心用网</p>
            <p>需要技术支持、定制APP或运维合作？欢迎联系我们</p>
            <span class="more">
                了解更多
                <Icon name="open_in_new" />
            </span>
        </div>
    </Anchor>
</template>

<script setup lang="ts">
import logo from '@/asset/image/logo.webp'
import Anchor from '@/component/widget/Anchor.vue'
import Divider from '@/component/layout/Divider.vue'
import Icon from '@/component/widget/Icon.vue'
import Image from '@/component/widget/Image.vue'
</script>

<style scoped>
/* 外围是 Anchor，这里把它改成卡片：背景不随状态变，悬停和按下只改边框颜色，里面的排版自己管 */
.ad {
    flex-direction: column;
    align-items: stretch;
    justify-content: flex-start;
    gap: 0;
    padding: 9px 19px; /* 原来的 10px 20px，减掉 1px 边框 */
    border: 1px solid var(--background-100);
    border-radius: 20px;
    background-color: var(--background-300);
    box-shadow: var(--shadow);
    text-align: left;
    transition: border-color 0.3s;

    &:hover,
    &:focus-visible {
        border-color: var(--primary-200);
        background-color: var(--background-300);
    }

    &:active {
        border-color: var(--primary-400);
        background-color: var(--background-300);
    }

    .head {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;
        padding-bottom: 5px;

        hgroup {
            display: flex;
            flex-direction: column;

            h1 {
                color: var(--text-100);
                font-size: 14px;
            }

            h2 {
                color: var(--text-300);
                font-size: 12px;
                font-weight: normal;
            }
        }

        .logo {
            display: flex;
            align-items: center;
            gap: 5px;
            color: var(--text-100);
            font-size: 20px;
            font-weight: bold;

            img {
                width: 20px;
                height: 20px;
            }
        }
    }

    .body {
        display: flex;
        flex-direction: column;
        padding-top: 20px;
        color: var(--text-300);
        font-size: 12px;
        line-height: 1.5;

        .more {
            display: flex;
            align-items: center;
            gap: 5px;
            padding-top: 10px;
            color: var(--primary-100);

            .icon {
                font-size: 12px;
            }
        }
    }
}
</style>
