import { useTranslation } from "react-i18next";
import { useLocation } from "react-router-dom";
import { Routes } from "@/router";

const RouteHeaderImage = () => {
  const location = useLocation();
  const { i18n } = useTranslation();
  const path = location.pathname;

  // Detect if language is Chinese (simplified or traditional)
  const isZh = i18n.language.startsWith("zh");
  const suffix = isZh ? "-zh" : "";

  let headerName = "";

  if (path === Routes.HOME) headerName = "memos";
  else if (path.startsWith(Routes.EXPLORE)) headerName = "explore";
  else if (path.startsWith(Routes.ARCHIVED)) headerName = "memos";
  else if (path.startsWith("/u/"))
    headerName = "memos"; // Profile
  else if (path.startsWith(Routes.CHAT)) headerName = "ai";
  else if (path.startsWith(Routes.SCHEDULE)) headerName = "schedule";
  else if (path.startsWith(Routes.REVIEW)) headerName = "review";
  else if (path.startsWith(Routes.KNOWLEDGE_GRAPH)) headerName = "knowledge";
  else if (path.startsWith(Routes.ATTACHMENTS)) headerName = "files";
  else if (path.startsWith(Routes.INBOX)) headerName = "inbox";
  else if (path.startsWith(Routes.SETTING)) headerName = "memos";
  else if (path.startsWith("/memos/")) headerName = "memos"; // Detail

  if (!headerName) return null;

  const headerSrc = `/headers/header-${headerName}${suffix}.svg`;

  return <img src={headerSrc} alt="Page Header" className="h-8 w-auto object-contain select-none opacity-90 dark:opacity-100" />;
};

export default RouteHeaderImage;
