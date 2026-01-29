import { useEffect } from "react";
import { useAIChat } from "@/contexts/AIChatContext";

const GEEK_MODE_BODY_CLASS = "geek-mode";

/**
 * Geek Mode Theme Provider
 *
 * Applies `.geek-mode` class to document.body when Geek Mode is enabled.
 * This enables all Geek Mode visual styles defined in themes/geek-mode.css.
 *
 * Should be placed near the root of the app, within AIChatProvider.
 */
export function GeekModeThemeProvider() {
  const { state } = useAIChat();
  const { geekMode } = state;

  useEffect(() => {
    if (typeof document === "undefined") return;

    if (geekMode) {
      document.body.classList.add(GEEK_MODE_BODY_CLASS);
    } else {
      document.body.classList.remove(GEEK_MODE_BODY_CLASS);
    }

    // Cleanup on unmount
    return () => {
      document.body.classList.remove(GEEK_MODE_BODY_CLASS);
    };
  }, [geekMode]);

  return null;
}
