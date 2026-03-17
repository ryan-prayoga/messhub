import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { sharedExpensesServerApi, usersServerApi } from '$lib/api/server';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function canManage(role: string | undefined) {
  return role === 'admin' || role === 'treasurer';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
  await parent();

  if (!locals.token || !locals.user) {
    return {
      expenses: [],
      summary: null,
      payers: [],
      canManage: false,
      loadError: 'Sesi login tidak ditemukan.'
    };
  }

  try {
    const [expenseResult, payerResult] = await Promise.all([
      sharedExpensesServerApi.list(fetch, locals.token),
      canManage(locals.user.role) ? usersServerApi.list(fetch, locals.token) : Promise.resolve(null)
    ]);

    return {
      expenses: expenseResult.data.items,
      summary: expenseResult.data.summary,
      payers:
        payerResult?.data.filter((member) => !member.archived_at).map((member) => ({
          id: member.id,
          name: member.name,
          status: member.archived_at ? 'archived' : member.is_active ? 'active' : 'inactive'
        })) ?? [],
      canManage: canManage(locals.user.role),
      loadError: null
    };
  } catch (error) {
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Pengeluaran bersama belum dapat dimuat.');

    return {
      expenses: [],
      summary: null,
      payers: [],
      canManage: canManage(locals.user?.role),
      loadError: failure.message
    };
  }
};

export const actions: Actions = {
  createExpense: async ({ cookies, fetch, locals, request }) => {
    const values = {
      expense_date: '',
      category: '',
      description: '',
      amount: '',
      paid_by_user_id: '',
      status: 'fronted',
      notes: '',
      proof_url: ''
    };
    const formData = await request.formData();
    Object.assign(values, {
      expense_date: normalizeString(formData.get('expense_date')),
      category: normalizeString(formData.get('category')),
      description: normalizeString(formData.get('description')),
      amount: normalizeString(formData.get('amount')),
      paid_by_user_id: normalizeString(formData.get('paid_by_user_id')),
      status: normalizeString(formData.get('status')) || 'fronted',
      notes: normalizeString(formData.get('notes')),
      proof_url: normalizeString(formData.get('proof_url'))
    });

    const { token, user } = await requireServerUser({ cookies, fetch, locals });
    if (!canManage(user.role)) {
      return fail(403, {
        action: 'createExpense',
        message: 'Hanya admin dan bendahara yang bisa mencatat pengeluaran bersama.',
        values
      });
    }

    const amount = Number(values.amount);
    if (
      values.expense_date === '' ||
      values.category === '' ||
      values.description === '' ||
      values.paid_by_user_id === '' ||
      !Number.isInteger(amount) ||
      amount <= 0
    ) {
      return fail(400, {
        action: 'createExpense',
        message: 'Tanggal, kategori, deskripsi, nominal, dan pihak yang membayar wajib diisi dengan benar.',
        values
      });
    }

    try {
      await sharedExpensesServerApi.create(fetch, token, {
        expense_date: values.expense_date,
        category: values.category,
        description: values.description,
        amount,
        paid_by_user_id: values.paid_by_user_id,
        status: values.status as 'personal' | 'fronted' | 'partially_reimbursed' | 'reimbursed',
        ...(values.notes ? { notes: values.notes } : {}),
        ...(values.proof_url ? { proof_url: values.proof_url } : {})
      });

      return {
        action: 'createExpense',
        success: 'Pengeluaran bersama berhasil dicatat.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Pengeluaran bersama belum dapat disimpan.');

      return fail(failure.status, {
        action: 'createExpense',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  updateExpense: async ({ cookies, fetch, locals, request }) => {
    const values = {
      expense_id: '',
      expense_date: '',
      category: '',
      description: '',
      amount: '',
      paid_by_user_id: '',
      status: '',
      notes: '',
      proof_url: ''
    };
    const formData = await request.formData();
    Object.assign(values, {
      expense_id: normalizeString(formData.get('expense_id')),
      expense_date: normalizeString(formData.get('expense_date')),
      category: normalizeString(formData.get('category')),
      description: normalizeString(formData.get('description')),
      amount: normalizeString(formData.get('amount')),
      paid_by_user_id: normalizeString(formData.get('paid_by_user_id')),
      status: normalizeString(formData.get('status')),
      notes: normalizeString(formData.get('notes')),
      proof_url: normalizeString(formData.get('proof_url'))
    });

    const { token, user } = await requireServerUser({ cookies, fetch, locals });
    if (!canManage(user.role)) {
      return fail(403, {
        action: 'updateExpense',
        message: 'Hanya admin dan bendahara yang bisa mengubah pengeluaran bersama.',
        values
      });
    }

    const amount = Number(values.amount);
    if (
      values.expense_id === '' ||
      values.expense_date === '' ||
      values.category === '' ||
      values.description === '' ||
      values.paid_by_user_id === '' ||
      values.status === '' ||
      !Number.isInteger(amount) ||
      amount <= 0
    ) {
      return fail(400, {
        action: 'updateExpense',
        message: 'Seluruh field utama pengeluaran bersama wajib diisi dengan benar.',
        values
      });
    }

    try {
      await sharedExpensesServerApi.update(fetch, token, values.expense_id, {
        expense_date: values.expense_date,
        category: values.category,
        description: values.description,
        amount,
        paid_by_user_id: values.paid_by_user_id,
        status: values.status as 'personal' | 'fronted' | 'partially_reimbursed' | 'reimbursed',
        notes: values.notes,
        proof_url: values.proof_url
      });

      return {
        action: 'updateExpense',
        success: 'Pengeluaran bersama berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Pengeluaran bersama belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'updateExpense',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
