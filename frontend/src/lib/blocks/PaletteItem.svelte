<script lang="ts">
  let {
    label,
    onClick,
  }: {
    label: string;
    onClick: () => void;
  } = $props();
</script>

<button
  onclick={onClick}
  class="w-full text-left rounded-md border
    border-border bg-card px-3 py-2 text-sm
    font-medium hover:bg-accent transition-colors"
>
  {label}
</button>
