import {
  type Column,
  type ColumnDef,
  type ColumnFiltersState,
  type SortingState,
  type VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsUpDown,
  ChevronUp,
  ChevronDown,
  Settings2,
} from "lucide-react";
import * as React from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { cn } from "@/lib/utils";

export type { ColumnDef };

export type DataTableProps<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  filterColumn?: string;
  filterPlaceholder?: string;
  enableColumnVisibility?: boolean;
  pageSize?: number;
  emptyText?: React.ReactNode;
  className?: string;
  borderless?: boolean;
  initialSorting?: SortingState;
};

export function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  className,
}: {
  column: Column<TData, TValue>;
  title: string;
  className?: string;
}) {
  if (!column.getCanSort()) {
    return <div className={cn("text-muted-foreground", className)}>{title}</div>;
  }

  const sorted = column.getIsSorted();
  const active = sorted !== false || column.getIsFiltered();

  return (
    <Button
      variant="ghost"
      size="sm"
      className={cn(
        "group -ml-2 h-8 group-hover:text-foreground",
        !active && "text-muted-foreground",
        className,
      )}
      onClick={() => column.toggleSorting(sorted === "asc")}
    >
      {title}
      {sorted === "asc" ? (
        <ChevronUp className="size-3.5" />
      ) : sorted === "desc" ? (
        <ChevronDown className="size-3.5" />
      ) : (
        <ChevronsUpDown className="size-3.5 opacity-0 transition-opacity group-hover:opacity-100" />
      )}
    </Button>
  );
}

function ColumnVisibilityToggle<TData>({ column }: { column: Column<TData, unknown> }) {
  const handleCheckedChange = React.useCallback(
    (value: boolean) => column.toggleVisibility(!!value),
    [column],
  );

  return (
    <DropdownMenuCheckboxItem
      className="capitalize"
      checked={column.getIsVisible()}
      onCheckedChange={handleCheckedChange}
    >
      {column.id}
    </DropdownMenuCheckboxItem>
  );
}

export function DataTable<TData, TValue>({
  columns,
  data,
  filterColumn,
  filterPlaceholder,
  enableColumnVisibility = false,
  pageSize = 10,
  emptyText,
  className,
  borderless = false,
  initialSorting = [],
}: DataTableProps<TData, TValue>) {
  const { t } = useTranslation();
  const [sorting, setSorting] = React.useState<SortingState>(initialSorting);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
  const [globalFilter, setGlobalFilter] = React.useState("");

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onGlobalFilterChange: setGlobalFilter,
    initialState: {
      pagination: {
        pageSize,
      },
    },
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      globalFilter,
    },
  });

  const handleFilterChange = React.useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      if (filterColumn) {
        table.getColumn(filterColumn)?.setFilterValue(event.target.value);
      } else {
        setGlobalFilter(event.target.value);
      }
    },
    [filterColumn, table],
  );

  const handlePreviousPage = React.useCallback(() => table.previousPage(), [table]);
  const handleNextPage = React.useCallback(() => table.nextPage(), [table]);

  return (
    <div className={cn("space-y-4", className)}>
      <div className="flex items-center gap-2">
        <Input
          placeholder={filterPlaceholder ?? t("common.table.filter")}
          value={
            filterColumn
              ? ((table.getColumn(filterColumn)?.getFilterValue() as string) ?? "")
              : globalFilter
          }
          onChange={handleFilterChange}
          className="max-w-sm"
          aria-label={t("common.table.filterLabel")}
        />
        {enableColumnVisibility && (
          <DropdownMenu>
            <DropdownMenuTrigger
              render={<Button variant="outline" size="sm" className="ml-auto" />}
            >
              <Settings2 className="size-4" />
              {t("common.table.columns")}
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>{t("common.table.toggleColumns")}</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {table
                .getAllColumns()
                .filter((column) => column.getCanHide())
                .map((column) => (
                  <ColumnVisibilityToggle key={column.id} column={column} />
                ))}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>

      <div className={cn("rounded-md", !borderless && "border")}>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow
                key={headerGroup.id}
                className={cn("hover:bg-transparent", borderless && "border-0")}
              >
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id} className={cn(borderless && "border-0")}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  {emptyText ?? t("common.table.noResults")}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {table.getPageCount() > 1 && (
        <div className="flex items-center justify-end">
          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={handlePreviousPage}
              disabled={!table.getCanPreviousPage()}
              aria-label={t("common.table.previousLabel")}
            >
              <ChevronLeft className="size-4" />
              {t("common.table.previous")}
            </Button>
            <span className="text-sm">
              {t("common.table.page", {
                page: table.getState().pagination.pageIndex + 1,
                total: table.getPageCount(),
              })}
            </span>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={handleNextPage}
              disabled={!table.getCanNextPage()}
              aria-label={t("common.table.nextLabel")}
            >
              {t("common.table.next")}
              <ChevronRight className="size-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
