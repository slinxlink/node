<template>
    <div class="datetime" :class="{ disabled }">
        <VueDatePicker
            ref="picker"
            :model-value="value || null"
            model-type="yyyy-MM-dd HH:mm:ss"
            :formats="{ input: 'yyyy-MM-dd HH:mm:ss' }"
            :time-config="{ enableSeconds: true }"
            :aria-labels="{ closeTimePicker: 'datetime-close-time', prevMonth: '上一月', nextMonth: '下一月' }"
            :config="{ monthChangeOnScroll: false }"
            :locale="zhCN"
            :dark="resolved === 'dark'"
            :placeholder="placeholder"
            :disabled="disabled"
            :input-attrs="{ hideInputIcon: true, clearable: false }"
            :six-weeks="true"
            :text-input="true"
            :teleport="true"
            @update:model-value="emit('update:value', $event ?? '')"
            @overlay-toggle="onOverlay"
            @internal-model-change="pending = $event"
            @closed="onClosed"
            @menu-mounted="(menu) => menu.addEventListener('wheel', onWheel, { passive: false })"
            @menu-unmounted="(menu) => menu.removeEventListener('wheel', onWheel)"
        >
            <template #action-extra>
                <div class="datetime-preview">{{ previewText }}</div>
            </template>
            <template #action-row="{ closePicker, selectDate }">
                <button type="button" class="dp--action-button datetime-switch" @click="picker?.switchView(timeOpen ? 'calendar' : 'time')">
                    <Icon :name="timeOpen ? 'calendar_month' : 'schedule'" />
                </button>
                <div class="dp--action-buttons">
                    <button type="button" class="dp--action-button dp--action-cancel" @click="closePicker">取消</button>
                    <button type="button" class="dp--action-button dp--action-cancel" @click="picker?.updateInternalModelValue(new Date())">此刻</button>
                    <button type="button" class="dp--action-button dp--action-select" @click="selectDate">确定</button>
                </div>
            </template>
        </VueDatePicker>
        <span class="trailing">
            <Button variant="text" :tabindex="-1" :disabled="disabled" @click="value ? emit('update:value', '') : picker?.openMenu()">
                <span class="swap">
                    <Icon name="calendar_month" :class="{ show: !value }" />
                    <Icon name="cancel" :class="{ show: !!value }" />
                </span>
            </Button>
        </span>
    </div>
</template>

<script setup lang="ts">
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import { resolved } from '@/composable/theme'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'

// 日期时间选择，精确到秒;值是 "2026-11-04 12:00:00" 这样的文本，空表示没有选
// 用库而不用浏览器自带的 datetime-local，是因为 Safari 的原生选择器选不了时间
// 库自带的日历图标和清除按钮都关掉，换成右边同一个位置的按钮：没有值是日历(点开面板)，有值变成清除
// 面板底部的操作栏：左下角是日期和时间的切换，右边是取消、此刻、确定；库自带的两个切换条用样式藏掉
const props = defineProps<{
    value: string
    placeholder?: string
    disabled?: boolean
}>()

const emit = defineEmits<{
    'update:value': [value: string]
}>()

const picker = ref<InstanceType<typeof VueDatePicker>>()

// 当前是不是在选时间：左下角的切换按钮据此显示成日历或时钟
const timeOpen = ref(false)

// 预览：还没点确定的当前选择(库的 internal-model-change 事件给出)；没动过就显示已有的值，没有就是空的
const pending = ref<Date | Date[] | null>()

const previewText = computed(() => {
    if (pending.value === undefined) return props.value

    return pending.value instanceof Date ? format(pending.value, 'yyyy-MM-dd HH:mm:ss') : ''
})

// 面板关了，当前选择作废，回到已有的值
function onClosed() {
    timeOpen.value = false
    pending.value = undefined
}

function onOverlay({ open, overlay }: { open: boolean; overlay: string }) {
    if (overlay === 'time') timeOpen.value = open
}

// 滚轮翻月：库自己的节流太松（滚得稍快就是每 100ms 翻一次，触摸板一划能翻好几个月），关掉它，自己按累计滚动量来
// 累计够 SCROLL_STEP 才翻一个月，翻完冷却 SCROLL_COOLDOWN 毫秒，冷却里的滚动（含触摸板的惯性尾巴）直接丢掉；向下滚是下个月
const SCROLL_STEP = 100
const SCROLL_COOLDOWN = 400

