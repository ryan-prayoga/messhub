import { redirect } from '@sveltejs/kit';
import { ApiError, authServerApi } from '$lib/api/server';
import { clearAuthCookies } from '$lib/auth/session';

type ServerAuthContext = {
  fetch: typeof fetch;
  locals: App.Locals;
  cookies: import('@sveltejs/kit').Cookies;
};

export async function resolveServerUser(context: ServerAuthContext) {
  const { fetch, locals, cookies } = context;

  if (locals.user) {
    return locals.user;
  }

  if (!locals.token) {
    return null;
  }

  try {
    const response = await authServerApi.me(fetch, locals.token);
    locals.user = response.data;
    return locals.user;
  } catch (error) {
    if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
      clearAuthCookies(cookies);
      locals.token = null;
      locals.user = null;
      return null;
    }

    throw error;
  }
}

export async function requireServerUser(context: ServerAuthContext) {
  const user = await resolveServerUser(context);

  if (!context.locals.token || !user) {
    throw redirect(303, '/login');
  }

  return {
    token: context.locals.token,
    user
  };
}
