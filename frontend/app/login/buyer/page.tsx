'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { loginBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';
import { Box, Text, Button, Input, Stack } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function BuyerLoginPage() {
    const [error, setError] = useState('');
    const router = useRouter();
    const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<BuyerLoginFormData>({
        resolver: zodResolver(buyerLoginSchema),
    });

    const onSubmit = async (data: BuyerLoginFormData) => {
        setError('');
        const buyer = await loginBuyer(data);
        if (buyer) {
            window.location.href = '/auctions';
        } else {
            setError('名前またはパスワードが間違っています');
        }
    };

    return (
        <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" py="12" px="4">
            <Box maxW="md" w="full">
                <Stack spacing="8">
                    <Box textAlign="center">
                        <Text as="h2" fontSize="3xl" fontWeight="extrabold" color="gray.900">
                            中買人ログイン
                        </Text>
                        <Text mt="2" fontSize="sm" color="gray.600">
                            セリに参加するにはログインしてください
                        </Text>
                    </Box>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <Stack spacing="6">
                            <Stack spacing="0">
                                <Box>
                                    <label htmlFor="email" className={css({ srOnly: true })}>メールアドレス</label>
                                    <Input
                                        id="email"
                                        type="email"
                                        {...register('email')}
                                        placeholder="メールアドレス"
                                        className={css({ borderBottomLeftRadius: '0', borderBottomRightRadius: '0' })}
                                    />
                                    {errors.email && <Text color="red.500" fontSize="xs" mt="1">{errors.email.message}</Text>}
                                </Box>
                                <Box>
                                    <label htmlFor="password" className={css({ srOnly: true })}>パスワード</label>
                                    <Input
                                        id="password"
                                        type="password"
                                        {...register('password')}
                                        placeholder="パスワード"
                                        className={css({ borderTopLeftRadius: '0', borderTopRightRadius: '0', borderTop: 'none' })}
                                    />
                                    {errors.password && <Text color="red.500" fontSize="xs" mt="1">{errors.password.message}</Text>}
                                </Box>
                            </Stack>

                            {error && <Text color="red.500" fontSize="sm" textAlign="center">{error}</Text>}

                            <Button
                                type="submit"
                                disabled={isSubmitting}
                                w="full"
                                className={css({ bg: 'indigo.600', _hover: { bg: 'indigo.700' }, color: 'white' })}
                            >
                                ログイン
                            </Button>
                            <Box textAlign="center">
                                <Link href="/signup" className={css({ fontSize: 'sm', color: 'indigo.600', _hover: { color: 'indigo.500' } })}>
                                    アカウントをお持ちでない方はこちら
                                </Link>
                            </Box>
                        </Stack>
                    </form>
                </Stack>
            </Box>
        </Box>
    );
}
