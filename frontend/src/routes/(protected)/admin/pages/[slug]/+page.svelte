<script lang="ts">
  import { blocksAPI } from '$lib/api/blocks';
  import type { Blocks, ResolvedBlock } from '$lib/types/blocks';
  import { untrack } from 'svelte';
  import {
    DndContext,
    closestCenter,
    useSensor,
    useSensors,
    PointerSensor,
  } from '@dnd-kit-svelte/core';
  import {
    SortableContext,
    verticalListSortingStrategy,
    arrayMove,
  } from '@dnd-kit-svelte/sortable';
  import SortableBlockItem from '$lib/blocks/SortableBlockItem.svelte';
  import BlockPicker from '$lib/blocks/BlockPicker.svelte';
  import { blockRegistry, type BlockCollection } from '$lib/blocks/registry';

  let { data } = $props();

  let pageBlocks = $state<ResolvedBlock[]>(untrack(() => data.pageBlocks));
  let pageId = $derived(data.page.id);

  let blocks = $derived<Blocks[]>(untrack(() => data.blocks));
  console.log('blocks', blocks);

  let blockIds = $derived(pageBlocks.map((b) => b.junction_id));

  let adding = $state(false);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 5 },
    }),
  );

  async function persistOrder() {
    await blocksAPI.reorder(
      pageId,
      pageBlocks.map((b) => ({
        id: b.junction_id,
        sort_order: b.sort_order,
      })),
    );
  }

  async function addBlock(collection: BlockCollection, atIndex?: number) {
    adding = true;
    try {
      const defaults: Record<BlockCollection, Record<string, unknown>> = {
        hero: {
          heading: 'New Hero',
          subheading: '',
          cta_label: '',
          cta_url: '',
        },
        richtext: {
          content: '',
          format: 'html',
        },
      };

      const { id: blockId } = await blocksAPI.createBlock(
        collection,
        defaults[collection],
      );

      const insertIndex = atIndex ?? pageBlocks.length;

      const junction = await blocksAPI.addToPage(
        pageId,
        blockId,
        collection,
        insertIndex,
      );

      const newBlock: ResolvedBlock = {
        junction_id: junction.id,
        sort_order: insertIndex,
        hide_block: false,
        collection,
        block_id: blockId,
        data: defaults[collection] as any,
      };

      const next = [...pageBlocks];
      next.splice(insertIndex, 0, newBlock);
      pageBlocks = next.map((b, i) => ({ ...b, sort_order: i }));

      if (insertIndex !== pageBlocks.length - 1) {
        await persistOrder();
      }
    } finally {
      adding = false;
    }
  }

  async function removeBlock(junctionId: number) {
    await blocksAPI.deleteFromPage(pageId, junctionId);
    pageBlocks = pageBlocks.filter((b) => b.junction_id !== junctionId);
  }

  async function handleDragEnd(event: any) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIndex = pageBlocks.findIndex((b) => b.junction_id === active.id);
    const newIndex = pageBlocks.findIndex((b) => b.junction_id === over.id);
    if (oldIndex === -1 || newIndex === -1) return;

    pageBlocks = arrayMove(pageBlocks, oldIndex, newIndex).map((b, i) => ({
      ...b,
      sort_order: i,
    }));

    await persistOrder();
  }
</script>

<div class="-m-6 flex h-[calc(100vh-3rem)]">
  <!-- Main canvas -->
  <div class="flex-1 overflow-y-auto p-6">
    <div class="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 class="text-2xl font-semibold">
          {data.page.title}
        </h1>
        <span
          class="text-xs text-muted-foreground
          uppercase tracking-wide"
        >
          {data.page.status}
        </span>
      </div>

      {#if pageBlocks.length === 0}
        <div
          class="border-2 border-dashed
          border-border rounded-lg p-16
          text-center text-muted-foreground space-y-4"
        >
          <p>No blocks yet.</p>
          <BlockPicker onAdd={(c) => addBlock(c, 0)} />
        </div>
      {:else}
        <DndContext
          {sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={blockIds}
            strategy={verticalListSortingStrategy}
          >
            <div class="space-y-4 bg-slate-100 p-4">
              {#each pageBlocks as block (block.junction_id)}
                {@const Component = blockRegistry[block.collection]}
                <SortableBlockItem
                  id={block.junction_id}
                  onRemove={() => removeBlock(block.junction_id)}
                >
                  {#if Component}
                    <Component {block} />
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
  </div>

  <!-- Right blocks panel -->
  <aside
    class="w-72 shrink-0 h-full flex flex-col
    bg-background border-l border-border"
  >
    <div class="h-12 flex items-center px-4 border-b border-border">
      <h2 class="text-sm font-semibold">Blocks</h2>
    </div>
    <div class="flex-1 overflow-y-auto p-4">
      <BlockPicker onAdd={(c) => addBlock(c)} expanded />
    </div>
  </aside>
</div>
