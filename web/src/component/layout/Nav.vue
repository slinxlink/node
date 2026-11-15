<template>
    <nav :class="{ open: sideOpen, menu: menuOpen }">
        <section class="header">
            <div class="brand">
                <Image :src="logo" />
                <hgroup class="collapse">
                    <strong>SLINX</strong>
                    <span>节点管理</span>
                </hgroup>
            </div>
            <Button class="menu-btn" @click="menuOpen = !menuOpen">
                <Icon :name="menuOpen ? 'close' : 'menu'" :class="{ active: menuOpen }" />
            </Button>
        </section>
        <section class="router">
            <template v-for="(item, index) in routes" :key="index">
                <Divider v-if="'divider' in item" />
                <router-link v-else :to="item.to" custom v-slot="{ navigate, href, isActive }">
                    <Anchor :href="href" :class="{ active: isActive }" @click="navigate">
                        <Icon :name="item.icon" />
                        <span class="collapse">{{ item.label }}</span>
                    </Anchor>
                </router-link>
            </template>
            <Button class="side-btn" @click="sideOpen = !sideOpen">
                <Icon name="chevron_right" />
            </Button>
        </section>
        <section class="user">
            <Button @click="userOpen = !userOpen">
                <Avatar :name="configStore.Username" />
                <span class="collapse">{{ configStore.Username }}</span>
                <Icon name="more_vert" class="collapse" />
            </Button>
        </section>
        <section class="footer" :class="{ open: userOpen }">
            <div>
                <Anchor @click.prevent="signOut">
                    <Icon name="logout" />
                    <span class="collapse">登出</span>
                </Anchor>
            </div>
        </section>
    </nav>
</template>

<script setup lang="ts">
import { configStore, loadConfig } from '@/store/config'
import { logout } from '@/api/auth'
import logo from '@/asset/image/logo.webp'
import Image from '@/component/widget/Image.vue'
import Button from '@/component/widget/Button.vue'
import Anchor from '@/component/widget/Anchor.vue'
import Divider from '@/component/layout/Divider.vue'
import Icon from '@/component/widget/Icon.vue'
import Avatar from '@/component/widget/Avatar.vue'

type NavLink = { to: string; icon: string; label: string } | { divider: true }

const routes = computed<NavLink[]>(() => [
    { to: '/dashboard', icon: 'dashboard', label: '仪表盘' },
    { divider: true },
    { to: '/inbound', icon: 'add_link', label: '入站' },
    { to: '/user', icon: 'rss_feed', label: '用户' },
    { to: '/endpoint', icon: 'cloud_sync', label: '端点' },
    ...(configStore.BoardEnable
        ? [{ to: '/board', icon: 'flight', label: '面板对接' }]
        : []),
    { divider: true },
    { to: '/detect', icon: 'travel_explore', label: '线路检测' },
    { divider: true },
    { to: '/cert', icon: 'verified_user', label: '证书' },
    { to: '/core', icon: 'handyman', label: '核心配置' },
    { to: '/config', icon: 'settings', label: '面板设置' },
])

const router = useRouter()
const route = useRoute()

const sideOpen = ref(false)
const menuOpen = ref(false)
const userOpen = ref(false)

watch(
    () => route.path,
    () => {
        menuOpen.value = false
    },
)

onMounted(() => loadConfig())

async function signOut() {
    await logout()
    localStorage.removeItem('token')
    router.push('/login')
}
</script>

