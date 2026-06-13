import type { User } from "$lib/types/types";

const state = $state({ user: null as User | null, loading: true });

export const auth = {
  get user() {
    return state.user;
  },
  get loading() {
    return state.loading;
  },

  async fetchMe() {
    try {
      const res = await fetch("/api/me", { credentials: "include" });
      state.user = res.ok ? await res.json() : null;
    } catch {
      state.user = null;
    } finally {
      state.loading = false;
    }
  },

  async logout() {
    await fetch("/api/auth/logout", { method: "POST", credentials: "include" });
    state.user = null;
  },
};
