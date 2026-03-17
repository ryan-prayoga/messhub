<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ImportCommitResult, WalletImportPreview } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  let pendingAction: 'preview' | 'commit' | null = null;
  let preview: WalletImportPreview | null = form?.preview ?? null;
  let committed: ImportCommitResult | null = form?.committed ?? null;

  $: if (form?.preview) {
    preview = form.preview;
    committed = null;
  }

  $: if (form?.committed) {
    committed = form.committed;
    preview = null;
  }

  function enhanceWithAction(actionName: 'preview' | 'commit'): SubmitFunction {
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

  function rowStatusClass(status: 'valid' | 'invalid') {
    return status === 'valid'
      ? 'badge bg-emerald-100 text-emerald-700'
      : 'badge bg-rose-100 text-rose-700';
  }
</script>

<div class="space-y-5 lg:space-y-6">
  <PageCard
    title="Impor Transaksi Kas"
    description="Unggah CSV kas lama, cek hasil parsing dan total nominal, lalu commit setelah preview."
  >
    {#if $navigating?.to?.url.pathname === '/admin/import/wallet' || pendingAction}
      <StatePanel
        tone="loading"
        title="Memuat"
        message={pendingAction === 'commit'
          ? 'Menyimpan hasil impor transaksi...'
          : 'Memproses preview CSV transaksi...'}
      />
    {/if}

    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Akses ditolak"
        message="Halaman impor hanya tersedia untuk admin mess."
      />
    {:else}
      {#if form?.message}
        <StatePanel
          tone="error"
          title="Gagal memproses"
          message={form.message}
          requestId={form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null}
        />
      {:else if form?.success}
        <div class="helper-box mb-4 border-emerald-200 bg-emerald-50/80">
          <p class="helper-label text-emerald-700">Berhasil</p>
          <p class="mt-2 text-sm leading-6 text-emerald-800">{form.success}</p>
        </div>
      {/if}

      <div class="grid gap-4 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)]">
        <section class="space-y-4">
          <article class="app-panel p-5">
            <p class="eyebrow">Langkah 1</p>
            <h2 class="section-title mt-1">Siapkan file CSV</h2>
            <p class="section-subtitle mt-2">
              Sistem menerima template MessHub dan juga header yang umum dipakai di sheet kas lama seperti Tanggal, Deskripsi, Pemasukan, Pengeluaran, dan Bukti.
            </p>

            <div class="mt-4 flex flex-wrap gap-3">
              <a href="/templates/wallet-import-template.csv" class="btn-secondary px-4 py-3" download>
                Unduh template transaksi
              </a>
              <a href="/admin/import" class="btn-secondary px-4 py-3">Kembali ke pusat impor</a>
            </div>

            <div class="mt-4 helper-box">
              <p class="helper-label">Aturan utama</p>
              <p class="mt-2 text-sm leading-6 text-slate-600">
                Isi salah satu kolom income atau expense. Saldo tidak diimpor karena akan dihitung ulang dari transaksi yang berhasil masuk.
              </p>
            </div>
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Langkah 2</p>
            <h2 class="section-title mt-1">Unggah dan preview</h2>

            <form
              method="POST"
              action="?/preview"
              enctype="multipart/form-data"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('preview')}
            >
              <label>
                <span class="field-label">File CSV transaksi</span>
                <input class="input-field file:mr-3 file:rounded-xl file:border-0 file:bg-slate-900 file:px-3 file:py-2 file:text-sm file:font-semibold file:text-white" type="file" name="file" accept=".csv,text/csv" required />
              </label>

              <button type="submit" class="btn-primary w-full px-4 py-3" disabled={pendingAction === 'preview'}>
                {pendingAction === 'preview' ? 'Memproses preview...' : 'Buat preview impor'}
              </button>
            </form>
          </article>
        </section>

        <section class="space-y-4">
          <article class="stat-card bg-slate-950 text-white">
            <p class="helper-label text-slate-300">Inferensi kategori</p>
            <p class="mt-2 text-sm leading-6 text-slate-200">
              Sistem akan mencoba mengenali kategori sederhana seperti wifi, hibah, galon, dan kebersihan. Sisanya masuk ke kategori lainnya.
            </p>
          </article>

          {#if committed}
            <article class="app-panel p-5">
              <p class="eyebrow">Hasil commit</p>
              <h2 class="section-title mt-1">Impor selesai</h2>

              <div class="mt-4 grid gap-3 sm:grid-cols-3">
                <div class="stat-card">
                  <p class="helper-label">Berhasil</p>
                  <p class="mt-2 text-3xl font-semibold text-ink">{committed.imported_rows}</p>
                </div>

                <div class="stat-card bg-white">
                  <p class="helper-label">Tidak masuk</p>
                  <p class="mt-2 text-3xl font-semibold text-ink">{committed.failed_rows}</p>
                </div>

                <div class="stat-card bg-white">
                  <p class="helper-label">Total baris</p>
                  <p class="mt-2 text-3xl font-semibold text-ink">{committed.total_rows}</p>
                </div>
              </div>

              <div class="mt-4 grid gap-3 sm:grid-cols-2">
                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Total pemasukan</p>
                  <p class="mt-2 text-sm font-semibold text-emerald-700">{formatCurrency(committed.total_income ?? 0)}</p>
                </div>

                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Total pengeluaran</p>
                  <p class="mt-2 text-sm font-semibold text-rose-700">{formatCurrency(committed.total_expense ?? 0)}</p>
                </div>
              </div>
            </article>
          {/if}
        </section>
      </div>

      {#if preview}
        <section class="mt-4 space-y-4">
          <article class="app-panel p-5">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <p class="eyebrow">Preview</p>
                <h2 class="section-title mt-1">{preview.file_name}</h2>
                <p class="section-subtitle mt-2">
                  Pastikan semua nominal, tipe transaksi, dan tanggal sudah sesuai sebelum commit.
                </p>
              </div>

              <span class={`badge ${preview.can_commit ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'}`}>
                {preview.can_commit ? 'Siap diimpor' : 'Perlu perbaikan'}
              </span>
            </div>

            <div class="mt-4 grid gap-3 sm:grid-cols-5">
              <div class="stat-card">
                <p class="helper-label">Total baris</p>
                <p class="mt-2 text-3xl font-semibold text-ink">{preview.summary.total_rows}</p>
              </div>

              <div class="stat-card bg-white">
                <p class="helper-label">Valid</p>
                <p class="mt-2 text-3xl font-semibold text-emerald-700">{preview.summary.valid_rows}</p>
              </div>

              <div class="stat-card bg-white">
                <p class="helper-label">Invalid</p>
                <p class="mt-2 text-3xl font-semibold text-rose-700">{preview.summary.invalid_rows}</p>
              </div>

              <div class="stat-card bg-white">
                <p class="helper-label">Pemasukan</p>
                <p class="mt-2 text-sm font-semibold text-emerald-700">{formatCurrency(preview.summary.total_income)}</p>
              </div>

              <div class="stat-card bg-white">
                <p class="helper-label">Pengeluaran</p>
                <p class="mt-2 text-sm font-semibold text-rose-700">{formatCurrency(preview.summary.total_expense)}</p>
              </div>
            </div>

            {#if preview.warnings.length > 0}
              <div class="mt-4 space-y-2">
                {#each preview.warnings as warning}
                  <div class="helper-box-brand">
                    <p class="helper-label text-sky-700">Peringatan</p>
                    <p class="mt-2 text-sm leading-6 text-slate-700">{warning.message}</p>
                  </div>
                {/each}
              </div>
            {/if}
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Langkah 3</p>
            <h2 class="section-title mt-1">Commit impor</h2>

            <form method="POST" action="?/commit" class="mt-4 space-y-4" use:enhance={enhanceWithAction('commit')}>
              <input type="hidden" name="job_id" value={preview.job_id} />

              <button type="submit" class="btn-primary w-full px-4 py-3" disabled={!preview.can_commit || pendingAction === 'commit'}>
                {pendingAction === 'commit' ? 'Menyimpan impor...' : 'Commit impor transaksi'}
              </button>
            </form>
          </article>

          <article class="space-y-3">
            {#each preview.rows as row}
              <div class="stat-card bg-white">
                <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                  <div class="min-w-0">
                    <div class="flex flex-wrap items-center gap-2">
                      <p class="text-sm font-semibold text-ink">Baris {row.row_number}</p>
                      <span class={rowStatusClass(row.status)}>
                        {row.status === 'valid' ? 'Valid' : 'Invalid'}
                      </span>
                    </div>

                    <p class="mt-2 text-sm text-slate-700">{row.description || '-'}</p>
                    <p class="mt-1 text-xs text-slate-500">
                      {row.normalized_transaction_date || row.transaction_date || '-'}
                    </p>
                  </div>

                  <div class="grid gap-2 text-sm text-slate-600 sm:text-right">
                    <p>Tipe: {row.type === 'income' ? 'Pemasukan' : row.type === 'expense' ? 'Pengeluaran' : '-'}</p>
                    <p>Nominal: {row.amount ? formatCurrency(row.amount) : '-'}</p>
                    <p>Kategori: {row.category || '-'}</p>
                  </div>
                </div>

                {#if row.proof}
                  <div class="mt-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Bukti</p>
                    <p class="mt-2 break-all text-sm text-slate-700">{row.proof}</p>
                  </div>
                {/if}

                {#if row.errors.length > 0}
                  <div class="mt-3 rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3">
                    <p class="helper-label text-rose-700">Alasan invalid</p>
                    <div class="mt-2 space-y-1">
                      {#each row.errors as item}
                        <p class="text-sm leading-6 text-rose-800">{item}</p>
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </article>
        </section>
      {/if}
    {/if}
  </PageCard>
</div>
