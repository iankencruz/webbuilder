import type { ResolvedBlock } from "$lib/types/blocks";

const BASE = "/api/pages";

export const blocksAPI = {
  async getResolved(fetch: typeof globalThis.fetch, pageId: number): Promise<ResolvedBlock[]> {
    const res = await fetch(`${BASE}/${pageId}/blocks/resolved`, { credentials: "include" });
    if (!res.ok) {
      throw new Error(`Failed to fetch resolved blocks for page ${pageId}: ${res.statusText}`);
    }
    return res.json();
  },

  async createBlock(collection: string, data: Record<string, unknown>): Promise<{ id: number }> {
    const res = await fetch(`/api/blocks/${collection}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ collection, data }),
    });
    if (!res.ok) {
      throw new Error(`Failed to create block in collection ${collection}: ${res.statusText}`);
    }
    return res.json();
  },

  async addToPage(pageId: number, blockId: number, collection: string, sortOrder: number) {
    const res = await fetch(`${BASE}/${pageId}/blocks`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        block_id: blockId,
        block_collection: collection,
        sort_order: sortOrder,
      }),
    });
    if (!res.ok) throw new Error("Failed to add block to page");
    return res.json();
  },

  async updateBlock(collection: string, blockId: number, data: Record<string, unknown>) {
    const res = await fetch(`/api/blocks/${collection}/${blockId}`, {
      method: "PUT",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to update block");
    return res.json();
  },

  async deleteFromPage(pageId: number, junctionId: number) {
    const res = await fetch(`${BASE}/${pageId}/blocks/${junctionId}`, {
      method: "DELETE",
      credentials: "include",
    });
    if (!res.ok) throw new Error("Failed to remove block from page");
  },

  async reorder(pageId: number, items: { id: number; sort_order: number }[]) {
    const res = await fetch(`${BASE}/${pageId}/blocks/reorder`, {
      method: "PATCH",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(items),
    });
    if (!res.ok) throw new Error("Failed to reorder blocks");
  },
};
