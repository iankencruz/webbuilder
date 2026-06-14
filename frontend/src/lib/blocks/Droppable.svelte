<script lang="ts">
  import { useDroppable } from '@dnd-kit-svelte/core';

  let {
    id,
    children,
  }: {
    id: string;
    children: import('svelte').Snippet;
  } = $props();

  const { setNodeRef, isOver } = useDroppable({ id: () => id });
</script>

<div
  {@attach setNodeRef}
  class="rounded-lg transition-colors
    {isOver.current ? 'ring-2 ring-primary/40 bg-primary/5' : ''}"
>
  {@render children()}
</div>
