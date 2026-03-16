import { env } from '$env/dynamic/public';

export const APP_NAME = env.PUBLIC_APP_NAME || 'MessHub';
export const API_BASE_URL = env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';
