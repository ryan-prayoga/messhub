<script lang="ts">
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { Notification } from '$lib/api/types';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;

  function formatDate(value: string) {
    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(new Date(value));
  }

  function notificationHref(item: Notification) {
    if (item.type === 'activity_created' || item.type === 'comment_created') {
      return '/feed';
    }

    if (item.type === 'wifi_bill_created' || item.type === 'wifi_payment_verified') {
      return '/wifi';
    }

    return '/dashboard';
  }

  function cardClass(item: Notification) {
    if (item.is_read) {
      return 'rounded-3xl border border-line bg-white/80 px-4 py-4';
    }

    return 'rounded-3xl border border-line bg-panel px-4 py-4 shadow-sm';
  }
</script>

<div class="space-y-4">
  <PageCard
    eyebrow="Inbox"
    icon="lucide:bell-ring"
    title="Notifications"
    description="Notifikasi dalam aplikasi untuk aktivitas baru, komentar, dan status wifi."
  >
    {#if form?.message}
      <FeedbackBanner tone="error" title="Belum bisa memproses" message={form.message} />
    {:else if form?.success}
      <FeedbackBanner tone="success" title="Berhasil" message={form.success} />
    {/if}

    {#if data.loadError}
      <StatePanel tone="error" title="Belum bisa memuat notifikasi" message={data.loadError} />
    {:else}
      <div class="mb-4 flex items-center justify-between rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3">
        <div>
          <p class="helper-label">Belum dibaca</p>
          <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
            {data.notificationSummary.unread_count}
          </p>
        </div>

        <form method="POST" action="?/markAllRead">
          <button
            type="submit"
            class="btn-secondary px-4 py-2.5"
            disabled={data.notificationSummary.unread_count === 0}
          >
            Tandai semua
          </button>
        </form>
      </div>

      {#if data.notificationSummary.items.length === 0}
        <StatePanel
          tone="empty"
          title="Inbox masih kosong"
          message="Belum ada notifikasi baru. Update penting dari wifi, feed, dan aktivitas mess akan muncul di sini."
          actionHref="/feed"
          actionLabel="Buka feed"
        />
      {:else}
        <div class="space-y-3">
          {#each data.notificationSummary.items as item}
            <article class={cardClass(item)}>
              <div class="flex items-start justify-between gap-4">
                <a href={notificationHref(item)} class="min-w-0 flex-1">
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="text-sm font-semibold text-ink">{item.title}</p>
                    {#if !item.is_read}
                      <span class="badge-brand">Baru</span>
                    {/if}
                  </div>
                  <p class="mt-2 text-sm leading-6 text-slate-600">{item.message}</p>
                  <p class="mt-2 text-xs text-slate-500">{formatDate(item.created_at)}</p>
                </a>

                <form method="POST" action="?/markOneRead">
                  <input type="hidden" name="notification_id" value={item.id} />
                  <button type="submit" class="btn-secondary px-3 py-2 text-xs" disabled={item.is_read}>
                    {item.is_read ? 'Sudah dibaca' : 'Tandai dibaca'}
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
