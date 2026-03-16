import { describe, expect, it } from "vitest";
import { screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ColumnDef } from "@tanstack/react-table";

import { render } from "@/test/test-utils";

import { DataTable } from "./index";

type TestRow = {
  name: string;
  age: number;
  address: string;
};

const columns: ColumnDef<TestRow>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "age",
    header: "Age",
  },
  {
    accessorKey: "address",
    header: "Address",
  },
];

function createRows(count: number): TestRow[] {
  return Array.from({ length: count }, (_, index) => ({
    name: `User ${index + 1}`,
    age: 20 + index,
    address: `Street ${index + 1}`,
  }));
}

function getBodyRows() {
  const table = screen.getByRole("table");
  const rows = within(table).getAllByRole("row");

  return rows.slice(1);
}

describe("DataTable", () => {
  it("renders rows from data", () => {
    render(
      <DataTable
        data={[
          { name: "Mike", age: 32, address: "10 Downing Street" },
          { name: "John", age: 42, address: "11 Downing Street" },
        ]}
        columns={columns}
      />,
    );

    expect(screen.getByText("Mike")).toBeInTheDocument();
    expect(screen.getByText("John")).toBeInTheDocument();
    expect(screen.getByText("42")).toBeInTheDocument();
  });

  it("filters rows with the filter input (global)", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        data={[
          { name: "Mike", age: 32, address: "10 Downing Street" },
          { name: "John", age: 42, address: "11 Downing Street" },
        ]}
        columns={columns}
      />,
    );

    await user.type(screen.getByLabelText("Filter table"), "john");

    expect(screen.queryByText("Mike")).not.toBeInTheDocument();
    expect(screen.getByText("John")).toBeInTheDocument();
  });

  it("filters rows by a specific column", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        data={[
          { name: "Mike", age: 32, address: "10 Downing Street" },
          { name: "John", age: 42, address: "11 Downing Street" },
        ]}
        columns={columns}
        filterColumn="name"
        filterPlaceholder="Filter names..."
      />,
    );

    await user.type(screen.getByLabelText("Filter table"), "John");

    expect(screen.queryByText("Mike")).not.toBeInTheDocument();
    expect(screen.getByText("John")).toBeInTheDocument();
  });

  it("sorts rows when clicking a column header", async () => {
    const user = userEvent.setup();

    const sortableColumns: ColumnDef<TestRow>[] = [
      { accessorKey: "name", header: "Name" },
      {
        accessorKey: "age",
        header: ({ column }) => (
          <button onClick={() => column.toggleSorting()}>Age</button>
        ),
      },
      { accessorKey: "address", header: "Address" },
    ];

    render(
      <DataTable
        data={[
          { name: "Mike", age: 32, address: "10 Downing Street" },
          { name: "John", age: 42, address: "11 Downing Street" },
        ]}
        columns={sortableColumns}
      />,
    );

    // TanStack auto-detects numeric columns and sorts descending first
    await user.click(screen.getByRole("button", { name: /^age$/i }));

    const descendingRows = getBodyRows();
    expect(within(descendingRows[0]).getByText("42")).toBeInTheDocument();
    expect(within(descendingRows[1]).getByText("32")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /^age$/i }));

    const ascendingRows = getBodyRows();
    expect(within(ascendingRows[0]).getByText("32")).toBeInTheDocument();
    expect(within(ascendingRows[1]).getByText("42")).toBeInTheDocument();
  });

  it("paginates rows on the client", async () => {
    const user = userEvent.setup();

    render(<DataTable data={createRows(12)} columns={columns} pageSize={5} />);

    expect(screen.getByText("User 1")).toBeInTheDocument();
    expect(screen.getByText("User 5")).toBeInTheDocument();
    expect(screen.queryByText("User 6")).not.toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /next page/i }));

    expect(screen.queryByText("User 1")).not.toBeInTheDocument();
    expect(screen.getByText("User 6")).toBeInTheDocument();
    expect(screen.getByText("Page 2 of 3")).toBeInTheDocument();
  });

  it("supports custom cell rendering", () => {
    const customColumns: ColumnDef<TestRow>[] = [
      { accessorKey: "name", header: "Name" },
      { accessorKey: "age", header: "Age" },
      {
        accessorKey: "address",
        header: "Address",
        cell: ({ getValue }) => <span>Location: {getValue<string>()}</span>,
      },
    ];

    render(
      <DataTable
        data={[{ name: "Mike", age: 32, address: "10 Downing Street" }]}
        columns={customColumns}
      />,
    );

    expect(screen.getByText("Location: 10 Downing Street")).toBeInTheDocument();
  });

  it("renders the empty state when no rows match", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        data={[{ name: "Mike", age: 32, address: "10 Downing Street" }]}
        columns={columns}
        emptyText="Nothing here"
      />,
    );

    await user.type(screen.getByLabelText("Filter table"), "missing");

    expect(screen.getByText("Nothing here")).toBeInTheDocument();
  });
});
