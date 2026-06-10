<script lang="ts">
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import type { Snippet } from 'svelte';

  type Props = {
    open: boolean;
    header: Snippet;
    footer: Snippet;
    children: Snippet;
  };

  let { open = $bindable(), header, footer, children }: Props = $props();
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content>
    <!-- Header -->
    {#if header}
      <AlertDialog.Header>
        {@render header()}
      </AlertDialog.Header>
    {/if}
    <!-- children -->
    {#if children}
      <div class="py-2 text-sm text-muted-foreground">
        {@render children()}
      </div>
    {/if}

    <!-- Footer -->
    {#if footer}
      <AlertDialog.Footer>
        {@render footer()}
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
