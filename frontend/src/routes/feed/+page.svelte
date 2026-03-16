<script lang="ts">
  import { enhance } from '$app/forms';
  import { invalidateAll } from '$app/navigation';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import PullToRefresh from '$lib/components/PullToRefresh.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ActivityFeedItem, ActivityType } from '$lib/api/types';
  import { queueCreateActivity, queueCreateComment } from '$lib/pwa/feed-sync';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  let pendingAction: string | null = null;
  let offlineQueueMessage: string | null = null;

  const activityLabels: Record<ActivityType, string> = {
    contribution: 'Kontribusi',
    food: 'Makanan',
    rice: 'Nasi',
    announcement: 'Pengumuman',
    other: 'Lainnya'
  };

  function formatDate(value: string | null, withTime = true) {
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

  function activityBadgeClass(type: ActivityType) {
    if (type === 'contribution') {
      return 'badge-success';
    }

    if (type === 'food') {
      return 'badge-warning';
    }

    if (type === 'rice') {
      return 'badge-info';
    }

    if (type === 'announcement') {
      return 'badge-strong';
    }

    return 'badge-muted';
  }

  function createActivityValue(field: 'type' | 'title' | 'content' | 'points') {
    if (form?.action === 'createActivity') {
      const values = form.values as Record<string, string> | undefined;
      const value = values?.[field];

      if (typeof value === 'string') {
        return value;
      }
    }

    return data.defaults[field];
  }

  function reactionState(item: ActivityFeedItem) {
    return item.reactions.find((reaction) => reaction.reaction_type === 'like') ?? null;
  }

  function hasClaimed(item: ActivityFeedItem) {
    return item.claims.some((claim) => claim.user_id === data.user?.id);
  }

  function hasRiceResponse(item: ActivityFeedItem) {
    return item.rice_responses.some((response) => response.user_id === data.user?.id);
  }

  function joinNames(values: { user_name: string }[]) {
    if (values.length === 0) {
      return '-';
    }

    return values.map((value) => value.user_name).join(', ');
  }

  function isExpired(item: ActivityFeedItem) {
    return !!item.activity.expires_at && new Date(item.activity.expires_at).getTime() < Date.now();
  }

  function enhanceCreateActivity(): SubmitFunction {
    return ({ formData, formElement, cancel }) => {
      pendingAction = 'createActivity';
      offlineQueueMessage = null;

      if (!navigator.onLine) {
        cancel();
        pendingAction = null;

        const type = String(formData.get('type') ?? '').trim() as ActivityType;
        const title = String(formData.get('title') ?? '').trim();
        const content = String(formData.get('content') ?? '').trim();
        const pointsValue = Number(String(formData.get('points') ?? '0').trim());

        if (type === 'contribution' && (!Number.isInteger(pointsValue) || pointsValue <= 0)) {
          offlineQueueMessage = 'Aktivitas kontribusi membutuhkan poin lebih dari 0 sebelum dimasukkan ke outbox.';
          return;
        }

        void queueCreateActivity({
          type,
          title,
          content,
          ...(type === 'contribution' ? { points: pointsValue } : {})
        })
          .then(() => {
            formElement.reset();
            offlineQueueMessage = 'Aktivitas disimpan di outbox dan akan dikirim otomatis saat online.';
          })
          .catch((error) => {
            offlineQueueMessage =
              error instanceof Error ? error.message : 'Gagal menyimpan aktivitas ke outbox.';
          });

        return;
      }

      return async ({ update }) => {
        await update();
        pendingAction = null;
      };
    };
  }

  function enhanceComment(activityID: string): SubmitFunction {
    return ({ formData, formElement, cancel }) => {
      pendingAction = `comment-${activityID}`;
      offlineQueueMessage = null;

      if (!navigator.onLine) {
        cancel();
        pendingAction = null;

        void queueCreateComment({
          activityID,
          comment: String(formData.get('comment') ?? '').trim()
        })
          .then(() => {
            formElement.reset();
            offlineQueueMessage = 'Komentar masuk ke outbox dan akan dikirim otomatis saat online.';
          })
          .catch((error) => {
            offlineQueueMessage =
              error instanceof Error ? error.message : 'Gagal menyimpan komentar ke outbox.';
          });

        return;
      }

      return async ({ update }) => {
        await update();
        pendingAction = null;
      };
    };
  }

  async function refreshPage() {
    await invalidateAll();
  }
</script>

<PullToRefresh onRefresh={refreshPage}>
<div class="space-y-4">
  <PageCard
    eyebrow="Feed"
    icon="lucide:newspaper"
    title="Feed"
    description="Aktivitas mess, klaim makanan, rencana nasi, komentar, dan reaksi dalam satu alur yang lebih hidup."
  >
    {#if $navigating?.to?.url.pathname === '/feed' || pendingAction}
      <PageSkeleton statCards={1} rows={4} />
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

    {#if offlineQueueMessage}
      <FeedbackBanner tone="info" title="Outbox" message={offlineQueueMessage} />
    {/if}

    {#if data.loadError}
      <StatePanel tone="error" title="Belum bisa memuat feed" message={data.loadError} />
    {:else}
      <article class="app-panel p-5">
        <p class="eyebrow">Aktivitas Baru</p>
        <h2 class="section-title mt-1">Tambah aktivitas baru</h2>
        <p class="section-subtitle mt-2">
          Kontribusi akan masuk ke leaderboard. Post makanan dan nasi otomatis punya masa aktif sementara.
        </p>

        <form method="POST" action="?/createActivity" class="mt-4 space-y-4" use:enhance={enhanceCreateActivity()}>
          <label>
            <span class="field-label">Jenis aktivitas</span>
            <select class="input-field" name="type">
              <option value="food" selected={createActivityValue('type') === 'food'}>Makanan</option>
              <option value="rice" selected={createActivityValue('type') === 'rice'}>Nasi</option>
              <option value="contribution" selected={createActivityValue('type') === 'contribution'}>
                Kontribusi
              </option>
              <option value="announcement" selected={createActivityValue('type') === 'announcement'}>
                Pengumuman
              </option>
              <option value="other" selected={createActivityValue('type') === 'other'}>Lainnya</option>
            </select>
          </label>

          <label>
            <span class="field-label">Judul</span>
            <input
              class="input-field"
              type="text"
              name="title"
              required
              placeholder="Contoh: Ada ayam goreng di dapur"
              value={createActivityValue('title')}
            />
          </label>

          <label>
            <span class="field-label">Isi</span>
            <textarea
              class="input-field min-h-[120px]"
              name="content"
              required
              placeholder="Tulis detail singkat supaya mudah ditindaklanjuti."
            >{createActivityValue('content')}</textarea>
          </label>

          <label>
            <span class="field-label">Poin</span>
            <input
              class="input-field"
              type="number"
              name="points"
              min="1"
              step="1"
              value={createActivityValue('points')}
            />
            <p class="mt-2 text-xs text-slate-500">Dipakai saat jenis aktivitas = kontribusi. Tipe lain akan mengabaikan nilai ini.</p>
          </label>

          <button type="submit" class="btn-primary w-full px-4 py-3" disabled={pendingAction === 'createActivity'}>
            {pendingAction === 'createActivity' ? 'Menerbitkan...' : 'Terbitkan aktivitas'}
          </button>
        </form>
      </article>

      {#if data.activities.length === 0}
        <StatePanel tone="empty" title="Belum ada aktivitas" message="Belum ada aktivitas. Terbitkan item pertama dari form di atas." />
      {:else}
        <div class="mt-4 space-y-4">
          {#each data.activities as item}
            <article class="app-panel p-5">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span class={activityBadgeClass(item.activity.type)}>
                      {activityLabels[item.activity.type]}
                    </span>
                    {#if item.activity.type === 'contribution'}
                      <span class="badge bg-emerald-50 text-emerald-700">{item.activity.points} poin</span>
                    {/if}
                    {#if isExpired(item)}
                      <span class="badge-muted">Sudah berakhir</span>
                    {/if}
                  </div>

                  <h2 class="section-title mt-3">{item.activity.title}</h2>
                  <p class="mt-2 text-sm leading-6 text-slate-600">{item.activity.content}</p>
                </div>

                <div class="text-right text-xs text-slate-500">
                  <p class="font-semibold text-slate-700">{item.activity.created_by_name}</p>
                  <p class="mt-1">{formatDate(item.activity.created_at)}</p>
                  {#if item.activity.expires_at}
                    <p class="mt-1">Aktif sampai {formatDate(item.activity.expires_at)}</p>
                  {/if}
                </div>
              </div>

              {#if item.activity.type === 'food'}
                <div class="helper-box mt-4">
                  <p class="helper-label">Klaim makanan</p>
                  <p class="mt-2 text-sm text-slate-600">Sudah diambil oleh: {joinNames(item.claims)}</p>
                </div>
              {/if}

              {#if item.activity.type === 'rice'}
                <div class="helper-box mt-4">
                  <p class="helper-label">Respons nasi</p>
                  <p class="mt-2 text-sm text-slate-600">Yang ikut: {joinNames(item.rice_responses)}</p>
                </div>
              {/if}

              <div class="mt-4 flex flex-wrap gap-2">
                <form method="POST" action="?/react">
                  <input type="hidden" name="activity_id" value={item.activity.id} />
                  <button
                    type="submit"
                    class={reactionState(item)?.reacted ? 'btn-primary px-4 py-2.5' : 'btn-secondary px-4 py-2.5'}
                  >
                    {reactionState(item)?.reacted ? 'Disukai' : 'Suka'} ({reactionState(item)?.count ?? 0})
                  </button>
                </form>

                {#if item.activity.type === 'food'}
                  <form method="POST" action="?/claim">
                    <input type="hidden" name="activity_id" value={item.activity.id} />
                    <button
                      type="submit"
                      class="btn-secondary px-4 py-2.5"
                      disabled={hasClaimed(item)}
                    >
                      {hasClaimed(item) ? 'Sudah diambil' : 'Ambil makanan'}
                    </button>
                  </form>
                {/if}

                {#if item.activity.type === 'rice'}
                  <form method="POST" action="?/riceResponse">
                    <input type="hidden" name="activity_id" value={item.activity.id} />
                    <button
                      type="submit"
                      class="btn-secondary px-4 py-2.5"
                      disabled={hasRiceResponse(item)}
                    >
                      {hasRiceResponse(item) ? 'Sudah respon' : 'Ikut makan nasi'}
                    </button>
                  </form>
                {/if}
              </div>

              <div class="mt-4 space-y-3">
                <div class="flex items-center justify-between">
                  <p class="helper-label">Komentar</p>
                  <span class="badge-muted">{item.comments.length}</span>
                </div>

                {#if item.comments.length === 0}
                  <StatePanel tone="empty" title="Belum ada komentar" message="Belum ada komentar untuk aktivitas ini." />
                {:else}
                  <div class="space-y-2">
                    {#each item.comments as comment}
                      <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                        <div class="flex items-center justify-between gap-3">
                          <p class="text-sm font-semibold text-ink">{comment.user_name}</p>
                          <p class="text-xs text-slate-500">{formatDate(comment.created_at)}</p>
                        </div>
                        <p class="mt-2 text-sm leading-6 text-slate-600">{comment.comment}</p>
                      </div>
                    {/each}
                  </div>
                {/if}

                <form method="POST" action="?/comment" class="space-y-3" use:enhance={enhanceComment(item.activity.id)}>
                  <input type="hidden" name="activity_id" value={item.activity.id} />
                  <textarea
                    class="input-field min-h-[96px]"
                    name="comment"
                    placeholder="Tambahkan komentar singkat..."
                    required
                  ></textarea>
                  <button
                    type="submit"
                    class="btn-secondary w-full px-4 py-3"
                    disabled={pendingAction === `comment-${item.activity.id}`}
                  >
                    {pendingAction === `comment-${item.activity.id}` ? 'Menyimpan...' : 'Kirim komentar'}
                  </button>
                </form>
              </div>
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
</PullToRefresh>
