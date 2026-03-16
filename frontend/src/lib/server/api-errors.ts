import { redirect } from '@sveltejs/kit';
import type { Cookies } from '@sveltejs/kit';
import { clearAuthCookies } from '$lib/auth/session';
import { ApiError, isUnauthorizedError } from '$lib/api/http';

export type ApiFailureState = {
  status: number;
  message: string;
  code?: string;
  requestId?: string;
  kind: 'forbidden' | 'network' | 'error';
};

export function throwIfUnauthorized(error: unknown, cookies: Cookies) {
  if (!isUnauthorizedError(error)) {
    return;
  }

  clearAuthCookies(cookies);
  throw redirect(303, '/login');
}

export function toApiFailureState(error: unknown, fallbackMessage: string): ApiFailureState {
  if (error instanceof ApiError) {
    const apiError = error;

    return {
      status: apiError.status || 503,
      message: apiError.message || fallbackMessage,
      code: apiError.code,
      requestId: apiError.requestId,
      kind: apiError.status === 403 ? 'forbidden' : apiError.isNetworkError ? 'network' : 'error'
    };
  }

  return {
    status: 503,
    message: error instanceof Error ? error.message : fallbackMessage,
    kind: 'error'
  };
}
