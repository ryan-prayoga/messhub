<script lang="ts">
  import { enhance } from '$app/forms';
  import type { SubmitFunction } from '@sveltejs/kit';

  type LoginFormState = {
    message?: string;
    values?: {
      email?: string;
    };
  };

  export let form: LoginFormState | undefined;

  let isSubmitting = false;

  const enhanceLogin: SubmitFunction = () => {
    isSubmitting = true;

    return async ({ update }) => {
      await update();
      isSubmitting = false;
    };
  };
</script>

<div class="app-shell">
  <div class="mx-auto flex min-h-screen w-full max-w-5xl items-center justify-center px-4 py-10 sm:px-6">
    <section class="app-panel w-full max-w-md p-6 sm:p-8">
      <div class="mb-8">
        <div class="flex items-center justify-between gap-3">
          <span class="badge-brand">MessHub</span>
          <span class="badge-muted">Internal PWA</span>
        </div>

        <h1 class="mt-4 text-3xl font-semibold tracking-[-0.03em] text-ink">Masuk ke MessHub</h1>
        <p class="mt-2 text-sm leading-6 text-slate-500">
          Gunakan akun yang sudah disiapkan admin untuk mengakses operasional harian mess.
        </p>
      </div>

      <form
        method="POST"
        use:enhance={enhanceLogin}
        class="space-y-4"
        aria-busy={isSubmitting}
      >
        <label class="block">
          <span class="field-label">Email</span>
          <input
            name="email"
            type="email"
            autocomplete="username"
            placeholder="admin@messhub.local"
            class="input-field"
            value={form?.values?.email ?? ''}
            disabled={isSubmitting}
          />
        </label>

        <label class="block">
          <span class="field-label">Password</span>
          <input
            name="password"
            type="password"
            autocomplete="current-password"
            placeholder="Masukkan password"
            class="input-field"
            disabled={isSubmitting}
          />
        </label>

        {#if form?.message}
          <div class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {form.message}
          </div>
        {/if}

        <button type="submit" class="btn-primary w-full" disabled={isSubmitting}>
          {#if isSubmitting}
            Memproses...
          {:else}
            Sign in
          {/if}
        </button>

        <p class="text-center text-sm text-slate-500">
          Jika login gagal, cek kredensial seed atau koneksi frontend ke backend auth.
        </p>
      </form>

      <div class="mt-6 helper-box-brand">
        <p class="helper-label text-sky-700">Akun seed default</p>
        <p class="mt-2 text-sm font-medium text-slate-900">Email: <code>admin@messhub.local</code></p>
        <p class="mt-1 text-sm leading-6 text-slate-600">
          Password mengikuti nilai di <code>.env</code> atau <code>.env.example</code> backend.
        </p>
      </div>

      <div class="mt-4 helper-box">
        <p class="helper-label">Catatan</p>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Halaman ini dibuat ringan dan mobile-first agar enak dipakai di HP tanpa mengubah flow auth inti.
        </p>
      </div>
    </section>
  </div>
</div>
