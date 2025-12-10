'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { loginBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';
import { Box, Text, Button, Input, Stack, Card } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

export default function BuyerLoginPage() {
    const [isLoading, setIsLoading] = useState(false);
    const [loginError, setLoginError] = useState('');
    const router = useRouter();
    const t = useTranslations();
    const { register, handleSubmit, formState: { errors } } = useForm<BuyerLoginFormData>({
        resolver: zodResolver(buyerLoginSchema),
    });

    const onSubmit = async (data: BuyerLoginFormData) => {
        setIsLoading(true);
        setLoginError('');

        try {
            const buyer = await loginBuyer(data);
            if (buyer) {
                router.push('/auctions');
            } else {
                setLoginError(t('Public.Login.error_credentials'));
            }
        } catch {
            setLoginError(t('Public.Login.error_credentials'));
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box display="flex" minH="screen" alignItems="center" justifyContent="center" bg="gray.50">
            <Card width="full" maxW="md" padding="lg" shadow="md">
                <Box textAlign="center" mb="8">
                    <Stack spacing="4">
                        <Text as="h1" variant="h3" fontWeight="bold" className={css({ color: 'blue.700' })}>
                            {t('Public.Login.title')}
                        </Text>
                        <Text className={css({ color: 'gray.600' })}>
                            {t('Public.Login.description')}
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
                            {errors.email && <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">{errors.email.message}</Text>}
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
                            {errors.password && <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">{errors.password.message}</Text>}
                        </Box>
                        {loginError && (
                            <Box bg="red.50" p="3" borderRadius="md">
                                <Text className={css({ color: 'red.600', fontSize: 'sm', textAlign: 'center' })}>{loginError}</Text>
                            </Box>
                        )}
                        <Button
                            type="submit"
                            width="full"
                            size="lg"
                            disabled={isLoading}
                            variant="primary"
                        >
                            {isLoading ? t('Common.loading') : t('Public.Login.submit')}
                        </Button>
                        <Box textAlign="center">
                            <Link href="/signup" className={css({ fontSize: 'sm', color: 'blue.600', _hover: { textDecoration: 'underline' }, cursor: 'pointer' })}>
                                {t('Public.Login.signup_link')}
                            </Link>
                        </Box>
                    </Stack>
                </form>
            </Card>
        </Box>
    );
}
