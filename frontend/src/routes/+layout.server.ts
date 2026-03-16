import { error as svelteError } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { ApiError, authServerApi } from '$lib/api/server';
import { clearAuthCookies } from '$lib/auth/session';
import { requireAuth } from '$lib/auth/guard';

export const load: LayoutServerLoad = async ({ cookies, fetch, locals, url }) => {
  let user = null;

  if (locals.token) {
    try {
      const response = await authServerApi.me(fetch, locals.token);
      user = response.data;
    } catch (error) {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        clearAuthCookies(cookies);
        locals.token = null;
      } else {
        throw svelteError(503, {
          message: error instanceof Error ? error.message : 'Failed to verify session'
        });
      }
    }
  }

  locals.user = user;
  requireAuth(url.pathname, user);

  return {
    user
  };
};
