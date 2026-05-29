<script lang="ts">
  import { auth } from '$lib/auth.svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { navGroups } from '$lib/data/dashboard';

  let { children } = $props();
  let sidebarOpen = $state(true);

  $effect(() => {
    if (!auth.loading && !auth.user) {
      goto('/login');
    }
  });

  const isActive = (href: string) => page.url.pathname === href;
</script>

{#if auth.user}
  <div class="flex h-screen bg-background overflow-hidden">
    <!-- Sidebar -->
    <aside
      class="flex flex-col bg-sidebar border-r border-sidebar-border transition-all duration-300 ease-in-out shrink-0 {sidebarOpen
        ? 'w-56'
        : 'w-16'}"
    >
      <!-- Logo -->
      <div
        class="flex items-center gap-3 px-4 h-14 border-b border-sidebar-border shrink-0"
      >
        <div
          class="w-7 h-7 rounded-lg bg-sidebar-primary flex items-center justify-center shrink-0"
        >
          <span class="text-sidebar-primary-foreground text-xs font-bold"
            >W</span
          >
        </div>
        {#if sidebarOpen}
          <span
            class="font-semibold text-sidebar-foreground text-sm tracking-tight"
            >Webbuilder</span
          >
        {/if}
      </div>

      <!-- Nav -->
      <nav class="flex-1 overflow-y-auto py-4 px-2 space-y-6">
        {#each navGroups as group}
          <div>
            {#if sidebarOpen}
              <p
                class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground px-2 mb-1"
              >
                {group.label}
              </p>
            {/if}
            <ul class="space-y-0.5">
              {#each group.items as item}
                <li>
                  <a
                    href={item.href}
                    class="flex items-center gap-3 px-2 py-2 rounded-md text-sm transition-colors
                    {isActive(item.href)
                      ? 'bg-sidebar-accent text-sidebar-accent-foreground font-medium'
                      : 'text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'}"
                    class:justify-center={!sidebarOpen}
                  >
                    <span class="text-base shrink-0 w-5 text-center"
                      >{item.icon}</span
                    >
                    {#if sidebarOpen}
                      <span>{item.label}</span>
                    {/if}
                  </a>
                </li>
              {/each}
            </ul>
          </div>
        {/each}
      </nav>

      <!-- User -->
      <div class="border-t border-sidebar-border p-3 shrink-0">
        <div class="flex items-center gap-3">
          <div
            class="w-8 h-8 rounded-full bg-sidebar-accent flex items-center justify-center shrink-0"
          >
            <span class="text-sidebar-accent-foreground text-xs font-semibold">
              {auth.user.first_name?.[0] ?? auth.user.email[0].toUpperCase()}
            </span>
          </div>
          {#if sidebarOpen}
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-sidebar-foreground truncate">
                {auth.user.first_name}
                {auth.user.last_name}
              </p>
              <p class="text-xs text-muted-foreground truncate">
                {auth.user.email}
              </p>
            </div>
          {/if}
        </div>
      </div>
    </aside>

    <!-- Main -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Topbar -->
      <header
        class="h-14 bg-card border-b border-border flex items-center gap-4 px-6 shrink-0"
      >
        <button
          onclick={() => (sidebarOpen = !sidebarOpen)}
          class="text-muted-foreground hover:text-foreground transition-colors"
          aria-label="Toggle sidebar"
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
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>

        <!-- Breadcrumb -->
        <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
          <span>Admin</span>
          <span>›</span>
          <span class="text-foreground font-medium capitalize">
            {page.url.pathname.split('/').at(-1)}
          </span>
        </div>

        <div class="ml-auto flex items-center gap-3">
          <!-- Search -->
          <div class="relative hidden sm:block">
            <input
              type="text"
              placeholder="Search..."
              class="w-48 pl-8 pr-3 py-1.5 text-sm bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent text-foreground placeholder:text-muted-foreground"
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

          <!-- Notifications -->
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

          <!-- Logout -->
          <button
            onclick={() => auth.logout().then(() => goto('/login'))}
            class="w-8 h-8 rounded-full bg-primary flex items-center justify-center hover:opacity-90 transition-opacity"
            aria-label="Logout"
          >
            <span class="text-primary-foreground text-xs font-semibold">
              {auth.user.first_name?.[0] ?? '?'}
            </span>
          </button>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-6">
        {@render children()}
      </main>
    </div>
  </div>
{/if}
