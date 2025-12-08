'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerSignupSchema, BuyerSignupFormData } from '@/src/models/schemas/buyer_auth';
import { buyerSchema, BuyerFormData } from '@/src/models/schemas/admin';
import { signupBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';
import { css } from 'styled-system/css';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { useTranslations } from 'next-intl';

export default function SignupPage() {
    const t = useTranslations();
    const [error, setError] = useState('');
    const [message, setMessage] = useState('');
    const router = useRouter();
    const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<BuyerSignupFormData>({
        resolver: zodResolver(buyerSignupSchema),
    });

    const onSubmit = async (data: BuyerSignupFormData) => {
        setError('');
        try {
            await signupBuyer(data);
            router.push('/login/buyer');
        } catch (e: any) {
            if (e.response && e.response.status === 409) {
                setError('登録に失敗しました。名前が既に使用されている可能性があります。');
            } else if (e.response && e.response.status >= 500) {
                setError('この操作の実行中にエラーが発生しました。運営にお問い合わせください');
            } else {
                setError('登録に失敗しました。入力内容をご確認ください。');
            }
        }
    };

    return (
        <Box className={css({ minH: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', bg: 'gray.50', py: '12', px: { base: '4', sm: '6', lg: '8' } })}>
            <Box className={css({ maxW: 'md', w: 'full', spaceY: '8' })}>
                <Box>
                    <Text variant="h2" className={css({ textAlign: 'center', fontSize: '3xl', fontWeight: 'extrabold', color: 'gray.900', mt: '6' })}>
                        中買人登録
                    </Text>
                    <Text className={css({ mt: '2', textAlign: 'center', fontSize: 'sm', color: 'gray.600' })}>
                        セリに参加するにはアカウントを作成してください
                    </Text>
                </Box>
                <form className={css({ mt: '8', spaceY: '6' })} onSubmit={handleSubmit(onSubmit)}>
                    <Box className={css({ rounded: 'md', shadow: 'sm' })}>
                        <Stack spacing="0" className={css({ '& > *': { position: 'relative', _focusWithin: { zIndex: 10 } } })}>
                            <Box>
                                <label htmlFor="name" className={css({ srOnly: true })}>名前</label>
                                <Input
                                    id="name"
                                    type="text"
                                    {...register('name')}
                                    className={css({ borderBottomLeftRadius: '0', borderBottomRightRadius: '0' })}
                                    placeholder="名前"
                                    error={!!errors.name}
                                />
                                {errors.name && <Text className={css({ color: 'red.500', fontSize: 'xs', mt: '1' })}>{errors.name.message}</Text>}
                            </Box>
                            <Box>
                                <label htmlFor="email" className={css({ srOnly: true })}>メールアドレス</label>
                                <Input
                                    id="email"
                                    type="email"
                                    {...register('email')}
                                    className={css({ borderRadius: '0', mt: '-1px' })}
                                    placeholder="メールアドレス"
                                    error={!!errors.email}
                                />
                                {errors.email && <Text className={css({ color: 'red.500', fontSize: 'xs', mt: '1' })}>{errors.email.message}</Text>}
                            </Box>
                            <Box>
                                <label htmlFor="organization" className={css({ srOnly: true })}>所属組織</label>
                                <Input
                                    id="organization"
                                    type="text"
                                    {...register('organization')}
                                    className={css({ borderRadius: '0', mt: '-1px' })}
                                    placeholder="所属組織"
                                    error={!!errors.organization}
                                />
                                {errors.organization && <Text className={css({ color: 'red.500', fontSize: 'xs', mt: '1' })}>{errors.organization.message}</Text>}
                            </Box>
                            <Box>
                                <label htmlFor="contact_info" className={css({ srOnly: true })}>連絡先</label>
                                <Input
                                    id="contact_info"
                                    type="text"
                                    {...register('contact_info')}
                                    className={css({ borderRadius: '0', mt: '-1px' })}
                                    placeholder="連絡先"
                                    error={!!errors.contact_info}
                                />
                                {errors.contact_info && <Text className={css({ color: 'red.500', fontSize: 'xs', mt: '1' })}>{errors.contact_info.message}</Text>}
                            </Box>
                            <Box>
                                <label htmlFor="password" className={css({ srOnly: true })}>パスワード</label>
                                <Input
                                    id="password"
                                    type="password"
                                    {...register('password')}
                                    className={css({ borderTopLeftRadius: '0', borderTopRightRadius: '0', mt: '-1px' })}
                                    placeholder="パスワード"
                                    error={!!errors.password}
                                />
                                {errors.password && <Text className={css({ color: 'red.500', fontSize: 'xs', mt: '1' })}>{errors.password.message}</Text>}
                            </Box>
                        </Stack>
                    </Box>

                    {error && <Box className={css({ color: 'red.500', fontSize: 'sm', textAlign: 'center' })}>{error}</Box>}

                    <Box>
                        <Button
                            type="submit"
                            disabled={isSubmitting}
                            className={css({ w: 'full' })}
                        >
                            {isSubmitting ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register)}
                        </Button>
                    </Box>
                    <Box className={css({ textAlign: 'center' })}>
                        <Link href="/login/buyer" className={css({ fontSize: 'sm', color: 'primary.600', _hover: { color: 'primary.500' } })}>
                            すでにアカウントをお持ちの方はこちら
                        </Link>
                    </Box>
                </form>
            </Box>
        </Box>
    );
}
