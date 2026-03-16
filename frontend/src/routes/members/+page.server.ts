import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, usersServerApi } from '$lib/api/server';
import type { UserRole } from '$lib/api/types';

function canManage(role: string | undefined) {
  return role === 'admin';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

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
      canManage: false,
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
      canManage: canManage(locals.user?.role),
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
        canManage: false,
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
      canManage: canManage(locals.user?.role),
      accessDenied: false,
      loadError: error instanceof Error ? error.message : 'Failed to load members'
    };
  }
};

export const actions: Actions = {
  updateRole: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      member_id: normalizeString(formData.get('member_id')),
      role: normalizeString(formData.get('role'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'updateRole',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'updateRole',
        message: 'Only admin can update member roles',
        values
      });
    }

    if (
      values.member_id === '' ||
      (values.role !== 'admin' && values.role !== 'treasurer' && values.role !== 'member')
    ) {
      return fail(400, {
        action: 'updateRole',
        message: 'Member reference and role are required',
        values
      });
    }

    try {
      await usersServerApi.update(fetch, locals.token, values.member_id, {
        role: values.role as UserRole
      });

      return {
        action: 'updateRole',
        success: 'Member role updated.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'updateRole',
        message: error instanceof Error ? error.message : 'Failed to update member role',
        values
      });
    }
  },
  toggleActive: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      member_id: normalizeString(formData.get('member_id')),
      is_active: normalizeString(formData.get('is_active'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'toggleActive',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'toggleActive',
        message: 'Only admin can update member activation',
        values
      });
    }

    if (values.member_id === '' || (values.is_active !== 'true' && values.is_active !== 'false')) {
      return fail(400, {
        action: 'toggleActive',
        message: 'Member reference and activation state are required',
        values
      });
    }

    try {
      const isActive = values.is_active === 'true';
      await usersServerApi.update(fetch, locals.token, values.member_id, {
        is_active: isActive
      });

      return {
        action: 'toggleActive',
        success: isActive ? 'Member activated.' : 'Member deactivated.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'toggleActive',
        message: error instanceof Error ? error.message : 'Failed to update member activation',
        values
      });
    }
  }
};
