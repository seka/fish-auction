'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerLoginSchema, BuyerLoginFormData } from '@/src/models/schemas/buyer_auth';
import { loginBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';

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
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                <div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                        中買人ログイン
                    </h2>
                    <p className="mt-2 text-center text-sm text-gray-600">
                        セリに参加するにはログインしてください
                    </p>
                </div>
                <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
                    <div className="rounded-md shadow-sm -space-y-px">
                        <div>
                            <label htmlFor="email" className="sr-only">メールアドレス</label>
                            <input
                                id="email"
                                type="email"
                                {...register('email')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="メールアドレス"
                            />
                            {errors.email && <p className="text-red-500 text-xs mt-1">{errors.email.message}</p>}
                        </div>
                        <div>
                            <label htmlFor="password" className="sr-only">パスワード</label>
                            <input
                                id="password"
                                type="password"
                                {...register('password')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="パスワード"
                            />
                            {errors.password && <p className="text-red-500 text-xs mt-1">{errors.password.message}</p>}
                        </div>
                    </div>

                    {error && <div className="text-red-500 text-sm text-center">{error}</div>}

                    <div>
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                        >
                            ログイン
                        </button>
                    </div>
                    <div className="text-center">
                        <Link href="/signup" className="text-sm text-indigo-600 hover:text-indigo-500">
                            アカウントをお持ちでない方はこちら
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    );
}
