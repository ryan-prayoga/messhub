import { API_BASE_URL } from '$lib/config/env';

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  token?: string;
  body?: Record<string, unknown>;
};

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  });

  if (!response.ok) {
    const fallback = { message: 'Request failed' };
    const payload = (await response.json().catch(() => fallback)) as { message?: string };
    throw new Error(payload.message || fallback.message);
  }

  return response.json() as Promise<T>;
}

export const authApi = {
  login: (payload: { email: string; password: string }) =>
    apiRequest<{
      data: {
        token: string;
        user: { id: string; name: string; email: string; role: 'admin' | 'treasurer' | 'member' };
      };
    }>('/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (token: string) =>
    apiRequest<{ data: { id: string; name: string; email: string; role: 'admin' | 'treasurer' | 'member' } }>(
      '/auth/me',
      { token }
    )
};
