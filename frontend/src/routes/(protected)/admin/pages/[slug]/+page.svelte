<script lang="ts">
  import { blocksAPI } from '$lib/api/blocks';
  import type { ResolvedBlock } from '$lib/types/blocks';
  import { DndContext, closestCenter } from '@dnd-kit-svelte/core';
  import {
    SortableContext,
    verticalListSortingStrategy,
    arrayMove,
  } from '@dnd-kit-svelte/sortable';
  import SortableBlockItem from '$lib/blocks/SortableBlockItem.svelte';

  import { blockRegistry, type BlockCollection } from '$lib/blocks/registry.js';
  import { untrack } from 'svelte';

  let { data } = $props();

  let blocks = $state<ResolvedBlock[]>(untrack(() => data.blocks));

  const blockComponents = blockRegistry;

  let blockIds = $derived(blocks.map((b) => b.junction_id));

  // -- Add Block --
  let adding = $state(false);

  async function addBlock(collection: BlockCollection) {
    adding = true;
    try {
      // 1. Create the block content row
      const defaults: Record<string, Record<string, unknown>> = {
        hero: {
          heading: 'New Hero',
          subheading: '',
          cta_label: '',
          cta_url: '',
        },
        rich_text: { content: '', format: 'html' },
      };
      const { id: blockId } = await blocksAPI.createBlock(
        collection,
        defaults[collection],
      );

      // 2. Attach it to the page junction
      const nextSortOrder = blocks.length;
      const junction = await blocksAPI.addToPage(
        data.page.id,
        blockId,
        collection,
        nextSortOrder,
      );

      // 3. Fetch the resolved block and append locally
      const newBlock: ResolvedBlock = {
        junction_id: junction.id,
        sort_order: nextSortOrder,
        hide_block: false,
        collection,
        block_id: blockId,
        data: defaults[collection] as any,
      };
      blocks = [...blocks, newBlock];
    } finally {
      adding = false;
    }
  }

  // --- Remove Block ---
  async function removeBlock(junctionId: number) {
    await blocksAPI.deleteFromPage(data.page.id, junctionId);
    blocks = blocks.filter((b) => b.junction_id !== junctionId);
  }

  // --- Reorder ---
  async function handleDragEnd(event: any) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIndex = blocks.findIndex((b) => b.junction_id === active.id);
    const newIndex = blocks.findIndex((b) => b.junction_id === over.id);

    // Optimistically reorder locally
    blocks = arrayMove(blocks, oldIndex, newIndex).map((b, i) => ({
      ...b,
      sort_order: i,
    }));

    // Persist to backend
    await blocksAPI.reorder(
      data.page.id,
      blocks.map((b) => ({ id: b.junction_id, sort_order: b.sort_order })),
    );
  }
</script>

<div class="max-w-3xl mx-auto py-8 px-4 space-y-6">
  <!-- Page header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-semibold">{data.page.title}</h1>
      <span class="text-xs text-muted-foreground uppercase tracking-wide"
        >{data.page.status}</span
      >
    </div>

    <!-- Add Block picker -->
    <div class="flex gap-2">
      <button
        onclick={() => addBlock('rich_text')}
        disabled={adding}
        class="px-3 py-1.5 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
      >
        + Rich Text
      </button>
      <button
        onclick={() => addBlock('hero')}
        disabled={adding}
        class="px-3 py-1.5 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
      >
        + Hero
      </button>
    </div>
  </div>

  <!-- Canvas -->
  {#if blocks.length === 0}
    <div
      class="border-2 border-dashed border-border rounded-lg p-16 text-center text-muted-foreground"
    >
      No blocks yet. Add one above.
    </div>
  {:else}
    <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <SortableContext items={blockIds} strategy={verticalListSortingStrategy}>
        <div class="space-y-3">
          {#each blocks as block (block.junction_id)}
            {@const Component = blockComponents[block.collection]}
            <SortableBlockItem
              id={block.junction_id}
              onRemove={() => removeBlock(block.junction_id)}
            >
              {#if Component}
                <Component {block} pageId={data.page.id} />
              {:else}
                <p class="text-sm text-muted-foreground">
                  Unknown block: {block.collection}
                </p>
              {/if}
            </SortableBlockItem>
          {/each}
        </div>
      </SortableContext>
    </DndContext>
  {/if}
</div>
