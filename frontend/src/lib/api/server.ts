import { env } from '$env/dynamic/private';
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
  ImportCommitResult,
  MemberImportPreview,
  MessSettings,
  MemberUser,
  NotificationList,
  Profile,
  SessionUser,
  SystemStatus,
  WalletImportPreview,
  WifiBillDetail,
  WifiBillWithSummary,
  WifiMyBill,
  WalletSummary,
  WalletTransactionPage,
  WalletTransactionType
} from '$lib/api/types';

const DEFAULT_PRIVATE_API_BASE_URL = 'http://127.0.0.1:4100/api/v1';

export { ApiError } from '$lib/api/http';

function resolveApiUrl(path: string) {
  const baseUrl = env.PRIVATE_API_BASE_URL || DEFAULT_PRIVATE_API_BASE_URL;
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  return `${baseUrl}${normalizedPath}`;
}

export async function apiServerRequest<T>(
  fetcher: typeof fetch,
  path: string,
  options: RequestOptions = {}
): Promise<ApiEnvelope<T>> {
  try {
    const response = await fetcher(resolveApiUrl(path), buildRequestInit(options));
    return await parseApiResponse<T>(response);
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    throw wrapNetworkError(error, 'Backend request failed');
  }
}

export const authServerApi = {
  login: (
    fetcher: typeof fetch,
    payload: { identifier: string; password: string; email?: string }
  ) =>
    apiServerRequest<{ token: string; user: SessionUser }>(fetcher, '/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<SessionUser>(fetcher, '/auth/me', {
      token
    })
};

export const usersServerApi = {
  list: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<MemberUser[]>(fetcher, '/users', {
      token
    }),
  update: (
    fetcher: typeof fetch,
    token: string,
    userID: string,
    payload: {
      role?: MemberUser['role'];
      is_active?: boolean;
      name?: string;
    }
  ) =>
    apiServerRequest<MemberUser>(fetcher, `/users/${userID}`, {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const profileServerApi = {
  get: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<Profile>(fetcher, '/profile', {
      token
    }),
  update: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      name?: string;
      phone?: string;
      avatar_url?: string;
    }
  ) =>
    apiServerRequest<Profile>(fetcher, '/profile', {
      method: 'PATCH',
      token,
      body: payload
    }),
  changePassword: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      current_password: string;
      new_password: string;
    }
  ) =>
    apiServerRequest<{ changed: boolean }>(fetcher, '/profile/password', {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const settingsServerApi = {
  get: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<MessSettings>(fetcher, '/settings', {
      token
    }),
  update: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      mess_name?: string;
      wifi_price?: number;
      wifi_deadline_day?: number;
      bank_account_name?: string;
      bank_account_number?: string;
    }
  ) =>
    apiServerRequest<MessSettings>(fetcher, '/settings', {
      method: 'PATCH',
      token,
      body: payload
    })
};

export const systemServerApi = {
  status: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<SystemStatus>(fetcher, '/system/status', {
      token
    })
};

export const activitiesServerApi = {
  list: (
    fetcher: typeof fetch,
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

    return apiServerRequest<ActivityFeedItem[]>(
      fetcher,
      `/activities${query ? `?${query}` : ''}`,
      {
        token
      }
    );
  },
  create: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      type: ActivityType;
      title: string;
      content: string;
      points?: number;
    }
  ) =>
    apiServerRequest<ActivityFeedItem>(fetcher, '/activities', {
      method: 'POST',
      token,
      body: payload
    }),
  listComments: (fetcher: typeof fetch, token: string, activityID: string) =>
    apiServerRequest<ActivityComment[]>(fetcher, `/activities/${activityID}/comments`, {
      token
    }),
  addComment: (
    fetcher: typeof fetch,
    token: string,
    activityID: string,
    payload: {
      comment: string;
    }
  ) =>
    apiServerRequest<ActivityFeedItem>(fetcher, `/activities/${activityID}/comments`, {
      method: 'POST',
      token,
      body: payload
    }),
  toggleReaction: (
    fetcher: typeof fetch,
    token: string,
    activityID: string,
    payload: {
      reaction_type: string;
    }
  ) =>
    apiServerRequest<ActivityFeedItem>(fetcher, `/activities/${activityID}/reactions`, {
      method: 'POST',
      token,
      body: payload
    }),
  claimFood: (fetcher: typeof fetch, token: string, activityID: string) =>
    apiServerRequest(fetcher, `/activities/${activityID}/claim`, {
      method: 'POST',
      token
    }),
  respondRice: (fetcher: typeof fetch, token: string, activityID: string) =>
    apiServerRequest(fetcher, `/activities/${activityID}/rice-response`, {
      method: 'POST',
      token
    })
};