let scrolled = 0
let flippedAt = 0

function onWheel(event: WheelEvent) {
    const menu = event.currentTarget as HTMLElement

    // 只管日历格子上的滚动；选时间、月、年的面板(在 DOM 里是日历的子元素)有自己的滚动，要放行
    const target = event.target as HTMLElement

    if (!target.closest('.dp--instance-calendar') || target.closest('.dp--overlay')) return

    event.preventDefault()

    if (Date.now() - flippedAt < SCROLL_COOLDOWN) {
        scrolled = 0
        return
    }

    // Firefox 的滚轮按行报告(一格只有 3 左右)，行和页的单位一格就够翻一个月
    scrolled += event.deltaMode === WheelEvent.DOM_DELTA_PIXEL ? event.deltaY : Math.sign(event.deltaY) * SCROLL_STEP

    if (Math.abs(scrolled) < SCROLL_STEP) return

    menu.querySelector<HTMLElement>(`[aria-label='${scrolled > 0 ? '下一月' : '上一月'}']`)?.click()

    scrolled = 0
    flippedAt = Date.now()
}
</script>

<style scoped>
.datetime {
    position: relative;
    width: 100%;

    &.disabled {
        opacity: 0.5;
    }

    /* 输入框和别的输入框一样：圆角胶囊、34px 高、没有边框 */
    &:deep(.dp--input) {
        height: 34px;
        padding: 0 35px 0 15px;
        border: none;
        border-radius: 999px;
        background-color: var(--background-200);
        font-size: 14px;
        line-height: 1;
        transition: background-color 0.3s ease-in-out;

        &:hover,
        &:focus {
            background-color: var(--background-100);
        }
    }

    /* 和 Input 右侧的按钮同一个位置 */
    .trailing {
        position: absolute;
        top: 50%;
        right: 10px;
        transform: translateY(-50%);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 20px;

        /* 日历和清除两个图标叠在同一个位置，有没有值切换时 0.3 秒交叉淡入 */
        .swap {
            position: relative;
            width: 20px;
            height: 20px;

            .icon {
                position: absolute;
                top: 0;
                left: 0;
                opacity: 0;
                transition: opacity 0.3s ease-in-out;

                &.show {
                    opacity: 1;
                }
            }
        }
    }
}
</style>

<!-- 弹出的面板挂在 body 下，scoped 样式够不到，这里用本项目的配色令牌盖掉库的主题变量 -->
<style>
body .dp--theme-dark,
body .dp--theme-light {
    --dp-font-family: inherit;
    --dp-font-size: 14px;
    --dp-border-radius: 20px;
    --dp-cell-border-radius: 10px;
    --dp-menu-min-width: 293px;
    --dp-background-color: var(--background-100);
    --dp-text-color: var(--text-100);
    --dp-icon-color: var(--text-300);
    --dp-secondary-color: var(--text-300);
    --dp-hover-color: var(--primary-500);
    --dp-hover-text-color: var(--text-100);
    --dp-hover-icon-color: var(--text-100);
    --dp-primary-color: var(--primary-300);
    --dp-primary-disabled-color: var(--primary-600);
    --dp-primary-text-color: var(--white);
    --dp-border-color: transparent;
    --dp-menu-border-color: transparent;
    --dp-border-color-hover: transparent;
    --dp-border-color-focus: transparent;
    --dp-disabled-color: var(--background-200);
    --dp-disabled-color-text: var(--text-300);
    --dp-scroll-bar-background: var(--background-100);
    --dp-scroll-bar-color: var(--text-300);
    --dp-highlight-color: color-mix(in srgb, var(--primary-300) 20%, transparent);

    /* 所有动画统一 0.3 秒：悬停变色、操作栏按钮、菜单出现和消失、月年面板滑动、上下月翻页 */
    --dp-common-transition: all 0.3s ease-in-out;
    --dp-action-row-transition: all 0.3s ease-in-out;
    --dp-animation-duration: 0.3s;
    --dp-transition-timing-general: 0.3s;
}

/* 日期格子横向也留 5px 间距，和上下行之间的 5px 一致：菜单加宽到 7 格(35)加 6 个间距(5)，格子两端对齐、中间等距；
   表头的星期用同样的排法，才对得上格子 */
body .dp--calendar-header,
body .dp--calendar-row {
    justify-content: space-between;
}

