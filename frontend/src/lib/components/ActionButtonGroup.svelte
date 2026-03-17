<script lang="ts">
  export let align: 'start' | 'end' | 'between' = 'start';
  export let bordered = false;

  const alignClass: Record<'start' | 'end' | 'between', string> = {
    start: 'action-button-group-start',
    end: 'action-button-group-end',
    between: 'action-button-group-between'
  };
</script>

<div class={`action-button-group ${alignClass[align]} ${bordered ? 'action-button-group-bordered' : ''}`}>
  <slot />
</div>
