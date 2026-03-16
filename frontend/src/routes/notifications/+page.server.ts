import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, notificationsServerApi } from '$lib/api/server';

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
  await parent();

  if (!locals.token) {
    return {
      notificationSummary: {
        items: [],
        unread_count: 0
      },
      loadError: 'Missing authenticated session'
    };
  }

  try {
    const response = await notificationsServerApi.list(fetch, locals.token, { limit: 40 });

    return {
      notificationSummary: response.data,
      loadError: null
    };
  } catch (error) {
    return {
      notificationSummary: {
        items: [],
        unread_count: 0
      },
      loadError: error instanceof Error ? error.message : 'Failed to load notifications'
    };
  }
};

export const actions: Actions = {
  markOneRead: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const id = normalizeString(formData.get('notification_id'));

    if (!locals.token) {
      return fail(401, {
        action: 'markOneRead',
        message: 'Missing authenticated session'
      });
    }

    if (id === '') {
      return fail(400, {
        action: 'markOneRead',
        message: 'Notification reference is required'
      });
    }

    try {
      await notificationsServerApi.markRead(fetch, locals.token, {
        ids: [id]
      });

      return {
        action: 'markOneRead',
        success: 'Notification ditandai sudah dibaca.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'markOneRead',
        message: error instanceof Error ? error.message : 'Failed to update notification'
      });
    }
  },
  markAllRead: async ({ fetch, locals }) => {
    if (!locals.token) {
      return fail(401, {
        action: 'markAllRead',
        message: 'Missing authenticated session'
      });
    }

    try {
      await notificationsServerApi.markRead(fetch, locals.token, {
        all: true
      });

      return {
        action: 'markAllRead',
        success: 'Semua notification ditandai sudah dibaca.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'markAllRead',
        message: error instanceof Error ? error.message : 'Failed to update notifications'
      });
    }
  }
};
