<!-- view/auth/Login.vue -->
<template>
    <Page class="login">
        <div class="box">
            <div class="brand">
                <img src="@/asset/image/logo.webp" alt="SLINX" draggable="false" />
                <h1>SLINX</h1>
            </div>
            <h2>嗨，你好</h2>
            <form @submit.prevent="submit">
                <Input v-model:value="form.username" placeholder="用户名" autocomplete="username" />
                <Input
                    v-model:value="form.password"
                    type="password"
                    placeholder="密码"
                    autocomplete="current-password"
                />
                <Button type="submit" variant="primary" :disabled="loading">
                    {{ loading ? '登录中...' : '登录' }}
                </Button>
            </form>
            <Divider>SLINX</Divider>
            <hgroup>
                <h3>自由互联 连接世界</h3>
            </hgroup>
        </div>
    </Page>
</template>

<script setup lang="ts">
import Page from '@/component/layout/Page.vue'
import Divider from '@/component/layout/Divider.vue'
import Input from '@/component/widget/Input.vue'
import Button from '@/component/widget/Button.vue'
import { login } from '@/api/auth'
import { toast } from '@/composable/toast'
import { bind, unbind } from '@/util/keyboard'

const router = useRouter()
const loading = ref(false)
const form = ref({ username: '', password: '' })

async function submit() {
    if (loading.value) return
    if (!form.value.username || !form.value.password) {
        toast.failure('请输入用户名和密码')
        return
    }
    loading.value = true
    try {
        const res = await login(form.value.username, form.value.password)
        localStorage.setItem('token', res.token)
        router.push('/')
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '登录失败')
    } finally {
        loading.value = false
    }
}

// 回车提交；表单里的 submit 按钮已经处理了输入框内的回车，这里只管焦点不在输入框的情况
const keys = bind({
    Enter: (e) => {
        if ((e.target as HTMLElement).tagName === 'BUTTON') return
        submit()
    },
})
onUnmounted(() => unbind(keys))
</script>

<style scoped>
.login {
    /* 内容少，登录框在页面中间；底部栏仍在最下面 */
    &:deep(main) {
        justify-content: center;
    }

    .box {
        display: flex;
        flex-direction: column;
        gap: 20px;
        width: 100%;
        max-width: 400px;
        min-width: 0;
        margin: 20px auto;
        padding: 30px 20px 20px;
        border-radius: 20px;
        background-color: var(--background-300);
        box-shadow: var(--shadow);
        font-size: 14px;
        color: var(--text-200);

        .brand {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;
            color: var(--text-100);

            h1 {
                font-size: 24px;
            }

            img {
                width: 40px;
                height: 40px;
            }
        }

        h2 {
            width: 100%;
            margin: 30px 0;
            font-size: 24px;
            font-weight: normal;
            color: var(--text-100);
            text-align: center;
        }

        hgroup {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 5px;
            margin-bottom: 10px;

            h3 {
                text-align: center;
                font-size: 16px;
            }
        }

        form {
            display: flex;
            flex-direction: column;
            gap: 20px;

            &:deep(input),
            button[type='submit'] {
                height: 50px;
            }
        }
    }
}
</style>
