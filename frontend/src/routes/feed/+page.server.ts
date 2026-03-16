import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { activitiesServerApi, ApiError } from '$lib/api/server';
import type { ActivityType } from '$lib/api/types';
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
      loadError: 'Missing authenticated session',
      defaults
    };
  }

  try {
    const response = await activitiesServerApi.list(fetch, locals.token, { limit: 20 });

    return {
      activities: response.data,
      loadError: null,
      defaults
    };
  } catch (error) {
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Failed to load activities');

    return {
      activities: [],
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

    if (!locals.token) {
      return fail(401, {
        action: 'createActivity',
        message: 'Missing authenticated session',
        values
      });
    }

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
          message: 'Contribution membutuhkan points lebih dari 0',
          values
        });
      }

      payload.points = points;
    }

    try {
      await activitiesServerApi.create(fetch, locals.token, payload);

      return {
        action: 'createActivity',
        success: 'Activity berhasil diposting.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to create activity');

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

    if (!locals.token) {
      return fail(401, {
        action: 'react',
        message: 'Missing authenticated session',
        values
      });
    }

    if (values.activity_id === '') {
      return fail(400, {
        action: 'react',
        message: 'Activity reference is required',
        values
      });
    }

    try {
      await activitiesServerApi.toggleReaction(fetch, locals.token, values.activity_id, {
        reaction_type: 'like'
      });

      return {
        action: 'react',
        success: 'Reaction diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to update reaction');

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

    if (!locals.token) {
      return fail(401, {
        action: 'comment',
        message: 'Missing authenticated session',
        values
      });
    }

    if (values.activity_id === '' || values.comment === '') {
      return fail(400, {
        action: 'comment',
        message: 'Activity dan comment wajib diisi',
        values
      });
    }

    try {
      await activitiesServerApi.addComment(fetch, locals.token, values.activity_id, {
        comment: values.comment
      });

      return {
        action: 'comment',
        success: 'Komentar ditambahkan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to add comment');

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

    if (!locals.token) {
      return fail(401, {
        action: 'claim',
        message: 'Missing authenticated session',
        values
      });
    }

    if (values.activity_id === '') {
      return fail(400, {
        action: 'claim',
        message: 'Activity reference is required',
        values
      });
    }

    try {
      await activitiesServerApi.claimFood(fetch, locals.token, values.activity_id);

      return {
        action: 'claim',
        success: 'Food claim tercatat.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to claim food');

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

    if (!locals.token) {
      return fail(401, {
        action: 'riceResponse',
        message: 'Missing authenticated session',
        values
      });
    }

    if (values.activity_id === '') {
      return fail(400, {
        action: 'riceResponse',
        message: 'Activity reference is required',
        values
      });
    }

    try {
      await activitiesServerApi.respondRice(fetch, locals.token, values.activity_id);

      return {
        action: 'riceResponse',
        success: 'Respons nasi tersimpan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to save rice response');

      return fail(failure.status, {
        action: 'riceResponse',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
