<script lang="ts">
  import { onMount } from 'svelte';
  import { networkState } from '$lib/pwa/network';
  import {
    getPushPromptState,
    requestPushPermissionAndSubscribe,
    type PushPromptState
  } from '$lib/pwa/push';

  let installEvent: BeforeInstallPromptEvent | null = null;
  let canInstall = false;
  let installing = false;
  let pushState: PushPromptState = 'unsupported';
  let pushBusy = false;
  let feedback: string | null = null;

  onMount(() => {
    const handleBeforeInstallPrompt = (event: Event) => {
      const promptEvent = event as BeforeInstallPromptEvent;
      promptEvent.preventDefault();
      installEvent = promptEvent;
      canInstall = true;
    };

    const handleAppInstalled = () => {
      installEvent = null;
      canInstall = false;
      feedback = 'MessHub sudah terpasang di homescreen.';
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    window.addEventListener('appinstalled', handleAppInstalled);
    void refreshPushState();

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  });

  async function installApp() {
    if (!installEvent) {
      return;
    }

    installing = true;
    await installEvent.prompt();
    const choice = await installEvent.userChoice.catch(() => null);
    installing = false;

    if (choice?.outcome === 'accepted') {
      canInstall = false;
      installEvent = null;
      feedback = 'MessHub sedang dipasang ke homescreen.';
      return;
    }

    feedback = 'Install ditunda untuk sekarang.';
  }

  async function enablePush() {
    pushBusy = true;

    try {
      const result = await requestPushPermissionAndSubscribe();
      pushState = result.state;
      feedback = result.message;
    } catch (error) {
      feedback = error instanceof Error ? error.message : 'Gagal mengaktifkan push notification.';
    } finally {
      pushBusy = false;
    }
  }

  async function refreshPushState() {
    pushState = await getPushPromptState();
  }

  $: showInstall = canInstall;
  $: showPushPrompt = pushState === 'permission-needed';
  $: showBar = !$networkState.online || showInstall || showPushPrompt || feedback;
</script>

{#if showBar}
  <div class="mx-auto w-full max-w-4xl px-4 pt-3 sm:px-6">
    <section class="app-panel border border-slate-200/80 bg-white/95 px-4 py-4">
      <div class="flex flex-col gap-3">
        {#if !$networkState.online}
          <div class="native-banner native-banner-offline">
            <div>
              <p class="helper-label text-amber-800">Offline mode</p>
              <p class="mt-1 text-sm leading-6 text-amber-950">
                Menampilkan data terakhir yang tersimpan. Outbox akan sync otomatis saat koneksi kembali.
              </p>
            </div>
          </div>
        {/if}

        {#if showInstall || showPushPrompt}
          <div class="grid gap-3 md:grid-cols-2">
            {#if showInstall}
              <div class="native-banner">
                <div>
                  <p class="helper-label">Homescreen</p>
                  <p class="mt-1 text-sm leading-6 text-slate-700">
                    Pasang MessHub agar terbuka seperti app native dari homescreen.
                  </p>
                </div>

                <button type="button" class="btn-primary px-4 py-3" on:click={installApp} disabled={installing}>
                  {installing ? 'Preparing...' : 'Install MessHub App'}
                </button>
              </div>
            {/if}

            {#if showPushPrompt}
              <div class="native-banner">
                <div>
                  <p class="helper-label">Notifications</p>
                  <p class="mt-1 text-sm leading-6 text-slate-700">
                    Aktifkan push notification untuk tagihan wifi dan update feed.
                  </p>
                </div>

                <button type="button" class="btn-secondary px-4 py-3" on:click={enablePush} disabled={pushBusy}>
                  {pushBusy ? 'Enabling...' : 'Aktifkan notifikasi'}
                </button>
              </div>
            {/if}
          </div>
        {/if}

        {#if feedback}
          <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-600">
            {feedback}
          </div>
        {/if}
      </div>
    </section>
  </div>
{/if}
