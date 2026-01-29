import { Terminal, Cpu } from "lucide-react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface MobileGeekModeToggleProps {
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
}

/**
 * Mobile Geek Mode Toggle - 移动端极客模式开关
 *
 * A floating action button (FAB) style toggle for mobile devices.
 * Shows Geek Mode status with visual feedback.
 *
 * When enabled:
 * - Green glow border
 * - Terminal icon instead of CPU
 * - Pulsing animation
 * - Scanline effect on touch
 */
export function MobileGeekModeToggle({
  enabled,
  onToggle,
}: MobileGeekModeToggleProps) {
  const { t } = useTranslation();

  return (
    <button
      onClick={() => onToggle(!enabled)}
      className={cn(
        // Base styles
        "fixed bottom-6 right-6 z-50",
        "w-14 h-14 rounded-full",
        "flex items-center justify-center",
        "transition-all duration-300",
        "active:scale-95",
        // Background
        "bg-background border-2",
        // Geek mode styles
        enabled && "geek-fab",
        !enabled && "border-border shadow-lg",
      )}
      title={enabled ? t("ai.geek_mode.enabled") : t("ai.geek_mode.tooltip")}
      aria-pressed={enabled}
      aria-label={enabled ? t("ai.geek_mode.disable_label") : t("ai.geek_mode.enable_label")}
    >
      {/* Scanline effect when enabled */}
      {enabled && (
        <span className="absolute inset-0 rounded-full overflow-hidden opacity-30 geek-scanlines" />
      )}

      {/* Icon */}
      <span className={cn("relative z-10", enabled && "geek-text-glow")}>
        {enabled ? (
          <Terminal className="w-6 h-6 animate-pulse" />
        ) : (
          <Cpu className="w-6 h-6 text-muted-foreground" />
        )}
      </span>

      {/* Status ring */}
      {enabled && (
        <span className="absolute -inset-1 rounded-full border border-primary/20 animate-ping" />
      )}
    </button>
  );
}
