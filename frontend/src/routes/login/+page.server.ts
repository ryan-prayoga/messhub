import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, authServerApi } from '$lib/api/server';
import { AUTH_COOKIE_KEYS, buildAuthCookieOptions } from '$lib/auth/session';

export const load: PageServerLoad = async ({ locals }) => {
  if (locals.user) {
    throw redirect(303, '/dashboard');
  }

  return {};
};

export const actions: Actions = {
  default: async ({ cookies, fetch, request, url }) => {
    const formData = await request.formData();
    const email = String(formData.get('email') || '').trim();
    const password = String(formData.get('password') || '');
    const values = { email };

    if (!email || !password) {
      return fail(400, { message: 'Email dan password wajib diisi.', values });
    }

    try {
      const payload = await authServerApi.login(fetch, { email, password });
      const { token } = payload.data;

      cookies.set(AUTH_COOKIE_KEYS.token, token, buildAuthCookieOptions(url));
    } catch (error) {
      if (error instanceof ApiError) {
        return fail(error.status, {
          message: error.message,
          values
        });
      }

      return fail(503, {
        message: 'Backend auth belum terjangkau. Cek service API di VPS atau env frontend.',
        values
      });
    }

    throw redirect(303, '/dashboard');
  }
};
