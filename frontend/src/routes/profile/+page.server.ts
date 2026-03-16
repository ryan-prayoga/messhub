import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { profileServerApi } from '$lib/api/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
  await parent();

  if (!locals.token) {
    return {
      profile: null,
      loadError: 'Sesi login tidak ditemukan.'
    };
  }

  try {
    const response = await profileServerApi.get(fetch, locals.token);

    return {
      profile: response.data,
      loadError: null
    };
  } catch (error) {
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Profil belum dapat dimuat.');

    return {
      profile: null,
      loadError: failure.message
    };
  }
};

export const actions: Actions = {
  updateProfile: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      name: normalizeString(formData.get('name')),
      phone: normalizeString(formData.get('phone')),
      avatar_url: normalizeString(formData.get('avatar_url'))
    };

    if (!locals.token) {
      return fail(401, {
        action: 'updateProfile',
        message: 'Sesi login tidak ditemukan.',
        values
      });
    }

    if (values.name === '') {
      return fail(400, {
        action: 'updateProfile',
        message: 'Nama wajib diisi.',
        values
      });
    }

    try {
      await profileServerApi.update(fetch, locals.token, values);

      return {
        action: 'updateProfile',
        success: 'Profil berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Profil belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'updateProfile',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  changePassword: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      current_password: normalizeString(formData.get('current_password')),
      new_password: normalizeString(formData.get('new_password')),
      confirm_password: normalizeString(formData.get('confirm_password'))
    };

    if (!locals.token) {
      return fail(401, {
        action: 'changePassword',
        message: 'Sesi login tidak ditemukan.'
      });
    }

    if (
      values.current_password === '' ||
      values.new_password === '' ||
      values.confirm_password === ''
    ) {
      return fail(400, {
        action: 'changePassword',
        message: 'Password lama, password baru, dan konfirmasi wajib diisi.'
      });
    }

    if (values.new_password !== values.confirm_password) {
      return fail(400, {
        action: 'changePassword',
        message: 'Konfirmasi password baru belum sama.'
      });
    }

    try {
      await profileServerApi.changePassword(fetch, locals.token, {
        current_password: values.current_password,
        new_password: values.new_password
      });

      return {
        action: 'changePassword',
        success: 'Password berhasil diganti.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Password belum dapat diganti.');

      return fail(failure.status, {
        action: 'changePassword',
        message: failure.message,
        requestId: failure.requestId
      });
    }
  }
};
