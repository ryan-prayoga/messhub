import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { activitiesServerApi, ApiError } from '$lib/api/server';
import type { ActivityType } from '$lib/api/types';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

const defaults = {
  type: 'food',
  title: '',
  content: '',
  points: '1'
} as const;

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
  await parent();

  if (!locals.token || !locals.user) {
    return {
      activities: [],
      expiredActivities: [],
      loadError: 'Sesi login tidak ditemukan.',
      defaults
    };
  }

  try {
    const [activeResponse, historyResponse] = await Promise.all([
      activitiesServerApi.list(fetch, locals.token, { limit: 20, status: 'active' }),
      activitiesServerApi.list(fetch, locals.token, { limit: 10, status: 'expired' })
    ]);

    return {
      activities: activeResponse.data,
      expiredActivities: historyResponse.data,
      loadError: null,
      defaults
    };
  } catch (error) {
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Aktivitas belum dapat dimuat.');

    return {
      activities: [],
      expiredActivities: [],
      loadError: failure.message,
      defaults
    };
  }
};

export const actions: Actions = {
  createActivity: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      type: normalizeString(formData.get('type')) || defaults.type,
      title: normalizeString(formData.get('title')),
      content: normalizeString(formData.get('content')),
      points: normalizeString(formData.get('points')) || defaults.points
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    const points = Number(values.points);
    const payload: {
      type: ActivityType;
      title: string;
      content: string;
      points?: number;
    } = {
      type: values.type as ActivityType,
      title: values.title,
      content: values.content
    };

    if (values.type === 'contribution') {
      if (!Number.isInteger(points) || points <= 0) {
        return fail(400, {
          action: 'createActivity',
          message: 'Aktivitas kontribusi membutuhkan poin lebih dari 0.',
          values
        });
      }

      payload.points = points;
    }

    try {
      await activitiesServerApi.create(fetch, token, payload);

      return {
        action: 'createActivity',
        success: 'Aktivitas berhasil diterbitkan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Aktivitas belum dapat diterbitkan.');

      return fail(failure.status, {
        action: 'createActivity',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  react: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      activity_id: normalizeString(formData.get('activity_id'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.activity_id === '') {
      return fail(400, {
        action: 'react',
        message: 'Referensi aktivitas wajib diisi.',
        values
      });
    }

    try {
      await activitiesServerApi.toggleReaction(fetch, token, values.activity_id, {
        reaction_type: 'like'
      });

      return {
        action: 'react',
        success: 'Reaksi berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Reaksi belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'react',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  comment: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      activity_id: normalizeString(formData.get('activity_id')),
      comment: normalizeString(formData.get('comment'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.activity_id === '' || values.comment === '') {
      return fail(400, {
        action: 'comment',
        message: 'Referensi aktivitas dan komentar wajib diisi.',
        values
      });
    }

    try {
      await activitiesServerApi.addComment(fetch, token, values.activity_id, {
        comment: values.comment
      });

      return {
        action: 'comment',
        success: 'Komentar ditambahkan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Komentar belum dapat ditambahkan.');

      return fail(failure.status, {
        action: 'comment',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  claim: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      activity_id: normalizeString(formData.get('activity_id'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.activity_id === '') {
      return fail(400, {
        action: 'claim',
        message: 'Referensi aktivitas wajib diisi.',
        values
      });
    }

    try {
      await activitiesServerApi.claimFood(fetch, token, values.activity_id);

      return {
        action: 'claim',
        success: 'Pengambilan makanan berhasil dicatat.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Pengambilan makanan belum dapat dicatat.');

      return fail(failure.status, {
        action: 'claim',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  riceResponse: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      activity_id: normalizeString(formData.get('activity_id'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.activity_id === '') {
      return fail(400, {
        action: 'riceResponse',
        message: 'Referensi aktivitas wajib diisi.',
        values
      });
    }

    try {
      await activitiesServerApi.respondRice(fetch, token, values.activity_id);

      return {
        action: 'riceResponse',
        success: 'Respons nasi tersimpan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Respons nasi belum dapat disimpan.');

      return fail(failure.status, {
        action: 'riceResponse',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
