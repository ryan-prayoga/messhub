import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { activitiesServerApi, ApiError } from '$lib/api/server';
import type { ActivityType } from '$lib/api/types';

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

const defaults = {
  type: 'food',
  title: '',
  content: '',
  points: '1'
} as const;

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
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
    return {
      activities: [],
      loadError: error instanceof Error ? error.message : 'Failed to load activities',
      defaults
    };
  }
};

export const actions: Actions = {
  createActivity: async ({ fetch, locals, request }) => {
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
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'createActivity',
        message: error instanceof Error ? error.message : 'Failed to create activity',
        values
      });
    }
  },
  react: async ({ fetch, locals, request }) => {
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
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'react',
        message: error instanceof Error ? error.message : 'Failed to update reaction',
        values
      });
    }
  },
  comment: async ({ fetch, locals, request }) => {
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
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'comment',
        message: error instanceof Error ? error.message : 'Failed to add comment',
        values
      });
    }
  },
  claim: async ({ fetch, locals, request }) => {
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
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'claim',
        message: error instanceof Error ? error.message : 'Failed to claim food',
        values
      });
    }
  },
  riceResponse: async ({ fetch, locals, request }) => {
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
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'riceResponse',
        message: error instanceof Error ? error.message : 'Failed to save rice response',
        values
      });
    }
  }
};
