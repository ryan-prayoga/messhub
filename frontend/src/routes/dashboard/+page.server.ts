import type { PageServerLoad } from './$types';
import {
  ApiError,
  contributionsServerApi,
  usersServerApi,
  walletServerApi,
  wifiServerApi
} from '$lib/api/server';

type MemberSummary = {
  total: number | null;
  active: number | null;
  inactive: number | null;
  state: 'ready' | 'restricted' | 'error';
  message: string | null;
};

type WalletSummary = {
  balance: number | null;
  totalIncome: number | null;
  totalExpense: number | null;
  state: 'ready' | 'error';
  message: string | null;
};

type WifiSummary = {
  monthLabel: string | null;
  verified: number | null;
  unpaid: number | null;
  pending: number | null;
  totalCollected: number | null;
  totalTarget: number | null;
  myStatus: string | null;
  deadline: string | null;
  state: 'ready' | 'empty' | 'error';
  message: string | null;
};

type LeaderboardSummary = {
  items: {
    rank: number;
    user_name: string;
    total_points: number;
    total_activities: number;
  }[];
  state: 'ready' | 'empty' | 'error';
  message: string | null;
};

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
  await parent();

  const summary: MemberSummary = {
    total: null,
    active: null,
    inactive: null,
    state: 'restricted',
    message: null
  };

  const walletSummary: WalletSummary = {
    balance: null,
    totalIncome: null,
    totalExpense: null,
    state: 'error',
    message: 'Wallet summary unavailable'
  };

  const wifiSummary: WifiSummary = {
    monthLabel: null,
    verified: null,
    unpaid: null,
    pending: null,
    totalCollected: null,
    totalTarget: null,
    myStatus: null,
    deadline: null,
    state: 'empty',
    message: 'No active wifi bill'
  };

  const leaderboardSummary: LeaderboardSummary = {
    items: [],
    state: 'empty',
    message: 'Belum ada kontribusi bulan ini'
  };

  if (locals.token) {
    const [walletResult, wifiResult, leaderboardResult] = await Promise.allSettled([
      walletServerApi.summary(fetch, locals.token),
      wifiServerApi.getActive(fetch, locals.token),
      contributionsServerApi.leaderboard(fetch, locals.token, 'month')
    ]);

    if (walletResult.status === 'fulfilled') {
      walletSummary.balance = walletResult.value.data.balance;
      walletSummary.totalIncome = walletResult.value.data.total_income;
      walletSummary.totalExpense = walletResult.value.data.total_expense;
      walletSummary.state = 'ready';
      walletSummary.message = null;
    } else {
      walletSummary.state = 'error';
      walletSummary.message =
        walletResult.reason instanceof Error
          ? walletResult.reason.message
          : 'Failed to load wallet summary';
    }

    if (wifiResult.status === 'fulfilled') {
      const wifiResponse = wifiResult.value;
      if (wifiResponse.data) {
        const bill = wifiResponse.data.bill;
        const member = wifiResponse.data.members[0] ?? null;

        wifiSummary.monthLabel = new Intl.DateTimeFormat('id-ID', {
          month: 'long',
          year: 'numeric'
        }).format(new Date(bill.year, bill.month - 1, 1));
        wifiSummary.verified = wifiResponse.data.summary.verified_count;
        wifiSummary.unpaid = wifiResponse.data.summary.unpaid_count;
        wifiSummary.pending = wifiResponse.data.summary.pending_count;
        wifiSummary.totalCollected = wifiResponse.data.summary.total_collected;
        wifiSummary.totalTarget = wifiResponse.data.summary.total_target;
        wifiSummary.deadline = bill.deadline_date;
        wifiSummary.myStatus = member?.payment_status ?? null;
        wifiSummary.state = 'ready';
        wifiSummary.message = null;
      }
    } else {
      wifiSummary.state = 'error';
      wifiSummary.message =
        wifiResult.reason instanceof Error ? wifiResult.reason.message : 'Failed to load wifi summary';
    }

    if (leaderboardResult.status === 'fulfilled') {
      leaderboardSummary.items = leaderboardResult.value.data;
      leaderboardSummary.state = leaderboardResult.value.data.length > 0 ? 'ready' : 'empty';
      leaderboardSummary.message =
        leaderboardResult.value.data.length > 0 ? null : 'Belum ada kontribusi bulan ini';
    } else {
      leaderboardSummary.state = 'error';
      leaderboardSummary.message =
        leaderboardResult.reason instanceof Error
          ? leaderboardResult.reason.message
          : 'Failed to load contribution leaderboard';
    }
  }

  if (locals.token && locals.user && ['admin', 'treasurer'].includes(locals.user.role)) {
    try {
      const response = await usersServerApi.list(fetch, locals.token);
      const members = response.data;
      const active = members.filter((member) => member.is_active).length;

      summary.total = members.length;
      summary.active = active;
      summary.inactive = members.length - active;
      summary.state = 'ready';
    } catch (error) {
      if (error instanceof ApiError && error.status === 403) {
        summary.state = 'restricted';
      } else {
        summary.state = 'error';
        summary.message = error instanceof Error ? error.message : 'Failed to load members summary';
      }
    }
  }

  return {
    authStatus: locals.user ? 'authenticated' : 'unauthenticated',
    memberSummary: summary,
    walletSummary,
    wifiSummary,
    leaderboardSummary
  };
};
