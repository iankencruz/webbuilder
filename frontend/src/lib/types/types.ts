export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
}

export type Page = {
  id: number;
  title: string;
  slug: string;
  status: 'draft' | 'published' | 'archived';
  seo_title: string;
  seo_description: string;
  created_at: string;
  updated_at: string;
};
