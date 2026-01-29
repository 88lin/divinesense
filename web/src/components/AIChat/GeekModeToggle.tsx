import { Terminal, X } from "lucide-react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface GeekModeToggleProps {
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
  disabled?: boolean;
  variant?: "header" | "toolbar" | "mobile";
}

/**
 * Geek Mode Toggle Component
 *
 * Clear state indication:
 * - Not enabled: Shows "极客" + Terminal icon (click to enter geek mode)
 * - Enabled: Shows "退出极客" + X icon (click to exit geek mode)
 *
 * Design:
 * - Minimal and clean, no overwhelming effects
 * - Green color when active to indicate state
 * - Consistent across all variants
 */
export function GeekModeToggle({ enabled, onToggle, disabled = false, variant = "header" }: GeekModeToggleProps) {
  const { t } = useTranslation();

  const isToolbar = variant === "toolbar";
  const isHeader = variant === "header";

  // Label based on state
  const label = enabled ? t("ai.geek_mode.disable_label") : t("ai.geek_mode.label");

  return (
    <button
      onClick={() => onToggle(!enabled)}
      disabled={disabled}
      className={cn(
        "relative flex items-center justify-center transition-all",
        "disabled:opacity-50 disabled:cursor-not-allowed",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        // Toolbar variant - compact, always shows label
        isToolbar && [
          "h-7 px-2 text-xs",
          !enabled
            ? "text-muted-foreground hover:text-foreground hover:bg-muted"
            : "text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300 hover:bg-green-500/10",
        ],
        // Header variant - shows label on larger screens
        isHeader && [
          "flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium",
          "hover:bg-muted",
          !enabled ? "text-muted-foreground" : "text-green-600 dark:text-green-400",
        ],
      )}
      title={enabled ? t("ai.geek_mode.enabled") : t("ai.geek_mode.tooltip")}
      aria-pressed={enabled}
    >
      {/* Icon - Terminal for entering, X for exiting */}
      <span className="shrink-0">
        {enabled ? (
          <X className={cn("w-4 h-4", isToolbar && "w-3.5 h-3.5")} />
        ) : (
          <Terminal className={cn("w-4 h-4", isToolbar && "w-3.5 h-3.5")} />
        )}
      </span>

      {/* Label */}
      {isToolbar || isHeader ? (
        <span className={cn(isHeader && "hidden sm:inline", isToolbar && "ml-1 whitespace-nowrap")}>{label}</span>
      ) : null}
    </button>
  );
}
