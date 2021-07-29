import {check, group, sleep} from 'k6';
import http from 'k6/http';
import {BASE_URL, COMMON_REQUEST_HEADERS, randomUsername, loadPhoto} from "./utils.js";

export let options = {
    stages: [
        {duration: '30s', target: 1000},
    ],
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'p(99.99)', 'count'],
};

export default function () {
    const headers = Object.assign({}, COMMON_REQUEST_HEADERS);;

    const rand = Math.floor(Math.random() * 5)

    group("login", function() {
        const res = http.post(
            `${BASE_URL}/auth/login`,
            `{"username":"${randomUsername(10000000)}", "password": "secret"}`,
            {
                tags: {name: '/login'},
                headers: COMMON_REQUEST_HEADERS,
            }
        );

        headers['X-Csrftoken'] = res.headers['X-Csrftoken'];

        check(res, {
            'can login': (res) => res.status === 200
        });
    });

    if (rand === 0) return;

    group("profile", function() {
        let res;
        if (rand === 1) {
            res = http.put(
                `${BASE_URL}/profile`,
                `{"nickname": "whatever"}`,
                {
                    tags: {name: '/profile'},
                    headers: headers,
                }
            );
            check(res, {
                'can change nickname': (res) => res.status === 200
            });
        }

        if (rand === 2) {
            const fd = loadPhoto();
            const formHeaders = Object.assign({}, headers);
            formHeaders['content-type'] = 'multipart/form-data; boundary=' + fd.boundary;
            res = http.post(
                `${BASE_URL}/profile-picture`, fd.body(),
                {
                    tags: {name: '/profile-picture'},
                    headers: formHeaders,
                }
            );
            check(res, {
                'can upload profile picture': (res) => res.status === 200
            });
        }

        if (rand === 3) {
            res = http.get(
                `${BASE_URL}/profile`,
                {
                    tags: {name: '/profile'},
                    headers: headers,
                }
            );
            check(res, {
                'can get profile': (res) => res.status === 200
            });
        }
    });

    group("logout", function() {
        if (rand !== 4) return;
        const res = http.post(
            `${BASE_URL}/auth/logout`,
            "",
            {
                tags: {name: '/logout'},
                headers: headers,
            }
        );
        check(res, {
            'can logout': (res) => res.status === 200
        });
    });
};