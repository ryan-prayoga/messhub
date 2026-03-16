import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { AUTH_COOKIE_KEYS } from '$lib/auth/session';

export const POST: RequestHandler = async ({ cookies }) => {
  const options = {
    path: '/',
    expires: new Date(0)
  };

  for (const key of Object.values(AUTH_COOKIE_KEYS)) {
    cookies.set(key, '', options);
  }

  throw redirect(303, '/login');
};
