import type { Handle } from '@sveltejs/kit';
import { AUTH_COOKIE_KEYS } from '$lib/auth/session';

export const handle: Handle = async ({ event, resolve }) => {
  event.locals.token = event.cookies.get(AUTH_COOKIE_KEYS.token) ?? null;
  event.locals.user = null;

  return resolve(event);
};
