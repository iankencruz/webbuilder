<script lang="ts">
  import PaletteItem from './PaletteItem.svelte';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import * as Popover from '$lib/components/ui/popover';
  import * as ScrollArea from '$lib/components/ui/scroll-area';
  import PlusIcon from '@lucide/svelte/icons/plus';
  import type { BlockCollection } from './registry';

  let {
    onAdd,
    compact = false,
    expanded = false,
  }: {
    onAdd: (collection: BlockCollection) => void;
    compact?: boolean;
    expanded?: boolean;
  } = $props();

  const paletteItems: { collection: BlockCollection; label: string }[] = [
    { collection: 'hero', label: 'Hero' },
    { collection: 'rich_text', label: 'Rich Text' },
  ];

  let search = $state('');
  let open = $state(false);

  let filtered = $derived(
    paletteItems.filter((item) =>
      item.label.toLowerCase().includes(search.toLowerCase()),
    ),
  );

  function pick(collection: BlockCollection) {
    onAdd(collection);
    open = false;
  }
</script>

{#if expanded}
  <div class="space-y-3">
    <Input type="text" bind:value={search} placeholder="Search blocks…" />
    <div class="space-y-2">
      {#each filtered as item (item.collection)}
        <PaletteItem label={item.label} onClick={() => pick(item.collection)} />
      {:else}
        <p class="text-sm text-muted-foreground px-1">No blocks found.</p>
      {/each}
    </div>
  </div>
{:else}
  <Popover.Root bind:open>
    <Popover.Trigger>
      {#snippet child({ props })}
        {#if compact}
          <Button
            {...props}
            variant="outline"
            size="icon"
            class="size-6 rounded-full"
          >
            <PlusIcon class="size-3.5" />
          </Button>
        {:else}
          <Button {...props} variant="outline" size="sm" class="w-full">
            <PlusIcon class="size-4" />
            Add Block
          </Button>
        {/if}
      {/snippet}
    </Popover.Trigger>

    <Popover.Content class="w-64 p-0" align="center">
      <div class="p-3 border-b border-border">
        <Input type="text" bind:value={search} placeholder="Search blocks…" />
      </div>

      <ScrollArea.Root class="h-64">
        <div class="p-3 space-y-2">
          {#each filtered as item (item.collection)}
            <PaletteItem
              label={item.label}
              onClick={() => pick(item.collection)}
            />
          {:else}
            <p class="text-sm text-muted-foreground px-1">No blocks found.</p>
          {/each}
        </div>
      </ScrollArea.Root>
    </Popover.Content>
  </Popover.Root>
{/if}
