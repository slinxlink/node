<template>
    <nav :class="{ collapse: Shrink }">
        <div class="header">
            <img src="@/asset/image/logo.webp" alt="SLINX" />
            <label class="title shrink">
                <strong>SLINX</strong>
                <span>节点管理</span>
            </label>
            <button class="menu" @click="Menu = !Menu">
                <i class="icon" :class="{ active: Menu }">{{ Menu ? 'close' : 'menu' }}</i>
            </button>
        </div>
        <ul class="link" :class="{ visible: Menu}">
            <li>
                <button class="shrink-btn circle" @click="Shrink = !Shrink">
                    <i class="icon">chevron_right</i>
                </button>
            </li>
            <template v-for="link in links" :key="link.to ?? link.divider">
                <li v-if="link.divider">
                    <div class="divider-notext"></div>
                </li>
                <li v-else>
                    <router-link :to="link.to!">
                        <i class="icon">{{ link.icon }}</i>
                        <span class="shrink">{{ link.label }}</span>
                    </router-link>
                </li>
            </template>
        </ul>
        <div class="footer" :class="{ visible: Menu}">
            <button @click="navFooter = !navFooter">
                <span class="circle lg">{{ configStore.Username.charAt(0) }}</span>
                <span class="shrink">{{ configStore.Username }}</span>
                <i class="icon shrink">more_vert</i>
            </button>
            <div class="content" :class="{ active: navFooter }">
                <a @click.prevent="logout">
                    <i class="icon">logout</i>
                    <span class="shrink">登出</span>
                </a>
            </div>
        </div>
    </nav>
</template>

<script setup lang="ts">
import { configStore, loadConfig } from '@/store/config'
import { logout as logoutApi } from '@/api/auth'

const router = useRouter()
const Shrink = ref(false)
const Menu = ref(false)
const navFooter = ref(false)
const route = useRoute()

const links = computed(() => [
    { to: '/panel/dashboard', icon: 'dashboard',      label: '仪表盘' },
    { divider: true },
    { to: '/panel/inbound',   icon: 'add_link',       label: '入站' },
    { to: '/panel/user',      icon: 'rss_feed',       label: '用户' },
    ...(configStore.BoardEnable ? [{ to: '/panel/board', icon: 'flight', label: '面板对接' }] : []),
    { divider: true },
    { to: '/panel/detect',    icon: 'travel_explore', label: 'IP检测' },
    { to: '/panel/core',      icon: 'handyman',       label: '核心配置' },
    { to: '/panel/config',    icon: 'settings',       label: '面板设置' },
])

onMounted(() => loadConfig())

async function logout() {
    await logoutApi()
    localStorage.removeItem('token')
    router.push('/login')
}

watch(() => route.path, () => {
    Menu.value = false
})
</script>

<style scoped>
nav {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--color-bg-dark);
    z-index: 99;
    box-shadow: var(--box-shadow);

    &.collapse {
        .shrink {
            max-width: 999px;
            visibility: visible;
        }
        .shrink-btn {
            transform: rotate(180deg);
        }
    }
    &:not(.collapse) {
        .header {
            gap: 0;
        }
        a {
            gap: 0;
        }
        .footer button {
            gap: 0;
        }
    }

    .shrink {
        display: flex;
        visibility: hidden;
        max-width: 0;
        transition: max-width 0.3s ease;
        overflow: hidden;
        white-space: nowrap;
    }

    .header {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 20px;
        height: 80px;
        white-space: nowrap;
        border-bottom: 1px solid var(--color-bg);
        transition: gap 0.3s ease;

        img {
            width: 40px;
            height: 40px;
        }

        .title {
            display: flex;
            flex-direction: row;
            gap: 10px;
            align-items: baseline;

            strong {
                color: var(--color-text-light);
                font-size: var(--font-size-xl);
                font-weight: bold;
            }
            span {
                color: var(--color-text-dark);
                font-size: var(--font-size-sm);
            }
        }
        .menu {
            display: none;
            margin-left: auto;
            padding: 10px;
            border-radius: 5px;
            background-color: var(--color-bg-dark);

            .icon {
                transition: 0.3s ease;

                &.active {
                    transform: rotate(180deg);
                }
            }

            &:hover {
                background-color: var(--color-primary);
            }
            &:active {
                background-color: var(--color-primary-dark);
            }
        }
    }

    .shrink-btn {
        margin-left: 15px;
        transition: 0.3s ease;
        background-color: var(--color-bg-light);
        padding: 5px;
    }

    ul {
        display: flex;
        flex-direction: column;
        white-space: nowrap;
        list-style: none;
        width: 100%;
        overflow-y: auto;
        padding-top: 10px;
        gap: 10px;

        li {
            padding: 0 10px;
        }
    }

    .footer {
        padding: 10px;
        margin-top: auto;
        border-top: 1px solid var(--color-bg);
        white-space: nowrap;

        button {
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
            color: var(--color-text);
            font-size: var(--font-size-md);
            width: 100%;
            padding: 10px;
            border-radius: 10px;
            background-color: var(--color-bg-dark);
            transition: gap 0.3s ease, background-color 0.3s ease;
            gap: 10px;

            .icon {
                margin-left: auto;
            }

            label {
                flex: 1;
                display: flex;
                flex-direction: row;
                align-items: center;
            }

            &:hover {
                background-color: var(--color-bg-light);
            }
            &:active {
                background-color: var(--color-bg);
                color: var(--color-text-light);
            }
        }

        .content {
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.3s ease-out;

            &.active {
                max-height: 100px;
            }
        }
    }

    a {
        display: flex;
        text-align: left;
        align-items: center;
        padding: 10px 20px;
        color: var(--color-text);
        gap: 10px;
        border-radius: 10px;
        transition: gap 0.3s ease, background-color 0.3s ease;

        &:hover {
            background-color: var(--color-bg-light);
            color: var(--color-text);
        }
        &:active {
            background-color: var(--color-bg);
            color: var(--color-text-light);
        }
        &.router-link-active {
            background-color: var(--color-bg);
            pointer-events: none;
            color: var(--color-text-light);
        }
    }

    @media (max-width: 768px) {
        width: 100%;
        height: auto;

        .shrink {
            max-width: 999px;
            visibility: visible;
        }

        &:not(.collapse) {
            .header {
                gap: 10px;

                .menu {
                    display: block;
                }
            }
            a {
                gap: 10px;
            }
            .footer button {
                gap: 10px;
            }
        }

        .shrink-btn {
            display: none;
        }

        ul, .footer {
            position: fixed;
            left: 0;
            right: 0;
            background-color: var(--color-bg-dark);
            opacity: 0;
            visibility: hidden;
            transition: opacity 0.3s ease, visibility 0s 0.3s;
            pointer-events: none;

            &.visible {
                opacity: 1;
                visibility: visible;
                transition: opacity 0.3s ease, visibility 0s;
                pointer-events: auto;
            }
        }

        ul {
            top: 80px;
            height: calc(100% - 80px);
        }

        .footer {
            bottom: 0;
        }
    }
}
</style>