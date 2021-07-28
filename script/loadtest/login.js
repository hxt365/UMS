import {check} from 'k6';
import http from 'k6/http';
import {BASE_URL, COMMON_REQUEST_HEADERS, randomUsername} from "./utils.js";

export let options = {
    stages: [
        {duration: '3m', target: 10},
    ],
};

export default function () {
    const res = http.post(
        `${BASE_URL}/auth/login`,
        `{"username":"${randomUsername()}", "password": "secret"}`,
        {
            tags: {name: '/login'},
            headers: COMMON_REQUEST_HEADERS,
        }
    );
    check(res, {
        'can login': (res) => res.status === 200
    });
};