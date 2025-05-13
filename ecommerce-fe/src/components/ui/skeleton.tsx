import { cn } from "@/lib/utils"
import * as React from "react"

const Skeleton = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    className={cn(
      "animate-pulse rounded-md bg-aqua-spring/80 transition-all duration-500 ease-in-out shadow-sm hover:bg-aqua-spring hover:shadow-md transform hover:scale-[1.005]",
      className
    )}
    ref={ref}
    {...props}
  />
))
Skeleton.displayName = "Skeleton"

export { Skeleton }
