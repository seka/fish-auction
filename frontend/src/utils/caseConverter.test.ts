import { describe, it, expect } from 'vitest';
import { toCamelCase, toSnakeCase } from './caseConverter';

describe('CaseConverter', () => {
    describe('toCamelCase', () => {
        it('converts snake_case object keys to camelCase', () => {
            const input = { first_name: 'John', last_name: 'Doe' };
            const output = toCamelCase(input);
            expect(output).toEqual({ firstName: 'John', lastName: 'Doe' });
        });

        it('converts nested objects', () => {
            const input = { user_info: { user_id: 1 } };
            const output = toCamelCase(input);
            expect(output).toEqual({ userInfo: { userId: 1 } });
        });

        it('converts array of objects', () => {
            const input = [{ user_id: 1 }, { user_id: 2 }];
            const output = toCamelCase(input);
            expect(output).toEqual([{ userId: 1 }, { userId: 2 }]);
        });
    });

    describe('toSnakeCase', () => {
        it('converts camelCase object keys to snake_case', () => {
            const input = { firstName: 'John', lastName: 'Doe' };
            const output = toSnakeCase(input);
            expect(output).toEqual({ first_name: 'John', last_name: 'Doe' });
        });

        it('converts nested objects', () => {
            const input = { userInfo: { userId: 1 } };
            const output = toSnakeCase(input);
            expect(output).toEqual({ user_info: { user_id: 1 } });
        });
    });
});
