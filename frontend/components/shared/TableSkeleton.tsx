import { Skeleton } from '@/components/ui/skeleton'

export default function TableSkeleton({ cols = 5, rows = 8 }: { cols?: number; rows?: number }) {
  return (
    <div className="rounded-lg border border-[#EBEBEB] bg-white overflow-hidden">
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className={`flex gap-4 px-4 py-3 ${i < rows - 1 ? 'border-b border-[#F5F5F5]' : ''}`}>
          {Array.from({ length: cols }).map((_, j) => (
            <Skeleton key={j} className="h-3.5 flex-1 bg-[#F2F2F2]" style={{ maxWidth: j === 0 ? 180 : undefined }} />
          ))}
        </div>
      ))}
    </div>
  )
}
