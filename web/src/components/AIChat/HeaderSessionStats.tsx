/**
 * HeaderSessionStats - Header 内嵌的会话统计组件
 *
 * PC 端显示三指标：
 * 1. ⏱ 会话持续时间
 * 2. 📦 Block 数量
 * 3. ⚡ Input/Output Token 总数（计费公式）
 *
 * 计费公式：
 * - 计费 Input = 普通 Input + Cache Write + Cache Read × 0.1
 * - 计费 Output = 普通 Output
 *
 * @see docs/dev-guides/AI_CHAT_INTERFACE.md
 */

import { Clock, Package } from "lucide-react";
import { memo, useMemo } from "react";
// import { useTranslation } from "react-i18next"; // Unused - kept for potential future i18n
import { cn } from "@/lib/utils";
import type { AIMode } from "@/types/aichat";
import type { Block as AIBlock } from "@/types/block";

interface SessionStatsData {
  duration: string;
  blockCount: number;
  inputTokens: number;
  outputTokens: number;
  billedTokens: number;
}

interface HeaderSessionStatsProps {
  /** Blocks to aggregate stats from */
  blocks?: AIBlock[];
  /** Current AI mode */
  mode?: AIMode;
  /** Compact display for mobile */
  compact?: boolean;
  /** className for styling */
  className?: string;
}

/**
 * Calculate aggregated session statistics from blocks
 */
const MAX_BLOCKS_TO_PROCESS = 100;

function calculateSessionStats(blocks: AIBlock[] | undefined): SessionStatsData | null {
  if (!blocks || blocks.length === 0) return null;

  // Only process the most recent blocks for performance
  const blocksToProcess = blocks.slice(-MAX_BLOCKS_TO_PROCESS);

  let totalDuration = 0;
  let inputTokens = 0;
  let outputTokens = 0;
  let cacheWriteTokens = 0;
  let cacheReadTokens = 0;

  for (const block of blocksToProcess) {
    // Track total duration
    const created = Number(block.createdTs);
    const updated = Number(block.updatedTs);
    totalDuration += updated - created;

    // Add tokens from sessionStats
    if (block.sessionStats) {
      inputTokens += block.sessionStats.inputTokens || 0;
      outputTokens += block.sessionStats.outputTokens || 0;
      cacheWriteTokens += block.sessionStats.cacheWriteTokens || 0;
      cacheReadTokens += block.sessionStats.cacheReadTokens || 0;
    }
  }

  // Format time
  const formatTime = (ms: number) => {
    if (ms < 1000) return "<1s";
    const s = Math.floor(ms / 1000);
    if (s < 60) return `${s}s`;
    const m = Math.floor(s / 60);
    return `${m}m`;
  };

  // Calculate billed tokens using the formula:
  // Input = 普通 Input + Cache Write + Cache Read × 0.1
  // Output = 普通 Output
  const billedInputTokens = inputTokens + cacheWriteTokens + Math.round(cacheReadTokens * 0.1);
  const billedOutputTokens = outputTokens;
  const totalBilledTokens = billedInputTokens + billedOutputTokens;

  return {
    duration: formatTime(totalDuration),
    blockCount: blocks.length,
    inputTokens: billedInputTokens,
    outputTokens: billedOutputTokens,
    billedTokens: totalBilledTokens,
  };
}

export const HeaderSessionStats = memo(function HeaderSessionStats({ blocks, compact = false, className }: HeaderSessionStatsProps) {
  // mode parameter kept for interface compatibility
  // const { t } = useTranslation();

  const stats = useMemo(() => calculateSessionStats(blocks), [blocks]);

  if (!stats) return null;

  // Desktop: Full stats row with three indicators
  if (!compact) {
    return (
      <div
        className={cn(
          "hidden lg:flex items-center gap-3 text-[11px] font-mono opacity-70 bg-muted/30 px-2 py-1 rounded border border-border/50",
          className,
        )}
      >
        {/* ⏱ 会话持续时间 */}
        <span className="flex items-center gap-1" title="会话持续时间">
          <Clock className="w-3 h-3" />
          <span className="font-medium">{stats.duration}</span>
        </span>

        {/* 📦 Block 数量 */}
        <span className="flex items-center gap-1" title="Block 数量">
          <Package className="w-3 h-3" />
          <span className="font-medium">{stats.blockCount}</span>
        </span>

        {/* ⚡ Token 总数（计费） */}
        <span className="flex items-center gap-1" title="计费 Token 总数">
          <span className="text-amber-500">⚡</span>
          <span className="font-medium">{stats.billedTokens > 0 ? `${(stats.billedTokens / 1000).toFixed(1)}k` : "0"}</span>
        </span>

        {/* 详细 Token 分解（悬停显示） */}
        {stats.billedTokens > 0 && (
          <span className="text-muted-foreground/60 text-[10px]" title={`In: ${stats.inputTokens} / Out: ${stats.outputTokens}`}>
            In/Out
          </span>
        )}
      </div>
    );
  }

  // Mobile: Simplified indicator (仅显示 Token)
  return (
    <div className={cn("flex items-center gap-1 text-[10px] font-mono opacity-80", className)}>
      {stats.billedTokens > 0 && (
        <span className="flex items-center gap-0.5 text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 px-1.5 py-0.5 rounded">
          <span>⚡</span>
          <span className="font-medium">{stats.billedTokens > 0 ? `${(stats.billedTokens / 1000).toFixed(1)}k` : "0"}</span>
        </span>
      )}
    </div>
  );
});

HeaderSessionStats.displayName = "HeaderSessionStats";

export default HeaderSessionStats;
