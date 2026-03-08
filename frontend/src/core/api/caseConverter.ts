export const toCamelCase = (obj: unknown): unknown => {
  if (Array.isArray(obj)) {
    return obj.map((v) => toCamelCase(v));
  } else if (obj !== null && typeof obj === 'object' && obj.constructor === Object) {
    const record = obj as Record<string, unknown>;
    return Object.keys(record).reduce(
      (result, key) => ({
        ...result,
        [key.replace(/_([a-z])/g, (g) => g[1].toUpperCase())]: toCamelCase(record[key]),
      }),
      {},
    );
  }
  return obj;
};

export const toSnakeCase = (obj: unknown): unknown => {
  if (Array.isArray(obj)) {
    return obj.map((v) => toSnakeCase(v));
  } else if (obj !== null && typeof obj === 'object' && obj.constructor === Object) {
    const record = obj as Record<string, unknown>;
    return Object.keys(record).reduce(
      (result, key) => ({
        ...result,
        [key.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`)]: toSnakeCase(record[key]),
      }),
      {},
    );
  }
  return obj;
};
