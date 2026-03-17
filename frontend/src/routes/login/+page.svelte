<script lang="ts">
  import { enhance } from '$app/forms';
  import type { SubmitFunction } from '@sveltejs/kit';
  import AppIcon from '$lib/components/AppIcon.svelte';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import Spinner from '$lib/components/Spinner.svelte';

  type LoginFormState = {
    message?: string;
    values?: {
      identifier?: string;
    };
  };

  export let form: LoginFormState | undefined;
  export let data: { notice?: string | null };

  const loginHighlights = [
    {
      icon: 'lucide:house',
      title: 'Warm, mobile-first workspace',
      description: 'Terbuka nyaman di browser desktop, mobile, dan mode PWA.'
    },
    {
      icon: 'lucide:shield-check',
      title: 'Role & auth tetap aman',
      description: 'Akses tetap mengikuti role existing tanpa mengubah alur izin yang sudah berjalan.'
    },
    {
      icon: 'lucide:sparkles',
      title: 'Masuk lebih fleksibel',
      description: 'Gunakan email atau username yang sudah disiapkan admin mess.'
    }
  ];

  let isSubmitting = false;
  let isRedirecting = false;

  const enhanceLogin: SubmitFunction = () => {
    isSubmitting = true;
    isRedirecting = false;

    return async ({ result, update }) => {
      if (result.type === 'redirect') {
        isRedirecting = true;
      }

      await update();

      if (result.type !== 'redirect') {
        isSubmitting = false;
      }
    };
  };
</script>

<div class="app-shell">
  <div class="mx-auto flex min-h-screen w-full max-w-6xl items-center px-4 py-10 sm:px-6 lg:px-8">
    <div class="grid w-full gap-5 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
      <section class="app-panel hidden lg:block">
        <p class="eyebrow">Mess Traspac Menyala</p>
        <h1 class="mt-4 text-5xl font-semibold tracking-[-0.06em] text-ink">
          MessHub untuk operasional harian yang terasa lebih rapi dan hidup.
        </h1>
        <p class="mt-5 max-w-2xl text-base leading-8 text-muted">
          Satu tempat untuk memantau kas, wifi, feed, anggota, dan pengaturan mess dengan tampilan yang hangat,
          ringan, serta tetap cepat dipakai setiap hari.
        </p>

        <div class="mt-8 grid gap-4">
          {#each loginHighlights as item}
            <article class="stat-card bg-white/70">
              <div class="flex items-start gap-4">
                <div class="nav-link-icon mt-1">
                  <AppIcon icon={item.icon} className="h-5 w-5" />
                </div>

                <div>
                  <h2 class="text-base font-semibold text-ink">{item.title}</h2>
                  <p class="mt-2 text-sm leading-6 text-muted">{item.description}</p>
                </div>
              </div>
            </article>
          {/each}
        </div>
      </section>

      <section class="app-panel w-full max-w-[34rem] justify-self-end">
        <div class="flex items-center justify-between gap-3">
          <span class="badge-brand">MessHub</span>
          <span class="badge-muted">PWA Internal</span>
        </div>

        <div class="mt-6">
          <p class="eyebrow">Selamat datang kembali</p>
          <h1 class="mt-3 text-3xl font-semibold tracking-[-0.05em] text-ink">Masuk ke akun Anda</h1>
          <p class="mt-3 text-sm leading-7 text-muted">
            Gunakan email atau username yang sudah dibuatkan admin untuk membuka dashboard dan fitur operasional mess.
          </p>
        </div>

        <form
          method="POST"
          use:enhance={enhanceLogin}
          class="mt-8 space-y-4"
          aria-busy={isSubmitting}
        >
          <label class="block">
            <span class="field-label">Email atau Username</span>
            <div class="relative">
              <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-dusty">
                <AppIcon icon="lucide:at-sign" className="h-4 w-4" />
              </span>
              <input
                name="identifier"
                type="text"
                autocomplete="username"
                placeholder="ryan atau admin@messhub.local"
                class="input-field pl-11"
                value={form?.values?.identifier ?? ''}
                disabled={isSubmitting}
              />
            </div>
          </label>

          <label class="block">
            <span class="field-label">Password</span>
            <div class="relative">
              <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-dusty">
                <AppIcon icon="lucide:key-round" className="h-4 w-4" />
              </span>
              <input
                name="password"
                type="password"
                autocomplete="current-password"
                placeholder="Masukkan password"
                class="input-field pl-11"
                disabled={isSubmitting}
              />
            </div>
          </label>

          {#if data.notice}
            <FeedbackBanner tone="info" title="Sesi diperbarui" message={data.notice} />
          {/if}

          {#if form?.message}
            <FeedbackBanner tone="error" title="Belum bisa masuk" message={form.message} />
          {/if}

          {#if isRedirecting}
            <FeedbackBanner
              tone="success"
              title="Masuk berhasil"
              message="Membuka dashboard MessHub untuk Anda..."
            />
          {/if}

          <button type="submit" class="btn-primary w-full" disabled={isSubmitting}>
            {#if isSubmitting}
              <Spinner className="h-4 w-4" />
              <span>{isRedirecting ? 'Membuka dashboard...' : 'Memeriksa akun...'}</span>
            {:else}
              <AppIcon icon="lucide:arrow-right" className="h-4 w-4" />
              <span>Masuk</span>
            {/if}
          </button>
        </form>

        <div class="mt-6 rounded-[24px] border border-line bg-white/60 p-4">
          <p class="text-sm leading-6 text-muted">
            Belum punya akses? Hubungi admin mess untuk dibuatkan akun aktif beserta role yang sesuai.
          </p>
        </div>
      </section>
    </div>
  </div>
</div>
