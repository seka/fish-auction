import { cva, type RecipeVariantProps } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

const buttonRecipe = cva({
    base: {
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontWeight: 'semibold',
        cursor: 'pointer',
        transition: 'all 0.2s',
        borderRadius: '0.5rem',
        outline: 'none',
        _disabled: {
            opacity: 0.5,
            cursor: 'not-allowed',
        },
    },
    variants: {
        variant: {
            primary: {
                bg: 'primary.600',
                color: 'white',
                _hover: {
                    bg: 'primary.700',
                },
            },
            secondary: {
                bg: 'secondary.600',
                color: 'white',
                _hover: {
                    bg: 'secondary.700',
                },
            },
            outline: {
                bg: 'transparent',
                border: '1px solid',
                borderColor: 'gray.300',
                color: 'gray.700',
                _hover: {
                    bg: 'gray.50',
                },
            },
        },
        size: {
            sm: {
                px: '3',
                py: '1.5',
                fontSize: '0.75rem', // 12px
            },
            md: {
                px: '4',
                py: '2',
                fontSize: '0.875rem', // 14px
            },
            lg: {
                px: '6',
                py: '3',
                fontSize: '1rem', // 16px
            },
        },
    },
    defaultVariants: {
        variant: 'primary',
        size: 'md',
    },
});

export type ButtonProps = RecipeVariantProps<typeof buttonRecipe> &
    React.ButtonHTMLAttributes<HTMLButtonElement>;

export const Button = styled('button', buttonRecipe);
