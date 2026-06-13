<template>
    <div class="user-form">
        <div class="section">
            <div class="section-body">
                <div class="form-row">
                    <span class="form-label">用户名</span>
                    <Input v-model="form.Name" />
                </div>
                <div class="form-row">
                    <span class="form-label">Token</span>
                    <RefreshBtn @click="genToken">
                        <Input v-model="form.Token" />
                    </RefreshBtn>
                </div>
                <div class="form-row">
                    <span class="form-label">UUID</span>
                    <RefreshBtn @click="genUUID">
                        <Input v-model="form.UUID" />
                    </RefreshBtn>
                </div>
                <div class="form-row">
                    <span class="form-label">密码</span>
                    <RefreshBtn @click="genPassword">
                        <Input v-model="form.Password" />
                    </RefreshBtn>
                </div>
                <div class="form-row">
                    <span class="form-label">关联入站</span>
                    <MultiSelect
                        v-model="form.Inbounds"
                        :options="inboundOptions"
                    />
                </div>
                <div class="form-row">
                    <span class="form-label">启用</span>
                    <Toggle v-model="form.Enable" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Input from '@/component/ui/Input.vue'
import MultiSelect from '@/component/ui/MultiSelect.vue'
import RefreshBtn from '@/component/ui/RefreshBtn.vue'
import { generateToken, generateUUID, generatePassword } from '@/api/generate'

const props = defineProps<{
    inbounds: any[]
}>()

const form = defineModel<any>({ default: () => ({
    Enable: true,
    Name: '',
    Token: '',
    UUID: '',
    Password: '',
    Inbounds: '[]',
}) })

const inboundOptions = computed(() =>
    props.inbounds.map(ib => ({
        label: String(ib.Port),
        value: String(ib.ID),
    }))
)

async function genToken() {
    const res = await generateToken()
    form.value.Token = res.token
}

async function genUUID() {
    const res = await generateUUID()
    form.value.UUID = res.uuid
}

async function genPassword() {
    const res = await generatePassword()
    form.value.Password = res.password
}
</script>

<style scoped>
.user-form {
    interpolate-size: allow-keywords;
}

.section-body {
    padding: 16px 20px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.form-row {
    display: flex;
    align-items: center;
    gap: 16px;
}

.form-label {
    width: 100px;
    flex-shrink: 0;
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
    text-align: right;
}
</style>