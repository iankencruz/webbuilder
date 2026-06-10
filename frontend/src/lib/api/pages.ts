import type { Page } from '$lib/types';

const BASE = '/api/pages';

export const pagesAPI = {
  async list(fetch: typeof globalThis.fetch): Promise<Page[]> {
    const res = await fetch(BASE, { credentials: 'include' });
    if (!res.ok) throw new Error('Failed to fetch pages');
    return res.json();
  },

  async get(id: number): Promise<Page> {
    const res = await fetch(`${BASE}/${id}`, { credentials: 'include' });
    if (!res.ok) throw new Error('Failed to fetch page');
    return res.json();
  },

  async create(title: string, slug: string): Promise<Page> {
    const res = await fetch(BASE, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, slug, status: 'draft' }),
    });
    if (!res.ok) throw new Error('Failed to create page');
    return res.json();
  },

  async update(
    id: number,
    data: Partial<Omit<Page, 'id' | 'created_at' | 'updated_at'>>,
  ): Promise<Page> {
    const res = await fetch(`${BASE}/${id}`, {
      method: 'PUT',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to update page');
    return res.json();
  },

  async delete(slug: string): Promise<void> {
    const res = await fetch(`${BASE}/${slug}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!res.ok) throw new Error('Failed to delete page');
  },
};
