import {check} from 'k6';
import http from 'k6/http';
import {BASE_URL, COMMON_REQUEST_HEADERS, randomUsername} from "./utils.js";

export let options = {
    stages: [
        {duration: '1m', target: 1},
    ],
};

export default function () {
    const res = http.post(
        `${BASE_URL}/auth/login`,
        `{"username":"${randomUsername(1)}", "password": "secret"}`,
        {
            tags: {name: '/login'},
            headers: COMMON_REQUEST_HEADERS,
        }
    );
    check(res, {
        'can login': (res) => res.status === 200
    });
};