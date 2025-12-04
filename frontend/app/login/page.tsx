'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useLogin } from './_hooks/useAuth';
import { loginSchema, LoginFormData } from '@/src/models/schemas/auth';
import { css } from 'styled-system/css';
import { Box, Text, Button, Input } from '@/src/core/ui';

export default function LoginPage() {
    const [error, setError] = useState('');
    const router = useRouter();
    const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
    });
    const { login, isLoading } = useLogin();

    const onSubmit = async (data: LoginFormData) => {
        setError('');

        const success = await login(data.password);

        if (success) {
            window.location.href = '/admin';
        } else {
            setError('パスワードが間違っています');
        }
    };

    return (
        <Box className={css({ minH: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', bg: 'gray.50', py: '12', px: { base: '4', sm: '6', lg: '8' } })}>
            <Box className={css({ maxW: 'md', w: 'full', spaceY: '8' })}>
                <Box>
                    <Text variant="h2" className={css({ textAlign: 'center', fontSize: '3xl', fontWeight: 'extrabold', color: 'gray.900', mt: '6' })}>
                        管理者ログイン
                    </Text>
                    <Text className={css({ mt: '2', textAlign: 'center', fontSize: 'sm', color: 'gray.600' })}>
                        管理画面へアクセスするにはパスワードを入力してください
                    </Text>
                </Box>
                <form className={css({ mt: '8', spaceY: '6' })} onSubmit={handleSubmit(onSubmit)}>
                    <Box className={css({ rounded: 'md', shadow: 'sm' })}>
                        <Box>
                            <label htmlFor="password" className={css({ srOnly: true })}>
                                パスワード
                            </label>
                            <Input
                                id="password"
                                type="password"
                                {...register('password')}
                                placeholder="パスワード"
                                error={!!errors.password}
                            />
                            {errors.password && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.password.message}</Text>
                            )}
                        </Box>
                    </Box>

                    {error && (
                        <Box className={css({ color: 'red.500', fontSize: 'sm', textAlign: 'center', fontWeight: 'bold' })}>
                            {error}
                        </Box>
                    )}

                    <Box>
                        <Button
                            type="submit"
                            disabled={isLoading}
                            className={css({ w: 'full' })}
                        >
                            {isLoading ? 'ログイン中...' : 'ログイン'}
                        </Button>
                    </Box>
                </form>
            </Box>
        </Box>
    );
}