<style scoped>
nav {
    --gap: 0px;
    --header-h: 80px;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    width: 80px;
    height: 100%;
    transition: width 0.3s ease;
    background-color: var(--background-300);
    box-shadow: var(--shadow);
    white-space: nowrap;
    z-index: 79;

    &.open {
        --gap: 10px;
        width: 200px;

        .collapse {
            max-width: 999px;
            visibility: visible;
        }

        .side-btn {
            transform: rotate(180deg);
        }
    }

    button,
    a {
        background-color: transparent;
    }

    /* 触屏没有真正的悬停，点完按钮不要留着悬停底色 */
    @media (hover: none) {
        button:hover {
            background-color: transparent;
        }
    }

    /* 收起时隐藏的文字 */
    .collapse {
        display: flex;
        visibility: hidden;
        max-width: 0;
        transition: max-width 0.3s ease;
        overflow: hidden;
    }

    /* 头部和底部的药丸按钮：靠左，图标推到最右 */
    .brand,
    .user button {
        justify-content: flex-start;
        gap: var(--gap);
        overflow: hidden;
        padding: 10px;

        .icon {
            margin-left: auto;
            margin-right: -5px;
            color: var(--text-300);
        }
    }

    .header {
        display: flex;
        align-items: center;
        padding: 10px;
        height: var(--header-h);

        img {
            flex-shrink: 0;
            width: 40px;
            height: 40px;
        }

        hgroup {
            flex-direction: column;
            justify-content: center;
            text-align: left;

            strong {
                color: var(--text-100);
                font-size: 20px;
                font-weight: bold;
                line-height: 1.2;
            }

            span {
                color: var(--text-300);
                font-size: 12px;
            }
        }

        .brand {
            display: flex;
            align-items: center;
            flex: 1;
            min-width: 0;
            transition: gap 0.3s ease;
        }

        .menu-btn {
            display: none;
            width: 40px;
            height: 40px;
            margin-left: auto;
            padding: 10px;
            border-radius: 10px;

            .icon {
                transition: 0.3s ease;

                &.active {
                    transform: rotate(180deg);
                }
            }
        }
    }

    .router,
    .footer {
        a {
            justify-content: start;
            gap: var(--gap);
            padding: 10px 20px;
            transition:
                background-color 0.3s,
                color 0.3s,
                border-color 0.3s,
                gap 0.3s ease;
            color: var(--text-200);

            &:hover {
                background-color: var(--background-100);
                color: var(--text-200);
            }

            &:active {
                background-color: var(--background-200);
                color: var(--text-100);
            }
        }
    }

    /* 底部菜单：用 grid 行高过渡展开收起 */
    .footer {
        display: grid;
        grid-template-rows: 0fr;
        transition:
            grid-template-rows 0.3s ease,
            margin-bottom 0.3s ease,
            padding 0.3s ease;

        > div {
            display: flex;
            flex-direction: column;
            gap: 10px;
            min-height: 0;
            overflow: hidden;
        }

        &.open {
            grid-template-rows: 1fr;
        }
    }

    .footer {
        padding: 0 10px;

        &.open {
            padding-bottom: 10px;
        }
    }

    .router {
        display: flex;
        flex: 1;
        flex-direction: column;
        gap: 10px;
        min-height: 0;
        padding: 10px;
        overflow-y: auto;
        border-block: 1px solid var(--background-200);

        a.active {
            background-color: var(--background-200);
            pointer-events: none;
            color: var(--text-100);
        }

        .side-btn {
            width: auto;
            margin: auto auto 0 15px;
            padding: 5px;
            transition: 0.3s ease;
            background-color: var(--background-200);

            &:hover {
                background-color: var(--background-100);
            }

            &:active {
                background-color: var(--primary-600);
            }
        }
    }

    .user {
        padding: 10px;
    }

    @media (max-width: 768px) {
        --gap: 10px;
        position: fixed;
        inset: 0 0 auto;
        width: auto;
        height: auto;
        background-color: color-mix(in srgb, var(--background-300) 50%, transparent);
        -webkit-backdrop-filter: blur(5px);
        backdrop-filter: blur(5px);

        &.open {
            width: auto;
        }

        /* 菜单展开时整个导航铺满屏幕 */
        &.menu {
            inset: 0;
            height: 100%;
        }

        /* 菜单没展开时只留头部 */
        &:not(.menu) > :not(.header) {
            display: none;
        }

        .collapse {
            max-width: 250px;
            visibility: visible;
        }

        .header {
            .menu-btn {
                display: block;
            }

            .brand {
                flex: none;
                padding-right: 20px;
            }
        }

        .router {
            a.active {
                background-color: transparent;
            }

            .side-btn {
                display: none;
            }
        }
    }
}
</style>
