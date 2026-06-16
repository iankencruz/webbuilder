export interface HeroBlock {
  id: number;
  heading: string;
  subheading: string | null;
  cta_label: string | null;
  cta_url: string | null;
  created_at: string;
  updated_at: string;
}

export interface RichTextBlock {
  id: number;
  content: string;
  format: string;
  created_at: string;
  updated_at: string;
}

export type BlockData = HeroBlock | RichTextBlock;

export interface ResolvedBlock {
  junction_id: number;
  sort_order: number;
  hide_block: boolean;
  collection: BlockCollection;
  block_id: number;
  data: BlockData;
}

export interface Blocks {
  code: string;
  label: string;
  description: string;
  created_at: string;
}
export type BlockCollection = string;
