<script lang="ts" generics="TData extends Record<string, unknown>">
  import type { Table } from '@tanstack/table-core';
  import { Button } from '$lib/components/ui/button/index.js';

  let { table }: { table: Table<TData> } = $props();
</script>

<div class="flex items-center justify-between px-2">
  <!-- selected rows count -->
  <p class="text-sm text-muted-foreground">
    {table.getFilteredSelectedRowModel().rows.length} of{' '}
    {table.getFilteredRowModel().rows.length} row(s) selected
  </p>

  <!-- page controls -->
  <div class="flex items-center gap-4">
    <p class="text-sm text-muted-foreground">
      Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
    </p>
    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        onclick={() => table.previousPage()}
        disabled={!table.getCanPreviousPage()}
      >
        Previous
      </Button>
      <Button
        variant="outline"
        size="sm"
        onclick={() => table.nextPage()}
        disabled={!table.getCanNextPage()}
      >
        Next
      </Button>
    </div>
  </div>
</div>
