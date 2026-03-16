<script lang="ts">
  import AppIcon from '$lib/components/AppIcon.svelte';

  export let tone: 'success' | 'error' | 'info' | 'warning' = 'info';
  export let title = '';
  export let message = '';

  const toneClasses = {
    success: 'feedback-banner feedback-banner-success',
    error: 'feedback-banner feedback-banner-error',
    info: 'feedback-banner feedback-banner-info',
    warning: 'feedback-banner feedback-banner-warning'
  } as const;

  const toneIcons = {
    success: 'lucide:badge-check',
    error: 'lucide:triangle-alert',
    info: 'lucide:sparkles',
    warning: 'lucide:shield-alert'
  } as const;
</script>

<div class={toneClasses[tone]}>
  <div class="feedback-banner-icon">
    <AppIcon icon={toneIcons[tone]} className="h-5 w-5" />
  </div>

  <div class="min-w-0">
    {#if title}
      <p class="feedback-banner-title">{title}</p>
    {/if}
    <p class="feedback-banner-message">{message}</p>
  </div>
</div>
