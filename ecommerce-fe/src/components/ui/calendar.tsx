import * as React from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"

export type CalendarProps = {
  className?: string
  classNames?: Record<string, string>
  showOutsideDays?: boolean
  mode?: "single" | "range" | "multiple"
  selected?: Date | null
  onSelect?: (date: Date | null) => void
  disabled?: boolean
  footer?: React.ReactNode
}

/**
 * Basic calendar component that replaces react-day-picker to avoid dependency issues.
 * This is a simplified version that just renders a static calendar UI without actual date picking functionality.
 */
function Calendar({
  className,
  classNames,
  showOutsideDays = true,
  ...props
}: CalendarProps) {
  // Current month information
  const today = new Date()
  const currentMonth = today.getMonth()
  const currentYear = today.getFullYear()

  // Month name
  const monthName = new Intl.DateTimeFormat('en-US', { month: 'long' }).format(today)

  // Day headers
  const weekdays = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"]

  return (
    <div className={cn("p-3", className)}>
      <div className="flex justify-center pt-1 relative items-center">
        <button className={cn(
          buttonVariants({ variant: "outline" }),
          "h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100 absolute left-1"
        )}>
          <ChevronLeft className="h-4 w-4" />
        </button>
        <div className="text-sm font-medium">
          {monthName} {currentYear}
        </div>
        <button className={cn(
          buttonVariants({ variant: "outline" }),
          "h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100 absolute right-1"
        )}>
          <ChevronRight className="h-4 w-4" />
        </button>
      </div>

      <div className="mt-4">
        {/* Weekday headers */}
        <div className="flex w-full">
          {weekdays.map(day => (
            <div key={day} className="text-muted-foreground rounded-md w-8 font-normal text-[0.8rem] text-center">
              {day}
            </div>
          ))}
        </div>

        {/* Calendar grid - simplified static version */}
        {Array.from({ length: 5 }).map((_, weekIndex) => (
          <div key={weekIndex} className="flex w-full mt-2">
            {Array.from({ length: 7 }).map((_, dayIndex) => {
              const date = weekIndex * 7 + dayIndex + 1
              const isCurrentDay = date === today.getDate() && currentMonth === today.getMonth()
              const isOutsideMonth = date > 31

              return (
                <div
                  key={dayIndex}
                  className={cn(
                    "relative p-0 text-center text-sm",
                    isOutsideMonth && "invisible"
                  )}
                >
                  <button
                    className={cn(
                      buttonVariants({ variant: "ghost" }),
                      "h-8 w-8 p-0 font-normal",
                      isCurrentDay && "bg-accent text-accent-foreground"
                    )}
                    disabled={isOutsideMonth || props.disabled}
                  >
                    {isOutsideMonth ? "" : date}
                  </button>
                </div>
              )
            })}
          </div>
        ))}
      </div>

      {props.footer && (
        <div className="mt-4">
          {props.footer}
        </div>
      )}
    </div>
  )
}

Calendar.displayName = "Calendar"

export { Calendar } 