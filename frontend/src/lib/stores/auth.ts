import { writable } from 'svelte/store';
import type { AuthStatus, SessionUser } from '$lib/api/types';

type AuthState = {
  status: AuthStatus;
  user: SessionUser | null;
};

function createAuthState() {
  const { subscribe, set } = writable<AuthState>({
    status: 'loading',
    user: null
  });

  return {
    subscribe,
    sync(user: SessionUser | null) {
      set({
        status: user ? 'authenticated' : 'unauthenticated',
        user
      });
    }
  };
}

export const authState = createAuthState();
