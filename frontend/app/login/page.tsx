'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useLogin } from './_hooks/useAuth';
import { loginSchema, LoginFormData } from '@/src/models/schemas/auth';
import { css } from 'styled-system/css';
import { Box, Text, Button, Input, Card, Stack } from '@/src/core/ui';
import { useTranslations } from 'next-intl';

export default function LoginPage() {
    const [error, setError] = useState('');
    const router = useRouter();
    const t = useTranslations();
    const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
    });
    const { login, isLoading } = useLogin();

    const onSubmit = async (data: LoginFormData) => {
        setError('');

        const success = await login(data.email, data.password);

        if (success) {
            router.push('/admin');
        } else {
            setError(t('Admin.Login.error_invalid_password'));
        }
    };

    return (
        <Box display="flex" minH="screen" alignItems="center" justifyContent="center" bg="gray.100">
            <Card width="full" maxW="md" padding="lg" shadow="lg">
                <Box textAlign="center" mb="8">
                    <Stack spacing="4">
                        <Text as="h1" variant="h3" fontWeight="bold" className={css({ color: 'indigo.700' })}>
                            {t('Admin.Login.title')}
                        </Text>
                        <Text className={css({ color: 'gray.600' })}>
                            {t('Admin.Login.description')}
                        </Text>
                    </Stack>
                </Box>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Stack spacing="6">
                        <Box>
                            <label htmlFor="email" className={css({ srOnly: true })}>{t('Common.email')}</label>
                            <Input
                                id="email"
                                type="email"
                                {...register('email')}
                                placeholder={t('Common.email')}
                                bg="white"
                            />
                            {errors.email && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.email.message}</Text>
                            )}
                        </Box>
                        <Box>
                            <label htmlFor="password" className={css({ srOnly: true })}>{t('Common.password')}</label>
                            <Input
                                id="password"
                                type="password"
                                {...register('password')}
                                placeholder={t('Common.password')}
                                bg="white"
                            />
                            {errors.password && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.password.message}</Text>
                            )}
                        </Box>
                        {error && (
                            <Box bg="red.50" p="3" borderRadius="md">
                                <Text className={css({ color: 'red.600', fontSize: 'sm', textAlign: 'center' })}>{error}</Text>
                            </Box>
                        )}
                        <Button
                            type="submit"
                            width="full"
                            size="lg"
                            disabled={isLoading}
                            variant="primary"
                        >
                            {isLoading ? t('Admin.Login.logging_in') : t('Common.submit')}
                        </Button>
                        <Box textAlign="center">
                            <Link href="/login/admin/forgot_password" className={css({ fontSize: 'sm', color: 'gray.500', _hover: { textDecoration: 'underline' }, cursor: 'pointer', display: 'block', mb: '2' })}>
                                パスワードをお忘れの方はこちら
                            </Link>
                        </Box>
                    </Stack>
                </form>
            </Card>
        </Box>
    );
}
