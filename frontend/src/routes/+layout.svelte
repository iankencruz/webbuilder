<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
  import { onMount } from 'svelte';
  import { auth } from '$lib/auth.svelte';

	let { children } = $props();


  onMount(() => auth.fetchMe())
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

{#if auth.loading}
  <div>
  <span>Loading...</span>
  </div>
{:else}
{@render children()}
{/if}
