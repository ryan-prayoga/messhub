import type { PageServerLoad } from './$types';
import { ApiError, usersServerApi } from '$lib/api/server';

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
  await parent();

  if (!locals.token) {
    return {
      members: [],
      summary: {
        total: 0,
        active: 0,
        inactive: 0
      },
      accessDenied: false,
      loadError: 'Missing auth token'
    };
  }

  try {
    const response = await usersServerApi.list(fetch, locals.token);
    const members = response.data;
    const active = members.filter((member) => member.is_active).length;

    return {
      members,
      summary: {
        total: members.length,
        active,
        inactive: members.length - active
      },
      accessDenied: false,
      loadError: null
    };
  } catch (error) {
    if (error instanceof ApiError && error.status === 403) {
      return {
        members: [],
        summary: {
          total: 0,
          active: 0,
          inactive: 0
        },
        accessDenied: true,
        loadError: null
      };
    }

    return {
      members: [],
      summary: {
        total: 0,
        active: 0,
        inactive: 0
      },
      accessDenied: false,
      loadError: error instanceof Error ? error.message : 'Failed to load members'
    };
  }
};
