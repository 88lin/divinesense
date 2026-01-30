import { MenuIcon } from "lucide-react";
import { useState } from "react";
import { Outlet } from "react-router-dom";
import { AIChatSidebar } from "@/components/AIChat/AIChatSidebar";
import { GeekModeThemeProvider } from "@/components/AIChat/GeekModeThemeProvider";
import NavigationDrawer from "@/components/NavigationDrawer";
import RouteHeaderImage from "@/components/RouteHeaderImage";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { AIChatProvider, useAIChat } from "@/contexts/AIChatContext";
import useMediaQuery from "@/hooks/useMediaQuery";
import { cn } from "@/lib/utils";

/**
 * AI Chat Layout - 优化的聊天布局
 *
 * UX/UI 改进：
 * - 优化移动端和桌面端的布局切换
 * - 统一间距和边框样式
 * - 改进侧边栏和主内容的视觉层次
 * - 移动端 Header 支持极客模式视觉反馈
 */
const AIChatLayoutContent = () => {
  const lg = useMediaQuery("lg");
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const { state } = useAIChat();
  const geekMode = state.geekMode || false;
  const immersiveMode = state.immersiveMode || false;

  return (
    <section className="@container w-full h-screen flex flex-col lg:h-screen overflow-hidden bg-zinc-50 dark:bg-zinc-950">
      {/* Mobile Header */}
      <div
        className={cn(
          "lg:hidden flex-none relative flex items-center justify-center px-4 h-14 shrink-0 border-b transition-colors",
          // Normal mode styles
          !geekMode && "border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900",
          // Geek mode styles
          geekMode && [
            "border-green-500/30 bg-green-950/20 dark:bg-green-950/40",
            // Bottom accent line
            "after:content-[''] after:absolute after:bottom-0 after:left-0 after:right-0 after:h-[2px] after:bg-gradient-to-r after:from-transparent via-green-500 to-transparent",
          ],
        )}
      >
        {/* Left - Navigation Drawer */}
        <div className="absolute left-4 top-0 bottom-0 flex items-center">
          <NavigationDrawer />
        </div>

        {/* Center - Title with geek mode styling */}
        <div className={cn("flex items-center gap-2", geekMode && "font-mono")}>
          {/* Geek mode indicator */}
          {geekMode && (
            <span className="flex items-center gap-1 text-xs text-green-600 dark:text-green-400">
              <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
              <span className="hidden sm:inline">GEEK</span>
            </span>
          )}
          <RouteHeaderImage />
        </div>

        {/* Right - Sidebar Toggle */}
        <div className="absolute right-0 top-0 bottom-0 px-3 flex items-center">
          <Sheet open={mobileSidebarOpen} onOpenChange={setMobileSidebarOpen}>
            <SheetContent
              side="right"
              className={cn(
                "w-80 max-w-full [&_.absolute.top-4.right-4]:hidden border-l",
                !geekMode && "bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800",
                geekMode && "bg-green-950/10 dark:bg-green-950/20 border-green-500/30",
                "gap-0",
              )}
            >
              <SheetHeader className="p-0">
                <SheetTitle className="sr-only">AI Assistant</SheetTitle>
              </SheetHeader>
              <AIChatSidebar className="h-full" onClose={() => setMobileSidebarOpen(false)} />
            </SheetContent>
          </Sheet>
          <Button variant="ghost" size="icon" onClick={() => setMobileSidebarOpen(true)} aria-label="Open sidebar" className="h-9 w-9">
            <MenuIcon className={cn("w-5 h-5", geekMode && "text-green-600 dark:text-green-400")} />
          </Button>
        </div>
      </div>

      {/* Desktop Sidebar - Hidden in immersive mode */}
      {lg && !immersiveMode && (
        <div
          className={cn(
            "fixed top-0 left-16 shrink-0 h-svh border-r backdrop-blur-sm w-72 overflow-hidden pt-2 transition-colors",
            !geekMode && "border-zinc-200/80 dark:border-zinc-800/80 bg-zinc-50/95 dark:bg-zinc-900/95",
            geekMode && "border-green-500/20 bg-green-950/10 dark:bg-green-950/20",
          )}
        >
          <AIChatSidebar className="h-full" />
        </div>
      )}

      {/* Main Content */}
      <div
        className={cn(
          "flex-1 min-h-0 overflow-hidden transition-all duration-300",
          !geekMode && "bg-white dark:bg-zinc-900",
          geekMode && "bg-green-50/30 dark:bg-green-950/10",
          lg && !immersiveMode ? "pl-72" : "",
        )}
      >
        <Outlet />
      </div>
    </section>
  );
};

const AIChatLayout = () => {
  return (
    <AIChatProvider>
      <GeekModeThemeProvider />
      <AIChatLayoutContent />
    </AIChatProvider>
  );
};

export default AIChatLayout;
