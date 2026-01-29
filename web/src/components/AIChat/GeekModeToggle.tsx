import { Cpu } from "lucide-react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface GeekModeToggleProps {
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
  disabled?: boolean;
  variant?: "header" | "toolbar";
}

/**
 * Geek Mode Toggle Component
 *
 * Allows users to enable/disable Geek Mode for code-related tasks.
 * When enabled, the backend can use Claude Code CLI headless mode.
 *
 * Variants:
 * - "header": Shown in ChatHeader (default)
 * - "toolbar": Shown in input toolbar
 */
export function GeekModeToggle({
  enabled,
  onToggle,
  disabled = false,
  variant = "header",
}: GeekModeToggleProps) {
  const { t } = useTranslation();

  const isHeader = variant === "header";

  return (
    <button
      onClick={() => onToggle(!enabled)}
      disabled={disabled}
      className={cn(
        "flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium transition-all",
        "hover:bg-muted active:scale-[0.98]",
        "disabled:opacity-50 disabled:cursor-not-allowed",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        enabled ? "bg-primary/10 text-primary" : "text-muted-foreground",
      )}
      title={t("ai.geek_mode.tooltip")}
      aria-pressed={enabled}
      aria-label={t("ai.geek_mode.aria_label")}
    >
      <Cpu className={cn("w-4 h-4", enabled && "animate-pulse")} />
      <span className={cn(isHeader ? "hidden sm:inline" : "")}>
        {t("ai.geek_mode.label")}
      </span>
      {enabled && <span className="w-1.5 h-1.5 rounded-full bg-current" />}
    </button>
  );
}
