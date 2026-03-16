import { error as svelteError } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { ApiError, notificationsServerApi } from '$lib/api/server';
import type { NotificationList } from '$lib/api/types';
import { clearAuthCookies } from '$lib/auth/session';
import { requireAuth } from '$lib/auth/guard';
import { resolveServerUser } from '$lib/auth/server';
import { toApiFailureState } from '$lib/server/api-errors';

export const load: LayoutServerLoad = async ({ cookies, fetch, locals, url }) => {
  let user = null;
  let notificationSummary: NotificationList = {
    items: [],
    unread_count: 0
  };

  if (locals.token) {
    try {
      user = await resolveServerUser({ cookies, fetch, locals });
    } catch (error) {
      const failure = toApiFailureState(error, 'Sesi belum dapat diverifikasi.');
      throw svelteError(503, {
        message: failure.message
      });
    }
  }

  locals.user = user;
  requireAuth(url.pathname, user);

  if (locals.token && user) {
    try {
      const response = await notificationsServerApi.list(fetch, locals.token, { limit: 8 });
      notificationSummary = response.data;
    } catch (error) {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        clearAuthCookies(cookies);
        locals.token = null;
        locals.user = null;
      } else {
        const failure = toApiFailureState(error, 'Ringkasan notifikasi belum dapat dimuat.');
        console.error('notification summary failed', failure);
      }
    }
  }

  return {
    user,
    notificationSummary
  };
};
