import { env } from '$env/dynamic/public';

export const APP_NAME = env.PUBLIC_APP_NAME || 'MessHub';
export const API_BASE_URL = env.PUBLIC_API_BASE_URL || '/api/v1';
