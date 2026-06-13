<template>
    <a class="ad" href="https://www.slinx.ink" target="_blank" rel="noopener noreferrer">
        <div class="header">
            <div class="title">
                <h1>www.slinx.link</h1>
                <h2>独享节点 · 独立IP · 专属域名 · 无限设备</h2>
            </div>
            <div class="logo">
                <img src="@/asset/image/logo.webp" alt="SLINX" />
                <span>SLINX</span>
            </div>
        </div>
        <div class="divider-notext"></div>
        <div class="body">
            <p>厌倦了共享机场？自建节点太麻烦？交给我们，省心用网</p>
            <p>需要技术支持、定制APP或运维合作？欢迎联系我们</p>
            <a>了解跟多<i class="icon">open_in_new</i></a>
        </div>
    </a>
</template>

<script setup lang="ts">

</script>

<style scoped>
.ad {
    display: flex;
    flex-direction: column;
    background-color: var(--color-bg-dark);
    border-radius: 20px;
    width: calc(50% - var(--gap-size-half));
    padding: 10px 20px;
    color: var(--color-text-light);
    border-color: var(--color-bg-dark);
    transition: border-color 0.3s;
    cursor: pointer;

    &:hover {
        border: 1px solid var(--color-primary-light);
        padding: 9px 19px;
    }

    &:active {
        border: 1px solid var(--color-primary-dark);
        padding: 9px 19px;
    }

    @media (max-width: 768px) {
        width: 100%;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-bottom: 5px;
        
        .title {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            h1 {
                font-size: var(--font-size);
            }
            h2 {
                font-size: var(--font-size-sm);
                color: var(--color-text-dark);
            }
        }

        .logo {
            display: flex;
            justify-content: row;
            margin-left: auto;
            font-size: var(--font-size-lg);
            font-weight: bold;
            align-items: center;
            gap: 5px;
            img {
                height: 20px;
                width: 20px;
            }
        }
    }
    .body {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 0px !important;
        padding-top: 20px;

        p {
            font-size: var(--font-size-sm);
            color: var(--color-text-dark);
            line-height: 1.5;
        }

        a {
            display: flex;
            flex-direction: row;
            align-items: center;
            gap: 5px;
            font-size: var(--font-size-sm);
            padding-top: 10px;
            .icon {
                font-size: var(--font-size-sm);
            }
        }
    }
}
</style>