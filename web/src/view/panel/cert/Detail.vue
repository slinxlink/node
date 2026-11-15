<!-- view/panel/cert/Detail.vue -->
<template>
    <div class="detail">
        <Field>
            <template #prefix>域名</template>
            {{ cert.Domain }}
        </Field>
        <Field>
            <template #prefix>申请方式</template>
            {{ display.mode(cert.Mode) }}
        </Field>
        <Field>
            <template #prefix>状态</template>
            <Stamp :color="status.color">{{ status.text }}</Stamp>
        </Field>
        <Field>
            <template #prefix>过期时间</template>
            {{ formatTime(cert.ExpireAt) }}
        </Field>
        <Field v-if="cert.Mode !== 'manual'">
            <template #prefix>自动续签</template>
            <Toggle :value="cert.AutoRenew" disabled />
        </Field>
        <Field v-if="cert.Mode !== 'manual' && acme">
            <template #prefix>ACME</template>
            {{ acme.Email }} [{{ display.acmeProvider(acme.Provider) }}]
        </Field>
        <Field v-if="cert.Mode === 'dns' && dns">
            <template #prefix>DNS</template>
            {{ dns.Name }} [{{ display.dnsProvider(dns.Provider) }}]
        </Field>
        <Field>
            <template #prefix>证书内容</template>
            <Textarea :value="content.cert" readonly />
        </Field>
        <Field>
            <template #prefix>私钥内容</template>
            <Textarea :value="content.key" readonly />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Textarea from '@/component/widget/Textarea.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { getCertContent } from '@/api/cert'
import type { Acme, Cert, DnsAccount } from '@/api/cert'
import { formatTime } from '@/util/format'
import { display } from '@/view/panel/cert/display'

// 证书的详情：基本信息、用的账号、证书和私钥的内容(只读)
const props = defineProps<{
    cert: Cert
    acmes: Acme[]
    dnsList: DnsAccount[]
}>()

const content = ref({ cert: '', key: '' })

const status = computed(() => display.status(props.cert.ExpireAt))
const acme = computed(() => props.acmes.find((item) => item.ID === props.cert.Acme))
const dns = computed(() => props.dnsList.find((item) => item.ID === props.cert.Dns))

onMounted(async () => {
    content.value = await getCertContent(props.cert.ID)
})
</script>

<style scoped>
.detail {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    text-align: left;

    :deep(.prefix) {
        flex-basis: 100px;
        align-self: flex-start;
        min-height: 40px;
    }

    :deep(textarea) {
        height: 140px;
    }
}
</style>
