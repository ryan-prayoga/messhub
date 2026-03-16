import type { Handle } from '@sveltejs/kit';
import { AUTH_COOKIE_KEYS } from '$lib/auth/session';

const contentSecurityPolicy = [
  "default-src 'self'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
  "object-src 'none'",
  "script-src 'self' 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data: https: http:",
  "font-src 'self' data:",
  "connect-src 'self' https: http:",
  "manifest-src 'self'",
  "worker-src 'self' blob:"
].join('; ');

export const handle: Handle = async ({ event, resolve }) => {
  event.locals.token = event.cookies.get(AUTH_COOKIE_KEYS.token) ?? null;
  event.locals.user = null;

  const response = await resolve(event);
  response.headers.set('X-Frame-Options', 'DENY');
  response.headers.set('X-Content-Type-Options', 'nosniff');
  response.headers.set('X-XSS-Protection', '1; mode=block');
  response.headers.set('Content-Security-Policy', contentSecurityPolicy);

  return response;
};
