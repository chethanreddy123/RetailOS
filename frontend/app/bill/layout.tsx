import { Providers } from '@/app/providers'

export default function BillLayout({ children }: { children: React.ReactNode }) {
  return <Providers>{children}</Providers>
}
