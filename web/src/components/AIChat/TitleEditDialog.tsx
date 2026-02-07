import { Wand2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { aiServiceClient } from "@/connect";
import { cn } from "@/lib/utils";
import type { ConversationSummary } from "@/types/aichat";

interface TitleEditDialogProps {
  conversation: ConversationSummary;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: (newTitle: string) => void;
}

export function TitleEditDialog({ conversation, open, onOpenChange, onSuccess }: TitleEditDialogProps) {
  const { t } = useTranslation();
  const [title, setTitle] = useState(conversation.title);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const handleGenerateTitle = async () => {
    setIsGenerating(true);
    try {
      const response = await aiServiceClient.generateConversationTitle({
        id: Number(conversation.id),
      });
      setTitle(response.title);
    } catch (error) {
      console.error("Failed to generate title:", error);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await aiServiceClient.updateAIConversation({
        id: Number(conversation.id),
        title,
      });
      onSuccess?.(title);
      onOpenChange(false);
    } catch (error) {
      console.error("Failed to update title:", error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !isSaving) {
      handleSave();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-[28rem]">
        <DialogHeader>
          <DialogTitle>{t("ai.aichat.title-edit.title")}</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Title Input with Generate Button */}
          <div className="flex gap-2">
            <Input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={t("ai.aichat.title-edit.placeholder")}
              className="flex-1"
              autoFocus
            />
            <Button
              variant="outline"
              size="icon"
              onClick={handleGenerateTitle}
              disabled={isGenerating}
              title={t("ai.aichat.title-edit.generate-tooltip")}
            >
              <Wand2 className={cn("w-4 h-4", isGenerating && "animate-spin")} />
            </Button>
          </div>

          {/* Hint Text */}
          <p className="text-xs text-muted-foreground">{t("ai.aichat.title-edit.hint")}</p>

          {/* Actions */}
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={() => onOpenChange(false)}>
              {t("common.cancel")}
            </Button>
            <Button onClick={handleSave} disabled={isSaving || !title.trim()}>
              {isSaving ? t("common.saving") : t("common.save")}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
