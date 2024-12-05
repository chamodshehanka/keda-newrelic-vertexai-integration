import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Define custom metrics
export let errorRate = new Rate('errors');

export let options = {
    stages: [
        { duration: '30s', target: 10 }, // Ramp-up to 10 users over 30 seconds
        { duration: '1m', target: 10 },  // Stay at 10 users for 1 minute
        { duration: '30s', target: 0 },  // Ramp-down to 0 users over 30 seconds
    ],
    thresholds: {
        errors: ['rate<0.01'], // Error rate should be less than 1%
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    },
};

export default function () {
    const baseURL = __ENV.BASE_URL || "http://localhost:8080";

    // Test the root endpoint
    let res1 = http.get(baseURL);
    check(res1, {
        'status is 200': (r) => r.status === 200,
        'response time is less than 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);

    // Test the /books endpoint
    let res2 = http.get(`${baseURL}/books`);
    check(res2, {
        'status is 200': (r) => r.status === 200,
        'response time is less than 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);

    sleep(1); // Sleep for 1 second between iterations
}