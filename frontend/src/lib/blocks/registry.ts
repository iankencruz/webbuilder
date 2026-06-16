// registry.ts
import type { Component } from "svelte";
import HeroBlockEditor from "$lib/blocks/HeroBlockEditor.svelte";
import RichTextBlockEditor from "$lib/blocks/RichTextBlockEditor.svelte";

export const blockRegistry: Record<string, Component<any>> = {
  hero: HeroBlockEditor,
  richtext: RichTextBlockEditor,
};
