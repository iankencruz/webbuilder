<script lang="ts">
  import * as Item from '$lib/components/ui/item/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import BadgeCheckIcon from '@lucide/svelte/icons/badge-check';
  import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
  let {
    label,
    description,
    onClick,
  }: {
    label: string;
    description: string;
    onClick: () => void;
  } = $props();
</script>

<Button
  variant="outline"
  onclick={onClick}
  class="w-full text-ellipsis  p-4 py-8"
>
  <Item.Root>
    <Item.Content>
      <Item.Title>{label}</Item.Title>
      <Item.Description class="w-ffull">{description}</Item.Description>
    </Item.Content>
  </Item.Root>
</Button>