export const contributionsServerApi = {
  leaderboard: (fetcher: typeof fetch, token: string, period: 'month' | 'all' = 'month') =>
    apiServerRequest<ContributionLeaderboardEntry[]>(
      fetcher,
      `/contributions/leaderboard?period=${period}`,
      {
        token
      }
    )
};

export const notificationsServerApi = {
  list: (
    fetcher: typeof fetch,
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

    return apiServerRequest<NotificationList>(
      fetcher,
      `/notifications${query ? `?${query}` : ''}`,
      {
        token
      }
    );
  },
  markRead: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      ids?: string[];
      all?: boolean;
    }
  ) =>
    apiServerRequest<{ updated_count: number }>(fetcher, '/notifications/read', {
      method: 'POST',
      token,
      body: payload
    })
};

export const walletServerApi = {
  summary: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WalletSummary>(fetcher, '/wallet', {
      token
    }),
  listTransactions: (
    fetcher: typeof fetch,
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

    return apiServerRequest<WalletTransactionPage>(
      fetcher,
      `/wallet/transactions${query ? `?${query}` : ''}`,
      {
        token
      }
    );
  },
  createTransaction: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      type: WalletTransactionType;
      category: string;
      amount: number;
      description: string;
    }
  ) =>
    apiServerRequest(fetcher, '/wallet/transactions', {
      method: 'POST',
      token,
      body: payload
    })
};

export const importsServerApi = {
  previewMembers: (fetcher: typeof fetch, token: string, formData: FormData) =>
    apiServerRequest<MemberImportPreview>(fetcher, '/import/members/preview', {
      method: 'POST',
      token,
      body: formData
    }),
  commitMembers: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      job_id: string;
      duplicate_strategy: 'skip' | 'fail';
      temporary_password: string;
    }
  ) =>
    apiServerRequest<ImportCommitResult>(fetcher, '/import/members/commit', {
      method: 'POST',
      token,
      body: payload
    }),
  previewWallet: (fetcher: typeof fetch, token: string, formData: FormData) =>
    apiServerRequest<WalletImportPreview>(fetcher, '/import/wallet/preview', {
      method: 'POST',
      token,
      body: formData
    }),
  commitWallet: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      job_id: string;
    }
  ) =>
    apiServerRequest<ImportCommitResult>(fetcher, '/import/wallet/commit', {
      method: 'POST',
      token,
      body: payload
    })
};

export const wifiServerApi = {
  listBills: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiBillWithSummary[]>(fetcher, '/wifi/bills', {
      token
    }),
  getBill: (fetcher: typeof fetch, token: string, billID: string) =>
    apiServerRequest<WifiBillDetail>(fetcher, `/wifi/bills/${billID}`, {
      token
    }),
  getActive: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiBillDetail | null>(fetcher, '/wifi/active', {
      token
    }),
  getMyBills: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiMyBill[]>(fetcher, '/wifi/my', {
      token
    }),
  createBill: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      month: number;
      year: number;
      nominal_per_person?: number;
      deadline_date?: string;
      status?: 'draft' | 'active' | 'closed';
    }
  ) =>
    apiServerRequest<WifiBillDetail>(fetcher, '/wifi/bills', {
      method: 'POST',
      token,
      body: payload
    }),
  submitProof: (
    fetcher: typeof fetch,
    token: string,
    billID: string,
    payload: {
      proof_url: string;
      note?: string;
    }
  ) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/submit`, {
      method: 'POST',
      token,
      body: payload
    }),
  verifyPayment: (fetcher: typeof fetch, token: string, billID: string, memberID: string) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/verify/${memberID}`, {
      method: 'PATCH',
      token
    }),
  rejectPayment: (
    fetcher: typeof fetch,
    token: string,
    billID: string,
    memberID: string,
    payload: {
      reason: string;
    }
  ) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/reject/${memberID}`, {
      method: 'PATCH',
      token,
      body: payload
    })
};
