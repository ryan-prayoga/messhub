<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import PageCard from '$lib/components/PageCard.svelte';
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

  function formatUptime(seconds: number) {
    if (!seconds || seconds < 60) {
      return `${seconds || 0}s`;
    }

    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (hours === 0) {
      return `${minutes}m`;
    }

    return `${hours}h ${minutes}m`;
  }
</script>

<div class="space-y-4">
  <PageCard
    title="Settings"
    description="Admin panel untuk konfigurasi mess, nominal wifi, deadline, dan info transfer."
  >
    {#if $navigating?.to?.url.pathname === '/settings' || pendingAction}
      <StatePanel
        tone="loading"
        title="Loading"
        message={pendingAction ? 'Menyimpan pengaturan...' : 'Memuat ulang data pengaturan...'}
      />
    {/if}

    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Forbidden"
        message={`Role ${data.user?.role} tidak memiliki akses ke admin settings.`}
      />
    {:else}
      {#if form?.message}
        <StatePanel
          tone="error"
          title="Error"
          message={form.message}
          requestId={form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null}
        />
      {:else if form?.success}
        <div class="helper-box mb-4 border-emerald-200 bg-emerald-50/80">
          <p class="helper-label text-emerald-700">Success</p>
          <p class="mt-2 text-sm leading-6 text-emerald-800">{form.success}</p>
        </div>
      {/if}

      {#if data.loadError || !data.settings}
        <StatePanel tone="error" title="Error" message={data.loadError ?? 'Settings unavailable'} />
      {:else}
        <div class="grid gap-4 xl:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
          <article class="app-panel p-5">
            <p class="eyebrow">Mess config</p>
            <h2 class="section-title mt-1">Pengaturan operasional</h2>

            <form
              method="POST"
              action="?/updateSettings"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('updateSettings')}
            >
              <label>
                <span class="field-label">Mess name</span>
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
                  <span class="field-label">Wifi price</span>
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
                  <span class="field-label">Wifi deadline day</span>
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
                <span class="field-label">Bank account name</span>
                <input
                  class="input-field"
                  type="text"
                  name="bank_account_name"
                  value={valueFor('bank_account_name')}
                  required
                />
              </label>

              <label>
                <span class="field-label">Bank account number</span>
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
                {pendingAction === 'updateSettings' ? 'Saving...' : 'Save settings'}
              </button>
            </form>
          </article>

          <section class="space-y-4">
            <article class="stat-card bg-slate-950 text-white">
              <p class="helper-label text-slate-300">Current transfer target</p>
              <p class="mt-2 text-2xl font-semibold tracking-[-0.03em]">{data.settings.mess_name}</p>
              <p class="mt-2 text-sm text-slate-300">
                Wifi {formatCurrency(data.settings.wifi_price)} per orang, deadline tanggal {data.settings.wifi_deadline_day}.
              </p>
              <p class="mt-4 text-sm text-slate-200">
                {data.settings.bank_account_name} • {data.settings.bank_account_number}
              </p>
            </article>

            <article class="app-panel p-5">
              <p class="eyebrow">System status</p>
              <h2 class="section-title mt-1">Runtime ringkas</h2>

              {#if data.systemStatus}
                <div class="mt-4 space-y-3">
                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">App status</p>
                    <p class="mt-2 text-sm font-semibold text-ink">{data.systemStatus.status}</p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Database</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {data.systemStatus.database_status} ({data.systemStatus.database_reachable ? 'reachable' : 'unreachable'})
                    </p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Server time</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {formatDateTime(data.systemStatus.server_time)}
                    </p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">App version</p>
                    <p class="mt-2 text-sm font-semibold text-ink">{data.systemStatus.app_version}</p>
                  </div>

                  <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <p class="helper-label">Uptime</p>
                    <p class="mt-2 text-sm font-semibold text-ink">
                      {formatUptime(data.systemStatus.uptime_seconds)}
                    </p>
                  </div>
                </div>
              {:else}
                <StatePanel tone="empty" title="Empty" message="System status belum tersedia dari backend." />
              {/if}
            </article>
          </section>
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
