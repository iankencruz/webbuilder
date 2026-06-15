export const ssr = false;

// src/routes/(admin)/+layout.ts
import { redirect } from "@sveltejs/kit";
import { auth } from "$lib/hooks/auth.svelte";

export const load = async () => {
  if (!auth.user) {
    await auth.fetchMe();
  }

  if (!auth.user) {
    throw redirect(302, "/login");
  }

  return { user: auth.user };
};
