// util/qrcode.ts
import QRCodeStyling from 'qr-code-styling'
import panelLogo from '@/asset/image/logo.webp'

interface QRCodeOptions {
    value: string
    logo?: string
    size?: number
}

export async function buildQRCode({ value, logo = panelLogo, size = 220 }: QRCodeOptions): Promise<string> {
    const scale = window.devicePixelRatio || 2

    const qr = new QRCodeStyling({
        width: size * scale,
        height: size * scale,
        margin: 5 * scale,
        data: value,
        image: logo,
        dotsOptions: { color: '#000000', type: 'square', roundSize: false },
        backgroundOptions: { color: '#ffffff' },
        imageOptions: { crossOrigin: 'anonymous', margin: 5 * scale, imageSize: 0.2 },
        qrOptions: { errorCorrectionLevel: 'H' },
    })

    const blob = await qr.getRawData('png')
    return URL.createObjectURL(blob as Blob)
}