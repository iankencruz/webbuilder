<script lang="ts">
  import './layout.css';
  import favicon from '$lib/assets/favicon.svg';
  import { onMount } from 'svelte';
  import { Toaster } from '$lib/components/ui/sonner';
  import { auth } from '$lib/hooks/auth.svelte';

  let { children } = $props();

  onMount(() => auth.fetchMe());
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<Toaster />
{#if auth.loading}
  <div>
    <span>Loading...</span>
  </div>
{:else}
  {@render children()}
{/if}
