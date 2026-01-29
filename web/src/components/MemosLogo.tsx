import { useTranslation } from "react-i18next";
import { useInstance } from "@/contexts/InstanceContext";
import { cn } from "@/lib/utils";
import UserAvatar from "./UserAvatar";

interface Props {
  className?: string;
  collapsed?: boolean;
}


function MemosLogo(props: Props) {
  const { collapsed } = props;
  const { generalSetting: instanceGeneralSetting } = useInstance();
  const { i18n } = useTranslation();
  const title = instanceGeneralSetting.customProfile?.title || "Memos";
  const avatarUrl = instanceGeneralSetting.customProfile?.logoUrl || "/logo.webp";

  const isZh = i18n.language.startsWith("zh");
  const suffix = isZh ? "-zh" : "";

  return (
    <div className={cn("relative w-full h-auto shrink-0", props.className)}>
      <div className={cn("w-auto flex flex-row justify-start items-center text-foreground", collapsed ? "px-1" : "px-2")}>
        {collapsed ? (
          <UserAvatar className="shrink-0" avatarUrl={avatarUrl} />
        ) : (
          <>
            <img src={`/full-logo-light${suffix}.svg`} alt={title} className="h-10 w-auto object-contain dark:hidden" />
            <img src={`/full-logo-dark${suffix}.svg`} alt={title} className="h-10 w-auto object-contain hidden dark:block" />
          </>
        )}
      </div>
    </div>
  );
}

export default MemosLogo;
