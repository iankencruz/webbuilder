<script lang="ts">
  import { goto } from '$app/navigation';
  import { API_BASE } from '$lib/config';
  import { auth } from '$lib/hooks/auth.svelte';

  if (auth.user) {
    console.log('Authenticated: redirecting to /dashboard');
    goto('/admin/dashboard');
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
  <div
    class="w-full max-w-sm p-8 bg-white rounded-2xl shadow-sm border border-gray-100"
  >
    <div class="mb-8 text-center">
      <h1 class="text-2xl font-semibold text-gray-900">Welcome back</h1>
      <p class="mt-1 text-sm text-gray-500">Sign in to your account</p>
    </div>

    <a
      href="{API_BASE}/api/auth/zitadel/login"
      class="flex items-center justify-center gap-3 w-full px-4 py-2.5 rounded-lg border border-gray-200 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
    >
      <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none">
        <path
          d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
      Continue with Zitadel
    </a>
  </div>
</div>
