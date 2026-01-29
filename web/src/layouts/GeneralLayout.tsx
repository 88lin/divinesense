import { Outlet } from "react-router-dom";
import NavigationDrawer from "@/components/NavigationDrawer";
import RouteHeaderImage from "@/components/RouteHeaderImage";
import useMediaQuery from "@/hooks/useMediaQuery";

const GeneralLayout = () => {
    const sm = useMediaQuery("sm");

    return (
        <section className="w-full h-full flex flex-col justify-start items-center overflow-hidden">
            {!sm && (
                <div className="w-full flex items-center justify-center px-4 py-3 border-b bg-background sticky top-0 z-10 shrink-0 h-14 overflow-hidden relative">
                    <div className="absolute left-4 top-0 bottom-0 flex items-center">
                        <NavigationDrawer />
                    </div>
                    <RouteHeaderImage />
                </div>
            )}
            <div className="w-full h-full overflow-y-auto">
                <div className="w-full min-h-full">
                    <Outlet />
                </div>
            </div>
        </section>
    );
};

export default GeneralLayout;
