<script lang="ts">
  import { blocksAPI } from '$lib/api/blocks';
  import type { ResolvedBlock, RichTextBlock } from '$lib/types/blocks';
  import { untrack } from 'svelte';

  let { block }: { block: ResolvedBlock } = $props();

  let content = $state(
    untrack(() => (block.data as RichTextBlock).content ?? ''),
  );
  let format = $state(
    untrack(() => (block.data as RichTextBlock).format ?? 'html'),
  );

  type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';
  let saveStatus = $state<SaveStatus>('idle');
  let debounceTimer: ReturnType<typeof setTimeout>;

  $effect(() => {
    // Track these reactive values
    void content;

    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(async () => {
      saveStatus = 'saving';
      try {
        (await blocksAPI.updateBlock('richtext', block.block_id, {
          content,
          format,
        }),
          (saveStatus = 'saved'));
        setTimeout(() => (saveStatus = 'idle'), 2000);
      } catch (e) {
        saveStatus = 'error';
      }
    }, 800);

    return () => clearTimeout(debounceTimer);
  });
</script>

<div class="space-y-2">
  <div class="flex items-center justify-between">
    <span
      class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
      >Rich Text</span
    >
    {#if saveStatus === 'saving'}
      <span class="text-xs text-muted-foreground">Saving…</span>
    {:else if saveStatus === 'saved'}
      <span class="text-xs text-green-500">Saved</span>
    {:else if saveStatus === 'error'}
      <span class="text-xs text-destructive">Error saving</span>
    {/if}
  </div>

  <textarea
    bind:value={content}
    rows={4}
    placeholder="Enter content…"
    class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
  ></textarea>
</div>
