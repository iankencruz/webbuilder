<script lang="ts">
  import { auth } from '$lib/hooks/auth.svelte';
  import { page } from '$app/state';
  import * as Sidebar from '$lib/components/ui/sidebar/index';
  import AppSidebar from '$lib/components/app-sidebar.svelte';

  let { children } = $props();
</script>

{#if auth.user}
  <Sidebar.Provider>
    <AppSidebar />
    <Sidebar.Inset>
      <header
        class="border-b h-12 bg-card flex items-center gap-4 px-4 shrink-0"
      >
        <Sidebar.Trigger />
        <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
          <span>Admin</span>
          <span>›</span>
          <span class="text-foreground font-medium capitalize">
            {page.url.pathname.split('/').at(-1)}
          </span>
        </div>
        <div class="ml-auto flex items-center gap-3">
          <div class="relative hidden sm:block">
            <input
              type="text"
              placeholder="Search..."
              class="w-48 pl-8 pr-3 py-1.5 text-sm bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring text-foreground placeholder:text-muted-foreground"
            />
            <svg
              class="absolute left-2.5 top-2 w-3.5 h-3.5 text-muted-foreground"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0"
              />
            </svg>
          </div>
          <button
            type="button"
            aria-label="Notifications"
            class="relative text-muted-foreground hover:text-foreground transition-colors"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
              />
            </svg>
            <span
              class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-destructive rounded-full"
            ></span>
          </button>
        </div>
      </header>
      <main class="flex-1 overflow-y-auto p-6">
        {@render children()}
      </main>
    </Sidebar.Inset>
  </Sidebar.Provider>
{/if}