body .dp--calendar-header-item,
body .dp--calendar-item {
    flex-grow: 0;
}

/* 行与行之间的 5px 是每行的上下外边距叠出来的，最后一行的下外边距就成了多出来的一截；去掉 */
body .dp--calendar-row:last-child {
    margin-bottom: 0;
}

/* 库自带的日期/时间切换条(居中那一整条)不要，切换放到操作栏左边 */
body [data-test-id='open-time-picker-btn'],
body [aria-label='datetime-close-time'],
body .dp--overlay-action {
    display: none;
}

/* 预览在库的 action-extra 容器里：日历(或时间)的下面、操作栏的上面，居中
   上下各 14px：菜单内边距 6px、操作栏内边距 8px 已经占了一部分，这里补 8px 和 6px；
   高度固定 16px，没有值时也占着，免得第一次选了之后整个框往下撑 */
body .dp--action-extra {
    padding: 8px 0 6px;
}

body .datetime-preview {
    height: 16px;
    line-height: 16px;
    text-align: center;
    font-size: var(--dp-preview-font-size);
    font-variant-numeric: tabular-nums;
}

/* 面板底部的操作栏用库自己的类名画(dp--action-*),样式和库默认的一样；只有最左边的切换按钮是加的：图标取暗色，悬停变亮 */
body .datetime-switch {
    margin-inline-start: 0;
    padding-inline: 0;
    color: var(--dp-icon-color);
    font-size: 20px;
    transition: color 0.3s ease-in-out;

    &:hover {
        color: var(--dp-text-color);
    }
}

/* 按下的颜色：和别处一样用 --primary-600（悬停是 --primary-500），文字变白；库里有悬停效果的都加上 */
body .dp--date-hoverable:active,
body .dp--inner-nav:active,
body .dp--month-year-select-base:active,
body .dp--overlay-cell:active,
body .dp--inc-dec-button:active,
body .dp--time-display:active:enabled {
    background: var(--primary-600);
    color: var(--white);
}

/* 时、分、秒三个大按钮：库默认用菜单的 20px 圆角，按钮才 40 多像素，圆得像药丸；改成和日期格子一样的圆角 */
body .dp--time-display {
    border-radius: var(--dp-cell-border-radius);
}

/* 取消、此刻：悬停时库给的是主题色边框，按下就用同一个颜色填满 */
body .dp--action-cancel:active {
    background: var(--dp-primary-color);
    border-color: var(--dp-primary-color);
    color: var(--white);
}

/* 日期/时间切换是个图标按钮，按下只变图标的颜色 */
body .datetime-switch:active {
    color: var(--primary-300);
}

/* 确定按钮本来就是主题色底，按下变深 */
body .dp--action-buttons .dp--action-select:active {
    background: var(--primary-600);
}

/* 选时间、月、年的面板盖在菜单上面，自己没有圆角，上面两个角会从菜单的圆角里露出来；下面接着操作栏，保持直角 */
body .dp--overlay {
    overflow: hidden;
    padding: 5px 5px 0;
    border-radius: var(--dp-border-radius) var(--dp-border-radius) 0 0;
}

/* 和 Select 的下拉一样：面板留 5px 内边距，里面的滚动区自己有圆角，滚出去的内容沿着圆角线被裁掉，不会像被硬切一刀；
   滚动区的高度要减掉上面这 5px，不然底部会被顶出去
   选时间的面板也带着 dp--overlay-container 这个类，但它靠 height: 100% 撑满并把数字居中，不能改它 */
body .dp--overlay-container:not(.dp--time-picker-overlay-container) {
    height: calc(var(--dp-overlay-height) - 5px);
    border-radius: 15px 15px 0 0;
}

/* 选时间、月、年的面板进出只淡入淡出，不要库默认的 22px 上下位移：位移时面板的背景会从菜单边缘露出来 */
body .dp-slide-up-enter-from,
body .dp-slide-up-leave-to,
body .dp-slide-down-enter-from,
body .dp-slide-down-leave-to {
    transform: none;
}

/* 日期格子的悬停变色库里只在移进去时有 0.2 秒动画，移出是瞬间的；统一成进出都是 0.3 秒 */
body .dp--cell-inner,
body .dp--cell-inner:hover {
    transition:
        background-color 0.3s ease-in-out,
        color 0.3s ease-in-out,
        border-color 0.3s ease-in-out;
}
</style>
