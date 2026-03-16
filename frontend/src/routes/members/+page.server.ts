import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, usersServerApi } from '$lib/api/server';
import type { UserRole } from '$lib/api/types';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function canManage(role: string | undefined) {
  return role === 'admin';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

function normalizeOptionalString(value: FormDataEntryValue | null) {
  const normalized = normalizeString(value);
  return normalized === '' ? '' : normalized;
}

function isValidRole(value: string): value is UserRole {
  return value === 'admin' || value === 'treasurer' || value === 'member';
}

function buildMemberValues(formData: FormData, includePassword = false) {
  return {
    member_id: normalizeString(formData.get('member_id')),
    name: normalizeString(formData.get('name')),
    email: normalizeString(formData.get('email')),
    username: normalizeOptionalString(formData.get('username')),
    phone: normalizeOptionalString(formData.get('phone')),
    joined_at: normalizeOptionalString(formData.get('joined_at')),
    role: normalizeString(formData.get('role')),
    is_active: normalizeString(formData.get('is_active')) || 'true',
    ...(includePassword
      ? {
          password: normalizeString(formData.get('password')),
          confirm_password: normalizeString(formData.get('confirm_password'))
        }
      : {})
  };
}

function validateMemberValues(
  values: ReturnType<typeof buildMemberValues>,
  includePassword = false
) {
  if (
    values.name === '' ||
    values.email === '' ||
    !values.email.includes('@') ||
    !isValidRole(values.role) ||
    (values.is_active !== 'true' && values.is_active !== 'false')
  ) {
    return 'Nama, email valid, role, dan status aktif wajib diisi dengan benar.';
  }

  if (includePassword) {
    const passwordValues = values as ReturnType<typeof buildMemberValues> & {
      password: string;
      confirm_password: string;
    };

    if (
      passwordValues.password === '' ||
      passwordValues.confirm_password === '' ||
      passwordValues.password.length < 8
    ) {
      return 'Password minimal 8 karakter dan wajib diisi lengkap.';
    }

    if (passwordValues.password !== passwordValues.confirm_password) {
      return 'Konfirmasi password belum sama.';
    }
  }

  return null;
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
  createMember: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = buildMemberValues(formData, true);
    const passwordValues = values as ReturnType<typeof buildMemberValues> & {
      password: string;
      confirm_password: string;
    };
    const session = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(session.user.role)) {
      return fail(403, {
        action: 'createMember',
        message: 'Hanya admin yang bisa menambah anggota baru.',
        values
      });
    }

    const validationMessage = validateMemberValues(values, true);
    if (validationMessage) {
      return fail(400, {
        action: 'createMember',
        message: validationMessage,
        values
      });
    }

    try {
      await usersServerApi.create(fetch, session.token, {
        name: values.name,
        email: values.email,
        username: values.username || undefined,
        phone: values.phone === '' ? undefined : values.phone,
        password: passwordValues.password,
        role: values.role as UserRole,
        is_active: values.is_active === 'true',
        joined_at: values.joined_at || undefined
      });

      return {
        action: 'createMember',
        success: 'Anggota baru berhasil ditambahkan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Anggota baru belum dapat disimpan.');

      return fail(failure.status, {
        action: 'createMember',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  updateMember: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = buildMemberValues(formData);
    const session = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(session.user.role)) {
      return fail(403, {
        action: 'updateMember',
        message: 'Hanya admin yang bisa mengubah data anggota.',
        values
      });
    }

    if (values.member_id === '') {
      return fail(400, {
        action: 'updateMember',
        message: 'Referensi anggota tidak ditemukan.',
        values
      });
    }

    const validationMessage = validateMemberValues(values);
    if (validationMessage) {
      return fail(400, {
        action: 'updateMember',
        message: validationMessage,
        values
      });
    }

    try {
      await usersServerApi.update(fetch, session.token, values.member_id, {
        name: values.name,
        email: values.email,
        username: values.username || undefined,
        phone: values.phone || '',
        joined_at: values.joined_at || undefined,
        role: values.role as UserRole,
        is_active: values.is_active === 'true'
      });

      return {
        action: 'updateMember',
        success: 'Data anggota berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Data anggota belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'updateMember',
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
    const session = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(session.user.role)) {
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
      await usersServerApi.update(fetch, session.token, values.member_id, {
        is_active: isActive
      });

      return {
        action: 'toggleActive',
        success: isActive ? 'Anggota berhasil diaktifkan kembali.' : 'Anggota berhasil dinonaktifkan.'
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
  },
  resetPassword: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      member_id: normalizeString(formData.get('member_id')),
      new_password: normalizeString(formData.get('new_password')),
      confirm_password: normalizeString(formData.get('confirm_password'))
    };
    const session = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(session.user.role)) {
      return fail(403, {
        action: 'resetPassword',
        message: 'Hanya admin yang bisa mereset password anggota.',
        values
      });
    }

    if (values.member_id === '' || values.new_password === '' || values.confirm_password === '') {
      return fail(400, {
        action: 'resetPassword',
        message: 'Password baru dan konfirmasi wajib diisi.',
        values
      });
    }

    if (values.new_password.length < 8) {
      return fail(400, {
        action: 'resetPassword',
        message: 'Password baru minimal 8 karakter.',
        values
      });
    }

    if (values.new_password !== values.confirm_password) {
      return fail(400, {
        action: 'resetPassword',
        message: 'Konfirmasi password belum sama.',
        values
      });
    }

    try {
      await usersServerApi.resetPassword(fetch, session.token, values.member_id, {
        new_password: values.new_password
      });

      return {
        action: 'resetPassword',
        success: 'Password anggota berhasil direset.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Password anggota belum dapat direset.');

      return fail(failure.status, {
        action: 'resetPassword',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
