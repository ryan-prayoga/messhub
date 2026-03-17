<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { Proposal, ProposalStatus } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  let pendingAction: string | null = null;

  const statusLabels: Record<ProposalStatus, string> = {
    active: 'Aktif',
    closed: 'Ditutup',
    approved: 'Disetujui',
    rejected: 'Ditolak'
  };

  function enhanceWithAction(actionName: string): SubmitFunction {
    return () => {
      pendingAction = actionName;

      return async ({ update }) => {
        await update();
        pendingAction = null;
      };
    };
  }

  function statusBadgeClass(status: ProposalStatus) {
    if (status === 'approved') {
      return 'badge-success';
    }

    if (status === 'rejected') {
      return 'badge-danger';
    }

    if (status === 'closed') {
      return 'badge-strong';
    }

    return 'badge-brand';
  }

  function formatDateTime(value: string | null) {
    if (!value) {
      return '-';
    }

    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(new Date(value));
  }

  function createValue(field: 'title' | 'description' | 'voting_start' | 'voting_end') {
    if (form?.action !== 'createProposal') {
      return '';
    }

    const values = form.values as Record<string, string> | undefined;
    return values?.[field] ?? '';
  }

  function finalizeValue(proposal: Proposal, field: 'status' | 'final_decision_note') {
    if (form?.action === 'finalizeProposal') {
      const values = form.values as Record<string, string> | undefined;
      if (values?.proposal_id === proposal.id) {
        return values?.[field] ?? '';
      }
    }

    if (field === 'status') {
      return proposal.status === 'rejected' ? 'rejected' : 'approved';
    }

    return proposal.final_decision_note ?? '';
  }
</script>

