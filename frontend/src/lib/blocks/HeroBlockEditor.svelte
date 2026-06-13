<script lang="ts">
  import type { ResolvedBlock, HeroBlock } from '$lib/types/blocks';
  import { blocksAPI } from '$lib/api/blocks';
  import { untrack } from 'svelte';

  let { block }: { block: ResolvedBlock } = $props();

  let heading = $state(untrack(() => (block.data as HeroBlock).heading ?? ''));
  let subheading = $state(
    untrack(() => (block.data as HeroBlock).subheading ?? ''),
  );
  let ctaLabel = $state(
    untrack(() => (block.data as HeroBlock).cta_label ?? ''),
  );
  let ctaUrl = $state(untrack(() => (block.data as HeroBlock).cta_url ?? ''));

  type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';
  let saveStatus = $state<SaveStatus>('idle');
  let debounceTimer: ReturnType<typeof setTimeout>;

  $effect(() => {
    // Track all reactive values
    void (heading + subheading + ctaLabel + ctaUrl);

    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(async () => {
      saveStatus = 'saving';
      try {
        await blocksAPI.updateBlock('hero', block.block_id, {
          heading,
          subheading,
          cta_label: ctaLabel,
          cta_url: ctaUrl,
        });
        saveStatus = 'saved';
        setTimeout(() => (saveStatus = 'idle'), 2000);
      } catch {
        saveStatus = 'error';
      }
    }, 800);

    return () => clearTimeout(debounceTimer);
  });
</script>

<div class="space-y-3">
  <div class="flex items-center justify-between">
    <span
      class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
      >Hero</span
    >
    {#if saveStatus === 'saving'}
      <span class="text-xs text-muted-foreground">Saving…</span>
    {:else if saveStatus === 'saved'}
      <span class="text-xs text-green-500">Saved</span>
    {:else if saveStatus === 'error'}
      <span class="text-xs text-destructive">Error saving</span>
    {/if}
  </div>

  <div class="space-y-2">
    <input
      bind:value={heading}
      type="text"
      placeholder="Heading"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-ring"
    />
    <input
      bind:value={subheading}
      type="text"
      placeholder="Subheading"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
    />
    <div class="flex gap-2">
      <input
        bind:value={ctaLabel}
        type="text"
        placeholder="CTA Label"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
      />
      <input
        bind:value={ctaUrl}
        type="text"
        placeholder="CTA URL"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
      />
    </div>
  </div>
</div>
