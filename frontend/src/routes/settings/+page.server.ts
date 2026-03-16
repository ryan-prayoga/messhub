import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { settingsServerApi, systemServerApi } from '$lib/api/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

function isAdmin(role: string | undefined) {
  return role === 'admin';
}

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
  await parent();

  if (!locals.token) {
    return {
      accessDenied: false,
      settings: null,
      systemStatus: null,
      loadError: 'Missing auth token'
    };
  }

  if (!isAdmin(locals.user?.role)) {
    return {
      accessDenied: true,
      settings: null,
      systemStatus: null,
      loadError: null
    };
  }

  const [settingsResult, systemResult] = await Promise.allSettled([
    settingsServerApi.get(fetch, locals.token),
    systemServerApi.status(fetch, locals.token)
  ]);

  if (settingsResult.status === 'rejected') {
    throwIfUnauthorized(settingsResult.reason, cookies);
  }

  if (systemResult.status === 'rejected') {
    throwIfUnauthorized(systemResult.reason, cookies);
  }

  const settingsFailure =
    settingsResult.status === 'rejected'
      ? toApiFailureState(settingsResult.reason, 'Failed to load settings')
      : null;

  return {
    accessDenied: false,
    settings: settingsResult.status === 'fulfilled' ? settingsResult.value.data : null,
    systemStatus: systemResult.status === 'fulfilled' ? systemResult.value.data : null,
    loadError: settingsFailure?.message ?? null
  };
};

export const actions: Actions = {
  updateSettings: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      mess_name: normalizeString(formData.get('mess_name')),
      wifi_price: normalizeString(formData.get('wifi_price')),
      wifi_deadline_day: normalizeString(formData.get('wifi_deadline_day')),
      bank_account_name: normalizeString(formData.get('bank_account_name')),
      bank_account_number: normalizeString(formData.get('bank_account_number'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'updateSettings',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!isAdmin(locals.user.role)) {
      return fail(403, {
        action: 'updateSettings',
        message: 'Only admin can update settings',
        values
      });
    }

    const wifiPrice = Number(values.wifi_price);
    const wifiDeadlineDay = Number(values.wifi_deadline_day);

    if (
      values.mess_name === '' ||
      values.bank_account_name === '' ||
      values.bank_account_number === '' ||
      !Number.isInteger(wifiPrice) ||
      wifiPrice <= 0 ||
      !Number.isInteger(wifiDeadlineDay) ||
      wifiDeadlineDay < 1 ||
      wifiDeadlineDay > 31
    ) {
      return fail(400, {
        action: 'updateSettings',
        message: 'All settings fields are required and must be valid',
        values
      });
    }

    try {
      await settingsServerApi.update(fetch, locals.token, {
        mess_name: values.mess_name,
        wifi_price: wifiPrice,
        wifi_deadline_day: wifiDeadlineDay,
        bank_account_name: values.bank_account_name,
        bank_account_number: values.bank_account_number
      });

      return {
        action: 'updateSettings',
        success: 'Settings updated.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to update settings');

      return fail(failure.status, {
        action: 'updateSettings',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
