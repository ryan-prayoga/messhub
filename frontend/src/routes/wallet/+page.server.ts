import type { PageServerLoad } from './$types';
import { ApiError, walletServerApi } from '$lib/api/server';

const DEFAULT_PAGE_SIZE = 20;

function parsePage(value: string | null) {
  const page = Number(value);

  if (!Number.isFinite(page) || page < 1) {
    return 1;
  }

  return Math.floor(page);
}

export const load: PageServerLoad = async ({ fetch, locals, parent, url }) => {
  await parent();

  const page = parsePage(url.searchParams.get('page'));

  if (!locals.token) {
    return {
      summary: null,
      transactions: [],
      pagination: {
        page,
        page_size: DEFAULT_PAGE_SIZE,
        total_items: 0,
        total_pages: 0
      },
      canCreate: false,
      loadError: 'Missing auth token'
    };
  }

  try {
    const [summaryResponse, transactionsResponse] = await Promise.all([
      walletServerApi.summary(fetch, locals.token),
      walletServerApi.listTransactions(fetch, locals.token, {
        page,
        pageSize: DEFAULT_PAGE_SIZE
      })
    ]);

    return {
      summary: summaryResponse.data,
      transactions: transactionsResponse.data.items,
      pagination: transactionsResponse.data.pagination,
      canCreate: ['admin', 'treasurer'].includes(locals.user?.role ?? ''),
      loadError: null
    };
  } catch (error) {
    return {
      summary: null,
      transactions: [],
      pagination: {
        page,
        page_size: DEFAULT_PAGE_SIZE,
        total_items: 0,
        total_pages: 0
      },
      canCreate: ['admin', 'treasurer'].includes(locals.user?.role ?? ''),
      loadError:
        error instanceof ApiError || error instanceof Error
          ? error.message
          : 'Failed to load wallet'
    };
  }
};
