<script lang="ts">
  import { pagesAPI } from '$lib/api/pages';
  import { goto } from '$app/navigation';
  import type { Page } from '$lib/types/types.js';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import * as Dialog from '$lib/components/ui/dialog';
  import Label from '$lib/components/ui/label/label.svelte';
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import ActionDialog from '$lib/components/ActionDialog.svelte';
  import AlertDialogDescription from '$lib/components/ui/alert-dialog/alert-dialog-description.svelte';
  import DataTable from '$lib/components/DataTable/DataTable.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';

  let { data } = $props();

  // Handle local mutations using a mutable state tracking data updates
  let localPages = $derived<Page[]>(data.pages ?? []);

  // React to fresh parent data updates gracefully
  $effect(() => {
    localPages = data?.pages ?? [];
  });

  let open = $state(false);
  let title = $state('');
  let slug = $state('');
  let creating = $state(false);
  let error = $state('');

  // Delete dialog states
  let deleteDialogOpen = $state(false);
  let slugToDelete = $state<string | null>(null);
  let status = $state('');

  $effect(() => {
    slug = title
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '');
  });

  async function handleCreate() {
    if (!title || !slug) return;
    creating = true;
    error = '';
    try {
      const page = await pagesAPI.create(title, slug);
      localPages = [page, ...localPages];
      open = false;
      title = '';
      slug = '';
      goto(`/admin/pages/${page.slug}`);
    } catch (e) {
      error = e instanceof Error ? e.message : 'Something went wrong';
    } finally {
      creating = false;
    }
  }

  // Step 1: Open dialog and cache targeted record ID
  function triggerDeleteConfirm(slug: string) {
    slugToDelete = slug;
    deleteDialogOpen = true;
  }

  // Step 2: User confirmed via the custom modal trigger
  async function confirmDelete() {
    if (slugToDelete === null) return;

    const slugID = slugToDelete;
    // Reset state early to smoothly close modal
    deleteDialogOpen = false;
    slugToDelete = null;

    try {
      status = 'Deleting item...';
      await pagesAPI.delete(slugID);
      localPages = localPages.filter((p) => p.slug !== slugID);
      status = 'Item deleted successfully!';
    } catch (e) {
      status = 'Failed to delete item.';
    }
  }

  function cancelDelete() {
    deleteDialogOpen = false;
    slugToDelete = null;
    status = 'Item deletion cancelled.';
  }
</script>

{#snippet statusBadge({ value }: { value: unknown; row: Page })}
  {@const status = value as string}
  <Badge
    variant={status === 'published'
      ? 'success'
      : status === 'draft'
        ? 'draft'
        : status === 'archived'
          ? 'archived'
          : 'outline'}
  >
    {status}
  </Badge>
{/snippet}

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-semibold text-foreground">Pages</h1>
      <p class="text-sm text-muted-foreground mt-1">
        Manage your website pages
      </p>
    </div>
    <Button onclick={() => (open = true)}>New Page</Button>
  </div>

  <!-- Pages list -->
  {#if localPages.length === 0}
    <div class="flex flex-col items-center justify-center py-24 text-center">
      <p class="text-muted-foreground text-sm">No pages yet</p>
      <Button class="mt-4" onclick={() => (open = true)}>
        Create your first page
      </Button>
    </div>
  {:else}
    <DataTable
      data={localPages}
      hiddenColumns={['id', 'seo_title', 'seo_description']}
      filterColumn="title"
      filterPlaceholder="Search pages..."
      columnOverrides={{ status: { cellSnippet: statusBadge } }}
      onRowAction={(action, row) => {
        if (action === 'edit') {
          goto(`/admin/pages/${row.slug}`);
        } else if (action === 'delete') {
          triggerDeleteConfirm(row.slug as string);
        }
      }}
    />
  {/if}
</div>

<!-- Create page dialog -->
<ActionDialog bind:open>
  {#snippet header()}
    <Dialog.Title>Create New Page</Dialog.Title>
    <Dialog.Description>
      Enter a title for your new page. The slug will be generated automatically.
    </Dialog.Description>
  {/snippet}

  {#snippet children()}
    <div class="space-y-4 py-2">
      <div class="space-y-1.5">
        <Label for="title">Title</Label>
        <Input id="title" bind:value={title} placeholder="About Us" autofocus />
      </div>
      <div class="space-y-1.5">
        <Label for="slug">Slug</Label>
        <Input id="slug" bind:value={slug} placeholder="about-us" />
      </div>
      {#if error}
        <p class="text-sm text-destructive">{error}</p>
      {/if}
    </div>
  {/snippet}

  {#snippet footer()}
    <Button
      variant="outline"
      onclick={() => (open = false)}
      disabled={creating}
    >
      Cancel
    </Button>
    <Button onclick={handleCreate} disabled={!title || creating}>
      {creating ? 'Creating...' : 'Create'}
    </Button>
  {/snippet}
</ActionDialog>

<!-- Confirmation Alert Dialog -->
<ActionDialog bind:open={deleteDialogOpen}>
  {#snippet header()}
    <AlertDialog.Title>Confirm Deletion</AlertDialog.Title>
    <AlertDialogDescription>
      Are you sure you want to delete this page? This action cannot be undone.
    </AlertDialogDescription>
  {/snippet}

  {#snippet children()}
    <p class="text-sm text-muted-foreground">
      This will permanently delete the page with slug: <strong
        >{slugToDelete}</strong
      >.
    </p>
  {/snippet}

  {#snippet footer()}
    <AlertDialog.Cancel onclick={cancelDelete}>Cancel</AlertDialog.Cancel>
    <AlertDialog.Action
      class="bg-destructive text-white hover:bg-destructive/90"
      onclick={confirmDelete}
    >
      Delete
    </AlertDialog.Action>
  {/snippet}
</ActionDialog>
