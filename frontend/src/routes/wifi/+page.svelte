<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import PageCard from '$lib/components/PageCard.svelte';
  import type { WifiBill, WifiBillMember, WifiBillStatus } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  const billStatusLabels: Record<WifiBillStatus, string> = {
    draft: 'Draft',
    active: 'Active',
    closed: 'Closed'
  };

  let pendingAction: string | null = null;

  function enhanceWithAction(actionName: string): SubmitFunction {
    return () => {
      pendingAction = actionName;

      return async ({ update }) => {
        await update();
        pendingAction = null;
      };
    };
  }

  function formatCurrency(value: number) {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(value);
  }

  function formatDate(value: string | null, withTime = false) {
    if (!value) {
      return '-';
    }

    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      ...(withTime
        ? {
            hour: '2-digit',
            minute: '2-digit'
          }
        : {})
    }).format(new Date(value));
  }

  function formatMonthYear(month: number, year: number) {
    return new Intl.DateTimeFormat('id-ID', {
      month: 'long',
      year: 'numeric'
    }).format(new Date(year, month - 1, 1));
  }

  function paymentStatusLabel(status: WifiBillMember['payment_status']) {
    if (status === 'pending_verification') {
      return 'Pending';
    }

    if (status === 'verified') {
      return 'Verified';
    }

    if (status === 'rejected') {
      return 'Rejected';
    }

    return 'Unpaid';
  }

  function paymentStatusClass(status: WifiBillMember['payment_status']) {
    if (status === 'verified') {
      return 'badge bg-emerald-100 text-emerald-700';
    }

    if (status === 'pending_verification') {
      return 'badge bg-amber-100 text-amber-700';
    }

    if (status === 'rejected') {
      return 'badge bg-rose-100 text-rose-700';
    }

    return 'badge-muted';
  }

  function billStatusClass(status: WifiBillStatus) {
    if (status === 'active') {
      return 'badge-brand';
    }

    if (status === 'closed') {
      return 'badge bg-slate-950 text-white';
    }

    return 'badge bg-amber-100 text-amber-700';
  }

  function countLabel(value: number, singular: string, plural = `${singular}s`) {
    return `${value} ${value === 1 ? singular : plural}`;
  }

  function findOwnBillItem(bill: WifiBill | null | undefined) {
    if (!bill) {
      return null;
    }

    return data.myBills.find((item) => item.wifi_bill_id === bill.id) ?? null;
  }

  function createBillValue(
    field: 'month' | 'year' | 'nominal_per_person' | 'deadline_date' | 'status'
  ) {
    if (form?.action === 'createBill') {
      const values = form.values as Record<string, string> | undefined;
      const value = values?.[field];

      if (typeof value === 'string' && value !== '') {
        return value;
      }
    }

    return data.defaults[field];
  }

  function rejectReasonValue(memberID: string) {
    if (form?.action !== 'reject') {
      return '';
    }

    const values = form.values as Record<string, string> | undefined;
    if (values?.member_id !== memberID) {
      return '';
    }

    return values.reason ?? '';
  }

  function submitProofValue(field: 'proof_url' | 'note', fallback: string) {
    if (form?.action !== 'submitProof') {
      return fallback;
    }

    const values = form.values as Record<string, string> | undefined;
    return values?.[field] ?? fallback;
  }

  $: ownActiveBill = findOwnBillItem(data.activeBill?.bill);
</script>

