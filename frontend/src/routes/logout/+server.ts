import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { clearAuthCookies } from '$lib/auth/session';

export const POST: RequestHandler = async ({ cookies }) => {
  clearAuthCookies(cookies);

  throw redirect(303, '/login');
};
