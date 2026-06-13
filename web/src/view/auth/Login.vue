<template>
    <div class="view">
        <main class="auth">
            <div class="card">
                <div class="header">
                    <img src="@/asset/image/logo.webp" alt="SLINX" />
                    <span>SLINX</span>
                </div>
                <div class="body">
                    <span>嗨，你好</span>
                    <form @submit.prevent="handleLogin">
                        <input v-model="form.username" type="text" placeholder="用户名" autocomplete="username" />
                        <input v-model="form.password" type="password" placeholder="密码" autocomplete="current-password" />
                    </form>
                    <button type="submit" :disabled="loading" @click="handleLogin">
                        {{ loading ? '登录中...' : '登录' }}
                    </button>
                </div>
                <div class="divider">SLINX</div>
                <div class="footer">
                    <h1>自由互联 连接世界</h1>
                </div>
            </div>
        </main>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import { keyStack } from '@/composable/keyStack'
import Footer from '@/view/Footer.vue'
import { login } from '@/api/auth'

const router = useRouter()
const loading = ref(false)
const modal = inject<any>('modal')

const form = ref({ username: '', password: '' })

const handleLogin = async () => {
    if (!form.value.username || !form.value.password) {
        modal.value.show('error', '请输入用户名和密码')
        return
    }
    loading.value = true
    try {
        const res = await login(form.value.username, form.value.password)
        localStorage.setItem('token', res.token)
        router.push('/')
    } catch (err: any) {
        modal.value.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

keyStack(() => true, (e) => {
    if (e.key === 'Enter') handleLogin()
})
</script>