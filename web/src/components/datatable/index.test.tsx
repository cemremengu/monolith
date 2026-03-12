import { describe, expect, it } from "vitest";
import { screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { render } from "@/test/test-utils";

import { DataTable } from "./index";

type TestRow = {
  key: string;
  name: string;
  age: number;
  address: string;
};

const columns = [
  {
    title: "Name",
    dataIndex: "name",
    key: "name",
  },
  {
    title: "Age",
    dataIndex: "age",
    key: "age",
  },
  {
    title: "Address",
    dataIndex: "address",
    key: "address",
  },
] satisfies Array<{
  title: string;
  dataIndex: keyof TestRow;
  key: string;
}>;

function createRows(count: number): TestRow[] {
  return Array.from({ length: count }, (_, index) => ({
    key: String(index + 1),
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
  it("renders rows from dataSource", () => {
    render(
      <DataTable
        dataSource={[
          {
            key: "1",
            name: "Mike",
            age: 32,
            address: "10 Downing Street",
          },
          {
            key: "2",
            name: "John",
            age: 42,
            address: "11 Downing Street",
          },
        ]}
        columns={columns}
      />,
    );

    expect(screen.getByText("Mike")).toBeInTheDocument();
    expect(screen.getByText("John")).toBeInTheDocument();
    expect(screen.getByText("42")).toBeInTheDocument();
  });

  it("filters rows with the search input", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        dataSource={[
          {
            key: "1",
            name: "Mike",
            age: 32,
            address: "10 Downing Street",
          },
          {
            key: "2",
            name: "John",
            age: 42,
            address: "11 Downing Street",
          },
        ]}
        columns={columns}
      />,
    );

    await user.type(screen.getByLabelText("Search table"), "john");

    expect(screen.queryByText("Mike")).not.toBeInTheDocument();
    expect(screen.getByText("John")).toBeInTheDocument();
    expect(screen.getByText("1-1 of 1")).toBeInTheDocument();
  });

  it("sorts rows when clicking a column header", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        dataSource={[
          {
            key: "1",
            name: "Mike",
            age: 32,
            address: "10 Downing Street",
          },
          {
            key: "2",
            name: "John",
            age: 42,
            address: "11 Downing Street",
          },
        ]}
        columns={columns}
      />,
    );

    await user.click(screen.getByRole("button", { name: /^age$/i }));

    const ascendingRows = getBodyRows();
    expect(within(ascendingRows[0]).getByText("32")).toBeInTheDocument();
    expect(within(ascendingRows[1]).getByText("42")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /^age$/i }));

    const descendingRows = getBodyRows();
    expect(within(descendingRows[0]).getByText("42")).toBeInTheDocument();
    expect(within(descendingRows[1]).getByText("32")).toBeInTheDocument();
  });

  it("paginates rows on the client", async () => {
    const user = userEvent.setup();

    render(
      <DataTable dataSource={createRows(12)} columns={columns} pageSize={5} />,
    );

    expect(screen.getByText("User 1")).toBeInTheDocument();
    expect(screen.getByText("User 5")).toBeInTheDocument();
    expect(screen.queryByText("User 6")).not.toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /next page/i }));

    expect(screen.queryByText("User 1")).not.toBeInTheDocument();
    expect(screen.getByText("User 6")).toBeInTheDocument();
    expect(screen.getByText("Page 2 of 3")).toBeInTheDocument();
  });

  it("supports custom cell rendering", () => {
    render(
      <DataTable
        dataSource={[
          {
            key: "1",
            name: "Mike",
            age: 32,
            address: "10 Downing Street",
          },
        ]}
        columns={[
          ...columns.slice(0, 2),
          {
            title: "Address",
            dataIndex: "address",
            key: "address",
            render: (value) => <span>Location: {String(value)}</span>,
          },
        ]}
      />,
    );

    expect(screen.getByText("Location: 10 Downing Street")).toBeInTheDocument();
  });

  it("renders the empty state when no rows match", async () => {
    const user = userEvent.setup();

    render(
      <DataTable
        dataSource={[
          {
            key: "1",
            name: "Mike",
            age: 32,
            address: "10 Downing Street",
          },
        ]}
        columns={columns}
        emptyText="Nothing here"
      />,
    );

    await user.type(screen.getByLabelText("Search table"), "missing");

    expect(screen.getByText("Nothing here")).toBeInTheDocument();
  });
});
