<script lang="ts">
  export let title: string;
  export let description: string;
</script>

<section class="rounded-[28px] border border-line bg-panel p-5 shadow-sm">
  <div class="mb-4">
    <h2 class="text-xl font-bold text-ink">{title}</h2>
    <p class="mt-1 text-sm leading-6 text-slate-600">{description}</p>
  </div>
  <slot />
</section>
