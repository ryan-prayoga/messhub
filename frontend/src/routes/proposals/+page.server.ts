import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { proposalsServerApi } from '$lib/api/server';
import { requireServerUser } from '$lib/auth/server';
import { throwIfUnauthorized, toApiFailureState } from '$lib/server/api-errors';

function canFinalize(role: string | undefined) {
  return role === 'admin';
}

function normalizeString(value: FormDataEntryValue | null) {
  return typeof value === 'string' ? value.trim() : '';
}

export const load: PageServerLoad = async ({ cookies, fetch, locals, parent }) => {
  await parent();

  if (!locals.token || !locals.user) {
    return {
      proposals: [],
      canFinalize: false,
      loadError: 'Sesi login tidak ditemukan.'
    };
  }

  try {
    const response = await proposalsServerApi.list(fetch, locals.token);

    return {
      proposals: response.data,
      canFinalize: canFinalize(locals.user.role),
      loadError: null
    };
  } catch (error) {
    throwIfUnauthorized(error, cookies);
    const failure = toApiFailureState(error, 'Usulan belum dapat dimuat.');

    return {
      proposals: [],
      canFinalize: canFinalize(locals.user?.role),
      loadError: failure.message
    };
  }
};

export const actions: Actions = {
  createProposal: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      title: normalizeString(formData.get('title')),
      description: normalizeString(formData.get('description')),
      voting_start: normalizeString(formData.get('voting_start')),
      voting_end: normalizeString(formData.get('voting_end'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (values.title === '' || values.description === '') {
      return fail(400, {
        action: 'createProposal',
        message: 'Judul dan deskripsi usulan wajib diisi.',
        values
      });
    }

    try {
      await proposalsServerApi.create(fetch, token, {
        title: values.title,
        description: values.description,
        ...(values.voting_start ? { voting_start: values.voting_start } : {}),
        ...(values.voting_end ? { voting_end: values.voting_end } : {})
      });

      return {
        action: 'createProposal',
        success: 'Usulan berhasil dibuat.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Usulan belum dapat dibuat.');

      return fail(failure.status, {
        action: 'createProposal',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  voteProposal: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      proposal_id: normalizeString(formData.get('proposal_id')),
      vote_type: normalizeString(formData.get('vote_type'))
    };

    const { token } = await requireServerUser({ cookies, fetch, locals });

    if (
      values.proposal_id === '' ||
      (values.vote_type !== 'agree' && values.vote_type !== 'disagree')
    ) {
      return fail(400, {
        action: 'voteProposal',
        message: 'Referensi usulan dan tipe vote wajib diisi.',
        values
      });
    }

    try {
      await proposalsServerApi.vote(fetch, token, values.proposal_id, {
        vote_type: values.vote_type as 'agree' | 'disagree'
      });

      return {
        action: 'voteProposal',
        success: 'Vote berhasil disimpan.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Vote belum dapat disimpan.');

      return fail(failure.status, {
        action: 'voteProposal',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  closeProposal: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      proposal_id: normalizeString(formData.get('proposal_id'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });
    if (!canFinalize(user.role)) {
      return fail(403, {
        action: 'closeProposal',
        message: 'Hanya admin yang bisa menutup voting usulan.',
        values
      });
    }

    if (values.proposal_id === '') {
      return fail(400, {
        action: 'closeProposal',
        message: 'Referensi usulan tidak ditemukan.',
        values
      });
    }

    try {
      await proposalsServerApi.close(fetch, token, values.proposal_id);

      return {
        action: 'closeProposal',
        success: 'Voting usulan berhasil ditutup.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Voting usulan belum dapat ditutup.');

      return fail(failure.status, {
        action: 'closeProposal',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  },
  finalizeProposal: async ({ cookies, fetch, locals, request }) => {
    const formData = await request.formData();
    const values = {
      proposal_id: normalizeString(formData.get('proposal_id')),
      status: normalizeString(formData.get('status')),
      final_decision_note: normalizeString(formData.get('final_decision_note'))
    };

    const { token, user } = await requireServerUser({ cookies, fetch, locals });
    if (!canFinalize(user.role)) {
      return fail(403, {
        action: 'finalizeProposal',
        message: 'Hanya admin yang bisa memfinalisasi usulan.',
        values
      });
    }

    if (
      values.proposal_id === '' ||
      (values.status !== 'approved' && values.status !== 'rejected')
    ) {
      return fail(400, {
        action: 'finalizeProposal',
        message: 'Referensi usulan dan keputusan akhir wajib diisi.',
        values
      });
    }

    try {
      await proposalsServerApi.finalize(fetch, token, values.proposal_id, {
        status: values.status as 'approved' | 'rejected',
        ...(values.final_decision_note ? { final_decision_note: values.final_decision_note } : {})
      });

      return {
        action: 'finalizeProposal',
        success: 'Usulan berhasil difinalisasi.'
      };
    } catch (error) {
      throwIfUnauthorized(error, cookies);
      const failure = toApiFailureState(error, 'Usulan belum dapat difinalisasi.');

      return fail(failure.status, {
        action: 'finalizeProposal',
        message: failure.message,
        requestId: failure.requestId,
        values
      });
    }
  }
};
