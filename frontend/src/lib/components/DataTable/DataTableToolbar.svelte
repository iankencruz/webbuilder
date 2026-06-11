<script lang="ts" generics="TData extends Record<string, unknown>">
  import type { Table } from '@tanstack/table-core';
  import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

  let {
    table,
    filterColumn,
    filterPlaceholder = 'Search...',
  }: {
    table: Table<TData>;
    filterColumn?: string;
    filterPlaceholder?: string;
  } = $props();

  // fall back to first filterable column if none specified
  const activeFilter = $derived(
    filterColumn ?? table.getAllColumns().find((c) => c.getCanFilter())?.id,
  );
</script>

<div class="flex items-center gap-3">
  <!-- search input -->
  {#if activeFilter}
    <Input
      placeholder={filterPlaceholder}
      value={(table.getColumn(activeFilter)?.getFilterValue() as string) ?? ''}
      oninput={(e) =>
        table.getColumn(activeFilter)?.setFilterValue(e.currentTarget.value)}
      onchange={(e) =>
        table.getColumn(activeFilter)?.setFilterValue(e.currentTarget.value)}
      class="max-w-sm"
    />
  {/if}

  <!-- column visibility dropdown -->
  <DropdownMenu.Root>
    <DropdownMenu.Trigger>
      {#snippet child({ props })}
        <Button {...props} variant="outline" class="ms-auto">
          Columns <ChevronDownIcon class="ms-2 size-4" />
        </Button>
      {/snippet}
    </DropdownMenu.Trigger>
    <DropdownMenu.Content align="end">
      {#each table
        .getAllColumns()
        .filter((c) => c.getCanHide()) as column (column.id)}
        <DropdownMenu.CheckboxItem
          class="capitalize"
          bind:checked={
            () => column.getIsVisible(), (v) => column.toggleVisibility(!!v)
          }
        >
          {column.id.replace(/_/g, ' ')}
        </DropdownMenu.CheckboxItem>
      {/each}
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</div>
