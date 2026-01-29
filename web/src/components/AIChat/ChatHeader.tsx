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
 * UX/UI 设计原则：
 * - 仅展示助手信息和状态
 * - 工具按钮移至输入框工具栏
 * - 简洁清晰的视觉层次
 *
 * Geek Mode 变化：
 * - 头部背景变为终端深色
 * - 底部边框出现绿色渐变光效
 * - 助手名称使用等宽字体
 * - 状态指示器变为绿色闪烁点
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
        "flex items-center justify-between px-4 h-14 shrink-0",
        "border-b border-border/80",
        "bg-background/80 backdrop-blur-sm",
        // Geek mode header styles
        geekMode && "geek-header",
        className,
      )}
    >
      {/* Left Section */}
      <div className="flex items-center gap-2.5">
        <div className={cn(
          "w-9 h-9 flex items-center justify-center",
          geekMode && "geek-border rounded"
        )}>
          <img
            src="/assistant-avatar.webp"
            alt={assistantName}
            className={cn(
              "h-9 w-auto object-contain",
              geekMode && "brightness-110 contrast-125"
            )}
          />
        </div>
        <div className="flex flex-col">
          <h1 className={cn(
            "font-semibold text-foreground text-sm leading-tight",
            geekMode && "geek-mono"
          )}>
            {geekMode ? `[${assistantName}]` : assistantName}
          </h1>
          {/* 动作描述 - 替代能力徽章 */}
          {actionDescription ? (
            <span className={cn(
              "text-xs flex items-center gap-1",
              geekMode ? "geek-text" : "text-primary"
            )}>
              {geekMode ? (
                <span className="geek-status-dot active" />
              ) : (
                <span className="w-1.5 h-1.5 rounded-full bg-current animate-pulse" />
              )}
              {actionDescription}
            </span>
          ) : (
            <span className={cn(
              "text-xs",
              geekMode ? "geek-mono text-muted-foreground" : "text-muted-foreground"
            )}>
              {geekMode ? "$ ready" : t("ai.ready")}
            </span>
          )}
        </div>
      </div>

      {/* Right Section - Geek Mode Toggle + Status indicator */}
      <div className="flex items-center gap-2">
        <GeekModeToggle
          enabled={geekMode}
          onToggle={onGeekModeToggle ?? (() => {})}
          variant="header"
        />
        {isThinking && (
          <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <Sparkles className={cn(
              "w-4 h-4",
              geekMode ? "geek-text-glow animate-pulse" : "animate-pulse text-primary"
            )} />
          </div>
        )}
      </div>
    </header>
  );
}
