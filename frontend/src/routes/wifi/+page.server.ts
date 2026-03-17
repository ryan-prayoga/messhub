import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { ApiError, settingsServerApi, wifiServerApi } from '$lib/api/server';
import type { WifiBillStatus } from '$lib/api/types';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

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

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
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
      throwIfUnauthorized(error, cookies);
      console.error('wifi defaults failed', toApiFailureState(error, 'Default wifi belum dapat dimuat.'));
    }
  }

  if (!locals.token || !locals.user) {
    return {
      activeBill: null,
      bills: [],
      myBills: [],
      canManage: false,
      loadError: 'Sesi login tidak ditemukan.',
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
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Data wifi belum dapat dimuat.');

    return {
      activeBill: null,
      bills: [],
      myBills: [],
      canManage: canManage(locals.user?.role),
      loadError: failure.message,
      defaults
    };
  }
};

export const actions: Actions = {
  createBill: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      month: normalizeString(formData.get('month')),
      year: normalizeString(formData.get('year')),
      nominal_per_person: normalizeString(formData.get('nominal_per_person')),
      deadline_date: normalizeString(formData.get('deadline_date')),
      status: normalizeString(formData.get('status')) || 'active'
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(user.role)) {
      return fail(403, {
        action: 'createBill',
        message: 'Hanya admin dan bendahara yang bisa membuat tagihan wifi.',
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
        message: 'Bulan, tahun, nominal, dan jatuh tempo wajib diisi.',
        values
      });
    }

    try {
      await wifiServerApi.createBill(fetch, token, {
        month,
        year,
        nominal_per_person: nominalPerPerson,
        deadline_date: values.deadline_date,
        status: values.status as WifiBillStatus
      });

      return {
        action: 'createBill',
        success: 'Tagihan wifi berhasil dibuat dan anggota aktif sudah ditambahkan otomatis.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Tagihan wifi belum dapat dibuat.');

      return fail(failure.status, {
        action: 'createBill',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  submitProof: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      proof_url: normalizeString(formData.get('proof_url')),
      note: normalizeString(formData.get('note'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.bill_id === '' || values.proof_url === '') {
      return fail(400, {
        action: 'submitProof',
        message: 'Referensi bukti transfer wajib diisi.',
        values
      });
    }

    try {
      await wifiServerApi.submitProof(fetch, token, values.bill_id, {
        proof_url: values.proof_url,
        ...(values.note ? { note: values.note } : {})
      });

      return {
        action: 'submitProof',
        success: 'Bukti pembayaran berhasil dikirim untuk diverifikasi.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Bukti pembayaran belum dapat dikirim.');

      return fail(failure.status, {
        action: 'submitProof',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  verify: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      member_id: normalizeString(formData.get('member_id'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(user.role)) {
      return fail(403, {
        action: 'verify',
        message: 'Hanya admin dan bendahara yang bisa memverifikasi pembayaran wifi.',
        values
      });
    }

    if (values.bill_id === '' || values.member_id === '') {
      return fail(400, {
        action: 'verify',
        message: 'Referensi tagihan dan anggota wajib diisi.',
        values
      });
    }

    try {
      await wifiServerApi.verifyPayment(fetch, token, values.bill_id, values.member_id);

      return {
        action: 'verify',
        success: 'Pembayaran wifi berhasil diverifikasi.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Pembayaran wifi belum dapat diverifikasi.');

      return fail(failure.status, {
        action: 'verify',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  reject: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      member_id: normalizeString(formData.get('member_id')),
      reason: normalizeString(formData.get('reason'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(user.role)) {
      return fail(403, {
        action: 'reject',
        message: 'Hanya admin dan bendahara yang bisa menolak pembayaran wifi.',
        values
      });
    }

    if (values.bill_id === '' || values.member_id === '' || values.reason === '') {
      return fail(400, {
        action: 'reject',
        message: 'Referensi tagihan, anggota, dan alasan penolakan wajib diisi.',
        values
      });
    }

    try {
      await wifiServerApi.rejectPayment(fetch, token, values.bill_id, values.member_id, {
        reason: values.reason
      });

      return {
        action: 'reject',
        success: 'Pembayaran wifi ditolak dan anggota bisa mengirim ulang bukti baru.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Pembayaran wifi belum dapat ditolak.');

      return fail(failure.status, {
        action: 'reject',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  updateBillStatus: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      bill_id: normalizeString(formData.get('bill_id')),
      status: normalizeString(formData.get('status'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });

    if (!canManage(user.role)) {
      return fail(403, {
        action: 'updateBillStatus',
        message: 'Hanya admin dan bendahara yang bisa mengubah status tagihan wifi.',
        values
      });
    }

    if (
      values.bill_id === '' ||
      (values.status !== 'draft' && values.status !== 'active' && values.status !== 'closed')
    ) {
      return fail(400, {
        action: 'updateBillStatus',
        message: 'Referensi tagihan dan status baru wajib diisi dengan benar.',
        values
      });
    }

    try {
      await wifiServerApi.updateBillStatus(fetch, token, values.bill_id, {
        status: values.status as WifiBillStatus
      });

      return {
        action: 'updateBillStatus',
        success: 'Status tagihan wifi berhasil diperbarui.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Status tagihan wifi belum dapat diperbarui.');

      return fail(failure.status, {
        action: 'updateBillStatus',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
