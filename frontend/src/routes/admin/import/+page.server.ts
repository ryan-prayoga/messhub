import type { PageServerLoad } from './$types';

function isAdmin(role: string | undefined) {
  return role === 'admin';
}

export const load: PageServerLoad = async ({ parent, locals }) => {
  await parent();

  return {
    accessDenied: !isAdmin(locals.user?.role)
  };
};
