import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
  const id = event.cookies.get('mh_user_id');
  const email = event.cookies.get('mh_user_email');
  const name = event.cookies.get('mh_user_name');
  const role = event.cookies.get('mh_user_role') as 'admin' | 'treasurer' | 'member' | undefined;

  if (id && email && name && role) {
    event.locals.user = { id, email, name, role };
  } else {
    event.locals.user = null;
  }

  return resolve(event);
};
