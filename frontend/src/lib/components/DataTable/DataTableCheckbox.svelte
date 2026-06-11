<script lang="ts">
  import type { ComponentProps } from 'svelte';
  import Checkbox from '../ui/checkbox/checkbox.svelte';

  let {
    checked = false,
    onCheckedChange = (v: boolean) => (checked = v),
    ...restProps
  }: ComponentProps<typeof Checkbox> = $props();
</script>

<Checkbox
  bind:checked={() => checked, (v) => onCheckedChange(v)}
  {...restProps}
/>
