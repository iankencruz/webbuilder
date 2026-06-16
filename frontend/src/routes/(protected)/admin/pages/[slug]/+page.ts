import type { PageLoad } from "./$types";
import { pagesAPI } from "$lib/api/pages";
import { blocksAPI } from "$lib/api/blocks";

export const load: PageLoad = async ({ fetch, params }) => {
  const page = await pagesAPI.get(fetch, params.slug);
  const blocks = await blocksAPI.listBlocks(fetch);
  const pageBlocks = await blocksAPI.getResolved(fetch, page.id);
  return { page, blocks, pageBlocks };
};
