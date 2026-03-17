import type { SubmitFunction } from '@sveltejs/kit';
import { get, writable } from 'svelte/store';

export type ConfirmationConfig = {
  title: string;
  description: string;
  confirmLabel: string;
  cancelLabel?: string;
  icon?: string;
  destructive?: boolean;
};

export type ConfirmationDialogState = ConfirmationConfig & {
  actionKey: string;
  form: HTMLFormElement;
  submitter: HTMLElement | null;
};

type ControllerState = {
  dialog: ConfirmationDialogState | null;
  confirmedActionKey: string | null;
  requestingActionKey: string | null;
};

type ConfirmableControllerOptions = {
  setPendingAction: (value: string | null) => void;
};

export function createConfirmableSubmitController(options: ConfirmableControllerOptions) {
  const state = writable<ControllerState>({
    dialog: null,
    confirmedActionKey: null,
    requestingActionKey: null
  });

  function reset() {
    state.set({
      dialog: null,
      confirmedActionKey: null,
      requestingActionKey: null
    });
  }

  function enhance(actionKey: string, confirmation?: ConfirmationConfig): SubmitFunction {
    return ({ formElement, cancel, submitter }) => {
      const current = get(state);

      if (confirmation && current.confirmedActionKey !== actionKey) {
        cancel();
        state.set({
          dialog: {
            actionKey,
            form: formElement,
            submitter: submitter instanceof HTMLElement ? submitter : null,
            cancelLabel: confirmation.cancelLabel ?? 'Batal',
            ...confirmation
          },
          confirmedActionKey: null,
          requestingActionKey: null
        });
        return;
      }

      state.update((value) => ({
        ...value,
        confirmedActionKey: null,
        requestingActionKey: null
      }));
      options.setPendingAction(actionKey);

      return async ({ update }) => {
        try {
          await update();
        } finally {
          options.setPendingAction(null);
          const latest = get(state);

          if (latest.dialog?.actionKey === actionKey) {
            reset();
          } else {
            state.update((value) => ({
              ...value,
              requestingActionKey: null
            }));
          }
        }
      };
    };
  }

  function closeDialog() {
    if (get(state).requestingActionKey) {
      return;
    }

    reset();
  }

  function confirmDialog() {
    const current = get(state);
    if (!current.dialog || current.requestingActionKey === current.dialog.actionKey) {
      return;
    }

    state.set({
      ...current,
      confirmedActionKey: current.dialog.actionKey,
      requestingActionKey: current.dialog.actionKey
    });

    current.dialog.form.requestSubmit(
      current.dialog.submitter instanceof HTMLElement ? current.dialog.submitter : undefined
    );
  }

  return {
    state,
    enhance,
    closeDialog,
    confirmDialog
  };
}
