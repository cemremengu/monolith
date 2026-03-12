"use client";

import * as React from "react";
import {
  ArrowDown,
  ArrowUp,
  ArrowUpDown,
  ChevronLeft,
  ChevronRight,
  Search,
} from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
} from "@/components/ui/input-group";
import {
  Table as BaseTable,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

type DataTableRecord = Record<PropertyKey, unknown>;

type SortOrder = "asc" | "desc";

type Align = "left" | "center" | "right";

type SortState<TData extends DataTableRecord> = {
  column: DataTableColumn<TData>;
  order: SortOrder;
};

export type DataTableColumn<TData extends DataTableRecord> = {
  title: React.ReactNode;
  dataIndex: keyof TData | string;
  key: React.Key;
  render?: (value: unknown, record: TData, index: number) => React.ReactNode;
  sorter?: boolean | ((a: TData, b: TData) => number);
  searchable?: boolean;
  align?: Align;
  className?: string;
  headerClassName?: string;
  width?: React.CSSProperties["width"];
};

export type DataTableProps<TData extends DataTableRecord> = {
  dataSource: TData[];
  columns: DataTableColumn<TData>[];
  className?: string;
  emptyText?: React.ReactNode;
  pageSize?: number;
  rowKey?: keyof TData | ((record: TData) => React.Key);
  searchPlaceholder?: string;
};

function getValue(
  record: DataTableRecord,
  dataIndex: string | number | symbol,
) {
  if (typeof dataIndex === "symbol") {
    return record[dataIndex];
  }

  const resolvedDataIndex =
    typeof dataIndex === "number" ? String(dataIndex) : dataIndex;

  if (Object.hasOwn(record, resolvedDataIndex)) {
    return record[resolvedDataIndex];
  }

  return resolvedDataIndex.split(".").reduce<unknown>((value, key) => {
    if (value == null || typeof value !== "object") {
      return undefined;
    }

    return (value as DataTableRecord)[key];
  }, record);
}

function normalizeSearchValue(value: unknown): string {
  if (value == null) {
    return "";
  }

  if (Array.isArray(value)) {
    return value.map(normalizeSearchValue).join(" ");
  }

  if (value instanceof Date) {
    return value.toISOString();
  }

  if (typeof value === "object") {
    return JSON.stringify(value);
  }

  return String(value);
}

function compareValues(left: unknown, right: unknown) {
  if (left == null && right == null) {
    return 0;
  }

  if (left == null) {
    return 1;
  }

  if (right == null) {
    return -1;
  }

  if (typeof left === "number" && typeof right === "number") {
    return left - right;
  }

  if (typeof left === "boolean" && typeof right === "boolean") {
    return Number(left) - Number(right);
  }

  return String(left).localeCompare(String(right), undefined, {
    numeric: true,
    sensitivity: "base",
  });
}

function getCellAlignmentClassName(align: Align = "left") {
  switch (align) {
    case "center":
      return "text-center";
    case "right":
      return "text-right";
    default:
      return "text-left";
  }
}

function getPaginationLabel(
  totalItems: number,
  pageSize: number,
  currentPage: number,
) {
  if (totalItems === 0) {
    return "0 results";
  }

  const start = (currentPage - 1) * pageSize + 1;
  const end = Math.min(currentPage * pageSize, totalItems);

  return `${start}-${end} of ${totalItems}`;
}

export function DataTable<TData extends DataTableRecord>({
  dataSource,
  columns,
  className,
  emptyText = "No data",
  pageSize = 10,
  rowKey,
  searchPlaceholder = "Search",
}: DataTableProps<TData>) {
  const normalizedPageSize = pageSize > 0 ? pageSize : 10;
  const searchableColumns = columns.filter(
    (column) => column.searchable !== false,
  );
  const [searchQuery, setSearchQuery] = React.useState("");
  const [sortState, setSortState] = React.useState<SortState<TData> | null>(
    null,
  );
  const [page, setPage] = React.useState(1);

  const normalizedQuery = searchQuery.trim().toLocaleLowerCase();

  const filteredData = normalizedQuery
    ? dataSource.filter((record) =>
        searchableColumns.some((column) =>
          normalizeSearchValue(getValue(record, column.dataIndex))
            .toLocaleLowerCase()
            .includes(normalizedQuery),
        ),
      )
    : dataSource;

  const sortedData = sortState
    ? [...filteredData].sort((left, right) => {
        const sorter = sortState.column.sorter;
        const result =
          typeof sorter === "function"
            ? sorter(left, right)
            : compareValues(
                getValue(left, sortState.column.dataIndex),
                getValue(right, sortState.column.dataIndex),
              );

        return sortState.order === "asc" ? result : result * -1;
      })
    : filteredData;

  const totalPages = Math.max(
    1,
    Math.ceil(sortedData.length / normalizedPageSize),
  );
  const currentPage = Math.min(page, totalPages);
  const startIndex = (currentPage - 1) * normalizedPageSize;
  const paginatedData = sortedData.slice(
    startIndex,
    startIndex + normalizedPageSize,
  );

  const canSearch = searchableColumns.length > 0;

  function handleSort(column: DataTableColumn<TData>) {
    if (column.sorter === false) {
      return;
    }

    setPage(1);
    setSortState((currentSort) => {
      if (currentSort?.column.key !== column.key) {
        return {
          column,
          order: "asc",
        };
      }

      if (currentSort.order === "asc") {
        return {
          column,
          order: "desc",
        };
      }

      return null;
    });
  }

  function getRowKey(record: TData, index: number) {
    if (typeof rowKey === "function") {
      return rowKey(record);
    }

    if (rowKey) {
      return record[rowKey] as React.Key;
    }

    if (Object.hasOwn(record, "key")) {
      return record.key as React.Key;
    }

    return index;
  }

  return (
    <div className={cn("space-y-4", className)}>
      {canSearch ? (
        <InputGroup className="max-w-sm">
          <InputGroupAddon align="inline-start">
            <InputGroupText>
              <Search className="size-4" />
            </InputGroupText>
          </InputGroupAddon>
          <InputGroupInput
            value={searchQuery}
            onChange={(event) => {
              setSearchQuery(event.target.value);
              setPage(1);
            }}
            placeholder={searchPlaceholder}
            aria-label="Search table"
          />
        </InputGroup>
      ) : null}

      <div className="rounded-md border">
        <BaseTable>
          <TableHeader>
            <TableRow>
              {columns.map((column) => {
                const isSortable = column.sorter !== false;
                const isSorted = sortState?.column.key === column.key;
                const ariaSort =
                  isSorted && sortState?.order === "asc"
                    ? "ascending"
                    : isSorted && sortState?.order === "desc"
                      ? "descending"
                      : "none";

                return (
                  <TableHead
                    key={column.key}
                    className={cn(
                      getCellAlignmentClassName(column.align),
                      column.headerClassName,
                    )}
                    style={{ width: column.width }}
                    aria-sort={ariaSort}
                  >
                    {isSortable ? (
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className={cn(
                          "hover:text-foreground -mx-2 h-auto w-full justify-start px-2 py-1 font-medium",
                          column.align === "center" && "justify-center",
                          column.align === "right" && "justify-end",
                        )}
                        onClick={() => handleSort(column)}
                      >
                        <span>{column.title}</span>
                        {isSorted ? (
                          sortState?.order === "asc" ? (
                            <ArrowUp className="size-4" />
                          ) : (
                            <ArrowDown className="size-4" />
                          )
                        ) : (
                          <ArrowUpDown className="text-muted-foreground size-4" />
                        )}
                      </Button>
                    ) : (
                      column.title
                    )}
                  </TableHead>
                );
              })}
            </TableRow>
          </TableHeader>
          <TableBody>
            {paginatedData.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {emptyText}
                </TableCell>
              </TableRow>
            ) : (
              paginatedData.map((record, index) => (
                <TableRow key={getRowKey(record, startIndex + index)}>
                  {columns.map((column) => {
                    const value = getValue(record, column.dataIndex);

                    return (
                      <TableCell
                        key={column.key}
                        className={cn(
                          getCellAlignmentClassName(column.align),
                          column.className,
                        )}
                      >
                        {column.render
                          ? column.render(value, record, startIndex + index)
                          : normalizeSearchValue(value)}
                      </TableCell>
                    );
                  })}
                </TableRow>
              ))
            )}
          </TableBody>
        </BaseTable>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p className="text-muted-foreground text-sm">
          {getPaginationLabel(
            sortedData.length,
            normalizedPageSize,
            currentPage,
          )}
        </p>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => setPage((current) => Math.max(1, current - 1))}
            disabled={currentPage === 1}
            aria-label="Previous page"
          >
            <ChevronLeft className="size-4" />
            Previous
          </Button>
          <span className="text-sm">
            Page {currentPage} of {totalPages}
          </span>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              setPage((current) => Math.min(totalPages, current + 1))
            }
            disabled={currentPage === totalPages}
            aria-label="Next page"
          >
            Next
            <ChevronRight className="size-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
