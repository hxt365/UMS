import { FormData } from 'https://jslib.k6.io/formdata/0.0.2/index.js';
import http from 'k6/http';

export const BASE_URL = 'http://192.168.9.103:8000/api';

export let COMMON_REQUEST_HEADERS = {
    dnt: '1',
    'user-agent': 'Mozilla/5.0',
    'content-type': 'application/json',
    'X-Csrftoken': '',
    accept: '*/*',
    origin: BASE_URL,
    referer: BASE_URL
};

export function randomUsername(maxUsers) {
    return `user${Math.floor(Math.random() * maxUsers) + 1}`
};

const img = open('./../../assets/test_photo.png', 'b');

export function loadPhoto() {
    const fd = new FormData();
    fd.append('profilePicture', http.file(img, 'test_photo.png', 'image/png'));
    return fd;
}