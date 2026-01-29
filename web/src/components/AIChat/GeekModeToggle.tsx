import { Cpu, Terminal } from "lucide-react";
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
 * Allows users to enable/disable Geek Mode for code-related tasks.
 * When enabled, the backend can use Claude Code CLI headless mode.
 *
 * Visual Changes when enabled:
 * - Green terminal glow effect
 * - Animated border
 * - CPU icon transforms to Terminal icon
 * - Monospace font for label
 *
 * Variants:
 * - "header": Shown in ChatHeader (default)
 * - "toolbar": Shown in input toolbar
 * - "mobile": Mobile-optimized floating style
 */
export function GeekModeToggle({
  enabled,
  onToggle,
  disabled = false,
  variant = "header",
}: GeekModeToggleProps) {
  const { t } = useTranslation();

  const isHeader = variant === "header";
  const isMobile = variant === "mobile";

  return (
    <button
      onClick={() => onToggle(!enabled)}
      disabled={disabled}
      className={cn(
        "relative overflow-hidden flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium transition-all",
        "hover:bg-muted active:scale-[0.98]",
        "disabled:opacity-50 disabled:cursor-not-allowed",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        // Geek mode styles
        enabled && "geek-border geek-glow-pulse",
        // Base colors
        enabled ? "text-primary" : "text-muted-foreground",
        // Mobile specific
        isMobile && "px-3 py-2",
      )}
      title={enabled ? t("ai.geek_mode.enabled") : t("ai.geek_mode.tooltip")}
      aria-pressed={enabled}
      aria-label={t("ai.geek_mode.aria_label")}
    >
      {/* Scanline effect overlay when enabled */}
      {enabled && (
        <span className="absolute inset-0 pointer-events-none opacity-20 geek-scanlines" />
      )}

      {/* Icon - switches from Cpu to Terminal when enabled */}
      <span className={cn("relative z-10", enabled && "geek-text-glow")}>
        {enabled ? (
          <Terminal className={cn("w-4 h-4", enabled && "animate-pulse")} />
        ) : (
          <Cpu className="w-4 h-4" />
        )}
      </span>

      {/* Label - monospace in geek mode */}
      <span
        className={cn(
          "relative z-10",
          isHeader ? "hidden sm:inline" : "",
          isMobile ? "" : "hidden md:inline",
          enabled && "geek-mono"
        )}
      >
        {t("ai.geek_mode.label")}
      </span>

      {/* Status indicator - becomes a "terminal cursor" when enabled */}
      {enabled && (
        <span className="relative z-10 geek-cursor ml-0.5" />
      )}
    </button>
  );
}
