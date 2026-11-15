<template>
    <span ref="root" class="avatar" :class="{ empty: !name }" :style="style">
        <span v-if="name" ref="letterEl" class="letter"
            >{{ letter }}<i ref="marker" class="baseline"></i
        ></span>
    </span>
</template>

<script setup lang="ts">
const props = withDefaults(
    defineProps<{
        name: string
        size?: number
    }>(),
    {
        size: 40,
    },
)

const style = computed(() => ({
    width: `${props.size}px`,
    height: `${props.size}px`,
    fontSize: `${props.size * 0.45}px`,
}))

const letter = computed(() => props.name.charAt(0))

// 按字形实际墨迹位置补偿偏移，让任何字符垂直居中
// 零尺寸的 marker 底边就是基线；墨迹中心 = 基线 - 字形上下沿的中点，全程取自实际渲染，不依赖字体度量(中文走回退字体时度量对不上)
const root = ref<HTMLElement>()
const letterEl = ref<HTMLElement>()
const marker = ref<HTMLElement>()

function measure() {
    if (!root.value || !letterEl.value || !marker.value) return

    const ctx = document.createElement('canvas').getContext('2d')
    if (!ctx) return

    const css = getComputedStyle(root.value)
    ctx.font = `${css.fontWeight} ${css.fontSize} ${css.fontFamily}`

    const m = ctx.measureText(letter.value)

    letterEl.value.style.translate = '0 0'

    const baseline = marker.value.getBoundingClientRect().bottom
    const box = root.value.getBoundingClientRect()
    const ink = baseline - (m.actualBoundingBoxAscent - m.actualBoundingBoxDescent) / 2

    letterEl.value.style.translate = `0 ${box.top + box.height / 2 - ink}px`
}

onMounted(() => {
    measure()
    document.fonts?.ready.then(measure)
})

watch([letter, () => props.size], () => nextTick(measure))
</script>

<style scoped>
.avatar {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border-radius: 50%;
    font-weight: bold;
    color: var(--white);
    background: linear-gradient(135deg, var(--primary-300), var(--primary-500));
    -webkit-user-select: none;
    user-select: none;

    &.empty {
        background: transparent;
    }

    .baseline {
        display: inline-block;
        width: 0;
        height: 0;
    }
}
</style>