<div class="space-y-4">
  <PageCard
    eyebrow="Usulan"
    icon="lucide:file-text"
    title="Usulan & Voting"
    description="Buka usulan bersama, lihat hasil vote, dan biarkan keputusan bersama tercatat lebih transparan."
  >
    {#if $navigating?.to?.url.pathname === '/proposals' || pendingAction}
      <PageSkeleton statCards={2} rows={3} />
    {/if}

    {#if form?.message}
      <StatePanel
        tone="error"
        title="Belum bisa memproses"
        message={form.message}
        requestId={form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null}
      />
    {:else if form?.success}
      <FeedbackBanner tone="success" title="Berhasil" message={form.success} />
    {/if}

    {#if data.loadError}
      <StatePanel tone="error" title="Belum bisa memuat usulan" message={data.loadError} />
    {:else}
      <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <div class="stat-card bg-slate-950 text-white">
          <p class="helper-label text-slate-300">Usulan aktif</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em]">
            {data.proposals.filter((proposal) => proposal.status === 'active').length}
          </p>
          <p class="mt-2 text-sm text-slate-300">Masih membuka vote dari penghuni.</p>
        </div>
        <div class="stat-card bg-white">
          <p class="helper-label">Disetujui</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {data.proposals.filter((proposal) => proposal.status === 'approved').length}
          </p>
        </div>
        <div class="stat-card bg-white">
          <p class="helper-label">Ditolak</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {data.proposals.filter((proposal) => proposal.status === 'rejected').length}
          </p>
        </div>
        <div class="stat-card bg-white">
          <p class="helper-label">Voting selesai</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {data.proposals.filter((proposal) => proposal.status === 'closed').length}
          </p>
        </div>
      </div>

      <article class="app-panel mt-4 p-5">
        <p class="eyebrow">Usulan Baru</p>
        <h2 class="section-title mt-1">Buat usulan baru</h2>
        <form method="POST" action="?/createProposal" class="mt-4 space-y-4" use:enhance={enhanceWithAction('createProposal')}>
          <label>
            <span class="field-label">Judul</span>
            <input class="input-field" type="text" name="title" value={createValue('title')} placeholder="Contoh: Tambah dispenser baru" required />
          </label>
          <label>
            <span class="field-label">Deskripsi</span>
            <textarea class="input-field min-h-[120px]" name="description" required>{createValue('description')}</textarea>
          </label>
          <div class="grid gap-4 sm:grid-cols-2">
            <label>
              <span class="field-label">Mulai voting (opsional)</span>
              <input class="input-field" type="datetime-local" name="voting_start" value={createValue('voting_start')} />
            </label>
            <label>
              <span class="field-label">Batas voting (opsional)</span>
              <input class="input-field" type="datetime-local" name="voting_end" value={createValue('voting_end')} />
            </label>
          </div>
          <button type="submit" class="btn-primary w-full px-4 py-3" disabled={pendingAction === 'createProposal'}>
            {pendingAction === 'createProposal' ? 'Menyimpan...' : 'Buka usulan'}
          </button>
        </form>
      </article>

      {#if data.proposals.length === 0}
        <StatePanel
          tone="empty"
          title="Belum ada usulan"
          message="Usulan bersama yang membutuhkan voting akan muncul di sini agar proses keputusan lebih terbuka."
        />
      {:else}
        <div class="mt-4 space-y-4">
          {#each data.proposals as proposal}
            <article class="app-panel p-5">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span class={statusBadgeClass(proposal.status)}>{statusLabels[proposal.status]}</span>
                    {#if proposal.current_user_vote}
                      <span class="badge-muted">Vote saya: {proposal.current_user_vote === 'agree' ? 'Setuju' : 'Tidak setuju'}</span>
                    {/if}
                  </div>
                  <h2 class="section-title mt-3">{proposal.title}</h2>
                  <p class="mt-2 text-sm leading-6 text-slate-600">{proposal.description}</p>
                  <p class="mt-3 text-xs text-dusty">Dibuat oleh {proposal.created_by_name} • {formatDateTime(proposal.created_at)}</p>
                </div>

                <div class="grid min-w-[13rem] gap-2 text-right">
                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Hasil vote</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {proposal.agree_count} setuju • {proposal.disagree_count} tidak
                    </p>
                    <p class="mt-1 text-xs text-dusty">Total {proposal.total_votes} vote</p>
                  </div>
                </div>
              </div>

              <div class="mt-4 grid gap-3 sm:grid-cols-2">
                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Mulai voting</p>
                  <p class="mt-2 text-sm font-medium text-ink">{formatDateTime(proposal.voting_start)}</p>
                </div>
                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Batas voting</p>
                  <p class="mt-2 text-sm font-medium text-ink">{formatDateTime(proposal.voting_end)}</p>
                </div>
              </div>

              {#if proposal.final_decision_note}
                <div class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Catatan keputusan akhir</p>
                  <p class="mt-2 text-sm text-slate-700">{proposal.final_decision_note}</p>
                </div>
              {/if}

              {#if proposal.status === 'active' && !proposal.current_user_vote}
                <div class="mt-4 flex flex-wrap gap-2">
                  <form method="POST" action="?/voteProposal" use:enhance={enhanceWithAction(`vote-agree-${proposal.id}`)}>
                    <input type="hidden" name="proposal_id" value={proposal.id} />
                    <input type="hidden" name="vote_type" value="agree" />
                    <button type="submit" class="btn-primary px-4 py-2.5" disabled={pendingAction === `vote-agree-${proposal.id}`}>
                      {pendingAction === `vote-agree-${proposal.id}` ? 'Menyimpan...' : 'Setuju'}
                    </button>
                  </form>
                  <form method="POST" action="?/voteProposal" use:enhance={enhanceWithAction(`vote-disagree-${proposal.id}`)}>
                    <input type="hidden" name="proposal_id" value={proposal.id} />
                    <input type="hidden" name="vote_type" value="disagree" />
                    <button type="submit" class="btn-secondary px-4 py-2.5" disabled={pendingAction === `vote-disagree-${proposal.id}`}>
                      {pendingAction === `vote-disagree-${proposal.id}` ? 'Menyimpan...' : 'Tidak setuju'}
                    </button>
                  </form>
                </div>
              {/if}

              {#if data.canFinalize}
                <details class="mt-4 rounded-2xl border border-line bg-white/70 p-4">
                  <summary class="cursor-pointer text-sm font-semibold text-ink">Aksi admin</summary>
                  <div class="mt-4 space-y-4">
                    {#if proposal.status === 'active'}
                      <form method="POST" action="?/closeProposal" use:enhance={enhanceWithAction(`close-${proposal.id}`)}>
                        <input type="hidden" name="proposal_id" value={proposal.id} />
                        <button type="submit" class="btn-secondary w-full px-4 py-3" disabled={pendingAction === `close-${proposal.id}`}>
                          {pendingAction === `close-${proposal.id}` ? 'Menutup...' : 'Tutup voting'}
                        </button>
                      </form>
                    {/if}

                    {#if proposal.status === 'active' || proposal.status === 'closed'}
                      <form method="POST" action="?/finalizeProposal" class="space-y-4" use:enhance={enhanceWithAction(`finalize-${proposal.id}`)}>
                        <input type="hidden" name="proposal_id" value={proposal.id} />
                        <label>
                          <span class="field-label">Keputusan akhir</span>
                          <select class="input-field" name="status">
                            <option value="approved" selected={finalizeValue(proposal, 'status') === 'approved'}>Setujui</option>
                            <option value="rejected" selected={finalizeValue(proposal, 'status') === 'rejected'}>Tolak</option>
                          </select>
                        </label>
                        <label>
                          <span class="field-label">Catatan admin</span>
                          <textarea class="input-field min-h-[96px]" name="final_decision_note">{finalizeValue(proposal, 'final_decision_note')}</textarea>
                        </label>
                        <button type="submit" class="btn-primary w-full px-4 py-3" disabled={pendingAction === `finalize-${proposal.id}`}>
                          {pendingAction === `finalize-${proposal.id}` ? 'Menyimpan...' : 'Finalisasi usulan'}
                        </button>
                      </form>
                    {/if}
                  </div>
                </details>
              {/if}
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
