<script lang="ts">
  import { navigating } from '$app/stores';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ActivityFeedItem, ActivityType } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  const activityLabels: Record<ActivityType, string> = {
    contribution: 'Contribution',
    food: 'Food',
    rice: 'Rice',
    announcement: 'Announcement',
    other: 'Other'
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
      return 'badge bg-emerald-100 text-emerald-700';
    }

    if (type === 'food') {
      return 'badge bg-amber-100 text-amber-700';
    }

    if (type === 'rice') {
      return 'badge bg-sky-100 text-sky-700';
    }

    if (type === 'announcement') {
      return 'badge bg-slate-950 text-white';
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
</script>

<div class="space-y-4">
  <PageCard
    title="Feed"
    description="Aktivitas mess, klaim makanan, rencana nasi, komentar, dan reaction untuk Step 5."
  >
    {#if $navigating?.to?.url.pathname === '/feed'}
      <StatePanel tone="loading" title="Loading" message="Memuat ulang feed, komentar, dan reaction terbaru..." />
    {/if}

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

    {#if data.loadError}
      <StatePanel tone="error" title="Error" message={data.loadError} />
    {:else}
      <article class="app-panel p-5">
        <p class="eyebrow">Post Activity</p>
        <h2 class="section-title mt-1">Tambah aktivitas baru</h2>
        <p class="section-subtitle mt-2">
          Contribution akan masuk ke leaderboard. Food dan rice otomatis punya masa aktif sementara.
        </p>

        <form method="POST" action="?/createActivity" class="mt-4 space-y-4">
          <label>
            <span class="field-label">Type</span>
            <select class="input-field" name="type">
              <option value="food" selected={createActivityValue('type') === 'food'}>Food</option>
              <option value="rice" selected={createActivityValue('type') === 'rice'}>Rice</option>
              <option value="contribution" selected={createActivityValue('type') === 'contribution'}>
                Contribution
              </option>
              <option value="announcement" selected={createActivityValue('type') === 'announcement'}>
                Announcement
              </option>
              <option value="other" selected={createActivityValue('type') === 'other'}>Other</option>
            </select>
          </label>

          <label>
            <span class="field-label">Title</span>
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
            <span class="field-label">Content</span>
            <textarea
              class="input-field min-h-[120px]"
              name="content"
              required
              placeholder="Tulis detail singkat supaya mudah ditindaklanjuti."
            >{createActivityValue('content')}</textarea>
          </label>

          <label>
            <span class="field-label">Points</span>
            <input
              class="input-field"
              type="number"
              name="points"
              min="1"
              step="1"
              value={createActivityValue('points')}
            />
            <p class="mt-2 text-xs text-slate-500">Dipakai saat type = contribution. Tipe lain akan mengabaikan nilai ini.</p>
          </label>

          <button type="submit" class="btn-primary w-full px-4 py-3">Post activity</button>
        </form>
      </article>

      {#if data.activities.length === 0}
        <StatePanel tone="empty" title="Empty" message="Belum ada aktivitas. Posting item pertama dari form di atas." />
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
                      <span class="badge bg-emerald-50 text-emerald-700">{item.activity.points} pts</span>
                    {/if}
                    {#if isExpired(item)}
                      <span class="badge-muted">Expired</span>
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
                  <p class="helper-label">Food claims</p>
                  <p class="mt-2 text-sm text-slate-600">Claimed by: {joinNames(item.claims)}</p>
                </div>
              {/if}

              {#if item.activity.type === 'rice'}
                <div class="helper-box mt-4">
                  <p class="helper-label">Rice responses</p>
                  <p class="mt-2 text-sm text-slate-600">Responded by: {joinNames(item.rice_responses)}</p>
                </div>
              {/if}

              <div class="mt-4 flex flex-wrap gap-2">
                <form method="POST" action="?/react">
                  <input type="hidden" name="activity_id" value={item.activity.id} />
                  <button
                    type="submit"
                    class={reactionState(item)?.reacted ? 'btn-primary px-4 py-2.5' : 'btn-secondary px-4 py-2.5'}
                  >
                    {reactionState(item)?.reacted ? 'Liked' : 'Like'} ({reactionState(item)?.count ?? 0})
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
                      {hasClaimed(item) ? 'Sudah claim' : 'Claim food'}
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
                  <p class="helper-label">Comments</p>
                  <span class="badge-muted">{item.comments.length}</span>
                </div>

                {#if item.comments.length === 0}
                  <StatePanel tone="empty" title="Empty" message="Belum ada komentar." />
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

                <form method="POST" action="?/comment" class="space-y-3">
                  <input type="hidden" name="activity_id" value={item.activity.id} />
                  <textarea
                    class="input-field min-h-[96px]"
                    name="comment"
                    placeholder="Tambahkan komentar singkat..."
                    required
                  ></textarea>
                  <button type="submit" class="btn-secondary w-full px-4 py-3">Add comment</button>
                </form>
              </div>
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
