import { writable } from 'svelte/store';

export const authUser = writable<App.Locals['user']>(null);
