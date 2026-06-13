import type { Component } from "svelte";
import HeroBlockEditor from "$lib/blocks/HeroBlockEditor.svelte";
import RichTextBlockEditor from "$lib/blocks/RichTextBlockEditor.svelte";

const registry = {
  hero: HeroBlockEditor,
  rich_text: RichTextBlockEditor,
} satisfies Record<string, Component<any>>;

export const blockRegistry = registry;

export type BlockCollection = keyof typeof blockRegistry;
