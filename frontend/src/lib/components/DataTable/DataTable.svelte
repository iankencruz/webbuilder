<script lang="ts" generics="TData extends Record<string, unknown>">
  import {
    type ColumnDef,
    type PaginationState,
    type SortingState,
    type ColumnFiltersState,
    type VisibilityState,
    type RowSelectionState,
    getCoreRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    getFilteredRowModel,
  } from '@tanstack/table-core';
  import {
    createSvelteTable,
    FlexRender,
    renderComponent,
    renderSnippet,
  } from '$lib/components/ui/data-table/index.js';
  import { createRawSnippet, type Snippet } from 'svelte';
  import * as Table from '$lib/components/ui/table/index.js';
  import DataTableCheckbox from './DataTableCheckbox.svelte';
  import DataTableActions from './DataTableActions.svelte';
  import DataTableToolbar from './DataTableToolbar.svelte';
  import DataTablePagination from './DataTablePagination.svelte';

  type ColumnOverride = {
    label?: string;
    sortable?: boolean;
    cell?: (value: unknown, row: TData) => string;
    cellSnippet?: Snippet<[{ value: unknown; row: TData }]>;
  };

  type DataTableProps = {
    data: TData[];
    hiddenColumns?: string[];
    columnOverrides?: Record<string, ColumnOverride>;
    onRowAction?: (action: string, row: TData) => void;
    pageSize?: number;
    filterColumn: string;
    filterPlaceholder?: string;
  };

  let {
    data,
    hiddenColumns = [],
    columnOverrides = {},
    onRowAction,
    pageSize = 10,
    filterColumn,
    filterPlaceholder = 'Search...',
  }: DataTableProps = $props();

  const formatKey = (key: string): string => {
    if (key === 'id') return 'ID';
    return key.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
  };

  const formatValue = (key: string, value: unknown): string => {
    if (value === null || value === undefined) return '—';
    if (key.endsWith('_at') && typeof value === 'string') {
      return new Date(value).toLocaleDateString();
    }
    return String(value);
  };

  const buildColumns = (sample: TData): ColumnDef<TData>[] => {
    const keys = Object.keys(sample);

    const selectCol: ColumnDef<TData> = {
      id: 'select',
      enableHiding: false,
      enableSorting: false,
      header: ({ table }) =>
        renderComponent(DataTableCheckbox, {
          checked: table.getIsAllRowsSelected(),
          indeterminate:
            table.getIsSomePageRowsSelected() &&
            !table.getIsAllPageRowsSelected(),
          onCheckedChange: (v: boolean) => table.toggleAllPageRowsSelected(!!v),
          'aria-label': 'Select all',
        }),
      cell: ({ row }) =>
        renderComponent(DataTableCheckbox, {
          checked: row.getIsSelected(),
          onCheckedChange: (v: boolean) => row.toggleSelected(!!v),
          'aria-label': `Select row ${row.index + 1}`,
        }),
    };

    const dataCols: ColumnDef<TData>[] = keys.map((key) => {
      const override = columnOverrides[key] ?? {};
      const label = override.label ?? formatKey(key);
      const sortable = override.sortable ?? true;

      return {
        accessorKey: key,
        enableSorting: sortable,
        enableHiding: true,
        header: label,
        cell: ({ row }) => {
          const value = row.original[key];

          // custom snippet rendering
          if (override.cellSnippet) {
            return renderSnippet(override.cellSnippet, {
              value,
              row: row.original as TData,
            });
          }

          // default string rendering with optional custom cell function
          const rendered = override.cell
            ? override.cell(value, row.original as TData)
            : formatValue(key, value);

          // default - formatValue handles null, dates and string
          const snippet = createRawSnippet<[{ content: string }]>((get) => ({
            render: () => `<span>${get().content}</span>`,
          }));
          return renderSnippet(snippet, { content: rendered });
        },
      };
    });

    const actionsCol: ColumnDef<TData> = {
      id: 'actions',
      enableHiding: false,
      enableSorting: false,
      cell: ({ row }) =>
        renderComponent(DataTableActions, {
          onAction: (action: string) =>
            onRowAction?.(action, row.original as TData),
        }),
    };

    return [selectCol, ...dataCols, ...(onRowAction ? [actionsCol] : [])];
  };

  let pagination = $derived<PaginationState>({ pageIndex: 0, pageSize });
  let sorting = $state<SortingState>([]);
  let columnFilters = $state<ColumnFiltersState>([]);
  let columnVisibility = $derived<VisibilityState>(
    Object.fromEntries(hiddenColumns.map((k) => [k, false])),
  );
  let rowSelection = $state<RowSelectionState>({});

  const columns = $derived(data.length > 0 ? buildColumns(data[0]) : []);

  const table = createSvelteTable({
    get data() {
      return data;
    },
    get columns() {
      return columns;
    },
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: {
      get pagination() {
        return pagination;
      },
      get sorting() {
        return sorting;
      },
      get columnFilters() {
        return columnFilters;
      },
      get columnVisibility() {
        return columnVisibility;
      },
      get rowSelection() {
        return rowSelection;
      },
    },
    onPaginationChange: (u) => {
      pagination = typeof u === 'function' ? u(pagination) : u;
    },
    onSortingChange: (u) => {
      sorting = typeof u === 'function' ? u(sorting) : u;
    },
    onColumnFiltersChange: (u) => {
      columnFilters = typeof u === 'function' ? u(columnFilters) : u;
    },
    onColumnVisibilityChange: (u) => {
      columnVisibility = typeof u === 'function' ? u(columnVisibility) : u;
    },
    onRowSelectionChange: (u) => {
      rowSelection = typeof u === 'function' ? u(rowSelection) : u;
    },
  });
</script>

<div class="space-y-4">
  <DataTableToolbar {table} {filterColumn} {filterPlaceholder} />

  <div class="rounded-md border border-border">
    <Table.Root>
      <Table.Header>
        {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
          <Table.Row>
            {#each headerGroup.headers as header (header.id)}
              <Table.Head class="[&:has([role=checkbox])]:ps-3">
                {#if !header.isPlaceholder}
                  <FlexRender
                    content={header.column.columnDef.header}
                    context={header.getContext()}
                  />
                {/if}
              </Table.Head>
            {/each}
          </Table.Row>
        {/each}
      </Table.Header>

      <Table.Body>
        {#each table.getRowModel().rows as row (row.id)}
          <Table.Row data-state={row.getIsSelected() && 'selected'}>
            {#each row.getVisibleCells() as cell (cell.id)}
              <Table.Cell class="[&:has([role=checkbox])]:ps-3">
                <FlexRender
                  content={cell.column.columnDef.cell}
                  context={cell.getContext()}
                />
              </Table.Cell>
            {/each}
          </Table.Row>
        {:else}
          <Table.Row>
            <Table.Cell
              colspan={columns.length}
              class="h-24 text-center text-muted-foreground"
            >
              No results.
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  </div>

  <DataTablePagination {table} />
</div>
