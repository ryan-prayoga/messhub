<script lang="ts">
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import PageCard from '$lib/components/PageCard.svelte';
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

  $: avatarPreview = profileValue('avatar_url');
</script>

<div class="space-y-4">
  <PageCard
    title="Profile"
    description="Kelola identitas akun pribadi, kontak, avatar, dan password."
  >
    {#if $navigating?.to?.url.pathname === '/profile' || pendingAction}
      <div class="helper-box mb-4">
        <p class="helper-label">Loading</p>
        <p class="mt-2 text-sm text-slate-600">
          {pendingAction ? 'Memproses perubahan profil...' : 'Memuat ulang data profil...'}
        </p>
      </div>
    {/if}

    {#if form?.message}
      <div class="helper-box-brand mb-4">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{form.message}</p>
      </div>
    {:else if form?.success}
      <div class="helper-box mb-4 border-emerald-200 bg-emerald-50/80">
        <p class="helper-label text-emerald-700">Success</p>
        <p class="mt-2 text-sm leading-6 text-emerald-800">{form.success}</p>
      </div>
    {/if}

    {#if data.loadError || !data.profile}
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{data.loadError ?? 'Profile unavailable'}</p>
      </div>
    {:else}
      <div class="grid gap-4 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)]">
        <section class="space-y-4">
          <article class="stat-card bg-slate-950 text-white">
            <p class="helper-label text-slate-300">Current account</p>
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
                  <span class="badge bg-white/10 text-white">{data.profile.role}</span>
                  <span class={`badge ${data.profile.is_active ? 'bg-emerald-500/20 text-emerald-100' : 'bg-rose-500/20 text-rose-100'}`}>
                    {data.profile.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Account facts</p>
            <h2 class="section-title mt-1">Info saat ini</h2>

            <div class="mt-4 grid gap-3 sm:grid-cols-2">
              <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Phone</p>
                <p class="mt-2 text-sm font-medium text-ink">{data.profile.phone ?? '-'}</p>
              </div>

              <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Joined</p>
                <p class="mt-2 text-sm font-medium text-ink">{formatDate(data.profile.joined_at)}</p>
              </div>
            </div>
          </article>
        </section>

        <section class="space-y-4">
          <article class="app-panel p-5">
            <p class="eyebrow">Edit profile</p>
            <h2 class="section-title mt-1">Perbarui data pribadi</h2>

            <form
              method="POST"
              action="?/updateProfile"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('updateProfile')}
            >
              <label>
                <span class="field-label">Name</span>
                <input
                  class="input-field"
                  type="text"
                  name="name"
                  value={profileValue('name')}
                  required
                />
              </label>

              <label>
                <span class="field-label">Phone</span>
                <input
                  class="input-field"
                  type="tel"
                  name="phone"
                  value={profileValue('phone')}
                  placeholder="0812..."
                />
              </label>

              <label>
                <span class="field-label">Avatar URL</span>
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
                {pendingAction === 'updateProfile' ? 'Saving...' : 'Save profile'}
              </button>
            </form>
          </article>

          <article class="app-panel p-5">
            <p class="eyebrow">Security</p>
            <h2 class="section-title mt-1">Ganti password</h2>
            <p class="section-subtitle mt-2">
              Password baru harus minimal 8 karakter.
            </p>

            <form
              method="POST"
              action="?/changePassword"
              class="mt-4 space-y-4"
              use:enhance={enhanceWithAction('changePassword')}
            >
              <label>
                <span class="field-label">Current password</span>
                <input class="input-field" type="password" name="current_password" required />
              </label>

              <label>
                <span class="field-label">New password</span>
                <input class="input-field" type="password" name="new_password" minlength="8" required />
              </label>

              <label>
                <span class="field-label">Confirm new password</span>
                <input class="input-field" type="password" name="confirm_password" minlength="8" required />
              </label>

              <button
                type="submit"
                class="btn-secondary w-full px-4 py-3"
                disabled={pendingAction === 'changePassword'}
              >
                {pendingAction === 'changePassword' ? 'Updating...' : 'Change password'}
              </button>
            </form>
          </article>
        </section>
      </div>
    {/if}
  </PageCard>
</div>
