"use client";

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { requestPasswordReset, ResetPasswordRequest } from '@/src/api/auth_reset';
import { Box, Button, Text, Stack } from '@/src/core/ui';
import { css } from '@/styled-system/css';

export default function ForgotPasswordPage() {
    const t = useTranslations();
    const router = useRouter();
    const [isSubmitted, setIsSubmitted] = useState(false);
    const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<ResetPasswordRequest>();

    const onSubmit = async (data: ResetPasswordRequest) => {
        try {
            await requestPasswordReset(data);
            setIsSubmitted(true);
        } catch (error) {
            console.error('Failed to request password reset', error);
            setIsSubmitted(true); // Treat as success
        }
    };

    if (isSubmitted) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" p="4">
                <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
                    <Text variant="h3" textAlign="center" mb="6" className={css({ color: 'primary.600', fontWeight: 'bold' })}>
                        メール送信完了
                    </Text>
                    <Text textAlign="center" className={css({ color: 'gray.600' })}>
                        入力されたメールアドレスにパスワード再設定用のリンクを送信しました（登録がある場合）。
                        メールをご確認ください。
                    </Text>
                    <Button
                        mt="6"
                        w="full"
                        variant="secondary"
                        onClick={() => router.push('/login/buyer')}
                    >
                        ログイン画面に戻る
                    </Button>
                </Box>
            </Box>
        );
    }

    return (
        <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" p="4">
            <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
                <Text variant="h3" textAlign="center" mb="6" className={css({ color: 'primary.600', fontWeight: 'bold' })}>
                    パスワードをお忘れの方
                </Text>
                <Text mb="6" className={css({ color: 'gray.600', textAlign: 'center' })}>
                    登録したメールアドレスを入力してください。再設定用のリンクを送信します。
                </Text>

                <form onSubmit={handleSubmit(onSubmit)}>
                    <Stack gap="4">
                        <Box w="full">
                            <label className={css({ display: 'block', mb: '1.5', fontWeight: 'medium', color: 'gray.700' })}>
                                メールアドレス
                            </label>
                            <input
                                {...register('email', {
                                    required: "メールアドレスを入力してください",
                                    pattern: {
                                        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                                        message: "有効なメールアドレスを入力してください",
                                    }
                                })}
                                className={css({
                                    w: 'full',
                                    p: '2.5',
                                    border: '1px solid',
                                    borderColor: 'gray.300',
                                    borderRadius: 'md',
                                    _focus: { borderColor: 'primary.500', outline: 'none', ring: '2px', ringColor: 'primary.200' }
                                })}
                                placeholder="example@fish-auction.com"
                            />
                            {errors.email && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.email.message}</Text>
                            )}
                        </Box>

                        <Button
                            type="submit"
                            w="full"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? '送信中...' : '送信する'}
                        </Button>

                        <Button
                            variant="outline"
                            width="full"
                            onClick={() => router.push('/login/buyer')}
                            style={{ border: 'none' }}
                        >
                            キャンセル
                        </Button>
                    </Stack>
                </form>
            </Box>
        </Box>
    );
}
