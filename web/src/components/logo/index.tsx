import { SidebarTrigger, useSidebar } from "@/components/ui/sidebar";
// import logoImage from "@/assets/logo.png";

export function Logo() {
  const { state } = useSidebar();

  return (
    <div className="flex items-center justify-between gap-2 py-1">
      {state === "expanded" ? (
        <div className="flex items-center gap-2 px-0.5 text-lg font-semibold">
          <div className="flex h-8 w-8 items-center justify-center overflow-hidden rounded-lg">
            {/* <img
              src={logoImage}
              alt="Monolith Logo"
              className="h-full w-full object-contain"
            /> */}
          </div>
          <span className="truncate">Monolith</span>
        </div>
      ) : (
        <div className="group flex w-full justify-center">
          <div className="flex h-8 w-8 items-center justify-center overflow-hidden rounded-lg transition-opacity group-hover:opacity-0">
            {/* <img
              src={logoImage}
              alt="Monolith Logo"
              className="h-full w-full object-contain"
            /> */}
          </div>
          <div className="absolute flex h-8 w-8 items-center justify-center opacity-0 transition-opacity group-hover:opacity-100">
            <SidebarTrigger />
          </div>
        </div>
      )}
      {state === "expanded" && <SidebarTrigger className="-ml-1" />}
    </div>
  );
}
