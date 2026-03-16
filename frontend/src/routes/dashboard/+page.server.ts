import type { PageServerLoad } from './$types';
import { ApiError, usersServerApi, walletServerApi } from '$lib/api/server';

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

  if (locals.token) {
    try {
      const response = await walletServerApi.summary(fetch, locals.token);

      walletSummary.balance = response.data.balance;
      walletSummary.totalIncome = response.data.total_income;
      walletSummary.totalExpense = response.data.total_expense;
      walletSummary.state = 'ready';
      walletSummary.message = null;
    } catch (error) {
      walletSummary.state = 'error';
      walletSummary.message = error instanceof Error ? error.message : 'Failed to load wallet summary';
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
    walletSummary
  };
};
