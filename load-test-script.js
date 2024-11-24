import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '5s', target: 10000 },
        { duration: '1m', target: 10000 },
        { duration: '5s', target: 0 },    
    ],
};

export default function () {
    let id = Math.floor(Math.random() * 100000); 
    let res = http.get(`http://localhost:8080/api/verve/accept?id=${id}`);

    check(res, {
        'ok': (r) => r.body === 'ok\n',
    });
}