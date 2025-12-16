"use client";

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { verifyResetToken, confirmPasswordReset, ResetPasswordConfirmRequest } from '@/src/api/auth_reset';
import { Box, Button, Text, Stack } from '@/src/core/ui';
import { css } from '@/styled-system/css';

import { Suspense } from 'react';
// ... imports

function ResetPasswordForm() {
    const t = useTranslations();
    const router = useRouter();
    const searchParams = useSearchParams();
    const token = searchParams.get('token');

    const [isVerifying, setIsVerifying] = useState(true);
    const [isValidToken, setIsValidToken] = useState(false);
    const [isComplete, setIsComplete] = useState(false);

    const { register, handleSubmit, watch, formState: { errors, isSubmitting } } = useForm<Omit<ResetPasswordConfirmRequest, 'token'> & { confirm_password: string }>();
    const newPassword = watch('new_password');

    useEffect(() => {
        if (!token) {
            setIsVerifying(false);
            return;
        }

        const verify = async () => {
            try {
                await verifyResetToken({ token });
                setIsValidToken(true);
            } catch (error) {
                console.error("Invalid token", error);
                setIsValidToken(false);
            } finally {
                setIsVerifying(false);
            }
        };
        verify();
    }, [token]);

    const onSubmit = async (data: Omit<ResetPasswordConfirmRequest, 'token'> & { confirm_password: string }) => {
        if (!token) return;
        try {
            await confirmPasswordReset({ token, new_password: data.new_password });
            setIsComplete(true);
        } catch (error) {
            console.error('Failed to reset password', error);
            // Handle error (show message)
            alert("パスワードの再設定に失敗しました。トークンの有効期限が切れている可能性があります。");
        }
    };

    if (isVerifying) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
                <Text>検証中...</Text>
            </Box>
        );
    }

    if (!token || !isValidToken) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" p="4">
                <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
                    <Text variant="h3" textAlign="center" mb="6" className={css({ color: 'red.600', fontWeight: 'bold' })}>
                        無効なリンク
                    </Text>
                    <Text textAlign="center" className={css({ color: 'gray.600' })}>
                        このリンクは無効か、有効期限が切れています。
                        もう一度パスワード再設定リクエストを行ってください。
                    </Text>
                    <Button
                        mt="6"
                        w="full"
                        variant="secondary"
                        onClick={() => router.push('/login/forgot_password')}
                    >
                        再設定リクエストページへ
                    </Button>
                </Box>
            </Box>
        );
    }

    if (isComplete) {
        return (
            <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" p="4">
                <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
                    <Text variant="h3" textAlign="center" mb="6" className={css({ color: 'primary.600', fontWeight: 'bold' })}>
                        再設定完了
                    </Text>
                    <Text textAlign="center" className={css({ color: 'gray.600' })}>
                        パスワードの再設定が完了しました。
                        新しいパスワードでログインしてください。
                    </Text>
                    <Button
                        mt="6"
                        w="full"
                        onClick={() => router.push('/login/buyer')}
                    >
                        ログインページへ
                    </Button>
                </Box>
            </Box>
        );
    }

    return (
        <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50" p="4">
            <Box maxW="md" w="full" bg="white" p="8" borderRadius="lg" shadow="md">
                <Text variant="h3" textAlign="center" mb="6" className={css({ color: 'primary.600', fontWeight: 'bold' })}>
                    新しいパスワードの設定
                </Text>

                <form onSubmit={handleSubmit(onSubmit)}>
                    <Stack gap="4">
                        <Box w="full">
                            <label className={css({ display: 'block', mb: '1.5', fontWeight: 'medium', color: 'gray.700' })}>
                                新しいパスワード
                            </label>
                            <input
                                type="password"
                                {...register('new_password', {
                                    required: "パスワードを入力してください",
                                    minLength: { value: 8, message: "8文字以上で入力してください" }
                                })}
                                className={css({
                                    w: 'full', p: '2.5', border: '1px solid', borderColor: 'gray.300', borderRadius: 'md',
                                    _focus: { borderColor: 'primary.500', outline: 'none', ring: '2px', ringColor: 'primary.200' }
                                })}
                            />
                            {errors.new_password && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.new_password.message}</Text>
                            )}
                        </Box>

                        <Box w="full">
                            <label className={css({ display: 'block', mb: '1.5', fontWeight: 'medium', color: 'gray.700' })}>
                                パスワード（確認）
                            </label>
                            <input
                                type="password"
                                {...register('confirm_password', {
                                    validate: value => value === newPassword || "パスワードが一致しません"
                                })}
                                className={css({
                                    w: 'full', p: '2.5', border: '1px solid', borderColor: 'gray.300', borderRadius: 'md',
                                    _focus: { borderColor: 'primary.500', outline: 'none', ring: '2px', ringColor: 'primary.200' }
                                })}
                            />
                            {errors.confirm_password && (
                                <Text className={css({ color: 'red.500', fontSize: 'sm', mt: '1' })}>{errors.confirm_password.message}</Text>
                            )}
                        </Box>

                        <Button
                            type="submit"
                            w="full"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? '変更中...' : 'パスワードを変更する'}
                        </Button>
                    </Stack>
                </form>
            </Box>
        </Box>
    );
}

export default function ResetPasswordPage() {
    return (
        <Suspense fallback={<Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.100"><Text>読み込み中...</Text></Box>}>
            <ResetPasswordForm />
        </Suspense>
    );
}