<div class="space-y-4">
  <PageCard
    title="Wifi"
    description="Tagihan wifi bulanan, submit bukti transfer, dan verifikasi pembayaran untuk STEP 3."
  >
    {#if $navigating?.to?.url.pathname === '/wifi' || pendingAction}
      <div class="helper-box mb-4">
        <p class="helper-label">Loading</p>
        <p class="mt-2 text-sm text-slate-600">
          {pendingAction ? 'Memproses aksi wifi...' : 'Memuat ulang data wifi...'}
        </p>
      </div>
    {/if}

    {#if form?.message}
      <div class="helper-box-brand mb-4">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{form.message}</p>
      </div>
    {:else if form?.success}
      <div class="helper-box mb-4 border-emerald-200 bg-emerald-50/80">
        <p class="helper-label text-emerald-700">Success</p>
        <p class="mt-2 text-sm leading-6 text-emerald-800">{form.success}</p>
      </div>
    {/if}

    {#if data.loadError}
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{data.loadError}</p>
      </div>
    {:else}
      {#if data.activeBill}
        <section class="space-y-4">
          <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            <div class="stat-card bg-slate-950 text-white sm:col-span-2">
              <div class="flex flex-wrap items-center gap-2">
                <p class="helper-label text-slate-300">Active bill</p>
                <span class={billStatusClass(data.activeBill.bill.status)}>
                  {billStatusLabels[data.activeBill.bill.status]}
                </span>
              </div>

              <p class="mt-2 text-2xl font-semibold tracking-[-0.04em]">
                {formatMonthYear(data.activeBill.bill.month, data.activeBill.bill.year)}
              </p>
              <p class="mt-2 text-sm text-slate-300">
                Nominal {formatCurrency(data.activeBill.bill.nominal_per_person)} per orang, deadline{' '}
                {formatDate(data.activeBill.bill.deadline_date)}.
              </p>
            </div>

            <div class="stat-card bg-emerald-50">
              <p class="helper-label text-emerald-700">Verified</p>
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
                {data.activeBill.summary.verified_count}
              </p>
              <p class="mt-2 text-sm text-slate-600">
                {formatCurrency(data.activeBill.summary.total_collected)} terkumpul
              </p>
            </div>

            <div class="stat-card bg-sky-50">
              <p class="helper-label text-sky-700">Target</p>
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
                {formatCurrency(data.activeBill.summary.total_target)}
              </p>
              <p class="mt-2 text-sm text-slate-600">
                {countLabel(data.activeBill.summary.total_members, 'member')}
              </p>
            </div>
          </div>

          <div class="grid gap-3 sm:grid-cols-3">
            <div class="stat-card bg-white">
              <p class="helper-label">Pending</p>
              <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
                {data.activeBill.summary.pending_count}
              </p>
            </div>

            <div class="stat-card bg-white">
              <p class="helper-label">Unpaid</p>
              <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
                {data.activeBill.summary.unpaid_count}
              </p>
            </div>

            <div class="stat-card bg-white">
              <p class="helper-label">Rejected</p>
              <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
                {data.activeBill.summary.rejected_count}
              </p>
            </div>
          </div>
        </section>
      {:else}
        <div class="empty-state">
          Belum ada wifi bill aktif. Admin atau treasurer bisa membuat bill bulanan dari halaman ini.
        </div>
      {/if}

      {#if data.canManage}
        <div class="mt-4 grid gap-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,1.25fr)]">
          <section class="space-y-4">
            <article class="app-panel p-5">
              <div class="flex items-start justify-between gap-3">
                <div>
                  <p class="eyebrow">Create Bill</p>
                  <h2 class="section-title mt-1">Buat tagihan wifi bulanan</h2>
                  <p class="section-subtitle mt-2">
                    Record untuk member aktif akan digenerate otomatis dengan nominal default Rp20.000.
                  </p>
                </div>
              </div>

              <form
                method="POST"
                action="?/createBill"
                class="mt-4 space-y-4"
                use:enhance={enhanceWithAction('createBill')}
              >
                <div class="grid gap-4 sm:grid-cols-2">
                  <label>
                    <span class="field-label">Month</span>
                    <input
                      class="input-field"
                      type="number"
                      name="month"
                      min="1"
                      max="12"
                      required
                      value={createBillValue('month')}
                    />
                  </label>

                  <label>
                    <span class="field-label">Year</span>
                    <input
                      class="input-field"
                      type="number"
                      name="year"
                      min="2024"
                      required
                      value={createBillValue('year')}
                    />
                  </label>
                </div>

                <div class="grid gap-4 sm:grid-cols-2">
                  <label>
                    <span class="field-label">Nominal per person</span>
                    <input
                      class="input-field"
                      type="number"
                      name="nominal_per_person"
                      min="1"
                      step="1"
                      required
                      value={createBillValue('nominal_per_person')}
                    />
                  </label>

                  <label>
                    <span class="field-label">Deadline</span>
                    <input
                      class="input-field"
                      type="date"
                      name="deadline_date"
                      required
                      value={createBillValue('deadline_date')}
                    />
                  </label>
                </div>

                <label>
                  <span class="field-label">Status</span>
                  <select class="input-field" name="status">
                    <option value="active" selected={createBillValue('status') === 'active'}>Active</option>
                    <option value="draft" selected={createBillValue('status') === 'draft'}>Draft</option>
                    <option value="closed" selected={createBillValue('status') === 'closed'}>Closed</option>
                  </select>
                </label>

                <button
                  type="submit"
                  class="btn-primary w-full px-4 py-3"
                  disabled={pendingAction === 'createBill'}
                >
                  {pendingAction === 'createBill' ? 'Creating...' : 'Create wifi bill'}
                </button>
              </form>
            </article>

            <article class="app-panel p-5">
              <p class="eyebrow">Bill History</p>
              <h2 class="section-title mt-1">Daftar bill wifi</h2>
              <p class="section-subtitle mt-2">Satu bill per bulan-tahun, diurutkan dari yang terbaru.</p>

              {#if data.bills.length === 0}
                <div class="empty-state mt-4">
                  Belum ada bill wifi yang tercatat.
                </div>
              {:else}
                <div class="mt-4 space-y-3">
                  {#each data.bills as bill}
                    <article class="stat-card bg-white">
                      <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                        <div>
                          <div class="flex flex-wrap items-center gap-2">
                            <h3 class="text-base font-semibold text-ink">
                              {formatMonthYear(bill.month, bill.year)}
                            </h3>
                            <span class={billStatusClass(bill.status)}>
                              {billStatusLabels[bill.status]}
                            </span>
                          </div>
                          <p class="mt-2 text-sm text-slate-500">
                            Deadline {formatDate(bill.deadline_date)} • {formatCurrency(bill.nominal_per_person)} per orang
                          </p>
                        </div>

                        <div class="text-left sm:text-right">
                          <p class="text-lg font-semibold text-ink">
                            {formatCurrency(bill.summary.total_collected)}
                          </p>
                          <p class="mt-1 text-xs text-slate-500">
                            dari {formatCurrency(bill.summary.total_target)}
                          </p>
                        </div>
                      </div>

                      <div class="mt-4 grid gap-2 sm:grid-cols-4">
                        <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                          <p class="helper-label">Verified</p>
                          <p class="mt-2 text-sm font-semibold text-ink">{bill.summary.verified_count}</p>
                        </div>
                        <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                          <p class="helper-label">Pending</p>
                          <p class="mt-2 text-sm font-semibold text-ink">{bill.summary.pending_count}</p>
                        </div>
                        <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                          <p class="helper-label">Unpaid</p>
                          <p class="mt-2 text-sm font-semibold text-ink">{bill.summary.unpaid_count}</p>
                        </div>
                        <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                          <p class="helper-label">Rejected</p>
                          <p class="mt-2 text-sm font-semibold text-ink">{bill.summary.rejected_count}</p>
                        </div>
                      </div>
                    </article>
                  {/each}
                </div>
              {/if}
            </article>
          </section>

          <article class="app-panel p-5">
            <p class="eyebrow">Verification</p>
            <h2 class="section-title mt-1">Status pembayaran member</h2>
            <p class="section-subtitle mt-2">
              Endpoint verify/reject memakai <code>wifi_bill_members.id</code> sebagai <code>memberId</code>.
            </p>

            {#if !data.activeBill}
              <div class="empty-state mt-4">
                Bill aktif belum ada, jadi daftar pembayaran belum tersedia.
              </div>
            {:else if data.activeBill.members.length === 0}
              <div class="empty-state mt-4">
                Belum ada member record pada bill aktif ini.
              </div>
            {:else}
              <div class="mt-4 space-y-3">
                {#each data.activeBill.members as member}
                  <article class="stat-card bg-white">
                    <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                      <div class="min-w-0">
                        <div class="flex flex-wrap items-center gap-2">
                          <h3 class="text-base font-semibold text-ink">{member.user_name}</h3>
                          <span class={paymentStatusClass(member.payment_status)}>
                            {paymentStatusLabel(member.payment_status)}
                          </span>
                        </div>
                        <p class="mt-1 break-all text-sm text-slate-500">{member.user_email}</p>
                      </div>

                      <div class="text-left sm:text-right">
                        <p class="text-lg font-semibold text-ink">{formatCurrency(member.amount)}</p>
                        <p class="mt-1 text-xs text-slate-500">
                          Submit {formatDate(member.submitted_at, true)}
                        </p>
                      </div>
                    </div>

                    <div class="mt-4 grid gap-2">
                      <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                        <p class="helper-label">Proof</p>
                        <p class="mt-2 break-all text-sm text-slate-700">
                          {member.proof_url || 'Belum ada bukti transfer'}
                        </p>
                      </div>

                      {#if member.note}
                        <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                          <p class="helper-label">Note</p>
                          <p class="mt-2 text-sm text-slate-700">{member.note}</p>
                        </div>
                      {/if}

                      {#if member.rejection_reason}
                        <div class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3">
                          <p class="helper-label text-rose-700">Rejection reason</p>
                          <p class="mt-2 text-sm text-rose-700">{member.rejection_reason}</p>
                        </div>
                      {/if}
                    </div>

                    {#if member.payment_status === 'pending_verification'}
                      <div class="mt-4 grid gap-3">
                        <form
                          method="POST"
                          action="?/verify"
                          use:enhance={enhanceWithAction(`verify-${member.id}`)}
                        >
                          <input type="hidden" name="bill_id" value={data.activeBill.bill.id} />
                          <input type="hidden" name="member_id" value={member.id} />
                          <button
                            type="submit"
                            class="btn-primary w-full px-4 py-3"
                            disabled={pendingAction === `verify-${member.id}`}
                          >
                            {pendingAction === `verify-${member.id}` ? 'Verifying...' : 'Verify payment'}
                          </button>
                        </form>

                        <form
                          method="POST"
                          action="?/reject"
                          class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_auto]"
                          use:enhance={enhanceWithAction(`reject-${member.id}`)}
                        >
                          <input type="hidden" name="bill_id" value={data.activeBill.bill.id} />
                          <input type="hidden" name="member_id" value={member.id} />
                          <input
                            class="input-field"
                            type="text"
                            name="reason"
                            placeholder="Alasan reject"
                            value={rejectReasonValue(member.id)}
                            required
                          />
                          <button
                            type="submit"
                            class="btn-secondary px-4 py-3"
                            disabled={pendingAction === `reject-${member.id}`}
                          >
                            {pendingAction === `reject-${member.id}` ? 'Rejecting...' : 'Reject'}
                          </button>
                        </form>
                      </div>
                    {/if}
                  </article>
                {/each}
              </div>
            {/if}
          </article>
        </div>
      {:else}
        <div class="mt-4 grid gap-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,0.95fr)]">
          <section class="space-y-4">
            <article class="app-panel p-5">
              <p class="eyebrow">My Active Bill</p>
              <h2 class="section-title mt-1">Status pembayaran saya</h2>
              <p class="section-subtitle mt-2">
                Transfer dilakukan di luar aplikasi. Halaman ini hanya mencatat bukti dan status verifikasi.
              </p>

              {#if !data.activeBill || !ownActiveBill}
                <div class="empty-state mt-4">
                  Tidak ada bill aktif yang perlu Anda submit saat ini.
                </div>
              {:else}
                <div class="mt-4 grid gap-3 sm:grid-cols-2">
                  <div class="stat-card bg-slate-950 text-white">
                    <p class="helper-label text-slate-300">Bill</p>
                    <p class="mt-2 text-2xl font-semibold tracking-[-0.04em]">
                      {formatMonthYear(ownActiveBill.month, ownActiveBill.year)}
                    </p>
                    <p class="mt-2 text-sm text-slate-300">
                      Deadline {formatDate(ownActiveBill.deadline_date)}
                    </p>
                  </div>

                  <div class="stat-card bg-white">
                    <p class="helper-label">Current status</p>
                    <div class="mt-2">
                      <span class={paymentStatusClass(ownActiveBill.payment_status)}>
                        {paymentStatusLabel(ownActiveBill.payment_status)}
                      </span>
                    </div>
                    <p class="mt-2 text-lg font-semibold text-ink">
                      {formatCurrency(ownActiveBill.amount)}
                    </p>
                  </div>
                </div>

                <div class="mt-4 space-y-3">
                  {#if ownActiveBill.proof_url}
                    <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                      <p class="helper-label">Latest proof</p>
                      <p class="mt-2 break-all text-sm text-slate-700">{ownActiveBill.proof_url}</p>
                    </div>
                  {/if}

                  {#if ownActiveBill.note}
                    <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                      <p class="helper-label">Latest note</p>
                      <p class="mt-2 text-sm text-slate-700">{ownActiveBill.note}</p>
                    </div>
                  {/if}

                  {#if ownActiveBill.rejection_reason}
                    <div class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3">
                      <p class="helper-label text-rose-700">Rejection reason</p>
                      <p class="mt-2 text-sm text-rose-700">{ownActiveBill.rejection_reason}</p>
                    </div>
                  {/if}

                  <form
                    method="POST"
                    action="?/submitProof"
                    class="space-y-4"
                    use:enhance={enhanceWithAction('submitProof')}
                  >
                    <input type="hidden" name="bill_id" value={ownActiveBill.wifi_bill_id} />

                    <label>
                      <span class="field-label">Proof reference</span>
                      <input
                        class="input-field"
                        type="text"
                        name="proof_url"
                        placeholder="Link upload, path file, atau catatan bukti transfer"
                        value={submitProofValue('proof_url', ownActiveBill.proof_url || '')}
                        required
                      />
                    </label>

                    <label>
                      <span class="field-label">Optional note</span>
                      <textarea
                        class="input-field min-h-28 resize-y"
                        name="note"
                        placeholder="Contoh: transfer via BCA a.n. Ryan"
                      >{submitProofValue('note', ownActiveBill.note || '')}</textarea>
                    </label>

                    <button
                      type="submit"
                      class="btn-primary w-full px-4 py-3"
                      disabled={pendingAction === 'submitProof' || ownActiveBill.payment_status === 'verified'}
                    >
                      {ownActiveBill.payment_status === 'verified'
                        ? 'Already verified'
                        : pendingAction === 'submitProof'
                          ? 'Submitting...'
                          : 'Submit payment proof'}
                    </button>
                  </form>
                </div>
              {/if}
            </article>
          </section>

          <article class="app-panel p-5">
            <p class="eyebrow">History</p>
            <h2 class="section-title mt-1">Riwayat pembayaran wifi</h2>
            <p class="section-subtitle mt-2">Status terbaru submit Anda per bulan.</p>

            {#if data.myBills.length === 0}
              <div class="empty-state mt-4">
                Riwayat pembayaran wifi belum ada.
              </div>
            {:else}
              <div class="mt-4 space-y-3">
                {#each data.myBills as bill}
                  <article class="stat-card bg-white">
                    <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                      <div>
                        <div class="flex flex-wrap items-center gap-2">
                          <h3 class="text-base font-semibold text-ink">
                            {formatMonthYear(bill.month, bill.year)}
                          </h3>
                          <span class={paymentStatusClass(bill.payment_status)}>
                            {paymentStatusLabel(bill.payment_status)}
                          </span>
                        </div>
                        <p class="mt-2 text-sm text-slate-500">
                          Deadline {formatDate(bill.deadline_date)} • {formatCurrency(bill.amount)}
                        </p>
                      </div>

                      <div class="text-left sm:text-right">
                        <p class="text-sm font-medium text-ink">
                          Submitted {formatDate(bill.submitted_at, true)}
                        </p>
                        <p class="mt-1 text-xs text-slate-500">
                          Verified {formatDate(bill.verified_at, true)}
                        </p>
                      </div>
                    </div>

                    {#if bill.proof_url}
                      <div class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                        <p class="helper-label">Proof</p>
                        <p class="mt-2 break-all text-sm text-slate-700">{bill.proof_url}</p>
                      </div>
                    {/if}

                    {#if bill.rejection_reason}
                      <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3">
                        <p class="helper-label text-rose-700">Rejection reason</p>
                        <p class="mt-2 text-sm text-rose-700">{bill.rejection_reason}</p>
                      </div>
                    {/if}
                  </article>
                {/each}
              </div>
            {/if}
          </article>
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
