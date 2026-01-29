import { Sparkles } from "lucide-react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { CapabilityStatus, CapabilityType } from "@/types/capability";
import { GeekModeToggle } from "./GeekModeToggle";

interface ChatHeaderProps {
  isThinking?: boolean;
  className?: string;
  currentCapability?: CapabilityType;
  capabilityStatus?: CapabilityStatus;
  geekMode?: boolean;
  onGeekModeToggle?: (enabled: boolean) => void;
}

/**
 * Chat Header - 简洁状态显示
 *
 * 设计原则：
 * - 极客模式下保持清爽，不使用过度特效
 * - 通过颜色、字体、状态点传达状态
 * - 移动端和桌面端体验一致
 */

/**
 * 根据当前能力和状态获取动作描述
 */
function getActionDescription(capability: CapabilityType, status: CapabilityStatus, t: (key: string) => string): string | null {
  if (status === "idle") return null;

  if (status === "thinking") {
    return t("ai.thinking");
  }

  if (status === "processing") {
    switch (capability) {
      case CapabilityType.MEMO:
        return t("ai.parrot.status.searching-memos");
      case CapabilityType.SCHEDULE:
        return t("ai.parrot.status.querying-schedule");
      case CapabilityType.AMAZING:
        return t("ai.parrot.status.analyzing");
      default:
        return t("ai.processing");
    }
  }

  return null;
}

export function ChatHeader({
  isThinking = false,
  className,
  currentCapability = CapabilityType.AUTO,
  capabilityStatus = "idle",
  geekMode = false,
  onGeekModeToggle,
}: ChatHeaderProps) {
  const { t } = useTranslation();
  const assistantName = t("ai.assistant-name");
  const actionDescription = getActionDescription(currentCapability, capabilityStatus, t);

  return (
    <header
      className={cn(
        "flex items-center justify-between px-4 h-14 shrink-0 transition-all",
        "border-b border-border/80",
        "bg-background/80 backdrop-blur-sm",
        // Geek mode: subtle green border and background
        geekMode && "border-green-500/20 bg-green-50/50 dark:bg-green-950/20",
        className,
      )}
    >
      {/* Left Section */}
      <div className="flex items-center gap-2.5">
        {/* Avatar with subtle geek mode border */}
        <div className={cn(
          "w-9 h-9 flex items-center justify-center rounded-lg transition-all",
          geekMode && "border border-green-500/30 bg-green-500/10"
        )}>
          <img
            src="/assistant-avatar.webp"
            alt={assistantName}
            className="h-9 w-auto object-contain"
          />
        </div>
        <div className="flex flex-col">
          <h1 className={cn(
            "font-semibold text-foreground text-sm leading-tight",
            geekMode && "font-mono"
          )}>
            {assistantName}
          </h1>
          {/* Status */}
          {actionDescription ? (
            <span className={cn(
              "text-xs flex items-center gap-1.5",
              geekMode ? "text-green-600 dark:text-green-400" : "text-primary"
            )}>
              {geekMode && (
                <span className="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse" />
              )}
              {actionDescription}
            </span>
          ) : (
            <span className={cn(
              "text-xs text-muted-foreground",
              geekMode && "font-mono"
            )}>
              {geekMode ? "$ ready" : t("ai.ready")}
            </span>
          )}
        </div>
      </div>

      {/* Right Section - Geek Mode Toggle + Thinking indicator */}
      <div className="flex items-center gap-2">
        <GeekModeToggle
          enabled={geekMode}
          onToggle={onGeekModeToggle ?? (() => {})}
          variant="header"
        />
        {isThinking && (
          <div className="flex items-center gap-1.5 text-sm">
            <Sparkles className={cn(
              "w-4 h-4 animate-pulse",
              geekMode ? "text-green-600 dark:text-green-400" : "text-primary"
            )} />
          </div>
        )}
      </div>
    </header>
  );
}
