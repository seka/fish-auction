import { css, cva, type RecipeVariantProps } from 'styled-system/css';
import { styled } from 'styled-system/jsx';
import { ReactNode } from 'react';

// Box: 汎用コンテナ
const boxRecipe = cva({
    base: {},
});

export type BoxProps = RecipeVariantProps<typeof boxRecipe> &
    React.HTMLAttributes<HTMLDivElement>;

export const Box = styled('div', boxRecipe);

// Stack: 縦方向の配置
const stackRecipe = cva({
    base: {
        display: 'flex',
        flexDirection: 'column',
    },
    variants: {
        spacing: {
            '0': { gap: '0' },
            '1': { gap: '0.25rem' },
            '2': { gap: '0.5rem' },
            '3': { gap: '0.75rem' },
            '4': { gap: '1rem' },
            '6': { gap: '1.5rem' },
            '8': { gap: '2rem' },
        },
        align: {
            start: { alignItems: 'flex-start' },
            center: { alignItems: 'center' },
            end: { alignItems: 'flex-end' },
            stretch: { alignItems: 'stretch' },
        },
    },
    defaultVariants: {
        spacing: '4',
        align: 'stretch',
    },
});

export type StackProps = RecipeVariantProps<typeof stackRecipe> &
    React.HTMLAttributes<HTMLDivElement>;

export const Stack = styled('div', stackRecipe);

// HStack: 横方向の配置
const hstackRecipe = cva({
    base: {
        display: 'flex',
        flexDirection: 'row',
    },
    variants: {
        spacing: {
            '0': { gap: '0' },
            '1': { gap: '0.25rem' },
            '2': { gap: '0.5rem' },
            '3': { gap: '0.75rem' },
            '4': { gap: '1rem' },
            '6': { gap: '1.5rem' },
            '8': { gap: '2rem' },
        },
        align: {
            start: { alignItems: 'flex-start' },
            center: { alignItems: 'center' },
            end: { alignItems: 'flex-end' },
            stretch: { alignItems: 'stretch' },
        },
        justify: {
            start: { justifyContent: 'flex-start' },
            center: { justifyContent: 'center' },
            end: { justifyContent: 'flex-end' },
            between: { justifyContent: 'space-between' },
        },
    },
    defaultVariants: {
        spacing: '4',
        align: 'center',
        justify: 'start',
    },
});

export type HStackProps = RecipeVariantProps<typeof hstackRecipe> &
    React.HTMLAttributes<HTMLDivElement>;

export const HStack = styled('div', hstackRecipe);
