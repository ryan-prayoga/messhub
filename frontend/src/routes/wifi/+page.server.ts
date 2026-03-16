import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, settingsServerApi, wifiServerApi } from '$lib/api/server';
import type { WifiBillStatus } from '$lib/api/types';

function canManage(role: string | undefined) {
  return role === 'admin' || role === 'treasurer';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

function buildDefaultBillValues(wifiPrice = 20000, wifiDeadlineDay = 10) {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1;
  const lastDay = new Date(year, month, 0).getDate();
  const deadlineDay = Math.min(Math.max(wifiDeadlineDay, 1), lastDay);
  const deadlineDate = `${year}-${String(month).padStart(2, '0')}-${String(deadlineDay).padStart(2, '0')}`;

  return {
    month: String(month),
    year: String(year),
    nominal_per_person: String(wifiPrice),
    deadline_date: deadlineDate,
    status: 'active'
  };
}

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
  await parent();

  let defaults = buildDefaultBillValues();

  if (locals.token && canManage(locals.user?.role)) {
    try {
      const settingsResponse = await settingsServerApi.get(fetch, locals.token);
      defaults = buildDefaultBillValues(
        settingsResponse.data?.wifi_price ?? 20000,
        settingsResponse.data?.wifi_deadline_day ?? 10
      );
    } catch (error) {
      console.error('wifi defaults failed', error);
    }
  }

  if (!locals.token || !locals.user) {
    return {
      activeBill: null,
      bills: [],
      myBills: [],
      canManage: false,
      loadError: 'Missing authenticated session',
      defaults
    };
  }

  try {
    const [activeResponse, myBillsResponse, billsResponse] = await Promise.all([
      wifiServerApi.getActive(fetch, locals.token),
      wifiServerApi.getMyBills(fetch, locals.token),
      canManage(locals.user.role) ? wifiServerApi.listBills(fetch, locals.token) : Promise.resolve(null)
    ]);

    return {
      activeBill: activeResponse.data,
      bills: billsResponse?.data ?? [],
      myBills: myBillsResponse.data,
      canManage: canManage(locals.user.role),
      loadError: null,
      defaults
    };
  } catch (error) {
    return {
      activeBill: null,
      bills: [],
      myBills: [],
      canManage: canManage(locals.user?.role),
      loadError:
        error instanceof ApiError || error instanceof Error
          ? error.message
          : 'Failed to load wifi data',
      defaults
    };
  }
};

export const actions: Actions = {
  createBill: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      month: normalizeString(formData.get('month')),
      year: normalizeString(formData.get('year')),
      nominal_per_person: normalizeString(formData.get('nominal_per_person')),
      deadline_date: normalizeString(formData.get('deadline_date')),
      status: normalizeString(formData.get('status')) || 'active'
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'createBill',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'createBill',
        message: 'Only admin and treasurer can create wifi bills',
        values
      });
    }

    const month = Number(values.month);
    const year = Number(values.year);
    const nominalPerPerson = Number(values.nominal_per_person);

    if (
      !Number.isInteger(month) ||
      !Number.isInteger(year) ||
      !Number.isInteger(nominalPerPerson) ||
      nominalPerPerson <= 0 ||
      values.deadline_date === ''
    ) {
      return fail(400, {
        action: 'createBill',
        message: 'Month, year, nominal, and deadline are required',
        values
      });
    }

    try {
      await wifiServerApi.createBill(fetch, locals.token, {
        month,
        year,
        nominal_per_person: nominalPerPerson,
        deadline_date: values.deadline_date,
        status: values.status as WifiBillStatus
      });

      return {
        action: 'createBill',
        success: 'Wifi bill created and active members were generated automatically.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'createBill',
        message: error instanceof Error ? error.message : 'Failed to create wifi bill',
        values
      });
    }
  },
  submitProof: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      proof_url: normalizeString(formData.get('proof_url')),
      note: normalizeString(formData.get('note'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'submitProof',
        message: 'Missing authenticated session',
        values
      });
    }

    if (values.bill_id === '' || values.proof_url === '') {
      return fail(400, {
        action: 'submitProof',
        message: 'Proof reference is required',
        values
      });
    }

    try {
      await wifiServerApi.submitProof(fetch, locals.token, values.bill_id, {
        proof_url: values.proof_url,
        ...(values.note ? { note: values.note } : {})
      });

      return {
        action: 'submitProof',
        success: 'Payment proof submitted for verification.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'submitProof',
        message: error instanceof Error ? error.message : 'Failed to submit wifi payment proof',
        values
      });
    }
  },
  verify: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      member_id: normalizeString(formData.get('member_id'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'verify',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'verify',
        message: 'Only admin and treasurer can verify wifi payments',
        values
      });
    }

    if (values.bill_id === '' || values.member_id === '') {
      return fail(400, {
        action: 'verify',
        message: 'Bill and member reference are required',
        values
      });
    }

    try {
      await wifiServerApi.verifyPayment(fetch, locals.token, values.bill_id, values.member_id);

      return {
        action: 'verify',
        success: 'Wifi payment verified.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'verify',
        message: error instanceof Error ? error.message : 'Failed to verify wifi payment',
        values
      });
    }
  },
  reject: async ({ fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      member_id: normalizeString(formData.get('member_id')),
      reason: normalizeString(formData.get('reason'))
    };

    if (!locals.token || !locals.user) {
      return fail(401, {
        action: 'reject',
        message: 'Missing authenticated session',
        values
      });
    }

    if (!canManage(locals.user.role)) {
      return fail(403, {
        action: 'reject',
        message: 'Only admin and treasurer can reject wifi payments',
        values
      });
    }

    if (values.bill_id === '' || values.member_id === '' || values.reason === '') {
      return fail(400, {
        action: 'reject',
        message: 'Bill, member reference, and rejection reason are required',
        values
      });
    }

    try {
      await wifiServerApi.rejectPayment(fetch, locals.token, values.bill_id, values.member_id, {
        reason: values.reason
      });

      return {
        action: 'reject',
        success: 'Wifi payment rejected and can be resubmitted by the member.'
      };
    } catch (error) {
      return fail(error instanceof ApiError ? error.status : 500, {
        action: 'reject',
        message: error instanceof Error ? error.message : 'Failed to reject wifi payment',
        values
      });
    }
  }
};
