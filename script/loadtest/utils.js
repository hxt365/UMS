export const BASE_URL = 'http://localhost:8000/api';

export const COMMON_REQUEST_HEADERS = {
    dnt: '1',
    'user-agent': 'Mozilla/5.0',
    'content-type': 'application/json',
    accept: '*/*',
    origin: BASE_URL,
    referer: BASE_URL
};

export function randomUsername(maxUsers) {
    return `user${Math.floor(Math.random() * maxUsers) + 1}`
};