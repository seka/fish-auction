export class ApiClient {
    async get<T>(url: string): Promise<T> {
        const res = await fetch(url);
        if (!res.ok) {
            throw new Error(`GET ${url} failed: ${res.statusText}`);
        }
        return res.json();
    }

    async post<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });
        if (!res.ok) {
            throw new Error(`POST ${url} failed: ${res.statusText}`);
        }
        // Some APIs might not return JSON on success (e.g. 201 Created with empty body)
        // Adjust based on backend response. Assuming JSON or handling empty response if needed.
        const text = await res.text();
        return text ? JSON.parse(text) : ({} as T);
    }
}

export const apiClient = new ApiClient();
