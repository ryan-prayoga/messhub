import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { API_BASE_URL } from '$lib/config/env';
import { AUTH_COOKIE_KEYS } from '$lib/auth/session';

export const load: PageServerLoad = async ({ locals }) => {
  if (locals.user) {
    throw redirect(303, '/');
  }

  return {};
};

export const actions: Actions = {
  default: async ({ cookies, fetch, request }) => {
    const formData = await request.formData();
    const email = String(formData.get('email') || '').trim();
    const password = String(formData.get('password') || '');

    if (!email || !password) {
      return fail(400, { message: 'Email dan password wajib diisi.' });
    }

    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, password })
    });

    const payload = (await response.json().catch(() => null)) as
      | {
          data?: {
            token: string;
            user: { id: string; name: string; email: string; role: 'admin' | 'treasurer' | 'member' };
          };
          message?: string;
        }
      | null;

    if (!response.ok || !payload?.data) {
      return fail(response.status, {
        message: payload?.message || 'Login gagal. Cek backend atau kredensial.'
      });
    }

    const { token, user } = payload.data;
    const secure = false;
    const cookieOptions = {
      path: '/',
      httpOnly: true,
      sameSite: 'lax' as const,
      secure,
      maxAge: 60 * 60 * 72
    };

    cookies.set(AUTH_COOKIE_KEYS.token, token, cookieOptions);
    cookies.set(AUTH_COOKIE_KEYS.userId, user.id, cookieOptions);
    cookies.set(AUTH_COOKIE_KEYS.userEmail, user.email, cookieOptions);
    cookies.set(AUTH_COOKIE_KEYS.userName, user.name, cookieOptions);
    cookies.set(AUTH_COOKIE_KEYS.userRole, user.role, cookieOptions);

    throw redirect(303, '/');
  }
};
