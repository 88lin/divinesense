import { FileText, TrendingUp } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { cn } from "@/lib/utils";
import { MemoQueryResultData } from "@/types/parrot";

interface MemoQueryResultProps {
  result: MemoQueryResultData;
  className?: string;
}

export function MemoQueryResult({ result, className }: MemoQueryResultProps) {
  const { t } = useTranslation();
  const { memos, query, count } = result;

  // Rerank 已排好序，直接使用
  const sortedMemos = memos;

  if (count === 0) {
    return (
      <div className={cn("flex flex-col items-center justify-center py-8 px-4 rounded-lg bg-muted border border-border", className)}>
        <FileText className="w-12 h-12 text-muted-foreground mb-3" />
        <p className="text-sm font-medium text-foreground">{t("ai.memo-query.no-results")}</p>
        <p className="text-xs text-muted-foreground mt-1">
          {t("ai.memo-query.query-label")}: "{query}"
        </p>
      </div>
    );
  }

  return (
    <div className={cn("rounded-lg bg-muted border border-border overflow-hidden", className)}>
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 bg-card border-b border-border">
        <div className="flex items-center space-x-2">
          <FileText className="w-5 h-5 text-blue-600 dark:text-blue-400" />
          <div>
            <h3 className="font-semibold text-sm text-foreground">{t("ai.memo-query.results-title")}</h3>
            <p className="text-xs text-muted-foreground">
              {t("ai.memo-query.query-label")}: "{query}" · {t("ai.memo-query.found-count", { count })}
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-1 px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/20">
          <TrendingUp className="w-4 h-4 text-blue-600 dark:text-blue-400" />
          <span className="text-xs font-medium text-blue-700 dark:text-blue-300">{t("ai.memo-query.sorted")}</span>
        </div>
      </div>

      {/* Results List */}
      <div className="divide-y divide-border">
        {sortedMemos.map((memo, index) => (
          <MemoQueryResultItem key={memo.uid} memo={memo} rank={index + 1} />
        ))}
      </div>
    </div>
  );
}

interface MemoQueryResultItemProps {
  memo: {
    uid: string;
    content: string;
    score: number;
  };
  rank: number;
}

function MemoQueryResultItem({ memo, rank }: MemoQueryResultItemProps) {
  return (
    <Link to={`/memo/${memo.uid}`} className="block px-4 py-3 hover:bg-muted transition-colors">
      <div className="flex items-start justify-between space-x-3">
        {/* Rank Badge */}
        <div
          className={cn(
            "flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold",
            rank <= 3 ? "bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300" : "bg-muted text-muted-foreground",
          )}
        >
          {rank}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <p className="text-sm text-foreground line-clamp-2">{memo.content}</p>
        </div>
      </div>
    </Link>
  );
}
