'use client';

import { UseFormReturn } from 'react-hook-form';
import { BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { Box, Text, Stack, Input, Button } from '@atoms';
import { css } from 'styled-system/css';

interface BuyerLoginFormProps {
  loginForm: UseFormReturn<BuyerLoginFormData>;
  onSubmit: (data: BuyerLoginFormData) => Promise<void>;
  loginError: string;
  t: (key: string) => string;
}

export const BuyerLoginForm = ({ loginForm, onSubmit, loginError, t }: BuyerLoginFormProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = loginForm;

  return (
    <Box
      minH="screen"
      display="flex"
      alignItems="center"
      justifyContent="center"
      bg="gray.50"
      py="12"
      px="4"
    >
      <Box maxW="md" w="full">
        <Stack spacing="8">
          <Box textAlign="center">
            <Text
              as="h2"
              fontSize="3xl"
              fontWeight="extrabold"
              className={css({ color: 'gray.900' })}
            >
              {t('Public.AuctionDetail.login_title')}
            </Text>
            <Text mt="2" fontSize="sm" className={css({ color: 'gray.700' })}>
              {t('Public.AuctionDetail.login_description')}
            </Text>
          </Box>
          <form onSubmit={handleSubmit(onSubmit)}>
            <Stack spacing="6">
              <Stack spacing="0">
                <Box>
                  <label htmlFor="email" className={css({ srOnly: true })}>
                    {t('Common.email')}
                  </label>
                  <Input
                    id="email"
                    type="email"
                    {...register('email')}
                    placeholder={t('Common.email')}
                    className={css({ borderBottomLeftRadius: '0', borderBottomRightRadius: '0' })}
                  />
                  {errors.email && (
                    <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">
                      {errors.email.message}
                    </Text>
                  )}
                </Box>
                <Box>
                  <label htmlFor="password" className={css({ srOnly: true })}>
                    {t('Common.password')}
                  </label>
                  <Input
                    id="password"
                    type="password"
                    {...register('password')}
                    placeholder={t('Common.password')}
                    className={css({
                      borderTopLeftRadius: '0',
                      borderTopRightRadius: '0',
                      borderTop: 'none',
                    })}
                  />
                  {errors.password && (
                    <Text className={css({ color: 'red.500' })} fontSize="xs" mt="1">
                      {errors.password.message}
                    </Text>
                  )}
                </Box>
              </Stack>

              {loginError && (
                <Text className={css({ color: 'red.500' })} fontSize="sm" textAlign="center">
                  {loginError}
                </Text>
              )}

              <Button
                type="submit"
                disabled={isSubmitting}
                width="full"
                className={css({
                  bg: 'indigo.600',
                  _hover: { bg: 'indigo.700' },
                  color: 'white',
                })}
              >
                {t('Public.Login.submit')}
              </Button>
            </Stack>
          </form>
        </Stack>
      </Box>
    </Box>
  );
};
