<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { PageData } from './$types';

  export let data: PageData;

  const importCards = [
    {
      href: '/admin/import/members',
      title: 'Impor Anggota',
      description: 'Unggah CSV anggota, cek duplikasi email, lalu simpan hanya baris yang valid.',
      templateHref: '/templates/member-import-template.csv'
    },
    {
      href: '/admin/import/wallet',
      title: 'Impor Transaksi Kas',
      description: 'Unggah CSV kas lama, cek validasi nominal dan tanggal, lalu commit setelah preview.',
      templateHref: '/templates/wallet-import-template.csv'
    }
  ];
</script>

<div class="space-y-4">
  <PageCard
    title="Pusat Impor Data"
    description="Pindahkan data spreadsheet lama ke MessHub lewat alur CSV yang aman dan bisa ditinjau dulu."
  >
    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Akses ditolak"
        message="Halaman impor hanya tersedia untuk admin mess."
      />
    {:else}
      <div class="grid gap-4 md:grid-cols-2">
        {#each importCards as card}
          <article class="app-panel p-5">
            <p class="eyebrow">Admin import</p>
            <h2 class="section-title mt-1">{card.title}</h2>
            <p class="section-subtitle mt-2">{card.description}</p>

            <div class="mt-4 flex flex-wrap gap-3">
              <a href={card.href} class="btn-primary px-4 py-3">Buka halaman</a>
              <a href={card.templateHref} class="btn-secondary px-4 py-3" download>Unduh template</a>
            </div>
          </article>
        {/each}
      </div>

      <div class="mt-4 helper-box">
        <p class="helper-label">Alur aman</p>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Setiap impor akan dicatat lengkap dengan siapa yang mengunggah, kapan preview dibuat, dan berapa baris yang berhasil disimpan.
        </p>
      </div>
    {/if}
  </PageCard>
</div>
