<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { SharedExpense, SharedExpenseStatus } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  let pendingAction: string | null = null;

  const statusLabels: Record<SharedExpenseStatus, string> = {
    personal: 'Tanggung pribadi',
    fronted: 'Talangan',
    partially_reimbursed: 'Diganti sebagian',
    reimbursed: 'Sudah lunas'
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

  function formatCurrency(value: number) {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(value);
  }

  function formatDate(value: string | null) {
    if (!value) {
      return '-';
    }

    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    }).format(new Date(value));
  }

  function statusBadgeClass(status: SharedExpenseStatus) {
    if (status === 'reimbursed') {
      return 'badge-success';
    }

    if (status === 'partially_reimbursed') {
      return 'badge-warning';
    }

    if (status === 'fronted') {
      return 'badge-strong';
    }

    return 'badge-muted';
  }

  function createValue(field: string) {
    if (form?.action !== 'createExpense') {
      return field === 'status' ? 'fronted' : '';
    }

    const values = form.values as Record<string, string> | undefined;
    return values?.[field] ?? '';
  }

  function editValue(expense: SharedExpense, field: string) {
    if (form?.action === 'updateExpense') {
      const values = form.values as Record<string, string> | undefined;
      if (values?.expense_id === expense.id) {
        return values[field] ?? '';
      }
    }

    switch (field) {
      case 'expense_date':
        return expense.expense_date.slice(0, 10);
      case 'category':
        return expense.category;
      case 'description':
        return expense.description;
      case 'amount':
        return String(expense.amount);
      case 'paid_by_user_id':
        return expense.paid_by_user_id;
      case 'status':
        return expense.status;
      case 'notes':
        return expense.notes ?? '';
      case 'proof_url':
        return expense.proof_url ?? '';
      default:
        return '';
    }
  }
</script>

