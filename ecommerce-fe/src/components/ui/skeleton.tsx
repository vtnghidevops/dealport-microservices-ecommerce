import { cn } from "@/lib/utils"
import * as React from "react"

const Skeleton = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    className={cn(
      "animate-pulse rounded-md bg-primary/10 transition-opacity duration-200 ease-in-out",
      className
    )}
    ref={ref}
    {...props}
  />
))
Skeleton.displayName = "Skeleton"

export { Skeleton }
