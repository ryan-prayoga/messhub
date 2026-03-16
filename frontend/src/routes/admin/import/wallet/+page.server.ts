import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { importsServerApi } from '$lib/api/server';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function isAdmin(role: string | undefined) {
  return role === 'admin';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ parent, locals }) => {
  await parent();

  return {
    accessDenied: !isAdmin(locals.user?.role)
  };
};

export const actions: Actions = {
  preview: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const file = formData.get('file');

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!isAdmin(user.role)) {
      return fail(403, {
        action: 'preview',
        message: 'Halaman impor hanya tersedia untuk admin mess.'
      });
    }

    if (!(file instanceof File) || file.size === 0) {
      return fail(400, {
        action: 'preview',
        message: 'Pilih file CSV transaksi terlebih dahulu.'
      });
    }

    try {
      const response = await importsServerApi.previewWallet(fetch, token, formData);

      return {
        action: 'preview',
        preview: response.data
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Preview impor transaksi belum dapat dibuat.');

      return fail(failure.status, {
        action: 'preview',
        message: failure.message,
        requestId: failure.requestId
      });
    }
  },
  commit: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      job_id: normalizeString(formData.get('job_id'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!isAdmin(user.role)) {
      return fail(403, {
        action: 'commit',
        message: 'Halaman impor hanya tersedia untuk admin mess.',
        values
      });
    }

    if (values.job_id === '') {
      return fail(400, {
        action: 'commit',
        message: 'Preview impor yang dipilih tidak valid.',
        values
      });
    }

    try {
      const response = await importsServerApi.commitWallet(fetch, token, {
        job_id: values.job_id
      });

      return {
        action: 'commit',
        committed: response.data,
        success: 'Impor transaksi kas berhasil dijalankan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Impor transaksi kas belum dapat disimpan.');

      return fail(failure.status, {
        action: 'commit',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
