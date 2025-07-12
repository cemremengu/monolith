import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

export function DashboardPage() {
  return (
    <div className="flex flex-1 flex-col">
      <div className="flex items-center gap-2 px-4 py-2">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink href="#">Dashboard</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
            <BreadcrumbItem>
              <BreadcrumbPage>Overview</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </div>
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <div className="grid auto-rows-min gap-4 md:grid-cols-3">
          <div className="bg-muted/50 aspect-video rounded-xl" />
          <div className="bg-muted/50 aspect-video rounded-xl" />
          <div className="bg-muted/50 aspect-video rounded-xl" />
        </div>
        <div className="bg-muted/50 min-h-[100vh] flex-1 rounded-xl md:min-h-min" />
        <div className="space-y-2">
          <p className="font-light">Light (300)</p>
          <p className="font-normal">
            Whereas recognition of the inherent dignity Normal (400)
          </p>
          <p className="font-medium">Medium (500)</p>
          <p className="font-semibold">SemiBold (600)</p>
          <p className="font-bold">Bold (700)</p>
          <p style={{ fontWeight: 450 }}>Custom (450)</p>
        </div>
      </div>
    </div>
  );
}
