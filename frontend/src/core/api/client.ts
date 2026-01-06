import { toCamelCase, toSnakeCase } from '@/src/utils/caseConverter';

export class ApiError extends Error {
    constructor(public status: number, public statusText: string, public data?: any) {
        super(`API Error ${status}: ${statusText}`);
        this.name = 'ApiError';
        // Try to extraction message from standard error format
        if (data && typeof data === 'object' && data.message) {
            this.message = data.message;
        }
    }
}

export class ApiClient {
    private async handleResponse<T>(res: Response): Promise<T> {
        if (!res.ok) {
            let data;
            try {
                data = await res.json();
            } catch (e) {
                // ignore json parse error
            }
            throw new ApiError(res.status, res.statusText, data);
        }
        // Handle empty response (e.g. 204 No Content)
        if (res.status === 204) {
            return {} as T;
        }
        const text = await res.text();
        return text ? toCamelCase(JSON.parse(text)) : ({} as T);
    }

    async get<T>(url: string): Promise<T> {
        const res = await fetch(url, {
            credentials: 'include',
        });
        return this.handleResponse<T>(res);
    }

    async post<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        return this.handleResponse<T>(res);
    }

    async put<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        return this.handleResponse<T>(res);
    }

    async delete(url: string): Promise<void> {
        const res = await fetch(url, {
            method: 'DELETE',
            credentials: 'include',
        });
        if (!res.ok) {
            let data;
            try {
                data = await res.json();
            } catch (e) {
                // ignore
            }
            throw new ApiError(res.status, res.statusText, data);
        }
    }

    async patch<T>(url: string, body: any): Promise<T> {
        const res = await fetch(url, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toSnakeCase(body)),
            credentials: 'include',
        });
        return this.handleResponse<T>(res);
    }
}

export const apiClient = new ApiClient();
