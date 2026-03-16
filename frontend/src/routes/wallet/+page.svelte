<script lang="ts">
  import { navigating } from '$app/stores';
  import PageCard from '$lib/components/PageCard.svelte';
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
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(new Date(value));
  }
</script>

<div class="space-y-4">
  <PageCard
    title="Wallet"
    description="Ringkasan Kantong Duafa dan daftar transaksi terbaru. Semua role bisa melihat, admin dan treasurer bisa mencatat transaksi baru."
  >
    {#if $navigating?.to?.url.pathname === '/wallet'}
      <div class="helper-box mb-4">
        <p class="helper-label">Loading</p>
        <p class="mt-2 text-sm text-slate-600">Memuat ulang saldo dan transaksi wallet...</p>
      </div>
    {/if}

    {#if data.loadError}
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{data.loadError}</p>
      </div>
    {:else if data.summary}
      <div class="grid gap-3 sm:grid-cols-3">
        <div class="stat-card bg-slate-950 text-white">
          <p class="helper-label text-slate-300">Current balance</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em]">
            {formatCurrency(data.summary.balance)}
          </p>
          <p class="mt-2 text-sm text-slate-300">Saldo berjalan dari semua pemasukan dan pengeluaran.</p>
        </div>

        <div class="stat-card bg-emerald-50">
          <p class="helper-label text-emerald-700">Total income</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {formatCurrency(data.summary.total_income)}
          </p>
        </div>

        <div class="stat-card bg-rose-50">
          <p class="helper-label text-rose-700">Total expense</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {formatCurrency(data.summary.total_expense)}
          </p>
        </div>
      </div>

      <div class="mt-4 flex items-center justify-between gap-3">
        <div class="helper-box grow">
          <p class="helper-label">Audit trail</p>
          <p class="mt-2 text-sm leading-6 text-slate-600">
            Semua transaksi diurutkan dari yang terbaru. Total data saat ini: {data.pagination.total_items}.
          </p>
        </div>

        {#if data.canCreate}
          <a href="/wallet/new" class="btn-primary shrink-0 px-4 py-3">
            New transaction
          </a>
        {/if}
      </div>

      {#if data.transactions.length === 0}
        <div class="mt-4 empty-state">
          Belum ada transaksi wallet. Setelah admin atau treasurer mencatat transaksi pertama, daftar ini akan terisi.
        </div>
      {:else}
        <div class="mt-4 space-y-3">
          {#each data.transactions as transaction}
            <article class="stat-card bg-white">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-2">
                    <h3 class="text-base font-semibold text-ink">{transaction.category}</h3>
                    <span
                      class={transaction.type === 'income'
                        ? 'badge bg-emerald-100 text-emerald-700'
                        : 'badge bg-rose-100 text-rose-700'}
                    >
                      {transaction.type}
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
                  <p class="mt-1 text-xs text-slate-500">{formatDate(transaction.created_at)}</p>
                </div>
              </div>

              <div class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Created by</p>
                <p class="mt-2 text-sm font-medium text-ink">
                  {transaction.created_by_name || transaction.created_by}
                </p>
              </div>
            </article>
          {/each}
        </div>

        {#if data.pagination.total_pages > 1}
          <div class="mt-4 flex items-center justify-between gap-3">
            <a
              href={data.pagination.page > 1 ? `/wallet?page=${data.pagination.page - 1}` : '/wallet'}
              class={`btn-secondary px-4 py-3 ${data.pagination.page <= 1 ? 'pointer-events-none opacity-50' : ''}`}
            >
              Previous
            </a>

            <p class="text-sm text-slate-500">
              Page {data.pagination.page} / {data.pagination.total_pages}
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
              Next
            </a>
          </div>
        {/if}
      {/if}
    {/if}
  </PageCard>
</div>
