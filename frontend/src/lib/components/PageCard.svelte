<script lang="ts">
  import AppIcon from '$lib/components/AppIcon.svelte';

  export let title: string;
  export let description = '';
  export let eyebrow = '';
  export let icon: string | null = null;
</script>

<section class="app-panel">
  <div class="page-card-header">
    <div class="flex min-w-0 items-start gap-4">
      {#if icon}
        <div class="nav-link-icon mt-0.5">
          <AppIcon {icon} className="h-5 w-5" />
        </div>
      {/if}

      <div class="min-w-0 space-y-1.5">
        {#if eyebrow}
          <p class="eyebrow">{eyebrow}</p>
        {/if}

        <h2 class="section-title">{title}</h2>

        {#if description}
          <p class="section-subtitle">{description}</p>
        {/if}
      </div>
    </div>

    <slot name="actions" />
  </div>

  <div class="page-card-body">
    <slot />
  </div>
</section>
