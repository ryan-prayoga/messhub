<script lang="ts">
  import { networkState } from '$lib/pwa/network';

  function retry() {
    window.location.href = '/dashboard';
  }
</script>

<div class="min-h-screen px-4 py-8 sm:px-6">
  <div class="mx-auto flex min-h-[80vh] w-full max-w-xl items-center">
    <section class="app-panel w-full p-6 sm:p-8">
      <p class="eyebrow">Offline fallback</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">MessHub sedang offline</h1>
      <p class="mt-3 text-sm leading-7 text-slate-600">
        Koneksi internet belum tersedia dan halaman yang diminta belum ada di cache lokal. Coba buka dashboard lagi
        saat koneksi kembali.
      </p>

      <div class="mt-5 rounded-3xl border border-slate-200 bg-slate-50 p-4 text-sm leading-6 text-slate-600">
        Status jaringan: {$networkState.online ? 'online' : 'offline'}.
      </div>

      <div class="mt-6 flex flex-col gap-3 sm:flex-row">
        <button type="button" class="btn-primary px-4 py-3" on:click={retry}>
          Coba buka dashboard
        </button>
        <a href="/feed" class="btn-secondary px-4 py-3">Buka feed terakhir</a>
      </div>
    </section>
  </div>
</div>
