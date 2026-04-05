import { toCamelCase, toSnakeCase } from '@/src/core/api/caseConverter';

export class ApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public data?: unknown,
  ) {
    super(`API Error ${status}: ${statusText}`);
    this.name = 'ApiError';
    // Try to extraction message from standard error format
    if (
      data &&
      typeof data === 'object' &&
      'message' in data &&
      typeof (data as { message: unknown }).message === 'string'
    ) {
      this.message = (data as { message: string }).message;
    }
  }
}

const getBaseUrl = () => {
  if (typeof window === 'undefined') {
    // Server-side
    return process.env.API_BASE_URL || 'http://backend:8080';
  }
  // Client-side
  return process.env.NEXT_PUBLIC_API_URL || '';
};

export class ApiClient {
  private baseUrl = getBaseUrl();

  private getFullUrl(path: string): string {
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    return `${this.baseUrl}${cleanPath}`;
  }

  private async handleResponse<T>(res: Response): Promise<T> {
    if (!res.ok) {
      let data;
      try {
        data = await res.json();
      } catch {
        // ignore json parse error
      }
      throw new ApiError(res.status, res.statusText, data);
    }
    // Handle empty response (e.g. 204 No Content)
    if (res.status === 204) {
      return {} as T;
    }
    const text = await res.text();
    return text ? (toCamelCase(JSON.parse(text)) as T) : ({} as T);
  }

  async get<T>(url: string): Promise<T> {
    const res = await fetch(this.getFullUrl(url), {
      credentials: 'include',
    });
    return this.handleResponse<T>(res);
  }

  async post<T>(url: string, body: unknown): Promise<T> {
    const res = await fetch(this.getFullUrl(url), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(toSnakeCase(body)),
      credentials: 'include',
    });
    return this.handleResponse<T>(res);
  }

  async put<T>(url: string, body: unknown): Promise<T> {
    const res = await fetch(this.getFullUrl(url), {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(toSnakeCase(body)),
      credentials: 'include',
    });
    return this.handleResponse<T>(res);
  }

  async delete(url: string): Promise<void> {
    const res = await fetch(this.getFullUrl(url), {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!res.ok) {
      let data;
      try {
        data = await res.json();
      } catch {
        // ignore
      }
      throw new ApiError(res.status, res.statusText, data);
    }
  }

  async patch<T>(url: string, body: unknown): Promise<T> {
    const res = await fetch(this.getFullUrl(url), {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(toSnakeCase(body)),
      credentials: 'include',
    });
    return this.handleResponse<T>(res);
  }
}

export const apiClient = new ApiClient();
