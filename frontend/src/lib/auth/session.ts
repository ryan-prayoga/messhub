import type { Cookies } from '@sveltejs/kit';

export const AUTH_COOKIE_KEYS = {
  token: 'mh_access_token'
} as const;

export const PUBLIC_ROUTES = ['/login'];
export const SHELLLESS_ROUTES = ['/login', '/offline'];
export const ALWAYS_AVAILABLE_ROUTES = ['/offline'];

export function buildAuthCookieOptions(url: URL) {
  return {
    path: '/',
    httpOnly: true,
    sameSite: 'lax' as const,
    secure: url.protocol === 'https:',
    maxAge: 60 * 60 * 72
  };
}

export function clearAuthCookies(cookies: Cookies) {
  cookies.delete(AUTH_COOKIE_KEYS.token, {
    path: '/'
  });
}
