import type { LayoutServerLoad } from './$types';
import { requireAuth } from '$lib/auth/guard';

export const load: LayoutServerLoad = async ({ locals, url }) => {
  requireAuth(url.pathname, locals.user);

  return {
    user: locals.user
  };
};
