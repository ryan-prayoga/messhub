import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { authServerApi } from '$lib/api/server';
import { AUTH_COOKIE_KEYS, buildAuthCookieOptions } from '$lib/auth/session';
import { toApiFailureState } from '$lib/server/api-errors';

export const load: PageServerLoad = async ({ locals, url }) => {
  if (locals.user) {
    throw redirect(303, '/dashboard');
  }

  return {
    notice: url.searchParams.get('message') === 'password-changed'
      ? 'Password berhasil diganti. Silakan login lagi dengan password baru.'
      : null
  };
};

export const actions: Actions = {
  default: async ({ cookies, fetch, request, url }) => {
    const formData = await request.formData();
    const identifier = String(formData.get('identifier') || '').trim();
    const password = String(formData.get('password') || '');
    const values = { identifier };

    if (!identifier || !password) {
      return fail(400, { message: 'Email atau username dan password wajib diisi.', values });
    }

    try {
      const payload = await authServerApi.login(fetch, { identifier, password });
      const { token } = payload.data;

      cookies.set(AUTH_COOKIE_KEYS.token, token, buildAuthCookieOptions(url));
    } catch (error) {
      const failure = toApiFailureState(error, 'Tidak dapat memproses login saat ini.');
      return fail(failure.status, {
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }

    throw redirect(303, '/dashboard');
  }
};
