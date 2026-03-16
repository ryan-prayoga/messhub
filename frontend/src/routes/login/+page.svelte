<script lang="ts">
  import { enhance } from '$app/forms';

  export let form: { message?: string } | undefined;
</script>

<div class="mx-auto flex min-h-screen max-w-md items-center px-5 py-10">
  <section class="w-full rounded-[32px] border border-line bg-panel p-6 shadow-shell">
    <div class="mb-8">
      <p class="text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">MessHub</p>
      <h1 class="mt-2 text-3xl font-bold text-ink">Login</h1>
      <p class="mt-2 text-sm leading-6 text-slate-600">
        Masuk dengan akun admin, treasurer, atau member yang dibuat dari backend.
      </p>
    </div>

    <form method="POST" use:enhance class="space-y-4">
      <label class="block">
        <span class="mb-2 block text-sm font-semibold text-slate-700">Email</span>
        <input
          name="email"
          type="email"
          placeholder="admin@messhub.local"
          class="w-full rounded-2xl border border-line bg-slate-50 px-4 py-3 outline-none transition focus:border-ink"
        />
      </label>

      <label class="block">
        <span class="mb-2 block text-sm font-semibold text-slate-700">Password</span>
        <input
          name="password"
          type="password"
          placeholder="••••••••"
          class="w-full rounded-2xl border border-line bg-slate-50 px-4 py-3 outline-none transition focus:border-ink"
        />
      </label>

      {#if form?.message}
        <div class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {form.message}
        </div>
      {/if}

      <button
        type="submit"
        class="w-full rounded-2xl bg-ink px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-800"
      >
        Sign in
      </button>
    </form>

    <div class="mt-6 rounded-2xl bg-orange-50 p-4 text-sm text-orange-900">
      <p class="font-semibold">Seed default</p>
      <p class="mt-1">Email `admin@messhub.local` dengan password dari `.env` atau `.env.example`.</p>
    </div>
  </section>
</div>
