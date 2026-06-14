<script lang="ts">
  import BlockPicker from './BlockPicker.svelte';
  import type { BlockCollection } from './registry';

  let {
    onAdd,
  }: {
    onAdd: (collection: BlockCollection) => void;
  } = $props();
</script>

<div
  class="group relative h-3 -my-1.5 flex items-center py-4.5
  justify-center"
>
  <div
    class="absolute inset-x-0 top-1/2 h-px bg-border
    opacity-80 group-hover:opacity-100 transition-opacity"
  ></div>
  <div
    class="relative z-10 opacity-80
    group-hover:opacity-100 transition-opacity"
  >
    <BlockPicker {onAdd} compact />
  </div>
</div>