<div class="space-y-4">
  <PageCard
    eyebrow="Patungan"
    icon="lucide:receipt"
    title="Pengeluaran Bersama"
    description="Catat pengeluaran non-kas, siapa yang menalangi, dan progres penggantiannya tanpa mengganggu wallet utama."
  >
    {#if $navigating?.to?.url.pathname === '/shared-expenses' || pendingAction}
      <PageSkeleton statCards={3} rows={3} />
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
      <StatePanel tone="error" title="Belum bisa memuat pengeluaran bersama" message={data.loadError} />
    {:else}
      {#if data.summary}
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <div class="stat-card bg-slate-950 text-white">
            <p class="helper-label text-slate-300">Total tercatat</p>
            <p class="mt-2 text-2xl font-semibold tracking-[-0.04em]">{formatCurrency(data.summary.total_amount)}</p>
            <p class="mt-2 text-sm text-slate-300">{data.summary.total_count} transaksi pengeluaran bersama.</p>
          </div>
          <div class="stat-card bg-white">
            <p class="helper-label">Talangan aktif</p>
            <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">{data.summary.fronted_count}</p>
            <p class="mt-2 text-sm text-muted">{formatCurrency(data.summary.outstanding_amount)} belum lunas.</p>
          </div>
          <div class="stat-card bg-white">
            <p class="helper-label">Bulan ini</p>
            <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">{formatCurrency(data.summary.this_month_amount)}</p>
            <p class="mt-2 text-sm text-muted">Akumulasi pengeluaran bulan berjalan.</p>
          </div>
          <div class="stat-card bg-white">
            <p class="helper-label">Arah penggunaan</p>
            <p class="mt-2 text-base font-semibold text-ink">Terpisah dari wallet</p>
            <p class="mt-2 text-sm text-muted">Tetap transparan tanpa mengurangi saldo Kantong Duafa.</p>
          </div>
        </div>
      {/if}

      {#if data.canManage}
        <article class="app-panel mt-4 p-5">
          <p class="eyebrow">Catat Baru</p>
          <h2 class="section-title mt-1">Tambah pengeluaran bersama</h2>
          <form method="POST" action="?/createExpense" class="mt-4 space-y-4" use:enhance={enhanceWithAction('createExpense')}>
            <div class="grid gap-4 sm:grid-cols-2">
              <label>
                <span class="field-label">Tanggal</span>
                <input class="input-field" type="date" name="expense_date" value={createValue('expense_date')} required />
              </label>
              <label>
                <span class="field-label">Kategori</span>
                <input class="input-field" type="text" name="category" value={createValue('category')} placeholder="Contoh: Kebersihan" required />
              </label>
            </div>

            <label>
              <span class="field-label">Deskripsi</span>
              <input class="input-field" type="text" name="description" value={createValue('description')} placeholder="Contoh: Bayar alat pel dapur" required />
            </label>

            <div class="grid gap-4 sm:grid-cols-3">
              <label>
                <span class="field-label">Nominal</span>
                <input class="input-field" type="number" min="1" step="1" name="amount" value={createValue('amount')} required />
              </label>
              <label>
                <span class="field-label">Dibayar oleh</span>
                <select class="input-field" name="paid_by_user_id" required>
                  <option value="">Pilih anggota</option>
                  {#each data.payers as payer}
                    <option value={payer.id} selected={createValue('paid_by_user_id') === payer.id}>{payer.name}</option>
                  {/each}
                </select>
              </label>
              <label>
                <span class="field-label">Status</span>
                <select class="input-field" name="status">
                  {#each Object.entries(statusLabels) as [value, label]}
                    <option value={value} selected={createValue('status') === value}>{label}</option>
                  {/each}
                </select>
              </label>
            </div>

            <label>
              <span class="field-label">Catatan</span>
              <textarea class="input-field min-h-[96px]" name="notes" placeholder="Catatan tambahan, konteks penggantian, atau detail pelunasan.">{createValue('notes')}</textarea>
            </label>

            <label>
              <span class="field-label">Bukti (opsional)</span>
              <input class="input-field" type="text" name="proof_url" value={createValue('proof_url')} placeholder="Link struk, foto, atau referensi dokumen" />
            </label>

            <button type="submit" class="btn-primary w-full px-4 py-3" disabled={pendingAction === 'createExpense'}>
              {pendingAction === 'createExpense' ? 'Menyimpan...' : 'Simpan pengeluaran bersama'}
            </button>
          </form>
        </article>
      {/if}

      {#if data.expenses.length === 0}
        <StatePanel
          tone="empty"
          title="Belum ada pengeluaran bersama"
          message="Catatan pengeluaran non-kas bersama akan muncul di sini agar semua penghuni bisa melihat transparansi talangan dan penggantian."
        />
      {:else}
        <div class="mt-4 space-y-4">
          {#each data.expenses as expense}
            <article class="app-panel p-5">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span class={statusBadgeClass(expense.status)}>{statusLabels[expense.status]}</span>
                    <span class="badge-muted">{expense.category}</span>
                  </div>
                  <h2 class="section-title mt-3">{expense.description}</h2>
                  <p class="mt-2 text-sm text-muted">Dibayar oleh {expense.paid_by_user_name} • {formatDate(expense.expense_date)}</p>
                </div>

                <div class="text-right">
                  <p class="text-2xl font-semibold tracking-[-0.04em] text-ink">{formatCurrency(expense.amount)}</p>
                  <p class="mt-1 text-xs text-dusty">Dicatat {expense.created_by_name}</p>
                </div>
              </div>

              {#if expense.notes}
                <div class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Catatan</p>
                  <p class="mt-2 text-sm text-slate-700">{expense.notes}</p>
                </div>
              {/if}

              {#if expense.proof_url}
                <div class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Bukti</p>
                  <p class="mt-2 break-all text-sm text-slate-700">{expense.proof_url}</p>
                </div>
              {/if}

              {#if data.canManage}
                <details class="mt-4 rounded-2xl border border-line bg-white/70 p-4">
                  <summary class="cursor-pointer text-sm font-semibold text-ink">Ubah pengeluaran ini</summary>
                  <form method="POST" action="?/updateExpense" class="mt-4 space-y-4" use:enhance={enhanceWithAction(`update-${expense.id}`)}>
                    <input type="hidden" name="expense_id" value={expense.id} />
                    <div class="grid gap-4 sm:grid-cols-2">
                      <label>
                        <span class="field-label">Tanggal</span>
                        <input class="input-field" type="date" name="expense_date" value={editValue(expense, 'expense_date')} required />
                      </label>
                      <label>
                        <span class="field-label">Kategori</span>
                        <input class="input-field" type="text" name="category" value={editValue(expense, 'category')} required />
                      </label>
                    </div>
                    <label>
                      <span class="field-label">Deskripsi</span>
                      <input class="input-field" type="text" name="description" value={editValue(expense, 'description')} required />
                    </label>
                    <div class="grid gap-4 sm:grid-cols-3">
                      <label>
                        <span class="field-label">Nominal</span>
                        <input class="input-field" type="number" min="1" step="1" name="amount" value={editValue(expense, 'amount')} required />
                      </label>
                      <label>
                        <span class="field-label">Dibayar oleh</span>
                        <select class="input-field" name="paid_by_user_id" required>
                          {#each data.payers as payer}
                            <option value={payer.id} selected={editValue(expense, 'paid_by_user_id') === payer.id}>{payer.name}</option>
                          {/each}
                        </select>
                      </label>
                      <label>
                        <span class="field-label">Status</span>
                        <select class="input-field" name="status">
                          {#each Object.entries(statusLabels) as [value, label]}
                            <option value={value} selected={editValue(expense, 'status') === value}>{label}</option>
                          {/each}
                        </select>
                      </label>
                    </div>
                    <label>
                      <span class="field-label">Catatan</span>
                      <textarea class="input-field min-h-[96px]" name="notes">{editValue(expense, 'notes')}</textarea>
                    </label>
                    <label>
                      <span class="field-label">Bukti</span>
                      <input class="input-field" type="text" name="proof_url" value={editValue(expense, 'proof_url')} />
                    </label>
                    <button type="submit" class="btn-secondary w-full px-4 py-3" disabled={pendingAction === `update-${expense.id}`}>
                      {pendingAction === `update-${expense.id}` ? 'Menyimpan...' : 'Perbarui pengeluaran bersama'}
                    </button>
                  </form>
                </details>
              {/if}
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
