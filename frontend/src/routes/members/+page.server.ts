import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, usersServerApi } from '$lib/api/server';
import type { UserRole } from '$lib/api/types';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function canManage(role: string | undefined) {
  return role === 'admin';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
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
      loadError: 'Sesi login tidak ditemukan.'
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
    throwIfUnauthorized(error, cookies);

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

    const failure = toApiFailureState(error, 'Daftar anggota belum dapat dimuat.');

    return {
      members: [],
      summary: {
        total: 0,
        active: 0,
        inactive: 0
      },
      canManage: canManage(locals.user?.role),
      accessDenied: false,
      loadError: failure.message
    };
  }
};

export const actions: Actions = {
  updateRole: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      member_id: normalizeString(formData.get('member_id')),
      role: normalizeString(formData.get('role'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'updateRole',
        message: 'Sesi login tidak ditemukan.',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'updateRole',
        message: 'Hanya admin yang bisa mengubah role anggota.',
        values
      });
    }

    if (
      values.member_id === '' ||
      (values.role !== 'admin' && values.role !== 'treasurer' && values.role !== 'member')
    ) {
      return fail(400, {
        action: 'updateRole',
        message: 'Data anggota dan role wajib dipilih.',
        values
      });
    }

    try {
      await usersServerApi.update(fetch, locals.token, values.member_id, {
        role: values.role as UserRole
      });

      return {
        action: 'updateRole',
        success: 'Role anggota berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Role anggota belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'updateRole',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  toggleActive: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      member_id: normalizeString(formData.get('member_id')),
      is_active: normalizeString(formData.get('is_active'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'toggleActive',
        message: 'Sesi login tidak ditemukan.',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'toggleActive',
        message: 'Hanya admin yang bisa mengubah status anggota.',
        values
      });
    }

    if (values.member_id === '' || (values.is_active !== 'true' && values.is_active !== 'false')) {
      return fail(400, {
        action: 'toggleActive',
        message: 'Data anggota dan status aktif wajib diisi.',
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
        success: isActive ? 'Anggota berhasil diaktifkan.' : 'Anggota berhasil dinonaktifkan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Status anggota belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'toggleActive',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
