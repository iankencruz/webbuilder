import type { PageLoad } from './$types';
import { pagesAPI } from '$lib/api/pages';

export const load: PageLoad = async ({ fetch }) => {
  const pages = await pagesAPI.list(fetch);
  return { pages };
};
