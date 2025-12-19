import { toCamelCase, toSnakeCase } from '@/src/utils/caseConverter';

export class ApiClient {
    async get<T>(url: string): Promise<T> {
        const res = await fetch(url, {
            credentials: 'include',
        });
        if (!res.ok) {
            throw new Error(`GET ${url} failed: ${res.statusText}`);
        }
        return toCamelCase(await res.json());
    }

    async post<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        if (!res.ok) {
            throw new Error(`POST ${url} failed: ${res.statusText}`);
        }
        // Some APIs might not return JSON on success (e.g. 201 Created with empty body)
        // Adjust based on backend response. Assuming JSON or handling empty response if needed.
        const text = await res.text();
        return text ? toCamelCase(JSON.parse(text)) : ({} as T);
    }

    async put<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        if (!res.ok) {
            throw new Error(`PUT ${url} failed: ${res.statusText}`);
        }
        const text = await res.text();
        return text ? toCamelCase(JSON.parse(text)) : ({} as T);
    }

    async delete(url: string): Promise<void> {
        const res = await fetch(url, {
            method: 'DELETE',
            credentials: 'include',
        });
        if (!res.ok) {
            throw new Error(`DELETE ${url} failed: ${res.statusText}`);
        }
    }

    async patch<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        if (!res.ok) {
            throw new Error(`PATCH ${url} failed: ${res.statusText}`);
        }
        const text = await res.text();
        return text ? toCamelCase(JSON.parse(text)) : ({} as T);
    }
}

export const apiClient = new ApiClient();
