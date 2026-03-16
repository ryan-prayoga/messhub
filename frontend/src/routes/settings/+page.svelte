<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

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

  function valueFor(field: keyof NonNullable<PageData['settings']>) {
    if (form?.action === 'updateSettings') {
      const values = form.values as Record<string, string> | undefined;
      const value = values?.[field];
      if (typeof value === 'string') {
        return value;
      }
    }

    return data.settings?.[field]?.toString() ?? '';
  }

  function formatCurrency(value: number) {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(value);
  }

  function formatDateTime(value: string | null | undefined) {
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
</script>

<div class="space-y-4">
  <PageCard
    eyebrow="Settings"
    icon="lucide:settings-2"
    title="Pengaturan Mess"
    description="Kelola nama mess, nominal wifi, deadline, dan rekening tujuan dari satu tempat."
  >
    {#if $navigating?.to?.url.pathname === '/settings' || pendingAction}
      <PageSkeleton statCards={2} rows={2} />
    {/if}

    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Akses ditolak"
        message="Halaman ini hanya dapat dibuka oleh admin mess."
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
        <FeedbackBanner tone="success" title="Berhasil" message={form.success} />
      {/if}

      {#if data.loadError || !data.settings}
        <StatePanel tone="error" title="Gagal memuat" message={data.loadError ?? 'Pengaturan belum tersedia.'} />
      {:else}
        <div class="grid gap-4 xl:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
          <article class="app-panel p-5">
            <p class="eyebrow">Konfigurasi</p>
            <h2 class="section-title mt-1">Pengaturan operasional</h2>

            <form
              method="POST"
              action="?/updateSettings"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('updateSettings')}
            >
              <label>
                <span class="field-label">Nama mess</span>
                <input
                  class="input-field"
                  type="text"
                  name="mess_name"
                  value={valueFor('mess_name')}
                  required
                />
              </label>

              <div class="grid gap-4 sm:grid-cols-2">
                <label>
                  <span class="field-label">Nominal wifi</span>
                  <input
                    class="input-field"
                    type="number"
                    name="wifi_price"
                    min="1"
                    step="1"
                    inputmode="numeric"
                    value={valueFor('wifi_price')}
                    required
                  />
                </label>

                <label>
                  <span class="field-label">Tanggal deadline wifi</span>
                  <input
                    class="input-field"
                    type="number"
                    name="wifi_deadline_day"
                    min="1"
                    max="31"
                    step="1"
                    inputmode="numeric"
                    value={valueFor('wifi_deadline_day')}
                    required
                  />
                </label>
              </div>

              <label>
                <span class="field-label">Nama rekening</span>
                <input
                  class="input-field"
                  type="text"
                  name="bank_account_name"
                  value={valueFor('bank_account_name')}
                  required
                />
              </label>

              <label>
                <span class="field-label">Nomor rekening</span>
                <input
                  class="input-field"
                  type="text"
                  name="bank_account_number"
                  value={valueFor('bank_account_number')}
                  required
                />
              </label>

              <button
                type="submit"
                class="btn-primary w-full px-4 py-3"
                disabled={pendingAction === 'updateSettings'}
              >
                {pendingAction === 'updateSettings' ? 'Menyimpan...' : 'Simpan pengaturan'}
              </button>
            </form>
          </article>

          <section class="space-y-4">
            <article class="stat-card bg-slate-950 text-white">
              <p class="helper-label text-slate-300">Rekening tujuan saat ini</p>
              <p class="mt-2 text-2xl font-semibold tracking-[-0.03em]">{data.settings.mess_name}</p>
              <p class="mt-2 text-sm text-slate-300">
                Wifi {formatCurrency(data.settings.wifi_price)} per orang, deadline tanggal {data.settings.wifi_deadline_day}.
              </p>
              <p class="mt-4 text-sm text-slate-200">
                {data.settings.bank_account_name} • {data.settings.bank_account_number}
              </p>
            </article>

            <article class="app-panel p-5">
              <p class="eyebrow">Status layanan</p>
              <h2 class="section-title mt-1">Pantauan singkat</h2>

              {#if data.systemStatus}
                <div class="mt-4 grid gap-3 sm:grid-cols-2">
                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Status aplikasi</p>
                    <p class="mt-2 text-sm font-semibold text-ink">{data.systemStatus.status}</p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Koneksi data</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {data.systemStatus.database_reachable
                        ? 'Terhubung dengan baik'
                        : 'Perlu pengecekan'}
                    </p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Waktu server</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {formatDateTime(data.systemStatus.server_time)}
                    </p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Catatan</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      Gunakan halaman ini untuk memastikan pengaturan dan layanan utama tetap siap dipakai.
                    </p>
                  </div>
                </div>
              {:else}
                <StatePanel tone="empty" title="Belum ada data" message="Status layanan belum tersedia saat ini." />
              {/if}
            </article>
          </section>
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
