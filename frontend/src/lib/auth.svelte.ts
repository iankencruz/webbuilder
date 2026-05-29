type User = {
  id: number;
  email: string;
} | null;

const authState = $state({ user: null as User, loading: true });

export const auth = {
  get user() {
    return authState.user;
  },
  get loading() {
    return authState.loading;
  },

  async fetchMe() {
    try {
      const res = await fetch("/api/me", { credentials: "include" });
      authState.user = res.ok ? await res.json() : null;
    } catch {
      authState.user = null;
    } finally {
      authState.loading = false;
    }
  },

  async logout() {
    await fetch("/api/auth/logout", { method: "POST", credentials: "include" });
    authState.user = null;
  },
};
