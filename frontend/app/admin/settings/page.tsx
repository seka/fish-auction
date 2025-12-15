'use client';

import { useState } from 'react';
import { Box, Stack, Text, Button, Input } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function AdminSettingsPage() {
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setMessage(null);

        if (newPassword !== confirmPassword) {
            setMessage({ type: 'error', text: '新しいパスワードが一致しません。' });
            return;
        }

        if (newPassword.length < 8) {
            setMessage({ type: 'error', text: 'パスワードは8文字以上である必要があります。' });
            return;
        }

        setIsLoading(true);

        try {
            const res = await fetch('/api/proxy/api/admin/password', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    current_password: currentPassword,
                    new_password: newPassword,
                }),
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.error || 'パスワードの更新に失敗しました。');
            }

            setMessage({ type: 'success', text: 'パスワードを更新しました。' });
            setCurrentPassword('');
            setNewPassword('');
            setConfirmPassword('');
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box p="6">
            <Text as="h1" fontSize="2xl" fontWeight="bold" mb="6" className={css({ color: 'indigo.900' })}>
                設定
            </Text>

            <Box bg="white" p="6" borderRadius="lg" shadow="sm" maxW="md">
                <Text as="h2" fontSize="lg" fontWeight="semibold" mb="4" className={css({ color: 'gray.800' })}>
                    パスワード変更
                </Text>

                {message && (
                    <Box
                        p="3"
                        mb="4"
                        borderRadius="md"
                        bg={message.type === 'success' ? 'green.50' : 'red.50'}
                        color={message.type === 'success' ? 'green.700' : 'red.700'}
                        border="1px solid"
                        borderColor={message.type === 'success' ? 'green.200' : 'red.200'}
                    >
                        {message.text}
                    </Box>
                )}

                <form onSubmit={handleSubmit}>
                    <Stack spacing="4" alignItems="stretch">
                        <Box>
                            <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                現在のパスワード
                            </Text>
                            <Input
                                type="password"
                                value={currentPassword}
                                onChange={(e) => setCurrentPassword(e.target.value)}
                                required
                                className={css({ w: 'full' })}
                            />
                        </Box>

                        <Box>
                            <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                新しいパスワード
                            </Text>
                            <Input
                                type="password"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                                minLength={8}
                                className={css({ w: 'full' })}
                            />
                        </Box>

                        <Box>
                            <Text as="label" display="block" mb="1" fontSize="sm" fontWeight="medium" className={css({ color: 'gray.700' })}>
                                新しいパスワード（確認）
                            </Text>
                            <Input
                                type="password"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                                minLength={8}
                                className={css({ w: 'full' })}
                            />
                        </Box>

                        <Button
                            type="submit"
                            disabled={isLoading}
                            className={css({
                                bg: 'indigo.600',
                                color: 'white',
                                _hover: { bg: 'indigo.700' },
                                _disabled: { opacity: 0.6, cursor: 'not-allowed' }
                            })}
                        >
                            {isLoading ? '更新中...' : 'パスワードを変更する'}
                        </Button>
                    </Stack>
                </form>
            </Box>
        </Box>
    );
}
