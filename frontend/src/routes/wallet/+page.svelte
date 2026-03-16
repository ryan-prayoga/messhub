<script lang="ts">
  import { navigating } from '$app/stores';
  import AppIcon from '$lib/components/AppIcon.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { PageData } from './$types';

  export let data: PageData;

  function formatCurrency(value: number) {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(value);
  }

  function formatDate(value: string) {
    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    }).format(new Date(value));
  }

  function isLink(value: string | null) {
    return typeof value === 'string' && /^https?:\/\//i.test(value);
  }
</script>

<div class="space-y-4">
  <PageCard
    eyebrow="Wallet"
    icon="lucide:wallet"
    title="Kantong Duafa"
    description="Ringkasan saldo dan daftar transaksi kas. Semua role bisa melihat, admin dan bendahara bisa mencatat transaksi baru."
  >
    {#if $navigating?.to?.url.pathname === '/wallet'}
      <PageSkeleton statCards={3} rows={3} />
    {/if}

    {#if data.loadError}
      <StatePanel tone="error" title="Belum bisa memuat wallet" message={data.loadError} />
    {:else if data.summary}
      <div class="grid gap-3 sm:grid-cols-3">
        <div class="stat-card bg-slate-950 text-white">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="helper-label text-slate-300">Saldo saat ini</p>
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em]">
                {formatCurrency(data.summary.balance)}
              </p>
              <p class="mt-2 text-sm text-slate-300">Saldo berjalan dari semua pemasukan dan pengeluaran.</p>
            </div>

            <AppIcon icon="lucide:landmark" className="h-5 w-5 text-white/80" />
          </div>
        </div>

        <div class="stat-card bg-emerald-50">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="helper-label text-emerald-700">Total pemasukan</p>
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
                {formatCurrency(data.summary.total_income)}
              </p>
            </div>

            <AppIcon icon="lucide:trending-up" className="h-5 w-5 text-emerald-700" />
          </div>
        </div>

        <div class="stat-card bg-rose-50">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="helper-label text-rose-700">Total pengeluaran</p>
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
                {formatCurrency(data.summary.total_expense)}
              </p>
            </div>

            <AppIcon icon="lucide:trending-down" className="h-5 w-5 text-rose-700" />
          </div>
        </div>
      </div>

      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="helper-box grow">
          <p class="helper-label">Riwayat transaksi</p>
          <p class="mt-2 text-sm leading-6 text-slate-600">
            Semua transaksi diurutkan dari yang terbaru. Total data saat ini: {data.pagination.total_items}.
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
          {#if data.user?.role === 'admin'}
            <a href="/admin/import/wallet" class="btn-secondary shrink-0 px-4 py-3">Impor CSV</a>
          {/if}

          {#if data.canCreate}
            <a href="/wallet/new" class="btn-primary shrink-0 px-4 py-3">
              Tambah transaksi
            </a>
          {/if}
        </div>
      </div>

      {#if data.transactions.length === 0}
        <StatePanel
          tone="empty"
          title="Belum ada transaksi kas"
          message="Belum ada pemasukan atau pengeluaran yang tercatat. Tambahkan transaksi pertama atau impor CSV kas lama agar saldo wallet mulai terbentuk."
          actionHref={data.canCreate ? '/wallet/new' : null}
          actionLabel={data.canCreate ? 'Tambah transaksi pertama' : ''}
        />
      {:else}
        <div class="mt-4 space-y-3">
          {#each data.transactions as transaction}
            <article class="stat-card bg-white">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-2">
                    <div class="nav-link-icon h-10 w-10 rounded-[16px]">
                      <AppIcon
                        icon={transaction.type === 'income' ? 'lucide:arrow-down-left' : 'lucide:arrow-up-right'}
                        className="h-4 w-4"
                      />
                    </div>
                    <h3 class="text-base font-semibold text-ink">{transaction.category}</h3>
                    <span
                      class={transaction.type === 'income'
                        ? 'badge bg-emerald-100 text-emerald-700'
                        : 'badge bg-rose-100 text-rose-700'}
                    >
                      {transaction.type === 'income' ? 'Pemasukan' : 'Pengeluaran'}
                    </span>
                  </div>

                  <p class="mt-2 text-sm leading-6 text-slate-500">{transaction.description}</p>
                </div>

                <div class="text-left sm:text-right">
                  <p class={transaction.type === 'income'
                    ? 'text-lg font-semibold text-emerald-700'
                    : 'text-lg font-semibold text-rose-700'}>
                    {transaction.type === 'income' ? '+' : '-'}{formatCurrency(transaction.amount)}
                  </p>
                  <p class="mt-1 text-xs text-slate-500">{formatDate(transaction.transaction_date)}</p>
                </div>
              </div>

              <div class="mt-4 grid gap-3 sm:grid-cols-2">
                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Tanggal transaksi</p>
                  <p class="mt-2 text-sm font-medium text-ink">{formatDate(transaction.transaction_date)}</p>
                </div>

                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Dicatat oleh</p>
                  <p class="mt-2 text-sm font-medium text-ink">
                    {transaction.created_by_name || transaction.created_by}
                  </p>
                </div>
              </div>

              {#if transaction.proof_url}
                <div class="mt-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Bukti</p>
                  {#if isLink(transaction.proof_url)}
                    <a
                      href={transaction.proof_url}
                      class="mt-2 inline-flex items-center gap-2 text-sm font-medium text-ink underline"
                      target="_blank"
                      rel="noreferrer"
                    >
                      <AppIcon icon="lucide:external-link" className="h-4 w-4" />
                      Buka bukti transaksi
                    </a>
                  {:else}
                    <p class="mt-2 break-all text-sm font-medium text-ink">{transaction.proof_url}</p>
                  {/if}
                </div>
              {/if}
            </article>
          {/each}
        </div>

        {#if data.pagination.total_pages > 1}
          <div class="mt-4 flex items-center justify-between gap-3">
            <a
              href={data.pagination.page > 1 ? `/wallet?page=${data.pagination.page - 1}` : '/wallet'}
              class={`btn-secondary px-4 py-3 ${data.pagination.page <= 1 ? 'pointer-events-none opacity-50' : ''}`}
            >
              Sebelumnya
            </a>

            <p class="text-sm text-slate-500">
              Halaman {data.pagination.page} / {data.pagination.total_pages}
            </p>

            <a
              href={data.pagination.page < data.pagination.total_pages
                ? `/wallet?page=${data.pagination.page + 1}`
                : `/wallet?page=${data.pagination.page}`}
              class={`btn-secondary px-4 py-3 ${
                data.pagination.page >= data.pagination.total_pages
                  ? 'pointer-events-none opacity-50'
                  : ''
              }`}
            >
              Berikutnya
            </a>
          </div>
        {/if}
      {/if}
    {/if}
  </PageCard>
</div>
