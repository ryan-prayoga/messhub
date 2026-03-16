import { error as svelteError } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { ApiError, authServerApi, notificationsServerApi } from '$lib/api/server';
import type { NotificationList } from '$lib/api/types';
import { clearAuthCookies } from '$lib/auth/session';
import { requireAuth } from '$lib/auth/guard';

export const load: LayoutServerLoad = async ({ cookies, fetch, locals, url }) => {
  let user = null;
  let notificationSummary: NotificationList = {
    items: [],
    unread_count: 0
  };

  if (locals.token) {
    try {
      const response = await authServerApi.me(fetch, locals.token);
      user = response.data;
    } catch (error) {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        clearAuthCookies(cookies);
        locals.token = null;
      } else {
        throw svelteError(503, {
          message: error instanceof Error ? error.message : 'Failed to verify session'
        });
      }
    }
  }

  locals.user = user;
  requireAuth(url.pathname, user);

  if (locals.token && user) {
    try {
      const response = await notificationsServerApi.list(fetch, locals.token, { limit: 8 });
      notificationSummary = response.data;
    } catch (error) {
      if (!(error instanceof ApiError && (error.status === 401 || error.status === 403))) {
        console.error('notification summary failed', error);
      }
    }
  }

  return {
    user,
    notificationSummary
  };
};
