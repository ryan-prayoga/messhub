import { redirect } from '@sveltejs/kit';
import { PUBLIC_ROUTES } from '$lib/auth/session';

export function requireAuth(pathname: string, user: App.Locals['user']) {
  if (!PUBLIC_ROUTES.includes(pathname) && !user) {
    throw redirect(303, '/login');
  }

  if (PUBLIC_ROUTES.includes(pathname) && user) {
    throw redirect(303, '/dashboard');
  }
}
