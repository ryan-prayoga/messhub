import { redirect } from '@sveltejs/kit';
import { ALWAYS_AVAILABLE_ROUTES, PUBLIC_ROUTES } from '$lib/auth/session';

export function requireAuth(pathname: string, user: App.Locals['user']) {
  if (ALWAYS_AVAILABLE_ROUTES.includes(pathname)) {
    return;
  }

  if (!PUBLIC_ROUTES.includes(pathname) && !user) {
    throw redirect(303, '/login');
  }

  if (PUBLIC_ROUTES.includes(pathname) && user) {
    throw redirect(303, '/dashboard');
  }
}
