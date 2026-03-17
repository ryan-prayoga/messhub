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

  function profileValue(field: 'name' | 'phone' | 'avatar_url') {
    if (form?.action === 'updateProfile') {
      const values = form.values as Record<string, string> | undefined;
      const value = values?.[field];
      if (typeof value === 'string') {
        return value;
      }
    }

    if (!data.profile) {
      return '';
    }

    if (field === 'name') {
      return data.profile.name;
    }

    if (field === 'phone') {
      return data.profile.phone ?? '';
    }

    return data.profile.avatar_url ?? '';
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

  const roleLabels: Record<string, string> = {
    admin: 'Admin',
    treasurer: 'Bendahara',
    member: 'Anggota'
  };

  $: avatarPreview = profileValue('avatar_url');
</script>

<div class="space-y-5 lg:space-y-6">
  <PageCard
    eyebrow="Profile"
    icon="lucide:user-round"
    title="Profil"
    description="Kelola identitas akun pribadi, kontak, username, avatar, dan password."
  >
    {#if $navigating?.to?.url.pathname === '/profile' || pendingAction}
      <PageSkeleton statCards={2} rows={2} />
    {/if}

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

    {#if data.loadError || !data.profile}
      <StatePanel tone="error" title="Gagal memuat" message={data.loadError ?? 'Profil belum tersedia.'} />
    {:else}
      <div class="grid gap-5 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)]">
        <section class="space-y-5">
          <article class="stat-card bg-slate-950 text-white">
            <p class="helper-label text-slate-300">Akun aktif</p>
            <div class="mt-4 flex items-center gap-4">
              {#if avatarPreview}
                <img
                  src={avatarPreview}
                  alt={`${data.profile.name} avatar`}
                  class="h-20 w-20 rounded-[24px] border border-white/20 object-cover"
                />
              {:else}
                <div class="flex h-20 w-20 items-center justify-center rounded-[24px] border border-white/15 bg-white/10 text-2xl font-semibold uppercase">
                  {data.profile.name.slice(0, 1)}
                </div>
              {/if}

              <div class="min-w-0">
                <h2 class="text-xl font-semibold tracking-[-0.03em]">{data.profile.name}</h2>
                <p class="mt-1 break-all text-sm text-slate-300">{data.profile.email}</p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span class="badge bg-white/10 text-white">
                    {roleLabels[data.profile.role] ?? data.profile.role}
                  </span>
                  <span class={`badge ${data.profile.is_active ? 'bg-[#7D8A74]/20 text-[#F8F6F2]' : 'bg-[#B96F62]/20 text-[#F8F6F2]'}`}>
                    {data.profile.is_active ? 'Aktif' : 'Nonaktif'}
                  </span>
                </div>
              </div>
            </div>
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Info akun</p>
            <h2 class="section-title mt-1">Data saat ini</h2>

            <div class="mt-4 grid gap-3 sm:grid-cols-2">
              <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Nomor HP</p>
                <p class="mt-2 text-sm font-medium text-ink">{data.profile.phone ?? '-'}</p>
              </div>

              <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Mulai tinggal</p>
                <p class="mt-2 text-sm font-medium text-ink">{formatDate(data.profile.joined_at)}</p>
              </div>

              <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 sm:col-span-2">
                <p class="helper-label">Username login</p>
                <p class="mt-2 text-sm font-medium text-ink">@{data.profile.username}</p>
              </div>
            </div>
          </article>
        </section>

        <section class="space-y-5">
          <article class="app-panel p-5">
            <p class="eyebrow">Perbarui profil</p>
            <h2 class="section-title mt-1">Ubah data pribadi</h2>

            <form
              method="POST"
              action="?/updateProfile"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('updateProfile')}
            >
              <label>
                <span class="field-label">Nama</span>
                <input
                  class="input-field"
                  type="text"
                  name="name"
                  value={profileValue('name')}
                  required
                />
              </label>

              <label>
                <span class="field-label">Nomor HP</span>
                <input
                  class="input-field"
                  type="tel"
                  name="phone"
                  value={profileValue('phone')}
                  placeholder="0812..."
                />
              </label>

              <label>
                <span class="field-label">URL avatar</span>
                <input
                  class="input-field"
                  type="url"
                  name="avatar_url"
                  value={profileValue('avatar_url')}
                  placeholder="https://..."
                />
              </label>

              <button
                type="submit"
                class="btn-primary w-full px-4 py-3"
                disabled={pendingAction === 'updateProfile'}
              >
                {pendingAction === 'updateProfile' ? 'Menyimpan...' : 'Simpan profil'}
              </button>
            </form>
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Keamanan</p>
            <h2 class="section-title mt-1">Ganti password</h2>
            <p class="section-subtitle mt-2">Password baru harus minimal 8 karakter.</p>

            <form
              method="POST"
              action="?/changePassword"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('changePassword')}
            >
              <label>
                <span class="field-label">Password saat ini</span>
                <input class="input-field" type="password" name="current_password" required />
              </label>

              <label>
                <span class="field-label">Password baru</span>
                <input class="input-field" type="password" name="new_password" minlength="8" required />
              </label>

              <label>
                <span class="field-label">Konfirmasi password baru</span>
                <input class="input-field" type="password" name="confirm_password" minlength="8" required />
              </label>

              <button
                type="submit"
                class="btn-secondary w-full px-4 py-3"
                disabled={pendingAction === 'changePassword'}
              >
                {pendingAction === 'changePassword' ? 'Memperbarui...' : 'Ganti password'}
              </button>
            </form>
          </article>
        </section>
      </div>
    {/if}
  </PageCard>
</div>
