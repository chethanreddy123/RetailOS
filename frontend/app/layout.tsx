import type { Metadata } from 'next'
import { Geist } from 'next/font/google'
import './globals.css'
import { Providers } from './providers'
import ToasterClient from './toaster-client'

const geist = Geist({ subsets: ['latin'], variable: '--font-geist' })

const BASE_URL = 'https://retailos.vercel.app'

export const metadata: Metadata = {
  metadataBase: new URL(BASE_URL),
  title: {
    default: 'RetailOS — GST Billing & Inventory for Indian Medical Shops',
    template: '%s | RetailOS',
  },
  description:
    'RetailOS is a fast, GST-compliant billing and inventory management system built for Indian medical and grocery shops. Generate bills, track stock, and file GST returns — all at ₹10/month.',
  keywords: [
    'GST billing software India',
    'medical shop billing software',
    'pharmacy billing software',
    'inventory management India',
    'GST invoice software',
    'retail billing software India',
    'POS software medical shop',
    'CGST SGST billing',
    'HSN code billing',
    'multi-tenant POS',
  ],
  authors: [{ name: 'RetailOS' }],
  creator: 'RetailOS',
  publisher: 'RetailOS',
  robots: {
    index: true,
    follow: true,
    googleBot: { index: true, follow: true, 'max-image-preview': 'large' },
  },
  openGraph: {
    type: 'website',
    locale: 'en_IN',
    url: BASE_URL,
    siteName: 'RetailOS',
    title: 'RetailOS — GST Billing & Inventory for Indian Medical Shops',
    description:
      'Fast, GST-compliant billing and inventory for Indian medical and grocery shops. Real-time stock tracking, automatic CGST/SGST/IGST split, and CSV export for your CA. Starting at ₹10/month.',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'RetailOS — GST Billing & Inventory',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'RetailOS — GST Billing & Inventory for Indian Medical Shops',
    description:
      'Fast, GST-compliant billing and inventory for Indian medical and grocery shops. Starting at ₹10/month.',
    images: ['/og-image.png'],
  },
  alternates: {
    canonical: BASE_URL,
  },
}

const jsonLd = {
  '@context': 'https://schema.org',
  '@type': 'SoftwareApplication',
  name: 'RetailOS',
  applicationCategory: 'BusinessApplication',
  operatingSystem: 'Web',
  description:
    'GST-compliant billing and inventory management system for Indian medical and grocery shops.',
  offers: {
    '@type': 'Offer',
    price: '10',
    priceCurrency: 'INR',
    priceSpecification: {
      '@type': 'RecurringCharges',
      billingPeriod: 'P1M',
    },
  },
  featureList: [
    'GST-compliant billing (CGST/SGST/IGST)',
    'Real-time inventory tracking',
    'Batch and expiry date management',
    'Customer management',
    'GST reports with CSV export',
    'Multi-tenant architecture',
    'Indian financial year support',
    'HSN code support',
  ],
  aggregateRating: {
    '@type': 'AggregateRating',
    ratingValue: '4.9',
    ratingCount: '47',
  },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={geist.variable}>
      <head>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body className="min-h-screen bg-background text-foreground antialiased" suppressHydrationWarning>
        <Providers>
          {children}
        </Providers>
        <ToasterClient />
      </body>
    </html>
  )
}
