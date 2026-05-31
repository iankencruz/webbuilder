<script lang="ts">
  import LogOutIcon from '@lucide/svelte/icons/log-out';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/hooks/auth.svelte';
  import { navGroups } from '$lib/data/dashboard';
  import * as Sidebar from '$lib/components/ui/sidebar/index';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
  import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';

  const isActive = (href: string) => page.url.pathname === href;
</script>

<Sidebar.Root collapsible="icon">
  <Sidebar.Header class="p-0 border-b border-sidebar-border">
    <Sidebar.Menu class="p-2">
      <Sidebar.MenuItem
        class="group-data-[collapsible=icon]:flex group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
      >
        <Sidebar.MenuButton
          size="lg"
          class="rounded-none pointer-events-none hover:bg-transparent data-[active=true]:bg-transparent"
        >
          <div
            class="w-8 h-8 rounded-lg bg-foreground flex items-center justify-center shrink-0"
          >
            <span class="text-background text-sm font-bold">W</span>
          </div>
          <span class="font-semibold text-foreground text-base tracking-tight">
            Webbuilder
          </span>
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Header>
  <Sidebar.Content>
    {#each navGroups as group}
      <Sidebar.Group>
        <Sidebar.GroupLabel
          class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground group-data-[collapsible=icon]:hidden"
        >
          {group.label}
        </Sidebar.GroupLabel>
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each group.items as item (item.title)}
              {@const Icon = item.icon}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton isActive={isActive(item.href)}>
                  {#snippet child({ props })}
                    <a href={item.href} {...props}>
                      <Icon />
                      <span>{item.title}</span>
                    </a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/each}
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    {/each}
  </Sidebar.Content>

  <Sidebar.Footer class="border-t border-sidebar-border">
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            {#snippet child({ props })}
              <Sidebar.MenuButton
                size="lg"
                {...props}
                class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                <div
                  class="w-8 h-8 rounded-full bg-muted flex items-center justify-center shrink-0"
                >
                  <span class="text-foreground text-xs font-semibold">
                    {auth.user?.first_name?.[0] ??
                      auth.user?.email?.[0]?.toUpperCase() ??
                      '?'}
                  </span>
                </div>
                <div class="flex-1 min-w-0 text-left">
                  <p
                    class="text-sm font-medium text-foreground truncate leading-tight"
                  >
                    {auth.user?.first_name}
                    {auth.user?.last_name}
                  </p>
                  <p class="text-xs text-muted-foreground truncate">
                    {auth.user?.email}
                  </p>
                </div>
                <ChevronUpIcon class="ms-auto" />
              </Sidebar.MenuButton>
            {/snippet}
          </DropdownMenu.Trigger>
          <DropdownMenu.Content
            side="top"
            class="w-(--bits-dropdown-menu-anchor-width)"
          >
            <DropdownMenu.Item>
              <span>Profile</span>
            </DropdownMenu.Item>
            <DropdownMenu.Item>
              <span>Billing</span>
            </DropdownMenu.Item>
            <DropdownMenu.Separator />
            <DropdownMenu.Item
              onclick={() => auth.logout().then(() => goto('/login'))}
            >
              <LogOutIcon />
              <span>Log out</span>
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Footer>

  <Sidebar.Rail />
</Sidebar.Root>
