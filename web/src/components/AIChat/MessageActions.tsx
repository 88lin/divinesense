import { RefreshCwIcon, TrashIcon } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

interface MessageActionsProps {
  onRegenerate: () => void;
  onDelete: () => void;
}

const MessageActions = ({ onRegenerate, onDelete }: MessageActionsProps) => {
  const { t } = useTranslation();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 min-w-[44px] min-h-[44px] opacity-0 group-hover:opacity-50 hover:opacity-100 transition-opacity"
          aria-label={t("ai.more-options")}
        >
          <svg className="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="12" cy="12" r="1" />
            <circle cx="12" cy="5" r="1" />
            <circle cx="12" cy="19" r="1" />
          </svg>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" sideOffset={4}>
        <DropdownMenuItem onClick={onRegenerate} className="cursor-pointer">
          <RefreshCwIcon className="w-4 h-4 mr-2" />
          {t("ai.regenerate")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={onDelete} className="text-destructive focus:text-destructive cursor-pointer">
          <TrashIcon className="w-4 h-4 mr-2" />
          {t("common.delete")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default MessageActions;
