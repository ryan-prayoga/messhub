<script lang="ts">
  export let tone: 'loading' | 'error' | 'empty' | 'forbidden' = 'empty';
  export let title = '';
  export let message: string;
  export let requestId: string | null = null;

  const classes = {
    loading: 'helper-box',
    error: 'helper-box-brand',
    empty: 'empty-state',
    forbidden: 'empty-state'
  } as const;

  const labelClasses = {
    loading: 'helper-label',
    error: 'helper-label text-sky-700',
    empty: 'helper-label',
    forbidden: 'helper-label'
  } as const;
</script>

<div class={classes[tone]}>
  {#if title}
    <p class={labelClasses[tone]}>{title}</p>
  {/if}

  <p class={`text-sm leading-6 ${tone === 'error' ? 'mt-2 text-slate-700' : 'mt-2 text-slate-600'}`}>
    {message}
  </p>

  {#if requestId}
    <p class="mt-2 text-xs text-slate-400">Request ID: {requestId}</p>
  {/if}
</div>
