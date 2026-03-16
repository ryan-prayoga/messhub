import type { PageServerLoad } from './$types';
import { contributionsServerApi } from '$lib/api/server';

export const load: PageServerLoad = async ({ fetch, locals, parent }) => {
  await parent();

  if (!locals.token) {
    return {
      monthly: [],
      allTime: [],
      loadError: 'Sesi login tidak ditemukan.'
    };
  }

  try {
    const [monthlyResponse, allTimeResponse] = await Promise.all([
      contributionsServerApi.leaderboard(fetch, locals.token, 'month'),
      contributionsServerApi.leaderboard(fetch, locals.token, 'all')
    ]);

    return {
      monthly: monthlyResponse.data,
      allTime: allTimeResponse.data,
      loadError: null
    };
  } catch (error) {
    return {
      monthly: [],
      allTime: [],
      loadError: error instanceof Error ? error.message : 'Leaderboard belum dapat dimuat.'
    };
  }
};
