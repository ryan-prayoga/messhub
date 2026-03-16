import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { walletServerApi } from '$lib/api/server';
import type { WalletTransactionType } from '$lib/api/types';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function canCreate(role: string | undefined) {
  return role === 'admin' || role === 'treasurer';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ parent, locals }) => {
  await parent();

  return {
    accessDenied: !canCreate(locals.user?.role)
  };
};

export const actions: Actions = {
  default: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      type: normalizeString(formData.get('type')) || 'income',
      category: normalizeString(formData.get('category')),
      amount: normalizeString(formData.get('amount')),
      description: normalizeString(formData.get('description'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canCreate(locals.user.role)) {
      return fail(403, {
        message: 'Only admin and treasurer can create wallet transactions',
        values
      });
    }

    const amount = Number(values.amount);
    if (
      (values.type !== 'income' && values.type !== 'expense') ||
      values.category === '' ||
      values.description === '' ||
      !Number.isInteger(amount) ||
      amount <= 0
    ) {
      return fail(400, {
        message: 'Type, category, amount, and description are required',
        values
      });
    }

    try {
      await walletServerApi.createTransaction(fetch, locals.token, {
        type: values.type as WalletTransactionType,
        category: values.category,
        amount,
        description: values.description
      });
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Failed to create wallet transaction');

      return fail(failure.status, {
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }

    throw redirect(303, '/wallet');
  }
};
