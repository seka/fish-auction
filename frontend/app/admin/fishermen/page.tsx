'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { fishermanSchema, FishermanFormData } from '@/src/models/schemas/admin';
import { useFishermen, useFishermanMutations } from './_hooks/useFisherman';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function FishermenPage() {
    const [message, setMessage] = useState('');

    const { fishermen, isLoading } = useFishermen();
    const { createFisherman, isCreating } = useFishermanMutations();

    const { register, handleSubmit, reset, formState: { errors } } = useForm<FishermanFormData>({
        resolver: zodResolver(fishermanSchema),
    });

    const onSubmit = async (data: FishermanFormData) => {
        try {
            await createFisherman({ name: data.name });
            setMessage('漁師を登録しました');
            reset();
        } catch (e) {
            console.error(e);
            setMessage('登録に失敗しました');
        }
    };

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" color="gray.800" mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                漁師管理
            </Text>

            {message && (
                <Box bg="blue.50" borderLeft="4px solid" borderColor="blue.500" color="blue.700" p="4" mb="8" borderRadius="sm" shadow="sm" role="alert">
                    <Text fontWeight="bold">通知</Text>
                    <Text>{message}</Text>
                </Box>
            )}

            <Box display="grid" gridTemplateColumns={{ base: '1fr', md: '3fr 1fr' }} gap="8" className={css({ md: { gridTemplateColumns: '1fr 2fr' } })}>
                {/* Form Section */}
                <Box className={css({ md: { gridColumn: '1 / 2' } })}>
                    <Card p="md" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" color="indigo.900" fontWeight="bold">
                                新規漁師登録
                            </Text>
                        </HStack>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" color="gray.700" mb="1">
                                        氏名
                                    </Text>
                                    <Input
                                        type="text"
                                        {...register('name')}
                                        placeholder="例: 山田 太郎"
                                        error={!!errors.name}
                                    />
                                    {errors.name && (
                                        <Text color="red.500" fontSize="sm" mt="1">{errors.name.message}</Text>
                                    )}
                                </Box>

                                <Button
                                    type="submit"
                                    disabled={isCreating}
                                    width="full"
                                    variant="primary"
                                >
                                    {isCreating ? '登録中...' : '登録する'}
                                </Button>
                            </Stack>
                        </form>
                    </Card>
                </Box>

                {/* List Section */}
                <Box className={css({ md: { gridColumn: '2 / 3' } })}>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text as="h2" variant="h4" color="gray.800" fontWeight="bold">漁師一覧</Text>
                        </Box>
                        {isLoading ? (
                            <Box p="6" textAlign="center" color="gray.500">読み込み中...</Box>
                        ) : fishermen.length === 0 ? (
                            <Box p="6" textAlign="center" color="gray.500">漁師が登録されていません</Box>
                        ) : (
                            <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
                                {fishermen.map((fisherman) => (
                                    <Box as="li" key={fisherman.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
                                        <HStack justify="between" align="center">
                                            <Box>
                                                <Text as="h3" fontSize="lg" fontWeight="bold" color="indigo.900">{fisherman.name}</Text>
                                                <Text fontSize="sm" color="gray.500" mt="1">ID: {fisherman.id}</Text>
                                            </Box>
                                        </HStack>
                                    </Box>
                                ))}
                            </Stack>
                        )}
                    </Card>
                </Box>
            </Box>
        </Box>
    );
}
