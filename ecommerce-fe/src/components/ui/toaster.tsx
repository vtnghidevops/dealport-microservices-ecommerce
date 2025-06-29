import { useToast } from "@/hooks/use-toast"
import {
  Toast,
  ToastClose,
  ToastDescription,
  ToastProvider,
  ToastTitle,
  ToastViewport,
} from "@/components/ui/toast"
import {
  CheckCircle,
  AlertCircle,
  AlertTriangle,
  Info,
  X
} from "lucide-react"

export function Toaster() {
  const { toasts } = useToast()

  return (
    <ToastProvider>
      {toasts.map(function ({ id, title, description, action, variant, ...props }) {
        // Select icon based on variant
        let Icon = Info;
        let iconColor = "text-[#3b82f6]"; // Default color for info - Material Blue

        if (variant === "success") {
          Icon = CheckCircle;
          iconColor = "text-[#10b981]"; // Tailwind Emerald 500
        } else if (variant === "destructive" || variant === "error") {
          Icon = X;
          iconColor = "text-[#ef4444]"; // Tailwind Red 500
        } else if (variant === "warning") {
          Icon = AlertTriangle;
          iconColor = "text-[#f59e0b]"; // Tailwind Amber 500
        } else if (variant === "info") {
          Icon = Info;
          iconColor = "text-[#3b82f6]"; // Tailwind Blue 500
        } else if (variant === "default") {
          Icon = AlertCircle;
          iconColor = "text-[#6b7280]"; // Tailwind Gray 500
        }

        return (
          <Toast key={id} {...props} variant={variant}>
            <div className="flex gap-3 w-full">
              <div className="flex items-center">
                <Icon className={`h-5 w-5 ${iconColor} flex-shrink-0`} />
              </div>
              <div className="grid gap-1 flex-1">
                {title && <ToastTitle>{title}</ToastTitle>}
                {description && (
                  <ToastDescription>{description}</ToastDescription>
                )}
              </div>
            </div>
            {action}
            <ToastClose />
          </Toast>
        )
      })}
      <ToastViewport />
    </ToastProvider>
  )
}
