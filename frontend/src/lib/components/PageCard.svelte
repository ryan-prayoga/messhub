<script lang="ts">
  export let title: string;
  export let description: string;
</script>

<section class="app-panel p-5 sm:p-6">
  <div class="mb-4 space-y-1.5">
    <h2 class="section-title">{title}</h2>
    <p class="section-subtitle">{description}</p>
  </div>
  <slot />
</section>
