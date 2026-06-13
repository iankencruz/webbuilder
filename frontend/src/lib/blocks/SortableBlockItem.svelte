<script lang="ts">
  import { useSortable } from '@dnd-kit-svelte/sortable';
  import { CSS } from '@dnd-kit-svelte/utilities';

  let { id, onRemove, children } = $props();

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: () => id });

  let style = $derived(
    CSS.Transform.toString(transform.current)
      ? `transform: ${CSS.Transform.toString(transform.current)}; transition: ${transition.current};`
      : '',
  );
</script>

<div
  {@attach setNodeRef}
  {style}
  class="relative group rounded-lg border border-border bg-card {isDragging.current
    ? 'opacity-50 z-10'
    : ''}"
>
  <!-- Drag handle -->
  <div
    {...attributes.current}
    {...listeners.current}
    class="absolute left-2 top-3 cursor-grab opacity-0 group-hover:opacity-100 transition-opacity text-muted-foreground"
  >
    ⠿
  </div>

  <!-- Remove button -->
  <button
    onclick={onRemove}
    class="absolute right-2 top-2 opacity-0 group-hover:opacity-100 transition-opacity text-muted-foreground hover:text-destructive text-xs"
  >
    ✕
  </button>

  <!-- Block editor content -->
  <div class="pl-8 pr-8 py-3">
    {@render children()}
  </div>
</div>
