import { API_BASE_URL } from '$lib/config/env';
import type { ApiEnvelope, MemberUser, SessionUser } from '$lib/api/types';

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  token?: string;
  body?: Record<string, unknown>;
};

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(status: number, message: string, code?: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<ApiEnvelope<T>> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  });

  const fallback = { message: 'Request failed' };
  const payload = (await response.json().catch(() => fallback)) as
    | (ApiEnvelope<T> & { message?: string })
    | { message?: string };

  if (!response.ok) {
    throw new ApiError(
      response.status,
      payload.message || fallback.message,
      'error' in payload ? payload.error?.code : undefined
    );
  }

  return payload as ApiEnvelope<T>;
}

export const authApi = {
  login: (payload: { email: string; password: string }) =>
    apiRequest<{ token: string; user: SessionUser }>('/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (token: string) => apiRequest<SessionUser>('/auth/me', { token })
};

export const usersApi = {
  list: (token: string) => apiRequest<MemberUser[]>('/users', { token })
};
