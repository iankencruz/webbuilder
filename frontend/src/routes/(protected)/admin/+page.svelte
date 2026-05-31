<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/hooks/auth.svelte';

  $effect(() => {
    if (!auth.loading && !auth.user) {
      goto('/login');
    }
  });
</script>
