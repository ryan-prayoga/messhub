import { API_BASE_URL } from '$lib/config/env';
import {
  ApiError,
  buildRequestInit,
  parseApiResponse,
  type RequestOptions,
  wrapNetworkError
} from '$lib/api/http';
import type {
  ActivityComment,
  ActivityFeedItem,
  ActivityType,
  ApiEnvelope,
  ContributionLeaderboardEntry,
  MessSettings,
  MemberUser,
  NotificationList,
  Profile,
  SessionUser,
  SystemStatus,
  WifiBillDetail,
  WifiBillWithSummary,
  WifiMyBill,
  WalletSummary,
  WalletTransactionPage,
  WalletTransactionType
} from '$lib/api/types';

export { ApiError } from '$lib/api/http';

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<ApiEnvelope<T>> {
  try {
    const response = await fetch(`${API_BASE_URL}${path}`, buildRequestInit(options));
    return await parseApiResponse<T>(response);
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    throw wrapNetworkError(error);
  }
}

export const authApi = {
  login: (payload: { identifier: string; password: string; email?: string }) =>
    apiRequest<{ token: string; user: SessionUser }>('/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (token: string) => apiRequest<SessionUser>('/auth/me', { token })
};

export const usersApi = {
  list: (token: string) => apiRequest<MemberUser[]>('/users', { token }),
  create: (
    token: string,
    payload: {
      name: string;
      email: string;
      username?: string;
      phone?: string;
      password: string;
      role: MemberUser['role'];
      is_active?: boolean;
      joined_at?: string;
    }
  ) =>
    apiRequest<MemberUser>('/users', {
      method: 'POST',
      token,
      body: payload
    }),
  update: (
    token: string,
    userID: string,
    payload: {
      role?: MemberUser['role'];
      is_active?: boolean;
      name?: string;
      email?: string;
      username?: string;
      phone?: string;
      joined_at?: string;
    }
  ) =>
    apiRequest<MemberUser>(`/users/${userID}`, {
      method: 'PATCH',
      token,
      body: payload
    }),
  resetPassword: (
    token: string,
    userID: string,
    payload: {
      new_password: string;
    }
  ) =>
    apiRequest<{ changed: boolean }>(`/users/${userID}/password`, {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const profileApi = {
  get: (token: string) => apiRequest<Profile>('/profile', { token }),
  update: (
    token: string,
    payload: {
      name?: string;
      phone?: string;
      avatar_url?: string;
    }
  ) =>
    apiRequest<Profile>('/profile', {
      method: 'PATCH',
      token,
      body: payload
    }),
  changePassword: (
    token: string,
    payload: {
      current_password: string;
      new_password: string;
    }
  ) =>
    apiRequest<{ changed: boolean }>('/profile/password', {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const settingsApi = {
  get: (token: string) => apiRequest<MessSettings>('/settings', { token }),
  update: (
    token: string,
    payload: {
      mess_name?: string;
      wifi_price?: number;
      wifi_deadline_day?: number;
      bank_account_name?: string;
      bank_account_number?: string;
    }
  ) =>
    apiRequest<MessSettings>('/settings', {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const systemApi = {
  status: (token: string) => apiRequest<SystemStatus>('/system/status', { token })
};

export const activitiesApi = {
  list: (
    token: string,
    params: {
      limit?: number;
    } = {}
  ) => {
    const searchParams = new URLSearchParams();

    if (params.limit) {
      searchParams.set('limit', String(params.limit));
    }

    const query = searchParams.toString();

    return apiRequest<ActivityFeedItem[]>(`/activities${query ? `?${query}` : ''}`, {
      token
    });
  },
  create: (
    token: string,
    payload: {
      type: ActivityType;
      title: string;
      content: string;
      points?: number;
    }
  ) =>
    apiRequest<ActivityFeedItem>('/activities', {
      method: 'POST',
      token,
      body: payload
    }),
  listComments: (token: string, activityID: string) =>
    apiRequest<ActivityComment[]>(`/activities/${activityID}/comments`, { token }),
  addComment: (
    token: string,
    activityID: string,
    payload: {
      comment: string;
    }
  ) =>
    apiRequest<ActivityFeedItem>(`/activities/${activityID}/comments`, {
      method: 'POST',
      token,
      body: payload
    }),
  toggleReaction: (
    token: string,
    activityID: string,
    payload: {
      reaction_type: string;
    }
  ) =>
    apiRequest<ActivityFeedItem>(`/activities/${activityID}/reactions`, {
      method: 'POST',
      token,
      body: payload
    }),
  claimFood: (token: string, activityID: string) =>
    apiRequest(`/activities/${activityID}/claim`, {
      method: 'POST',
      token
    }),
  respondRice: (token: string, activityID: string) =>
    apiRequest(`/activities/${activityID}/rice-response`, {
      method: 'POST',
      token
    })
};

export const contributionsApi = {
  leaderboard: (token: string, period: 'month' | 'all' = 'month') =>
    apiRequest<ContributionLeaderboardEntry[]>(`/contributions/leaderboard?period=${period}`, {
      token
    })
};

export const notificationsApi = {
  list: (
    token: string,
    params: {
      limit?: number;
    } = {}
  ) => {
    const searchParams = new URLSearchParams();

    if (params.limit) {
      searchParams.set('limit', String(params.limit));
    }

    const query = searchParams.toString();

    return apiRequest<NotificationList>(`/notifications${query ? `?${query}` : ''}`, {
      token
    });
  },
  markRead: (
    token: string,
    payload: {
      ids?: string[];
      all?: boolean;
    }
  ) =>
    apiRequest<{ updated_count: number }>('/notifications/read', {
      method: 'POST',
      token,
      body: payload
    })
};

export const pushApi = {
  subscribe: (payload: {
    endpoint: string;
    keys: {
      p256dh: string;
      auth: string;
    };
  }) =>
    apiRequest<{
      id: string;
      user_id: string;
      endpoint: string;
      p256dh_key: string;
      auth_key: string;
      created_at: string;
    }>('/push/subscribe', {
      method: 'POST',
      body: payload
    }),
  unsubscribe: (endpoint: string) =>
    apiRequest<{ removed: boolean }>('/push/unsubscribe', {
      method: 'DELETE',
      body: {
        endpoint
      }
    })
};

export const walletApi = {
  summary: (token: string) => apiRequest<WalletSummary>('/wallet', { token }),
  listTransactions: (
    token: string,
    params: {
      page?: number;
      pageSize?: number;
    } = {}
  ) => {
    const searchParams = new URLSearchParams();

    if (params.page) {
      searchParams.set('page', String(params.page));
    }

    if (params.pageSize) {
      searchParams.set('page_size', String(params.pageSize));
    }

    const query = searchParams.toString();

    return apiRequest<WalletTransactionPage>(`/wallet/transactions${query ? `?${query}` : ''}`, {
      token
    });
  },
  createTransaction: (
    token: string,
    payload: {
      type: WalletTransactionType;
      category: string;
      amount: number;
      description: string;
    }
  ) =>
    apiRequest('/wallet/transactions', {
      method: 'POST',
      token,
      body: payload
    })
};

export const wifiApi = {
  listBills: (token: string) => apiRequest<WifiBillWithSummary[]>('/wifi/bills', { token }),
  getBill: (token: string, billID: string) => apiRequest<WifiBillDetail>(`/wifi/bills/${billID}`, { token }),
  getActive: (token: string) => apiRequest<WifiBillDetail | null>('/wifi/active', { token }),
  getMyBills: (token: string) => apiRequest<WifiMyBill[]>('/wifi/my', { token }),
  createBill: (
    token: string,
    payload: {
      month: number;
      year: number;
      nominal_per_person?: number;
      deadline_date?: string;
      status?: 'draft' | 'active' | 'closed';
    }
  ) =>
    apiRequest<WifiBillDetail>('/wifi/bills', {
      method: 'POST',
      token,
      body: payload
    }),
  submitProof: (
    token: string,
    billID: string,
    payload: {
      proof_url: string;
      note?: string;
    }
  ) =>
    apiRequest(`/wifi/bills/${billID}/submit`, {
      method: 'POST',
      token,
      body: payload
    }),
  verifyPayment: (token: string, billID: string, memberID: string) =>
    apiRequest(`/wifi/bills/${billID}/verify/${memberID}`, {
      method: 'PATCH',
      token
    }),
  rejectPayment: (
    token: string,
    billID: string,
    memberID: string,
    payload: {
      reason: string;
    }
  ) =>
    apiRequest(`/wifi/bills/${billID}/reject/${memberID}`, {
      method: 'PATCH',
      token,
      body: payload
    })
};
