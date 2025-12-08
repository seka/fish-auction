'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { venueSchema, VenueFormData } from '@/src/models/schemas/auction';
import { useVenues, useVenueMutations } from './_hooks/useVenue';
import { Venue } from '@/src/models/venue';
import { Box, Stack, HStack, Text, Card, Button, Input } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { useTranslations } from 'next-intl';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

// Textarea component with similar styling to Input
const Textarea = styled('textarea', {
    base: {
        display: 'block',
        width: 'full',
        px: '3',
        py: '2',
        bg: 'white',
        color: 'gray.900',
        border: '1px solid',
        borderColor: 'gray.300',
        borderRadius: 'md',
        fontSize: 'sm',
        outline: 'none',
        transition: 'border-color 0.2s',
        resize: 'vertical',
        _focus: {
            borderColor: 'indigo.500',
            ring: '1px',
            ringColor: 'indigo.500',
        },
        _disabled: {
            bg: 'gray.50',
            cursor: 'not-allowed',
        },
    }
});

export default function AdminVenuesPage() {
    const t = useTranslations();
    const [message, setMessage] = useState('');
    const [editingVenue, setEditingVenue] = useState<Venue | null>(null);

    const { venues, isLoading } = useVenues();
    const { createVenue, updateVenue, deleteVenue, isCreating, isUpdating, isDeleting } = useVenueMutations();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<VenueFormData>({
        resolver: zodResolver(venueSchema),
    });

    const onSubmit = async (data: VenueFormData) => {
        try {
            if (editingVenue) {
                await updateVenue({ id: editingVenue.id, data });
                setMessage('会場を更新しました');
                setEditingVenue(null);
            } else {
                await createVenue(data);
                setMessage('会場を作成しました');
            }
            reset();
        } catch (e) {
            console.error(e);
            setMessage('エラーが発生しました');
        }
    };

    const onEdit = (venue: Venue) => {
        setEditingVenue(venue);
        setValue('name', venue.name);
        setValue('location', venue.location || '');
        setValue('description', venue.description || '');
    };

    const onCancelEdit = () => {
        setEditingVenue(null);
        reset();
    };

    const onDelete = async (id: number) => {
        if (confirm('本当に削除しますか？')) {
            try {
                await deleteVenue(id);
                setMessage('会場を削除しました');
            } catch (e) {
                console.error(e);
                setMessage('削除に失敗しました');
            }
        }
    };

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                会場管理
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
                    <Card padding="lg" className={css({ position: 'sticky', top: '6' })}>
                        <HStack mb="6">
                            <Box w="2" h="6" bg="indigo.500" mr="3" borderRadius="full" />
                            <Text as="h2" variant="h4" className={css({ color: 'indigo.900' })} fontWeight="bold">
                                {editingVenue ? '会場編集' : '新規会場登録'}
                            </Text>
                        </HStack>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <Stack spacing="4">
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        会場名
                                    </Text>
                                    <Input
                                        type="text"
                                        {...register('name')}
                                        placeholder="例: 豊洲市場"
                                        error={!!errors.name}
                                    />
                                    {errors.name && (
                                        <Text className={css({ color: 'red.500' })} fontSize="sm" mt="1">{errors.name.message}</Text>
                                    )}
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        所在地
                                    </Text>
                                    <Input
                                        type="text"
                                        {...register('location')}
                                        placeholder="例: 東京都江東区..."
                                    />
                                </Box>
                                <Box>
                                    <Text as="label" display="block" fontSize="sm" fontWeight="bold" className={css({ color: 'gray.700' })} mb="1">
                                        説明
                                    </Text>
                                    <Textarea
                                        {...register('description')}
                                        rows={3}
                                        placeholder="会場の説明..."
                                    />
                                </Box>

                                <HStack spacing="2">
                                    <Button
                                        type="submit"
                                        disabled={isCreating || isUpdating}
                                        width="full"
                                        className={css({ flex: '1' })}
                                        variant="primary"
                                    >
                                        {editingVenue ? (isUpdating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.update)) : (isCreating ? t(COMMON_TEXT_KEYS.loading) : t(COMMON_TEXT_KEYS.register))}
                                    </Button>
                                    {editingVenue && (
                                        <Button
                                            type="button"
                                            onClick={onCancelEdit}
                                            variant="outline"
                                        >
                                            {t(COMMON_TEXT_KEYS.cancel)}
                                        </Button>
                                    )}
                                </HStack>
                            </Stack>
                        </form>
                    </Card>
                </Box>

                {/* List Section */}
                <Box className={css({ md: { gridColumn: '2 / 3' } })}>
                    <Card padding="none" overflow="hidden">
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text as="h2" variant="h4" className={css({ color: 'gray.800' })} fontWeight="bold">会場一覧</Text>
                        </Box>
                        {isLoading ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>読み込み中...</Box>
                        ) : venues.length === 0 ? (
                            <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t(COMMON_TEXT_KEYS.no_data)}</Box>
                        ) : (
                            <Stack as="ul" spacing="0" divideY="1px" divideColor="gray.200">
                                {venues.map((venue) => (
                                    <Box as="li" key={venue.id} p="6" _hover={{ bg: 'gray.50' }} transition="colors">
                                        <HStack justify="between" align="start">
                                            <Box>
                                                <Text as="h3" fontSize="lg" fontWeight="bold" className={css({ color: 'indigo.900' })}>{venue.name}</Text>
                                                {venue.location && (
                                                    <Text fontSize="sm" className={css({ color: 'gray.700' })} mt="1" display="flex" alignItems="center">
                                                        <span className={css({ mr: '2' })}>📍</span>
                                                        {venue.location}
                                                    </Text>
                                                )}
                                                {venue.description && (
                                                    <Text fontSize="sm" className={css({ color: 'gray.600' })} mt="2">{venue.description}</Text>
                                                )}
                                            </Box>
                                            <HStack spacing="2">
                                                <Button size="sm" variant="outline" onClick={() => onEdit(venue)}>
                                                    {t(COMMON_TEXT_KEYS.edit)}
                                                </Button>
                                                <Button size="sm" className={css({ bg: 'red.50', color: 'red.600', _hover: { bg: 'red.100' } })} onClick={() => onDelete(venue.id)}>
                                                    {t(COMMON_TEXT_KEYS.delete)}
                                                </Button>
                                            </HStack>
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
