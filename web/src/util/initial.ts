// util/initial.ts
// 取文字的首字符，用来按键定位：英文字母转大写，数字原样，中文取拼音首字母，其他返回空串

// 中文按拼音排序，落在哪个声母分界字之后就是哪个字母(没有 I、U、V)
const letters = 'ABCDEFGHJKLMNOPQRSTWXYZ'
const marks = '阿八嚓哒妸发旮哈击喀垃妈拿噢啪七然仨他挖夕丫帀'
const collator = new Intl.Collator('zh-CN')

export function initial(label: string): string {
    const char = [...label.trim()][0] ?? ''

    if (/^[a-z]$/i.test(char)) return char.toUpperCase()

    if (/^[0-9]$/.test(char)) return char

    if (!/^[一-龥]$/.test(char)) return ''

    let index = -1

    for (let i = 0; i < marks.length; i++) {
        if (collator.compare(char, marks[i]!) < 0) break

        index = i
    }

    return letters[index] ?? ''
}
